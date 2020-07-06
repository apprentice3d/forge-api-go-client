package dm

import (
	"io"
)

// UploadObject adds to specified bucket the given data (can originate from a multipart-form or direct file read).
// Return details on uploaded object, including the object URN. Check ObjectDetails struct.
func (api BucketAPI3L) UploadObject3L(bucketKey string, objectName string, reader io.Reader) (result ObjectDetails, err error) {
	if err = api.Token.RefreshTokenIfRequired(api.Auth); err != nil {
		return
	}

	path := api.Auth.Host + api.BucketsAPIPath
	return uploadObject(path, bucketKey, objectName, reader, api.Token.Bearer().AccessToken)
}

// DownloadObject returns the reader stream of the response body
// Don't forget to close it!
func (api BucketAPI3L) DownloadObject3L(bucketKey string, objectName string) (reader io.ReadCloser, err error) {
	if err = api.Token.RefreshTokenIfRequired(api.Auth); err != nil {
		return
	}

	path := api.Auth.Host + api.BucketsAPIPath
	return downloadObject(path, bucketKey, objectName, api.Token.Bearer().AccessToken)
}

// ListObjects returns the bucket contains along with details on each item.
func (api BucketAPI3L) ListObjects3L(bucketKey, limit, beginsWith, startAt string) (result BucketContent, err error) {
	if err = api.Token.RefreshTokenIfRequired(api.Auth); err != nil {
		return
	}

	path := api.Auth.Host + api.BucketsAPIPath
	return listObjects(path, bucketKey, limit, beginsWith, startAt, api.Token.Bearer().AccessToken)
}
