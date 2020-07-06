package md_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/apprentice3d/forge-api-go-client/dm"
	"github.com/apprentice3d/forge-api-go-client/md"
)

func TestAPI_TranslateToSVF(t *testing.T) {
	// prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		t.Skipf("No Forge credentials present; skipping test")
	}

	bucketAPI := dm.NewBucketAPIWithCredentials(clientID, clientSecret)
	mdAPI := md.NewAPIWithCredentials(clientID, clientSecret)

	tempBucketName := "go_testing_md_bucket"
	testFilePath := "../assets/HelloWorld.rvt"

	var testObject dm.ObjectDetails

	t.Run("Create a temporary bucket", func(t *testing.T) {
		_, err := bucketAPI.CreateBucket(tempBucketName, "transient")

		if err != nil {
			t.Errorf("Failed to create a bucket: %s\n", err.Error())
		}
	})

	t.Run("Get bucket details", func(t *testing.T) {
		_, err := bucketAPI.GetBucketDetails(tempBucketName)

		if err != nil {
			t.Fatalf("Failed to get bucket details: %s\n", err.Error())
		}
	})

	t.Run("Upload an object into temp bucket", func(t *testing.T) {
		file, err := os.Open(testFilePath)
		if err != nil {
			t.Fatal("Cannot open testfile for reading")
		}
		defer file.Close()
		data, err := ioutil.ReadAll(file)
		if err != nil {
			t.Fatal("Cannot read the testfile")
		}

		testObject, err = bucketAPI.UploadObject(tempBucketName, "temp_file.rvt", data)

		if err != nil {
			t.Fatal("Could not upload the test object, got: ", err.Error())
		}

		if testObject.Size == 0 {
			t.Fatal("The test object was uploaded but it is zero-sized")
		}
	})

	t.Run("Translate object into SVF", func(t *testing.T) {

		result, err := mdAPI.TranslateToSVF(testObject.ObjectID)

		if err != nil {
			t.Error("Could not translate the test object, got: ", err.Error())
		}

		if result.Result == "created" {
			t.Error("The test object was uploaded, but failed to create the translation job")
		}
	})

	t.Run("Delete the temporary bucket", func(t *testing.T) {
		err := bucketAPI.DeleteBucket(tempBucketName)

		if err != nil {
			t.Fatalf("Failed to delete bucket: %s\n", err.Error())
		}
	})
}

func TestAPI_TranslateToSVF2_JSON_Creation(t *testing.T) {

	params := md.TranslationSVFPreset
	params.Input.URN = base64.RawStdEncoding.EncodeToString([]byte("just a test urn"))

	output, err := json.Marshal(&params)
	if err != nil {
		t.Fatal("Could not marshal the preset into JSON: ", err.Error())
	}

	referenceExample := `
{
        "input": {
          "urn": "anVzdCBhIHRlc3QgdXJu"
        },
        "output": {
			"destination": {
        		"region": "us"
      		},
          	"formats": [
            {
              "type": "svf",
              "views": [
                "2d",
                "3d"
              ]
            }
          ]
        }
      }
`

	var example md.TranslationParams
	err = json.Unmarshal([]byte(referenceExample), &example)
	if err != nil {
		t.Fatal("Could not unmarshal the reference example: ", err.Error())
	}

	expected, err := json.Marshal(example)
	if err != nil {
		t.Fatal("Could not marshal the reference example into JSON: ", err.Error())
	}

	if bytes.Compare(expected, output) != 0 {
		t.Fatalf("The translation params are not correct:\nexpected: %s\n created: %s",
			string(expected),
			string(output))

	}

}
