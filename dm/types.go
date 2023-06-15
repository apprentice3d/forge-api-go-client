package dm

import (
	"time"

	"github.com/woweh/forge-api-go-client/oauth"
)

/* BUCKET API TYPES */

// BucketAPI holds the necessary data for making Bucket related calls to Forge Data Management service
type BucketAPI struct {
	Authenticator oauth.ForgeAuthenticator
	BucketAPIPath string // = "/oss/v2/buckets", populate in NewBucketAPI
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

type BucketList []BucketInfo

type BucketInfo struct {
	BucketKey   string `json:"bucketKey"`
	CreatedDate int64  `json:"createdDate"`
	PolicyKey   string `json:"policyKey"`
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

// signedUploadUrls reflects the response from the signedS3UploadEndpoint.
type signedUploadUrls struct {
	UploadKey        string    `json:"uploadKey"`
	UploadExpiration time.Time `json:"uploadExpiration"`
	UrlExpiration    time.Time `json:"urlExpiration"`
	Urls             []string  `json:"urls"`
}

// uploadJob provides information for uploading a file
type uploadJob struct {
	// api is a pointer to an instance of the BucketAPI.
	api *BucketAPI
	// bucketKey is the key (= name) of the bucket where the file shall be stored.
	bucketKey string
	// objectKey is the key (= name) of the file in the Autodesk cloud (OSS).
	objectKey string
	// fileToUpload is the path of the file to upload.
	fileToUpload string
	// minutesExpiration is the custom expiration time within a 1 to 60 minutes range.
	minutesExpiration int
	// fileSize is the size of the file to upload.
	fileSize int64
	// totalParts is the total number of parts to process.
	totalParts int
	// numberOfBatches indicates the number of 'batches' that must be processed, how often signed URLs must be requested.
	// If totalParts > maxParts, then we need to request signedUploadUrls multiple times.
	numberOfBatches int
	// uploadKey is the identifier of the upload session, to differentiate multiple attempts to upload data for the same object.
	// This must be provided when re-requesting chunk URLs for the same blob if they expire, and when calling the Complete Upload endpoint.
	uploadKey string
}

// signedDownloadUrl reflects the response from the "signeds3download" endpoint.
type signedDownloadUrl struct {
	Status string `json:"status"`
	Url    string `json:"url"`
	Params struct {
		ContentType        string `json:"content-type"`
		ContentDisposition string `json:"content-disposition"`
	} `json:"params"`
	Size int    `json:"size"`
	Sha1 string `json:"sha1"`
}
