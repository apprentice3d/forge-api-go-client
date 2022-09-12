package dm

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
)

const (
	// megaByte is 1048576 bytes
	megaByte = 1 << 20
	// the maximum number of parts returned by the "signeds3upload" endpoint
	maxParts = 25
	// the name/ending of the signeds3upload endpoint
	signedS3UploadEndpoint = "signeds3upload"
)

var (
	// defaultSize is the default size of download/upload chunks.
	defaultSize = int64(100 * megaByte)
)

// UploadObject adds to specified bucket the given data (can originate from a multipart-form or direct file read).
// Return details on uploaded object, including the object URN (> ObjectId). Check uploadOkResult struct.
func (api BucketAPI) UploadObject(bucketKey, objectName, fileToUpload string) (result UploadResult, err error) {

	// Instructions for the S3 update:
	// https://forge.autodesk.com/blog/data-management-oss-object-storage-service-migrating-direct-s3-approach
	/*
			Direct-to-S3 approach for Data Management OSS
			To upload and download files, applications must generate a signed URL, then upload or download the binary. Here are the steps (pseudo code):

			Upload
			========

			1. Calculate the number of parts of the file to upload
				Note: Each uploaded part except for the last one must be at least 5MB (1024 * 5)

			2. Generate up to 25 URLs for uploading specific parts of the file using the
			   GET buckets/:bucketKey/objects/:objectKey/signeds3upload?firstPart=<index of first part>&parts=<number of parts>
		       endpoint.
				a) The part numbers start with 1
				b) For example, to generate upload URLs for parts 10 through 15, set firstPart to 10 and parts to 6
				c) This endpoint also returns an uploadKey that is used later to request additional URLs or to finalize the upload

			3. Upload remaining parts of the file to their corresponding upload URLs
				a) Consider retrying (for example, with an exponential backoff) individual uploads when the response code is 100-199, 429, or 500-599
				b) If the response code is 403, the upload URLs have expired; go back to step #2
				c) If you have used up all the upload URLs and there are still parts that must be uploaded, go back to step #2

			4. Finalize the upload using the POST buckets/:bucketKey/objects/:objectKey/signeds3upload endpoint, using the uploadKey value from step #2
	*/

	// initialize the uploadJob
	job := uploadJob{}
	job.api = api
	job.bucketKey = bucketKey
	job.objectKey = objectName
	job.fileToUpload = fileToUpload
	job.minutesExpiration = 60

	// Steps 1 & 2: Calculate the number of parts & generate signed URLs
	err = job.calculatePartsAndGetSignedUrls()
	if err != nil {
		return
	}

	// Step 3, upload the file(s) to the signed url(s)
	err = job.uploadFile()
	if err != nil {
		return
	}

	// Step 4, complete the upload
	result, err = job.completeUpload()

	return
}

// ListObjects returns the bucket contains along with details on each item.
func (api BucketAPI) ListObjects(bucketKey, limit, beginsWith, startAt string) (result BucketContent, err error) {
	bearer, err := api.Authenticator.GetToken("data:read")
	if err != nil {
		return
	}

	return listObjects(api.getPath(), bucketKey, limit, beginsWith, startAt, bearer.AccessToken)
}

func (api BucketAPI) getPath() string {
	return api.Authenticator.GetHostPath() + api.BucketAPIPath
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

//region Support Functions

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

func (job *uploadJob) calculatePartsAndGetSignedUrls() (err error) {

	fileInfo, err := os.Stat(job.fileToUpload)
	if err != nil {
		return
	}

	job.fileSize = fileInfo.Size()
	job.totalParts = int((job.fileSize / defaultSize) + 1)
	job.numberOfBatches = (job.totalParts / maxParts) + 1
	job.batch = make([]signedUploadUrls, 0)

	partsCounter := 0
	for i := 0; i < job.numberOfBatches; i++ {
		// Step 1, generate signed S3 upload url(s)
		firstPart := (i * maxParts) + 1

		parts := maxParts
		if job.totalParts < (partsCounter + maxParts) {
			// Say totalParts = 20:  part[0]=20, firstPart[0]=1
			// Say totalParts = 30:  part[0]=25, firstPart[0]=1, part[1]= 5, firstPart[1]=26
			// Say totalParts = 40:  part[0]=25, firstPart[0]=1, part[1]=15, firstPart[1]=26
			// Say totalParts = 50:  part[0]=25, firstPart[0]=1, part[1]=25, firstPart[1]=26
			parts = job.totalParts - partsCounter
		}

		uploadKey := ""
		if i > 0 {
			uploadKey = job.batch[i-1].UploadKey
		}

		uploadUrls, err := job.getSignedUploadUrls(uploadKey, firstPart, parts)
		if err != nil {
			return fmt.Errorf("Error getting signed URLs for parts %v-%v :\n%w", firstPart, firstPart+parts-1, err)
		}
		job.batch = append(job.batch, uploadUrls)

		partsCounter += maxParts
	}

	return
}

// getSignedUploadUrls calls the signedS3UploadEndpoint
func (job *uploadJob) getSignedUploadUrls(uploadKey string, firstPart, parts int) (result signedUploadUrls, err error) {

	// - https://forge.autodesk.com/en/docs/data/v2/reference/http/buckets-:bucketKey-objects-:objectKey-signeds3upload-GET/

	// - https://forge.autodesk.com/en/docs/data/v2/tutorials/upload-file/#step-4-generate-a-signed-s3-url
	// - https://forge.autodesk.com/en/docs/data/v2/tutorials/app-managed-bucket/#step-2-initiate-a-direct-to-s3-multipart-upload

	// 1 - determine the required number of parts
	// In the examples, typically a chunk size of 5 or 10 MB is used.
	// In the old API, the boundary for multipart uploads was 100 MB.

	accessToken, err := job.authenticate()
	if err != nil {
		return
	}

	// request the signed urls
	req, err := http.NewRequest("GET", job.getSignedS3UploadPath(), nil)
	if err != nil {
		return
	}

	addOrSetHeader(req, "Authorization", "Bearer "+accessToken)

	// appending to existing query args
	q := req.URL.Query()
	if uploadKey != "" {
		q.Add("uploadKey", uploadKey)
	}
	q.Add("firstPart", strconv.Itoa(firstPart))
	q.Add("parts", strconv.Itoa(parts))
	q.Add("minutesExpiration", strconv.Itoa(job.minutesExpiration))
	// assign encoded query string to http request
	req.URL.RawQuery = q.Encode()

	task := http.Client{}
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

	err = json.NewDecoder(response.Body).Decode(&result)

	return
}

func (job *uploadJob) uploadFile() (err error) {

	// - https://forge.autodesk.com/en/docs/data/v2/tutorials/upload-file/#step-5-upload-a-file-to-the-signed-url
	// - https://forge.autodesk.com/en/docs/data/v2/tutorials/app-managed-bucket/#step-3-split-the-file-and-upload

	file, err := os.Open(job.fileToUpload)
	if err != nil {
		return
	}
	defer file.Close()

	// to calculate the sha1 checksum
	sha1 := sha1.New()

	for _, uploadUrls := range job.batch {
		for _, url := range uploadUrls.Urls {

			bytesSlice := make([]byte, defaultSize)

			bytesRead, err := file.Read(bytesSlice)
			if err != nil {
				if err != io.EOF {
					return err
				}
				break
			}

			if bytesRead > 0 {
				buffer := bytes.NewBuffer(bytesSlice[:bytesRead])
				sha1.Write(buffer.Bytes())
				uploadChunk(url, buffer)
			}
		}
	}

	job.checkSum = fmt.Sprintf("%x", sha1.Sum(nil))

	return
}

func uploadChunk(signedUrl string, buffer *bytes.Buffer) (err error) {

	req, err := http.NewRequest("PUT", signedUrl, buffer)
	if err != nil {
		return
	}

	l := buffer.Len()
	req.ContentLength = int64(l)
	addOrSetHeader(req, "Content-Type", "application/octet-stream")
	addOrSetHeader(req, "Content-Length", strconv.Itoa(l))

	task := http.Client{}
	response, err := task.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		return
	}

	content, _ := ioutil.ReadAll(response.Body)
	err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))

	return
}

func (job *uploadJob) completeUpload() (result UploadResult, err error) {

	// - https://forge.autodesk.com/en/docs/data/v2/reference/http/buckets-:bucketKey-objects-:objectKey-signeds3upload-POST/

	// - https://forge.autodesk.com/en/docs/data/v2/tutorials/upload-file/#step-6-complete-the-upload
	// - https://forge.autodesk.com/en/docs/data/v2/tutorials/app-managed-bucket/#step-4-complete-the-upload

	accessToken, err := job.authenticate()
	if err != nil {
		return
	}

	// size	integer: The expected size of the uploaded object.
	// If provided, OSS will check this against the blob in S3 and return an error if the size does not match.
	bodyData := struct {
		UploadKey string `json:"uploadKey"`
		Size      int    `json:"size"`
	}{
		UploadKey: job.batch[0].UploadKey,
		Size:      int(job.fileSize),
	}

	bodyJson, err := json.Marshal(bodyData)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", job.getSignedS3UploadPath(), bytes.NewBuffer(bodyJson))
	if err != nil {
		return
	}

	addOrSetHeader(req, "Authorization", "Bearer "+accessToken)
	addOrSetHeader(req, "Content-Type", "application/json")
	addOrSetHeader(req, "x-ads-meta-Content-Type", "application/octet-stream")

	task := http.Client{}
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

	err = json.NewDecoder(response.Body).Decode(&result)

	return
}

func addOrSetHeader(req *http.Request, key, value string) {
	if req.Header.Get(key) == "" {
		req.Header.Add(key, value)
	} else {
		req.Header.Set(key, value)
	}
}

func (job *uploadJob) getSignedS3UploadPath() string {
	// https://developer.api.autodesk.com/oss/v2/buckets/:bucketKey/objects/:objectKey/signeds3upload
	// :bucketKey/objects/:objectKey/signeds3upload
	return job.api.Authenticator.GetHostPath() + path.Join(job.api.BucketAPIPath, job.bucketKey, "objects", job.objectKey, signedS3UploadEndpoint)
}

func (job *uploadJob) authenticate() (accessToken string, err error) {
	bearer, err := job.api.Authenticator.GetToken("data:write data:read")
	if err != nil {
		return
	}
	accessToken = bearer.AccessToken
	return
}

//endregion

//region Old/Deprecated Support Functions

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

//endregion
