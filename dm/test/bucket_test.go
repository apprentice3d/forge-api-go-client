package dm_test

import (
	"testing"

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

func TestBucketAPI_CreateBucket(t *testing.T) {

	bucketAPI := getBucketAPI(t)

	bucketKey := "forge_unit_testing_create_bucket"

	t.Run(
		"Create a bucket", func(t *testing.T) {
			_, err := bucketAPI.GetBucketDetails(bucketKey)
			if err == nil {
				t.Skip("The temp bucket already exists.")
			}

			_, err = bucketAPI.CreateBucket(bucketKey, dm.PolicyTransient)
			if err != nil {
				t.Fatalf("Failed to create a bucket: %s\n", err.Error())
			}
		},
	)

	t.Run(
		"Delete created bucket", func(t *testing.T) {
			err := bucketAPI.DeleteBucket(bucketKey)

			if err != nil {
				t.Fatalf("Failed to delete bucket: %s\n", err.Error())
			}
		},
	)

	t.Run(
		"Create a bucket with invalid name", func(t *testing.T) {
			invalidBucketKey := "$Invalid@Bucket%Key!"
			_, err := bucketAPI.CreateBucket(invalidBucketKey, dm.PolicyTransient)

			if err == nil {
				t.Fatal("Should fail creating a bucket with invalid name: ", invalidBucketKey)
			}
		},
	)

	t.Run(
		"Create a bucket with invalid policyKey", func(t *testing.T) {
			_, err := bucketAPI.CreateBucket("all_lower_case_bucket_key", "invalidPolicy")

			if err == nil {
				t.Fatalf("Should fail creating a bucket with invalid name\n")
			}
		},
	)
}

func TestBucketAPI_GetBucketDetails(t *testing.T) {

	bucketAPI := getBucketAPI(t)

	bucketKey := "forge_unit_testing_get_bucket_details"

	t.Run(
		"Create a bucket", func(t *testing.T) {
			_, err := bucketAPI.GetBucketDetails(bucketKey)
			if err == nil {
				t.Skip("The temp bucket already exists.")
			}

			_, err = bucketAPI.CreateBucket(bucketKey, dm.PolicyTransient)
			if err != nil {
				t.Fatalf("Failed to create a bucket: %s\n", err.Error())
			}
		},
	)

	t.Run(
		"Get bucket details", func(t *testing.T) {
			_, err := bucketAPI.GetBucketDetails(bucketKey)

			if err != nil {
				t.Fatalf("Failed to get bucket details: %s\n", err.Error())
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

func TestBucketAPI_ListBuckets(t *testing.T) {

	bucketAPI := getBucketAPI(t)

	bucketKey := "forge_unit_testing_list_buckets"

	t.Run(
		"List available buckets", func(t *testing.T) {
			_, err := bucketAPI.ListBuckets("", "", "")
			if err != nil {
				t.Fatalf("Failed to list buckets: %s\n", err.Error())
			}
		},
	)

	t.Run(
		"Create a bucket and find it among listed", func(t *testing.T) {

			_, err := bucketAPI.GetBucketDetails(bucketKey)
			if err == nil {
				t.Log("The temp bucket already exists, try to delete it.")

				err = bucketAPI.DeleteBucket(bucketKey)
				if err != nil {
					t.Error("Could not delete temp bucket, got: ", err.Error())
				}
			}

			_, err = bucketAPI.CreateBucket(bucketKey, dm.PolicyTransient)
			if err != nil {
				t.Errorf("Failed to create a bucket: %s\n", err.Error())
			}

			list, err := bucketAPI.ListBuckets("", "", "")

			if err != nil {
				t.Errorf("Failed to list buckets: %s\n", err.Error())
			}

			found := false

			for _, bucket := range list {
				if bucket.BucketKey == bucketKey {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("Could not find the %s bucket\n", bucketKey)
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
