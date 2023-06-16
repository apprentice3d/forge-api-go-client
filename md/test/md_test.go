package md_test

/*
package md_test provides "blackbox" tests for the md package.
These tests are meant to test the public API of the md package.
*/

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/woweh/forge-api-go-client/dm"
	"github.com/woweh/forge-api-go-client/md"
	"github.com/woweh/forge-api-go-client/oauth"
)

/*
NOTE:
- Buckets can only be deleted by the user who created them.
  => You might want to change the bucketKey if the bucket already exists.

- A bucketKey (= bucket name) must be globally unique across all applications and regions

- You can only run these tests when you have a valid client ID and secret.
  => You probably want to run the tests locally, with your own credentials.
*/

func TestModelDerivativeAPI_HappyPath_AllFunctions_Default_US(t *testing.T) {

	// prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

	// check client ID and secret
	if clientID == "" || clientSecret == "" {
		t.Skip("Skipping tests because FORGE_CLIENT_ID and/or FORGE_CLIENT_SECRET env variables are not set")
	}

	authenticator := oauth.NewTwoLegged(clientID, clientSecret)
	bucketAPI := dm.NewBucketAPI(authenticator)
	mdAPI := md.NewMDAPI(authenticator)

	testFilePath := "../../assets/HelloWorld.rvt"

	tempBucketName := "forge_api_go_client_unit_testing_happy_path_default_us"

	var uploadResult dm.UploadResult
	var translationResult md.TranslationResult
	var manifest md.Manifest
	var masterViewGuid string

	t.Run(
		"Create a temporary bucket", func(t *testing.T) {
			bucketExists, err := bucketAPI.BucketExists(tempBucketName)
			if err != nil {
				t.Errorf("Failed to check if bucket exists: %s\n", err.Error())
			}
			if bucketExists {
				t.Skip("Bucket already exists, skipping bucket creation")
			}

			_, err = bucketAPI.CreateBucket(tempBucketName, "transient")
			if err != nil {
				t.Errorf("Failed to create a bucket: %s\n", err.Error())
			}
		},
	)

	t.Run(
		"Get bucket details", func(t *testing.T) {
			_, err := bucketAPI.GetBucketDetails(tempBucketName)

			if err != nil {
				t.Errorf("Failed to get bucket details: %s\n", err.Error())
			}
		},
	)

	t.Run(
		"Upload object into bucket", func(t *testing.T) {
			file, err := os.Open(testFilePath)
			if err != nil {
				t.Error("Cannot open test file for reading")
			}
			defer file.Close()

			uploadResult, err = bucketAPI.UploadObject(tempBucketName, "temp_file.rvt", testFilePath)
			if err != nil {
				t.Error("Could not upload the test object, got: ", err.Error())
			}

			if uploadResult.Size == 0 {
				t.Error("The test object was uploaded but it is zero-sized")
			}
		},
	)

	t.Run(
		"Translate object into SVF", func(t *testing.T) {
			var err error
			translationResult, err = mdAPI.TranslateToSVF(uploadResult.ObjectId)

			if err != nil {
				t.Error("Could not translate the test object, got: ", err.Error())
			}

			if translationResult.Result != "created" {
				t.Error("The test object was uploaded, but failed to create the translation job")
			}
		},
	)

	t.Run(
		"Get manifest of the object", func(t *testing.T) {

			timeToWait := 5 * time.Second
			translating := true

			for translating {
				manifest, err := mdAPI.GetManifest(translationResult.URN)
				if err != nil {
					t.Errorf("Problems getting the manifest for %s: %s", translationResult.URN, err.Error())
				}

				if manifest.Type != "manifest" {
					t.Error("Expecting 'manifest' type, got ", manifest.Type)
				}

				if manifest.URN != translationResult.URN {
					t.Errorf("URN not matching: translation=%s\tmanifest=%s", translationResult.URN, manifest.URN)
				}

				switch manifest.Status {
				case md.StatusPending:
					t.Log("Translation pending...")
					time.Sleep(timeToWait)
				case md.StatusInProgress:
					t.Log("Translation in progress...")
					time.Sleep(timeToWait)
				case md.StatusSuccess:
					translating = false
				case md.StatusFailed:
					t.Error("Translation failed")
				case md.StatusTimeout:
					t.Error("Translation timed out")
				default:
					t.Errorf("Got unexpected status: %s", manifest.Status)
				}
			}

			if len(manifest.Derivatives) != 2 {
				t.Errorf("Expecting to have 2 derivative, got %d", len(manifest.Derivatives))
			}

			outputType := manifest.Derivatives[0].OutputType
			if outputType != "svf" {
				t.Errorf("Expecting first derivative to be 'svf', got %s", outputType)
			}
		},
	)

	t.Run(
		"Download the properties database URN", func(t *testing.T) {
			propertiesDatabaseUrn := manifest.GetPropertiesDatabaseUrn()
			if propertiesDatabaseUrn == "" {
				t.Error("Expecting a non-empty URN")
			}

			_, err := mdAPI.GetDerivative(manifest.URN, propertiesDatabaseUrn)
			if err != nil {
				t.Error("Failed to download the properties database, got: ", err.Error())
			}
		},
	)

	t.Run(
		"Download metadata and get master view GUID", func(t *testing.T) {
			metaData, err := mdAPI.GetMetadata(manifest.URN, md.DefaultXAdsHeaders())
			if err != nil {
				t.Error("Failed to download the metadata, got: ", err.Error())
			}

			if metaData.Data.Type != "metadata" {
				t.Error("Expecting 'metadata' result type, got ", metaData.Data.Type)
			}

			masterViewGuid = metaData.GetMasterModelViewGuid()
			if masterViewGuid == "" {
				t.Error("Expecting a non-empty master view GUID")
			}
		},
	)

	t.Run(
		"Download all properties", func(t *testing.T) {
			bytes, err := mdAPI.GetModelViewProperties(manifest.URN, masterViewGuid, md.DefaultXAdsHeaders())
			if err != nil {
				t.Error("Failed to download the properties, got: ", err.Error())
			}

			if len(bytes) == 0 {
				t.Error("Properties data (byte array) empty")
			}

			// convert the bytes to JSON
			jsonProperties, err := json.Marshal(string(bytes))
			if err != nil {
				t.Error("Failed to convert the properties to JSON, got: ", err.Error())
			}

			if len(jsonProperties) == 0 {
				t.Error("Properties data (JSON) empty")
			}
		},
	)

	t.Run(
		"Download object tree", func(t *testing.T) {
			tree, err := mdAPI.GetObjectTree(manifest.URN, masterViewGuid, true, md.DefaultXAdsHeaders())
			if err != nil {
				t.Error("Failed to download the object tree, got: ", err.Error())
			}

			if tree.Data.Type != "objects" {
				t.Error("Expecting 'objects' result type, got ", tree.Data.Type)
			}

			if len(tree.Data.Objects) == 0 {
				t.Error("Object tree is empty")
			}
		},
	)

	t.Cleanup(
		func() {
			if ok, _ := bucketAPI.BucketExists(tempBucketName); ok {
				t.Log("Try to delete the temporary bucket...")
				err := bucketAPI.DeleteBucket(tempBucketName)
				if err != nil {
					t.Errorf("Failed to delete bucket: %s\n", err.Error())
				}
			}

		},
	)
}
