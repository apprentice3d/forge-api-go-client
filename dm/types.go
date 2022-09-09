package dm

import (
	"time"

	"github.com/apprentice3d/forge-api-go-client/oauth"
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
	ObjectID    string            `json:"objectID"` // => urn = base64.RawStdEncoding.EncodeToString([]byte(ObjectID))
	ObjectKey   string            `json:"objectKey"`
	SHA1        string            `json:"sha1"`
	Size        uint64            `json:"size"`
	ContentType string            `json:"contentType, omitempty"`
	Location    string            `json:"location"`
	BlockSizes  []int64           `json:"blockSizes, omitempty"`
	Deltas      map[string]string `json:"deltas, omitempty"`
}

// UploadResult provides the OK/200 result of the completeUpload POST.
type UploadResult struct {
	BucketKey   string `json:"bucketKey"`
	ObjectId    string `json:"objectId"` // => urn = base64.RawStdEncoding.EncodeToString([]byte(ObjectID))
	ObjectKey   string `json:"objectKey"`
	Sha1        string `json:"sha1"` // this is only shown in an example, it's not in the documentation?!?!
	Size        int    `json:"size"`
	ContentType string `json:"content-type"`
	Location    string `json:"location"`
}

// BucketContent reflects the response when query Data Management API for bucket content.
type BucketContent struct {
	Items []ObjectDetails `json:"items"`
	Next  string          `json:"next"`
}

// signedUploadUrls provides the response from the signedS3UploadEndpoint
type signedUploadUrls struct {
	UploadKey        string    `json:"uploadKey"`
	UploadExpiration time.Time `json:"uploadExpiration"`
	UrlExpiration    time.Time `json:"urlExpiration"`
	Urls             []string  `json:"urls"`
}

// uploadJob provides information for uploading a file
type uploadJob struct {
	// the instance of the BucketAPI
	api BucketAPI
	// the key (= name) of the bucket where the file shall be stored
	bucketKey string
	// the key (= name) of the file in the Autodesk cloud (OSS)
	objectKey string
	// the path of the file to upload
	fileToUpload string
	// The custom expiration time within the 1 to 60 minutes range.
	minutesExpiration int
	// the size of the file to upload
	fileSize int64
	// the total number of parts to process
	totalParts int
	// if totalParts > maxParts, then we need to request signedUploadUrls multiple times
	numberOfBatches int
	// sha1 checksum of the uploaded file
	checkSum string
	// the slice with all signedUploadUrls that need processing
	batch []signedUploadUrls
}
