package dm

import (
	"github.com/outer-labs/forge-api-go-client/oauth"
)

// BucketAPI holds the necessary data for making Bucket related calls to Forge Data Management service
type BucketAPI3L struct {
	Auth           oauth.ThreeLeggedAuth
	Token          TokenRefresher
	BucketsAPIPath string
}

// NewBucketAPIWithCredentials returns a Bucket API client with default configurations
func NewBucketAPI3LWithCredentials(auth oauth.ThreeLeggedAuth, token TokenRefresher) *BucketAPI3L {
	return &BucketAPI3L{
		Auth:           auth,
		Token:          token,
		BucketsAPIPath: "/oss/v2/buckets",
	}
}

// CreateBucket creates and returns details of created bucket, or an error on failure
func (api BucketAPI3L) CreateBucket3L(bucketKey, policyKey string) (result BucketDetails, err error) {
	if err = api.Token.RefreshTokenIfRequired(api.Auth); err != nil {
		return
	}

	path := api.Auth.Host + api.BucketsAPIPath
	result, err = createBucket(path, bucketKey, policyKey, api.Token.Bearer().AccessToken)

	return
}

// DeleteBucket deletes bucket given its key.
// 	WARNING: The bucket delete call is undocumented.
func (api BucketAPI3L) DeleteBucket3L(bucketKey string) error {
	if err := api.Token.RefreshTokenIfRequired(api.Auth); err != nil {
		return err
	}

	path := api.Auth.Host + api.BucketsAPIPath

	return deleteBucket(path, bucketKey, api.Token.Bearer().AccessToken)
}

// ListBuckets returns a list of all buckets created or associated with Forge secrets used for token creation
func (api BucketAPI3L) ListBuckets3L(region, limit, startAt string) (result ListedBuckets, err error) {
	if err = api.Token.RefreshTokenIfRequired(api.Auth); err != nil {
		return
	}

	path := api.Auth.Host + api.BucketsAPIPath

	return listBuckets(path, region, limit, startAt, api.Token.Bearer().AccessToken)
}

// GetBucketDetails returns information associated to a bucket. See BucketDetails struct.
func (api BucketAPI3L) GetBucketDetails3L(bucketKey string) (result BucketDetails, err error) {
	if err = api.Token.RefreshTokenIfRequired(api.Auth); err != nil {
		return
	}

	path := api.Auth.Host + api.BucketsAPIPath
	return getBucketDetails(path, bucketKey, api.Token.Bearer().AccessToken)
}
