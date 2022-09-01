package dm

import (
	"time"

	"github.com/apprentice3d/forge-api-go-client/oauth"
)

/* BUCKET API TYPES */

// BucketAPI holds the necessary data for making Bucket related calls to Forge Data Management service
type BucketAPI struct {
	Authenticator oauth.ForgeAuthenticator
	BucketAPIPath string
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
	CreateDate  int64  `json:"createDate"`
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
		CreatedDate int64  `json:"createdDate"`
		PolicyKey   string `json:"policyKey"`
	} `json:"items"`
	Next string `json:"next"`
}

// ObjectDetails reflects the data presented when uploading an object to a bucket or requesting details on object.
type ObjectDetails struct {
	BucketKey   string            `json:"bucketKey"`
	ObjectID    string            `json:"objectID"`
	ObjectKey   string            `json:"objectKey"`
	SHA1        string            `json:"sha1"`
	Size        uint64            `json:"size"`
	ContentType string            `json:"contentType, omitempty"`
	Location    string            `json:"location"`
	BlockSizes  []int64           `json:"blockSizes, omitempty"`
	Deltas      map[string]string `json:"deltas, omitempty"`
}

// BucketContent reflects the response when query Data Management API for bucket content.
type BucketContent struct {
	Items []ObjectDetails `json:"items"`
	Next  string          `json:"next"`
}

type PreSignedUploadUrls struct {
	UploadKey        string    `json:"uploadKey"`
	UploadExpiration time.Time `json:"uploadExpiration"`
	UrlExpiration    time.Time `json:"urlExpiration"`
	Urls             []string  `json:"urls"`
}
