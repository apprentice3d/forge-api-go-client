package dm

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/woweh/forge-api-go-client"
	"github.com/woweh/forge-api-go-client/oauth"
)

// NewOssApi returns an OSS API client with default configurations and populates the RelativePath.
func NewOssApi(authenticator oauth.ForgeAuthenticator, region forge.Region) OssAPI {
	return OssAPI{
		Authenticator: authenticator,
		relativePath:  "/oss/v2/buckets",
		region:        region,
	}
}

// Region of the OSS API.
func (a *OssAPI) Region() forge.Region {
	return a.region
}

// SetRegion sets the Region of the OSS API.
func (a *OssAPI) SetRegion(region forge.Region) {
	a.region = region
}

// RelativePath of the OSS API.
func (a *OssAPI) RelativePath() string {
	return a.relativePath
}

// BaseUrl of the OSS API.
func (a *OssAPI) BaseUrl() string {
	return a.Authenticator.HostPath() + a.relativePath
}

// RetentionPolicy applies to all objects that are stored in a bucket.
//   - This cannot be changed at a later time!
//
// When creating a bucket, specifically set the policyKey to transient, temporary, or persistent.
type RetentionPolicy string

const (
	// PolicyTransient - Think of this type of storage as a cache. Use it for ephemeral results.
	// For example, you might use this for objects that are part of producing other persistent artifacts, but otherwise are not required to be available later.
	// Objects older than 24 hours are removed automatically.
	// Each upload of an object is considered unique, so, for example, if the same rendering is uploaded multiple times, each of them will have its own retention period of 24 hours.
	PolicyTransient RetentionPolicy = "transient"

	// PolicyTemporary - This type of storage is suitable for artifacts produced for user-uploaded content where after some period of activity, the user may rarely access the artifacts.
	// When an object has reached 30 days of age, it is deleted.
	PolicyTemporary RetentionPolicy = "temporary"

	// PolicyPersistent - Persistent storage is intended for user data.
	// When a file is uploaded, the owner should expect this item to be available for as long as the owner account is active, or until he or she deletes the item.
	PolicyPersistent RetentionPolicy = "persistent"
)

// CreateBucket creates and returns details of created bucket, or an error on failure.
// The region is taken from the OssAPI instance.
//   - bucketKey: A unique name you assign to a bucket. It must be globally unique across all applications and regions, otherwise the call will fail.
//     Possible values: -_.a-z0-9 (between 3-128 characters in length).
//     Note that you cannot change a bucket key.
//   - policyKey: Data retention policy. Acceptable values: transient, temporary, persistent.
//     This cannot be changed at a later time. The retention policy on the bucket applies to all objects stored within.
//
// References:
//   - https://aps.autodesk.com/en/docs/data/v2/reference/http/buckets-POST/
func (a *OssAPI) CreateBucket(bucketKey string, policyKey RetentionPolicy) (result BucketDetails, err error) {

	bearer, err := a.Authenticator.GetToken("bucket:create")
	if err != nil {
		return
	}

	result, err = createBucket(a.BaseUrl(), bucketKey, policyKey, bearer.AccessToken, a.region)

	return
}

// DeleteBucket deletes bucket given its key.
//   - The bucket must be owned by the application.
//   - We recommend only deleting small buckets used for acceptance testing or prototyping, since it can take a long time for a bucket to be deleted.
//   - Note that the bucket name will not be immediately available for reuse.
//
// References:
//   - https://aps.autodesk.com/en/docs/data/v2/reference/http/buckets-:bucketKey-DELETE/
func (a *OssAPI) DeleteBucket(bucketKey string) error {
	bearer, err := a.Authenticator.GetToken("bucket:delete")
	if err != nil {
		return err
	}

	return deleteBucket(a.BaseUrl(), bucketKey, bearer.AccessToken)
}

// ListBuckets returns a list of all buckets created or associated with Forge secrets used for token creation
//
// References:
//   - https://aps.autodesk.com/en/docs/data/v2/reference/http/buckets-GET/
func (a *OssAPI) ListBuckets(region forge.Region, limit, startAt string) (result BucketList, err error) {

	// init the result
	result = BucketList{}

loop:
	bearer, err := a.Authenticator.GetToken("bucket:read")
	if err != nil {
		return
	}

	tmpResult, err := listBuckets(a.BaseUrl(), region, limit, startAt, bearer.AccessToken)
	if err != nil {
		return
	}

	// append the result
	result = append(result, tmpResult.Items...)

	// if there are more items, get them
	for tmpResult.Next != "" {
		// extract the startAt from the next link
		startAt, err = extractStartAt(tmpResult.Next)
		if err != nil {
			return
		}
		goto loop
	}

	return
}

func extractStartAt(nextUrl string) (startAt string, err error) {
	parsedUrl, err := url.Parse(nextUrl)
	if err != nil {
		return "", err
	}
	startAt = parsedUrl.Query().Get("startAt")
	if startAt == "" {
		return "", errors.New("startAt not found in next url")
	}
	return startAt, nil
}

// GetBucketDetails returns information associated to a bucket. See BucketDetails struct.
//
// References:
//   - https://aps.autodesk.com/en/docs/data/v2/reference/http/buckets-:bucketKey-details-GET/
func (a *OssAPI) GetBucketDetails(bucketKey string) (result BucketDetails, err error) {
	bearer, err := a.Authenticator.GetToken("bucket:read")
	if err != nil {
		return
	}

	return getBucketDetails(a.BaseUrl(), bucketKey, bearer.AccessToken)
}

/*
 *	SUPPORT FUNCTIONS
 */

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

	return
}

func createBucket(
	path, bucketKey string, policyKey RetentionPolicy, token string, region forge.Region,
) (result BucketDetails, err error) {

	task := http.Client{}

	body, err := json.Marshal(
		CreateBucketRequest{
			bucketKey,
			policyKey,
		},
	)
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
