package dm_test

import (
	"testing"
)

/*  NOTE:
Since tests are run in parallel, you should use a unique bucketKey per test.
*/

func TestBucketAPI_CreateBucket(t *testing.T) {

	bucketAPI := getBucketAPI(t)

	bucketKey := "forge_unit_testing_create_bucket"

	t.Run("Create a bucket", func(t *testing.T) {
		_, err := bucketAPI.GetBucketDetails(bucketKey)
		if err == nil {
			t.Log("The temp bucket already exists, try to delete it.")

			err := bucketAPI.DeleteBucket(bucketKey)
			if err != nil {
				t.Error("Could not delete temp bucket, got: ", err.Error())
			}
		}

		_, err = bucketAPI.CreateBucket(bucketKey, "transient")
		if err != nil {
			t.Fatalf("Failed to create a bucket: %s\n", err.Error())
		}
	})

	t.Run("Delete created bucket", func(t *testing.T) {
		err := bucketAPI.DeleteBucket(bucketKey)

		if err != nil {
			t.Fatalf("Failed to delete bucket: %s\n", err.Error())
		}
	})

	t.Run("Create a bucket with invalid name", func(t *testing.T) {
		_, err := bucketAPI.CreateBucket("goTestingBucket", "transient")

		if err == nil {
			t.Fatalf("Should fail creating a bucket with invalid name\n")
		}
	})

	t.Run("Create a bucket with invalid policyKey", func(t *testing.T) {
		_, err := bucketAPI.CreateBucket("goTestingBucket", "democracy")

		if err == nil {
			t.Fatalf("Should fail creating a bucket with invalid name\n")
		}
	})
}

func TestBucketAPI_GetBucketDetails(t *testing.T) {

	bucketAPI := getBucketAPI(t)

	bucketKey := "forge_unit_testing_get_bucket_details"

	t.Run("Create a bucket", func(t *testing.T) {
		_, err := bucketAPI.GetBucketDetails(bucketKey)
		if err == nil {
			t.Log("The temp bucket already exists, try to delete it.")

			err := bucketAPI.DeleteBucket(bucketKey)
			if err != nil {
				t.Error("Could not delete temp bucket, got: ", err.Error())
			}
		}

		_, err = bucketAPI.CreateBucket(bucketKey, "transient")
		if err != nil {
			t.Fatalf("Failed to create a bucket: %s\n", err.Error())
		}
	})

	t.Run("Get bucket details", func(t *testing.T) {
		_, err := bucketAPI.GetBucketDetails(bucketKey)

		if err != nil {
			t.Fatalf("Failed to get bucket details: %s\n", err.Error())
		}
	})

	t.Run("Delete created bucket", func(t *testing.T) {
		err := bucketAPI.DeleteBucket(bucketKey)

		if err != nil {
			t.Fatalf("Failed to delete bucket: %s\n", err.Error())
		}
	})

	t.Run("Get nonexistent bucket", func(t *testing.T) {
		_, err := bucketAPI.GetBucketDetails(bucketKey + "30091981")

		if err == nil {
			t.Fatalf("Should fail getting getting details for non-existing bucket\n")
		}
	})
}

func TestBucketAPI_ListBuckets(t *testing.T) {

	bucketAPI := getBucketAPI(t)

	t.Run("List available buckets", func(t *testing.T) {
		_, err := bucketAPI.ListBuckets("", "", "")
		if err != nil {
			t.Fatalf("Failed to list buckets: %s\n", err.Error())
		}
	})

	// TODO: fix ListBuckets to list all buckets (support paging).
	// Enable this test again when that is done.
	/*
		t.Run("Create a bucket and find it among listed", func(t *testing.T) {

			bucketKey := "forge_unit_testing_list_buckets"

			_, err := bucketAPI.GetBucketDetails(bucketKey)
			if err == nil {
				t.Log("The temp bucket already exists, try to delete it.")

				err := bucketAPI.DeleteBucket(bucketKey)
				if err != nil {
					t.Error("Could not delete temp bucket, got: ", err.Error())
				}
			}

			_, err = bucketAPI.CreateBucket(bucketKey, "transient")
			if err != nil {
				t.Errorf("Failed to create a bucket: %s\n", err.Error())
			}

			list, err := bucketAPI.ListBuckets("", "", "")

			if err != nil {
				t.Errorf("Failed to list buckets: %s\n", err.Error())
			}

			found := false

			for _, bucket := range list.Items {
				if bucket.BucketKey == bucketKey {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("Could not find the %s bucket\n", bucketKey)
			}

			if err = bucketAPI.DeleteBucket(bucketKey); err != nil {
				t.Errorf("Failed to delete bucket: %s\n", err.Error())
			}
		})
	*/
}
