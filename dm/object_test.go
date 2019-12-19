package dm_test

import (
	"io"
	"os"
	"testing"
	"github.com/outer-labs/forge-api-go-client/dm"
)

func TestBucketAPI_ListObjects(t *testing.T) {
	// prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

	bucketAPI := dm.NewBucketAPIWithCredentials(clientID, clientSecret)

	// testBucketName := "just_a_test_bucket"
	testBucketName := os.Getenv("FORGE_OSS_TEST_BUCKET_KEY")

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

	tempBucket := "some_temp_bucket_for_testing"
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
		// data, err := ioutil.ReadAll(file) // returns []byte
		data := io.Reader(file)
		if err != nil {
			t.Fatal("Cannot read the testfile")
		}

		result, err := bucketAPI.UploadObject(tempBucket, "temp_file.rvt", data) // doesn't want []byte as data

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

func TestBucketAPI_DownloadObject(t *testing.T) {
	// prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

	bucketAPI := dm.NewBucketAPIWithCredentials(clientID, clientSecret)

	tempBucket := "test_bucket_for_download"
	testFilePath := "../assets/TestFile.txt"
	bucketDetails, err := bucketAPI.GetBucketDetails(tempBucket)

	// Check if bucket is still hanging around
	if err != nil && bucketDetails.CreateDate == "" {
		_, err := bucketAPI.CreateBucket(tempBucket, "transient")
		if err != nil {
			t.Error("Could not create temp bucket, got: ", err.Error())
		}
		defer deleteBucket(bucketAPI, tempBucket, t)
	}

	file, err := os.Open(testFilePath)
	if err != nil {
		t.Fatal("Cannot open testfile for reading")
	}
	defer file.Close()

	data := io.Reader(file)
	if err != nil {
		t.Fatal("Cannot read the testfile")
	}

	result, err := bucketAPI.UploadObject(tempBucket, "temp_file.txt", data) // doesn't want []byte as data

	if err != nil {
		t.Fatal("Could not upload the test object, got: ", err.Error())
	}

	if result.Size == 0 {
		t.Fatal("The test object was uploaded but it is zero-sized")
	}

	reader, err := bucketAPI.DownloadObject(tempBucket, "temp_file.txt")
	defer reader.Close()
	if err != nil {
		t.Fatal("Could not download the test object, got: ", err.Error())
	}
	buf := make([]byte, 15)
	if _, err := io.ReadFull(reader, buf); err != nil {
		t.Fatal(err)
	}
	if string(buf) != "Test test 1 2 3" {
		t.Fatal("Test file contents do not match what was downloaded, got: ", string(buf))
	}

}

func deleteBucket(bucketAPI dm.BucketAPI, bucketKey string, t *testing.T) {
	err := bucketAPI.DeleteBucket(bucketKey)
	if err != nil {
		t.Error("Could not delete temp bucket, got: ", err.Error())
	}
}