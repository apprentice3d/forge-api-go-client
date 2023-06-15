package dm

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// ListObjects returns the bucket contains along with details on each item.
func (api *BucketAPI) ListObjects(bucketKey, limit, beginsWith, startAt string) (result BucketContent, err error) {

	bearer, err := api.Authenticator.GetToken("data:read")
	if err != nil {
		return
	}

	result, err = listObjects(api.getPath(), bucketKey, limit, beginsWith, startAt, bearer.AccessToken)

	return
}

// DownloadObject downloads an on object, given the URL-encoded object name.
func (api *BucketAPI) DownloadObject(bucketKey string, objectName string) (result []byte, err error) {

	bearer, err := api.Authenticator.GetToken("data:read")
	if err != nil {
		return
	}

	downloadUrl, err := getSignedDownloadUrl(api.getPath(), bucketKey, objectName, bearer.AccessToken)
	if err != nil {
		return
	}

	result, err = downloadObjectUsingSignedUrl(&downloadUrl)

	return
}

// UploadObject adds to specified bucket the given data (can originate from a multipart-form or direct file read).
// Return details on uploaded object, including the object URN (> ObjectId). Check uploadOkResult struct.
func (api *BucketAPI) UploadObject(bucketKey, objectName, fileToUpload string) (result UploadResult, err error) {

	job, err := newUploadJob(api, bucketKey, objectName, fileToUpload)
	if err != nil {
		return
	}

	result, err = job.uploadFile()

	return
}

/*
 *	SUPPORT FUNCTIONS
 */

func listObjects(path, bucketKey, limit, beginsWith, startAt, token string) (result BucketContent, err error) {

	task := http.Client{}

	req, err := http.NewRequest("GET", path+"/"+bucketKey+"/objects", nil)

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
		content, _ := io.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	err = json.NewDecoder(response.Body).Decode(&result)

	return
}

func getSignedDownloadUrl(path, bucketKey, objectName string, token string) (result signedDownloadUrl, err error) {

	req, err := http.NewRequest("GET", path+"/"+bucketKey+"/objects/"+objectName+"/signeds3download", nil)
	if err != nil {
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)

	task := http.Client{}
	response, err := task.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		err = json.NewDecoder(response.Body).Decode(&result)
	} else {
		content, _ := io.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
	}

	return
}

func downloadObjectUsingSignedUrl(s *signedDownloadUrl) (result []byte, err error) {

	req, err := http.NewRequest("GET", s.Url, nil)
	if err != nil {
		return
	}

	task := http.Client{}
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

	result, err = io.ReadAll(response.Body)
	if err != nil {
		return
	}

	receivedSize := len(result)
	if receivedSize != s.Size {
		err = fmt.Errorf("the file size doesn't match, expected %v, but received %v", s.Size, receivedSize)
	}

	return
}
