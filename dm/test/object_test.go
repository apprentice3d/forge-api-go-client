package dm_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/woweh/forge-api-go-client"
	"github.com/woweh/forge-api-go-client/oauth"

	"github.com/woweh/forge-api-go-client/dm"
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
	objectKey    string = "rst_basic_sample_project.rvt"
	testFilePath        = "../assets/" + objectKey
)

func getBucketAPI(t *testing.T) dm.OssAPI {

	// prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	if clientID == "" {
		t.Fatal("clientID is empty")
	}

	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")
	if clientSecret == "" {
		t.Fatal("clientSecret is empty")
	}

	authenticator := oauth.NewTwoLegged(clientID, clientSecret)
	if authenticator == nil {
		t.Fatal("Error authenticating, authenticator is nil.")
	}

	return dm.NewOssApi(authenticator, forge.US)
}

func TestBucketAPI_ListObjects(t *testing.T) {

	bucketAPI := getBucketAPI(t)

	bucketKey := "forge_unit_testing_list_objects"

	t.Run(
		"Create a temp bucket to store an object", func(t *testing.T) {

			_, err := bucketAPI.GetBucketDetails(bucketKey)
			if err == nil {
				t.Skip("The temp bucket already exists.")
			}

			_, err = bucketAPI.CreateBucket(bucketKey, dm.PolicyTransient)
			if err != nil {
				t.Error("Could not create temp bucket, got: ", err.Error())
			}
		},
	)

	t.Run(
		"List bucket content", func(t *testing.T) {
			content, err := bucketAPI.ListObjects(bucketKey, "", "", "")
			if err != nil {
				t.Fatalf("Failed to list bucket content: %s\n", err.Error())
			}

			t.Logf("%#v", content)
		},
	)

	t.Run(
		"List bucket content of non-existing bucket", func(t *testing.T) {
			tmpBucketKey := fmt.Sprintf("%v", time.Now().UnixNano())
			content, err := bucketAPI.ListObjects(tmpBucketKey, "", "", "")
			if err == nil {
				t.Fatalf("Expected to fail upon listing a non-existing bucket, but it didn't, got %#v", content)
			}
		},
	)

	t.Cleanup(
		func() {
			t.Log("Cleaning up the temp bucket")
			err := bucketAPI.DeleteBucket(bucketKey)
			if err != nil {
				t.Error("Could not delete temp bucket, got: ", err.Error())
			}
		},
	)
}

func TestBucketAPI_UploadObject(t *testing.T) {

	bucketAPI := getBucketAPI(t)

	bucketKey := "forge_unit_testing_upload_object"

	t.Run(
		"Create a temp bucket to store an object", func(t *testing.T) {
			_, err := bucketAPI.GetBucketDetails(bucketKey)
			if err == nil {
				t.Skip("The temp bucket already exists, try to delete it.")
			}

			_, err = bucketAPI.CreateBucket(bucketKey, dm.PolicyTransient)
			if err != nil {
				t.Error("Could not create temp bucket, got: ", err.Error())
			}
		},
	)

	t.Run(
		"List objects in temp bucket, to make sure it is empty", func(t *testing.T) {
			content, err := bucketAPI.ListObjects(bucketKey, "", "", "")
			if err != nil {
				t.Fatalf("Failed to list bucket content: %s\n", err.Error())
			}
			if len(content.Items) != 0 {
				t.Fatalf("temp bucket supposed to be empty, got %#v", content)
			}
		},
	)

	t.Run(
		"Upload an object into temp bucket", func(t *testing.T) {
			result, err := bucketAPI.UploadObject(bucketKey, objectKey, testFilePath)

			if err != nil {
				t.Error("Could not upload the test object, got: ", err.Error())
				t.Fatal("Could not upload the test object, got: ", err.Error())
			}

			if result.Size == 0 {
				t.Fatal("The test object was uploaded but it is zero-sized")
			}
		},
	)

	t.Cleanup(
		func() {
			t.Log("Cleaning up the temp bucket")
			err := bucketAPI.DeleteBucket(bucketKey)
			if err != nil {
				t.Error("Could not delete temp bucket, got: ", err.Error())
			}
		},
	)
}

func TestBucketAPI_DownloadObject(t *testing.T) {

	bucketAPI := getBucketAPI(t)

	bucketKey := "forge_unit_testing_upload_and_download_object"

	t.Run(
		"Create a temp bucket to store an object", func(t *testing.T) {
			_, err := bucketAPI.GetBucketDetails(bucketKey)
			if err == nil {
				t.Skip("The temp bucket already exists.")
			}

			_, err = bucketAPI.CreateBucket(bucketKey, dm.PolicyTransient)
			if err != nil {
				t.Error("Could not create temp bucket, got: ", err.Error())
			}
		},
	)

	t.Run(
		"List objects in temp bucket, to make sure it is empty", func(t *testing.T) {
			content, err := bucketAPI.ListObjects(bucketKey, "", "", "")
			if err != nil {
				t.Fatalf("Failed to list bucket content: %s\n", err.Error())
			}
			if len(content.Items) != 0 {
				t.Fatalf("temp bucket supposed to be empty, got %#v", content)
			}
		},
	)

	t.Run(
		"Upload an object into temp bucket", func(t *testing.T) {
			result, err := bucketAPI.UploadObject(bucketKey, objectKey, testFilePath)

			if err != nil {
				t.Fatal("Could not upload the test object, got: ", err.Error())
			}

			if result.Size == 0 {
				t.Fatal("The test object was uploaded but it is zero-sized")
			}
		},
	)

	t.Run(
		"Download an object from the temp bucket", func(t *testing.T) {
			result, err := bucketAPI.DownloadObject(bucketKey, objectKey)
			if err != nil {
				t.Errorf("Problems getting the object %s: %s", objectKey, err.Error())
			}

			if len(result) == 0 {
				t.Errorf("The object %s was downloaded sucessfully, but it is empty.", objectKey)
			}

		},
	)

	t.Cleanup(
		func() {
			t.Log("Cleaning up the temp bucket")
			err := bucketAPI.DeleteBucket(bucketKey)
			if err != nil {
				t.Error("Could not delete temp bucket, got: ", err.Error())
			}
		},
	)
}
