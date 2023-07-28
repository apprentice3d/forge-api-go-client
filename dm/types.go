package dm

import (
	"github.com/woweh/forge-api-go-client"
	"github.com/woweh/forge-api-go-client/oauth"
)

/* BUCKET API TYPES */

// OssAPI holds the necessary data for making Object Storage Service (OSS) related calls to the Forge Data Management API.
type OssAPI struct {
	Authenticator oauth.ForgeAuthenticator
	relativePath  string // = "/oss/v2/buckets", populate in NewOssApi
	region        forge.Region
}

// CreateBucketRequest contains the data necessary to be passed upon bucket creation
type CreateBucketRequest struct {
	BucketKey string          `json:"bucketKey"`
	PolicyKey RetentionPolicy `json:"policyKey"`
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
	PolicyKey RetentionPolicy `json:"policyKey"`
}

// ErrorResult reflects the body content when a request failed (g.e. Bad request or key conflict)
type ErrorResult struct {
	Reason string `json:"reason"`
}

type BucketList []BucketInfo

type BucketInfo struct {
	BucketKey   string          `json:"bucketKey"`
	CreatedDate int64           `json:"createdDate"`
	PolicyKey   RetentionPolicy `json:"policyKey"`
}

// ListedBuckets reflects the response when query Data Management API for buckets associated with current Forge secrets.
type ListedBuckets struct {
	Items BucketList `json:"items"`
	Next  string     `json:"next"`
}

// ObjectDetails reflects the data presented when uploading an object to a bucket or requesting details on object.
type ObjectDetails struct {
	BucketKey   string            `json:"bucketKey"`
	ObjectID    string            `json:"objectID"` // => urn = base64.RawStdEncoding.EncodeToString([]byte(ObjectID))
	ObjectKey   string            `json:"objectKey"`
	SHA1        string            `json:"sha1"`
	Size        uint64            `json:"size"`
	ContentType string            `json:"contentType,omitempty"`
	Location    string            `json:"location"`
	BlockSizes  []int64           `json:"blockSizes,omitempty"`
	Deltas      map[string]string `json:"deltas,omitempty"`
}

// UploadResult provides the OK/200 result of the completeUpload POST.
type UploadResult struct {
	BucketKey   string `json:"bucketKey"`
	ObjectId    string `json:"objectId"` // => urn = base64.RawStdEncoding.EncodeToString([]byte(ObjectID))
	ObjectKey   string `json:"objectKey"`
	Size        int    `json:"size"`
	ContentType string `json:"content-type"`
	Location    string `json:"location"`
}

// BucketContent reflects the response when query Data Management API for bucket content.
type BucketContent struct {
	Items []ObjectDetails `json:"items"`
	Next  string          `json:"next"`
}
