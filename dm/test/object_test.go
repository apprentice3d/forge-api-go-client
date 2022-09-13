package dm_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/woweh/forge-api-go-client/oauth"

	"github.com/woweh/forge-api-go-client/dm"
)

const (
	bucketKey    string = "forge_unit_testing"
	objectKey    string = "rst_basic_sample_project.rvt"
	testFilePath string = "../assets/" + objectKey
)

func getBucketAPI(t *testing.T) dm.BucketAPI {

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

	return dm.NewBucketAPI(authenticator)
}

func TestBucketAPI_ListObjects(t *testing.T) {

	bucketAPI := getBucketAPI(t)

	_, err := bucketAPI.GetBucketDetails(bucketKey)
	if err != nil {
		// bucket doesn't exist yet, create it
		t.Run("Create a temp bucket to store an object", func(t *testing.T) {
			_, err := bucketAPI.CreateBucket(bucketKey, "transient")
			if err != nil {
				t.Error("Could not create temp bucket, got: ", err.Error())
			}
		})
	}

	t.Run("List bucket content", func(t *testing.T) {
		content, err := bucketAPI.ListObjects(bucketKey, "", "", "")
		if err != nil {
			t.Fatalf("Failed to list bucket content: %s\n", err.Error())
		}

		t.Logf("%#v", content)
	})

	t.Run("List bucket content of non-existing bucket", func(t *testing.T) {
		tmpBucketKey := fmt.Sprintf("%v", time.Now().UnixNano())
		content, err := bucketAPI.ListObjects(tmpBucketKey, "", "", "")
		if err == nil {
			t.Fatalf("Expected to fail upon listing a non-existing bucket, but it didn't, got %#v", content)
		}
	})

}

func TestBucketAPI_UploadObject(t *testing.T) {

	bucketAPI := getBucketAPI(t)

	_, err := bucketAPI.GetBucketDetails(bucketKey)
	if err != nil {
		// bucket doesn't exist yet, create it
		t.Run("Create a temp bucket to store an object", func(t *testing.T) {
			_, err := bucketAPI.CreateBucket(bucketKey, "transient")
			if err != nil {
				t.Error("Could not create temp bucket, got: ", err.Error())
			}
		})
	}

	t.Run("List objects in temp bucket, to make sure it is empty", func(t *testing.T) {
		content, err := bucketAPI.ListObjects(bucketKey, "", "", "")
		if err != nil {
			t.Fatalf("Failed to list bucket content: %s\n", err.Error())
		}
		if len(content.Items) != 0 {
			t.Fatalf("temp bucket supposed to be empty, got %#v", content)
		}
	})

	t.Run("Upload an object into temp bucket", func(t *testing.T) {
		result, err := bucketAPI.UploadObject(bucketKey, objectKey, testFilePath)

		if err != nil {
			t.Fatal("Could not upload the test object, got: ", err.Error())
		}

		if result.Size == 0 {
			t.Fatal("The test object was uploaded but it is zero-sized")
		}
	})

	t.Run("Delete the temp bucket", func(t *testing.T) {
		err := bucketAPI.DeleteBucket(bucketKey)
		if err != nil {
			t.Error("Could not delete temp bucket, got: ", err.Error())
		}
	})
}

func TestBucketAPI_DownloadObject(t *testing.T) {

	bucketAPI := getBucketAPI(t)

	t.Run("Create a temp bucket to store an object", func(t *testing.T) {
		_, err := bucketAPI.CreateBucket(bucketKey, "transient")
		if err != nil {
			t.Error("Could not create temp bucket, got: ", err.Error())
		}
	})

	t.Run("List objects in temp bucket, to make sure it is empty", func(t *testing.T) {
		content, err := bucketAPI.ListObjects(bucketKey, "", "", "")
		if err != nil {
			t.Fatalf("Failed to list bucket content: %s\n", err.Error())
		}
		if len(content.Items) != 0 {
			t.Fatalf("temp bucket supposed to be empty, got %#v", content)
		}
	})

	t.Run("Upload an object into temp bucket", func(t *testing.T) {
		result, err := bucketAPI.UploadObject(bucketKey, objectKey, testFilePath)

		if err != nil {
			t.Fatal("Could not upload the test object, got: ", err.Error())
		}

		if result.Size == 0 {
			t.Fatal("The test object was uploaded but it is zero-sized")
		}
	})

	t.Run("Download an object from the temp bucket", func(t *testing.T) {
		result, err := bucketAPI.DownloadObject(bucketKey, objectKey)
		if err != nil {
			t.Errorf("Problems getting the object %s: %s", objectKey, err.Error())
		}

		if len(result) == 0 {
			t.Errorf("The object %s was downloaded sucessfully, but it is empty.", objectKey)
		}

	})

	t.Run("Delete the temp bucket", func(t *testing.T) {
		err := bucketAPI.DeleteBucket(bucketKey)
		if err != nil {
			t.Error("Could not delete temp bucket, got: ", err.Error())
		}
	})
}
