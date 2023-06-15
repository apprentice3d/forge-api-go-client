package dm

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/woweh/forge-api-go-client"
	"github.com/woweh/forge-api-go-client/oauth"
)

// NewBucketAPI returns a Bucket API client with default configurations and populates the BucketAPIPath.
func NewBucketAPI(authenticator oauth.ForgeAuthenticator) BucketAPI {
	return BucketAPI{
		authenticator,
		"/oss/v2/buckets",
	}
}

// CreateBucket creates and returns details of created bucket, or an error on failure.
//
// This method uses the default region (US).
// Use CreateBucketForRegion to specify a region (> EMEA).
//
// References:
//   - https://aps.autodesk.com/en/docs/data/v2/reference/http/buckets-POST/
//
// Deprecated: Use CreateBucketForRegion instead.
func (api *BucketAPI) CreateBucket(bucketKey, policyKey string) (result BucketDetails, err error) {

	bearer, err := api.Authenticator.GetToken("bucket:create")
	if err != nil {
		return
	}

	result, err = createBucket(api.getPath(), bucketKey, policyKey, bearer.AccessToken, forge.US)

	return
}

// CreateBucket creates and returns details of created bucket, or an error on failure.
//
// References:
//   - https://aps.autodesk.com/en/docs/data/v2/reference/http/buckets-POST/
func (api *BucketAPI) CreateBucketForRegion(bucketKey, policyKey string, region forge.Region) (result BucketDetails, err error) {

	bearer, err := api.Authenticator.GetToken("bucket:create")
	if err != nil {
		return
	}

	result, err = createBucket(api.getPath(), bucketKey, policyKey, bearer.AccessToken, region)

	return
}

// DeleteBucket deletes bucket given its key.
//
// References:
//   - https://aps.autodesk.com/en/docs/data/v2/reference/http/buckets-:bucketKey-DELETE/
func (api *BucketAPI) DeleteBucket(bucketKey string) error {
	bearer, err := api.Authenticator.GetToken("bucket:delete")
	if err != nil {
		return err
	}

	return deleteBucket(api.getPath(), bucketKey, bearer.AccessToken)
}

// ListBuckets returns a list of all buckets created or associated with Forge secrets used for token creation
//
// References:
//   - https://aps.autodesk.com/en/docs/data/v2/reference/http/buckets-GET/
func (api *BucketAPI) ListBuckets(region forge.Region, limit, startAt string) (result ListedBuckets, err error) {
	bearer, err := api.Authenticator.GetToken("bucket:read")
	if err != nil {
		return
	}

	return listBuckets(api.getPath(), region, limit, startAt, bearer.AccessToken)
}

// GetBucketDetails returns information associated to a bucket. See BucketDetails struct.
//
// References:
//   - https://aps.autodesk.com/en/docs/data/v2/reference/http/buckets-:bucketKey-details-GET/
func (api *BucketAPI) GetBucketDetails(bucketKey string) (result BucketDetails, err error) {
	bearer, err := api.Authenticator.GetToken("bucket:read")
	if err != nil {
		return
	}

	return getBucketDetails(api.getPath(), bucketKey, bearer.AccessToken)
}

// BucketExists returns true if the bucket exists, false otherwise.
// This function calls GetBucketDetails and checks if the bucket key matches.
func (api *BucketAPI) BucketExists(bucketKey string) (exists bool, err error) {
	bearer, err := api.Authenticator.GetToken("bucket:read")
	if err != nil {
		return false, err
	}

	result, err := getBucketDetails(api.getPath(), bucketKey, bearer.AccessToken)
	if err != nil {
		return false, err
	}
	return result.BucketKey == bucketKey, nil
}

/*
 *	SUPPORT FUNCTIONS
 */

// getPath gets the full bucket API path (= api.Authenticator.GetHostPath() + api.BucketAPIPath).
func (api *BucketAPI) getPath() string {
	return api.Authenticator.GetHostPath() + api.BucketAPIPath
}

func getBucketDetails(path, bucketKey, token string) (result BucketDetails, err error) {
	task := http.Client{}

	req, err := http.NewRequest("GET", path+"/"+bucketKey+"/details", nil)

	if err != nil {
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	response, err := task.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		content, _ := io.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&result)

	return
}

func listBuckets(path string, region forge.Region, limit, startAt, token string) (result ListedBuckets, err error) {
	task := http.Client{}

	req, err := http.NewRequest("GET", path, nil)

	if err != nil {
		return
	}

	params := req.URL.Query()
	if len(region) != 0 {
		params.Add("region", string(region))
	}
	if len(limit) != 0 {
		params.Add("limit", limit)
	}
	if len(startAt) != 0 {
		params.Add("startAt", startAt)
	}

	req.URL.RawQuery = params.Encode()

	req.Header.Set("Authorization", "Bearer "+token)
	response, err := task.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		content, _ := io.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	err = json.NewDecoder(response.Body).Decode(&result)

	/* TODO: address the pagination of buckets
	if result.Next != "" {
		// get the next batch
	}
	*/

	return
}

func createBucket(path, bucketKey, policyKey, token string, region forge.Region) (result BucketDetails, err error) {

	task := http.Client{}

	body, err := json.Marshal(
		CreateBucketRequest{
			bucketKey,
			policyKey,
		})
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", path, bytes.NewReader(body))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("x-ads-region", string(region))
	response, err := task.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		content, _ := io.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	decoder := json.NewDecoder(response.Body)

	err = decoder.Decode(&result)

	return
}

func deleteBucket(path, bucketKey, token string) (err error) {
	task := http.Client{}

	req, err := http.NewRequest("DELETE", path+"/"+bucketKey, nil)
	if err != nil {
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	response, err := task.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		content, _ := io.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	return
}
