package dm_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/apprentice3d/forge-api-go-client/dm"
)

func TestBucketAPI_ListObjects(t *testing.T) {
	// prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")
	bucketAPI := dm.NewBucketAPIWithCredentials(clientID, clientSecret)

	testBucketName := "just_a_test_bucket"

	t.Run("List bucket content", func(t *testing.T) {
		content, err := bucketAPI.ListObjects(testBucketName, "", "", "")
		if err != nil {
			t.Fatalf("Failed to list bucket content: %s\n", err.Error())
		}

		t.Logf("%#v", content)
	})

	t.Run("List bucket content of non-existing bucket", func(t *testing.T) {
		content, err := bucketAPI.ListObjects(testBucketName+"hz", "", "", "")
		if err == nil {
			t.Fatalf("Expected to fail upon listing a non-existing bucket, but it didn't, got %#v", content)
		}
	})

}

func TestBucketAPI_UploadObject(t *testing.T) {

	// prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")
	bucketAPI := dm.NewBucketAPIWithCredentials(clientID, clientSecret)

	tempBucket := "some_temp_bucket_for_testings"
	testFilePath := "../assets/HelloWorld.rvt"

	t.Run("Create a temp bucket to store an object", func(t *testing.T) {
		_, err := bucketAPI.CreateBucket(tempBucket, "transient")
		if err != nil {
			t.Error("Could not create temp bucket, got: ", err.Error())
		}
	})

	t.Run("List objects in temp bucket, to make sure it is empty", func(t *testing.T) {
		content, err := bucketAPI.ListObjects(tempBucket, "", "", "")
		if err != nil {
			t.Fatalf("Failed to list bucket content: %s\n", err.Error())
		}
		if len(content.Items) != 0 {
			t.Fatalf("temp bucket supposed to be empty, got %#v", content)
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

		result, err := bucketAPI.UploadObject(tempBucket, "temp_file.rvt", data)

		if err != nil {
			t.Fatal("Could not upload the test object, got: ", err.Error())
		}

		if result.Size == 0 {
			t.Fatal("The test object was uploaded but it is zero-sized")
		}
	})

	t.Run("Delete the temp bucket", func(t *testing.T) {
		err := bucketAPI.DeleteBucket(tempBucket)
		if err != nil {
			t.Error("Could not delete temp bucket, got: ", err.Error())
		}
	})
}
