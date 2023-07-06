package md_test

/*
package md_test provides "blackbox" tests for the md package.
These tests are meant to test the public API of the md package.
*/

import (
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	"github.com/woweh/forge-api-go-client"
	"github.com/woweh/forge-api-go-client/dm"
	"github.com/woweh/forge-api-go-client/md"
	"github.com/woweh/forge-api-go-client/oauth"
)

/*
NOTE:
- You can only run these tests when you have a valid client ID and secret.
  => You probably want to run the tests locally, with your own credentials.
- A bucketKey (= bucket name) must be globally unique across all applications and regions
- Rules for bucketKey names: -_.a-z0-9 (between 3-128 characters in length)
- Buckets can only be deleted by the user who created them.
  => You might want to change the bucketKey if the bucket already exists.
- A bucket name will not be immediately available for reuse after deletion.
  => Best use a unique bucket name for each subtest.
  => You can also use a timestamp to make sure the bucket name is unique.
*/

const (
	testFilePath = "../../dm/assets/rst_basic_sample_project.rvt"
)

var (
	backoffSchedule = []time.Duration{
		1 * time.Second,
		3 * time.Second,
		7 * time.Second,
		15 * time.Second,
		31 * time.Second,
	}
	usTestsFailed   bool
	emeaTestsFailed bool
)

func TestMain(m *testing.M) {

	log.Println("In TestMain()...")

	usTestsFailed = false
	emeaTestsFailed = false

	log.Println("Running tests...")
	exitCode := m.Run()

	log.Println("Tests finished, determining exit code")
	log.Println("- usTestsFailed: ", usTestsFailed)
	log.Println("- emeaTestsFailed: ", emeaTestsFailed)
	if usTestsFailed || emeaTestsFailed {
		exitCode = 1
	} else {
		exitCode = 0
	}

	log.Println("Exiting with code: ", exitCode)
	os.Exit(exitCode)
}

func TestModelDerivativeAPI_HappyPath_AllFunctions(t *testing.T) {

	type args struct {
		region         forge.Region
		tempBucketName string
		objectName     string
		error          *bool
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "Default-US",
			args: args{
				region:         forge.US,
				tempBucketName: "forge_api_go_client_unit_testing_happy_path_default_us",
				objectName:     "rst_basic_sample_project_us.rvt",
				error:          &usTestsFailed,
			},
		},
		{
			name: "EMEA",
			args: args{
				region:         forge.EMEA,
				tempBucketName: "forge_api_go_client_unit_testing_happy_path_default_emea",
				objectName:     "rst_basic_sample_project_emea.rvt",
				error:          &emeaTestsFailed,
			},
		},
	}

	// run the tests
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {

				trueVal := true

				// prepare the credentials
				clientID := os.Getenv("FORGE_CLIENT_ID")
				clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

				// check client ID and secret
				if clientID == "" || clientSecret == "" {
					t.Skip("Skipping tests because FORGE_CLIENT_ID and/or FORGE_CLIENT_SECRET env variables are not set")
				}

				authenticator := oauth.NewTwoLegged(clientID, clientSecret)
				_, err := authenticator.GetToken("bucket:create bucket:read data:read data:create data:write")
				if err != nil {
					// can't continue if we can't get a token
					tt.args.error = &trueVal
					t.Fatalf("Failed to get token: %s\n", err.Error())
				}

				ossAPI := dm.NewOssApi(authenticator, tt.args.region)
				mdAPI := md.NewMdApi(authenticator, tt.args.region)

				t.Log("Checking if bucket already exists...")
				bucketDetails, err := ossAPI.GetBucketDetails(tt.args.tempBucketName)
				if err == nil {

					t.Log("Bucket exists, no need to create it...")
					t.Log(bucketDetails)

				} else {

					t.Log("Bucket does not exist, try creating it...")
					_, err = ossAPI.CreateBucket(tt.args.tempBucketName, dm.PolicyTransient)
					if err != nil {
						// can't continue if bucket creation fails
						tt.args.error = &trueVal
						t.Fatalf("Failed to create a bucket: %s\n", err.Error())
					}

					t.Log("Verify that bucket exists...")
					for _, backoff := range backoffSchedule {
						bucketDetails, err = ossAPI.GetBucketDetails(tt.args.tempBucketName)
						if err != nil {
							t.Logf("Failed to get bucket details: %s\n", err.Error())
							t.Log("Trying again...")
							time.Sleep(backoff)
						} else {
							t.Log(bucketDetails)
							break
						}
					}
					if err != nil {
						// can't continue if bucket creation failed
						t.Log("Bucket does not exist, even after waiting for it to be created")
						tt.args.error = &trueVal
						t.Fatalf("Failed to get bucket details: %s\n", err.Error())
					}
				}

				t.Log("Checking if test file exists...")
				file, err := os.Open(testFilePath)
				if err != nil {
					// can't continue if file cannot be opened/found
					tt.args.error = &trueVal
					t.Fatal("Cannot open test file for reading")
				}
				defer file.Close()

				t.Log("Uploading test object...")
				uploadResult, err := ossAPI.UploadObject(tt.args.tempBucketName, tt.args.objectName, testFilePath)
				if err != nil {
					// can't continue if upload fails
					tt.args.error = &trueVal
					t.Fatal("Could not upload the test object, got: ", err.Error())
				}
				if uploadResult.Size == 0 {
					// can't continue if upload fails
					tt.args.error = &trueVal
					t.Fatal("The test object was uploaded but it is zero-sized")
				}
				t.Log("Uploaded object details: ", uploadResult)

				t.Log("Creating translation job...")
				params := mdAPI.DefaultTranslationParams(md.UrnFromObjectId(uploadResult.ObjectId))
				translationJob, err := mdAPI.StartTranslation(params, md.DefaultXAdsHeaders())
				if err != nil {
					// can't continue if translation job creation fails
					tt.args.error = &trueVal
					t.Fatal("Could not create the translation job, got: ", err.Error())
				}
				t.Log("Translation job: ", translationJob)
				if translationJob.Result != "created" && translationJob.Result != "success" {
					// can't continue if translation job creation fails
					tt.args.error = &trueVal
					t.Fatal(
						"The the translation job result is neither \"created\" nor \"success\": ",
						translationJob.Result,
					)
				}

				// make this a fixed value for now, to avoid golang test timeouts
				timeToWait := time.Duration(5) * time.Second

				t.Log("Initial wait for the translation to get started...")
				time.Sleep(timeToWait)

				var manifest md.Manifest

				seconds := 0
				timeout := float64(60 * 60) // 1 hour
				startTime := time.Now()
				errorCount := 0

			loopUntilTranslationIsFinished:
				for time.Since(startTime).Seconds() < timeout && manifest.Status != md.StatusSuccess {
					seconds++

					t.Log("Getting manifest...")
					manifest, err = mdAPI.GetManifest(translationJob.URN)
					if err != nil {
						errorCount++
						if errorCount > 10 {
							t.Errorf("Too many errors getting the manifest for %s: %s", translationJob.URN, err.Error())
						} else {
							t.Logf("Problems getting the manifest for %s: %s", translationJob.URN, err.Error())
							t.Log("Waiting a bit and trying again...")
							time.Sleep(timeToWait)
							continue loopUntilTranslationIsFinished
						}
					}

					switch manifest.Status {
					case md.StatusPending:
						t.Log("Translation pending...")
						time.Sleep(timeToWait)
						continue loopUntilTranslationIsFinished

					case md.StatusInProgress:
						t.Logf("Translation in progress: %s", manifest.Progress)
						time.Sleep(timeToWait)

					case md.StatusSuccess:
						t.Log("Translation completed")
						// break out of the loop
						break loopUntilTranslationIsFinished

					case md.StatusFailed:
						// can't continue if translation failed
						tt.args.error = &trueVal
						t.Fatal("Translation failed")

					case md.StatusTimeout:
						// can't continue if translation timed out
						tt.args.error = &trueVal
						t.Fatal("Translation timed out")

					default:
						t.Errorf("Got unexpected status: %s", manifest.Status)
					}
				}

				if manifest.Type != "manifest" {
					t.Error("Expecting 'manifest' type, got ", manifest.Type)
				}

				if manifest.URN != translationJob.URN {
					// can't continue if URN doesn't match
					tt.args.error = &trueVal
					t.Fatalf("URN not matching: translation=%s\tmanifest=%s", translationJob.URN, manifest.URN)
				}

				if len(manifest.Derivatives) != 2 {
					t.Errorf("Expecting to have 2 derivative, got %d", len(manifest.Derivatives))
				}

				outputType := manifest.Derivatives[0].OutputType
				if outputType != "svf" {
					t.Errorf("Expecting first derivative to be 'svf', got %s", outputType)
				}

				t.Log("Getting properties database URN...")
				propertiesDatabaseUrn := manifest.GetPropertiesDatabaseUrn()
				if propertiesDatabaseUrn == "" {
					t.Error("Expecting a non-empty URN")
				}

				t.Log("Downloading properties database...")
				_, err = mdAPI.GetDerivative(manifest.URN, propertiesDatabaseUrn)
				if err != nil {
					t.Error("Failed to download the properties database, got: ", err.Error())
				}

				t.Log("Downloading metadata...")
				metaData, err := mdAPI.GetMetadata(manifest.URN)
				if err != nil {
					// can't continue if metadata download fails
					tt.args.error = &trueVal
					t.Fatal("Failed to download the metadata, got: ", err.Error())
				}

				if metaData.Data.Type != "metadata" {
					t.Error("Expecting 'metadata' result type, got ", metaData.Data.Type)
				}

				masterViewGuid := metaData.GetMasterModelViewGuid()
				if masterViewGuid == "" {
					// can't continue if master view GUID is empty
					tt.args.error = &trueVal
					t.Fatal("Expecting a non-empty master view GUID")
				}

				t.Log("Downloading all properties for master view: ", masterViewGuid)
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

				t.Log("Downloading object tree for master view: ", masterViewGuid)
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

				t.Cleanup(
					func() {
						t.Log("Try to delete the temporary bucket...")
						err = ossAPI.DeleteBucket(tt.args.tempBucketName)
						if err != nil {
							t.Logf("Failed to delete bucket: %s\n", err.Error())
						}
					},
				)
			},
		)
	}
}
