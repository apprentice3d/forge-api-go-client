package dm

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

// UploadObject adds to specified bucket the given data (can originate from a multipart-form or direct file read).
// Return details on uploaded object, including the object URN. Check ObjectDetails struct.
func (api BucketAPI) UploadObject(bucketKey, objectName, fileToUpload string) (result ObjectDetails, err error) {
	bearer, err := api.Authenticator.GetToken("data:write data:read")
	if err != nil {
		return
	}
	path := api.Authenticator.GetHostPath() + api.BucketAPIPath

	// Step 1, generate signed S3 upload url(s)

	// Step 2, upload the file(s) to the signed url(s)
	// - https://forge.autodesk.com/en/docs/data/v2/tutorials/upload-file/#step-5-upload-a-file-to-the-signed-url
	// - https://forge.autodesk.com/en/docs/data/v2/tutorials/app-managed-bucket/#step-3-split-the-file-and-upload

	// Step 3, complete the upload
	// - https://forge.autodesk.com/en/docs/data/v2/tutorials/upload-file/#step-6-complete-the-upload
	// - https://forge.autodesk.com/en/docs/data/v2/tutorials/app-managed-bucket/#step-4-complete-the-upload
	// - https://forge.autodesk.com/en/docs/data/v2/reference/http/buckets-:bucketKey-objects-:objectKey-signeds3upload-POST/

	return
}

// ListObjects returns the bucket contains along with details on each item.
func (api BucketAPI) ListObjects(bucketKey, limit, beginsWith, startAt string) (result BucketContent, err error) {
	bearer, err := api.Authenticator.GetToken("data:read")
	if err != nil {
		return
	}
	path := api.Authenticator.GetHostPath() + api.BucketAPIPath

	return listObjects(path, bucketKey, limit, beginsWith, startAt, bearer.AccessToken)
}

// DownloadObject downloads an on object, given the URL-encoded object name.
func (api BucketAPI) DownloadObject(bucketKey string, objectName string) (result []byte, err error) {
	bearer, err := api.Authenticator.GetToken("data:read")
	if err != nil {
		return
	}
	path := api.Authenticator.GetHostPath() + api.BucketAPIPath

	return downloadObject(path, bucketKey, objectName, bearer.AccessToken)
}

/*
 *	SUPPORT FUNCTIONS
 */

func listObjects(path, bucketKey, limit, beginsWith, startAt, token string) (result BucketContent, err error) {
	task := http.Client{}

	req, err := http.NewRequest("GET",
		path+"/"+bucketKey+"/objects",
		nil,
	)

	if err != nil {
		return
	}

	params := req.URL.Query()
	if len(beginsWith) != 0 {
		params.Add("beginsWith", beginsWith)
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

func getSignedUploadUrls(path, bucketKey, objectName, fileToUpload string, minutesExpiration int) (result PreSignedUploadUrls, err error) {

	// - https://forge.autodesk.com/en/docs/data/v2/tutorials/upload-file/#step-4-generate-a-signed-s3-url
	// - https://forge.autodesk.com/en/docs/data/v2/tutorials/app-managed-bucket/#step-2-initiate-a-direct-to-s3-multipart-upload
	// - https://forge.autodesk.com/en/docs/data/v2/reference/http/buckets-:bucketKey-objects-:objectKey-signeds3upload-GET/

	// 1 - determine the required number of parts
	// In the examples, typically a chunk size of 5 or 10 MB is used.
	// In the old API, the boundary for multipart uploads was 100 MB.

	return
}

func uploadObject(path, bucketKey, objectName string, data []byte, token string) (result ObjectDetails, err error) {

	task := http.Client{}

	dataContent := bytes.NewReader(data)
	req, err := http.NewRequest("PUT",
		path+"/"+bucketKey+"/objects/"+objectName,
		dataContent)

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

func downloadObject(path, bucketKey, objectName string, token string) (result []byte, err error) {

	task := http.Client{}

	req, err := http.NewRequest("GET",
		path+"/"+bucketKey+"/objects/"+objectName,
		nil)

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

	result, err = ioutil.ReadAll(response.Body)

	return

}
