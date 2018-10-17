package dm

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/apprentice3d/forge-api-go-client/oauth"
)

// BucketAPI holds the necessary data for making Bucket related calls to Forge Data Management service
type BucketAPI struct {
	oauth.TwoLeggedAuth
	BucketAPIPath string
}

// NewBucketAPIWithCredentials returns a Bucket API client with default configurations
func NewBucketAPIWithCredentials(ClientID string, ClientSecret string) BucketAPI {
	return BucketAPI{
		oauth.NewTwoLeggedClient(ClientID, ClientSecret),
		"/oss/v2/buckets",
	}
}

// CreateBucketRequest contains the data necessary to be passed upon bucket creation
type CreateBucketRequest struct {
	BucketKey string `json:"bucketKey"`
	PolicyKey string `json:"policyKey"`
}

// BucketDetails reflects the body content received upon creation of a bucket
type BucketDetails struct {
	BucketKey   string `json:"bucketKey"`
	BucketOwner string `json:"bucketOwner"`
	CreateDate  int64 `json:"createDate"`
	Permissions []struct {
		AuthID string `json:"authId"`
		Access string `json:"access"`
	} `json:"permissions"`
	PolicyKey string `json:"policyKey"`
}

// ErrorResult reflects the body content when a request failed (g.e. Bad request or key conflict)
type ErrorResult struct {
	Reason string `json:"reason"`
}

// ListedBuckets reflects the response when query Data Management API for buckets associated with current Forge secrets.
type ListedBuckets struct {
	Items []struct {
		BucketKey   string `json:"bucketKey"`
		CreatedDate int64 `json:"createdDate"`
		PolicyKey   string `json:"policyKey"`
	} `json:"items"`
	Next string `json:"next"`
}


// CreateBucket creates and returns details of created bucket, or an error on failure
func (api BucketAPI) CreateBucket(bucketKey, policyKey string) (result BucketDetails, err error) {

	bearer, err := api.Authenticate("bucket:create")
	if err != nil {
		return
	}
	path := api.Host + api.BucketAPIPath
	result, err = createBucket(path, bucketKey, policyKey, bearer.AccessToken)

	return
}

// DeleteBucket deletes bucket given its key.
// 	WARNING: The bucket delete call is undocumented.
func (api BucketAPI) DeleteBucket(bucketKey string) error {
	bearer, err := api.Authenticate("bucket:delete")
	if err != nil {
		return err
	}
	path := api.Host + api.BucketAPIPath

	return deleteBucket(path, bucketKey, bearer.AccessToken)
}

// ListBuckets returns a list of all buckets created or associated with Forge secrets used for token creation
func (api BucketAPI) ListBuckets(region, limit, startAt string) (result ListedBuckets, err error) {
	bearer, err := api.Authenticate("bucket:read")
	if err != nil {
		return
	}
	path := api.Host + api.BucketAPIPath

	return listBuckets(path, region, limit, startAt, bearer.AccessToken)
}

// GetBucketDetails returns information associated to a bucket. See BucketDetails struct.
func (api BucketAPI) GetBucketDetails(bucketKey string) (result BucketDetails, err error) {
	bearer, err := api.Authenticate("bucket:read")
	if err != nil {
		return
	}
	path := api.Host + api.BucketAPIPath

	return getBucketDetails(path, bucketKey, bearer.AccessToken)
}





/*
 *	SUPPORT FUNCTIONS
 */
func getBucketDetails(path, bucketKey, token string) (result BucketDetails, err error) {
	task := http.Client{}

	req, err := http.NewRequest("GET",
		path+"/"+bucketKey+"/details",
		nil,
	)

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
		content, _ := ioutil.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&result)

	return
}

func listBuckets(path, region, limit, startAt, token string) (result ListedBuckets, err error) {
	task := http.Client{}

	req, err := http.NewRequest("GET",
		path,
		nil,
	)

	if err != nil {
		return
	}

	params := req.URL.Query()
	if len(region) != 0 {
		params.Add("region", region)
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
		content, _ := ioutil.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&result)

	return
}

func createBucket(path, bucketKey, policyKey, token string) (result BucketDetails, err error) {

	task := http.Client{}

	body, err := json.Marshal(
		CreateBucketRequest{
			bucketKey,
			policyKey,
		})
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST",
		path,
		bytes.NewReader(body),
	)

	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	response, err := task.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		content, _ := ioutil.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	decoder := json.NewDecoder(response.Body)

	err = decoder.Decode(&result)

	return
}

func deleteBucket(path, bucketKey, token string) (err error) {
	task := http.Client{}

	req, err := http.NewRequest("DELETE",
		path+"/"+bucketKey,
		nil,
	)

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
		content, _ := ioutil.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	return
}
