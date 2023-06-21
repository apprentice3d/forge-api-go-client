package dm

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

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

const (
	// megaByte is 1048576 byte
	megaByte = 1 << 20

	// maxParts is the maximum number of parts returned by the "signeds3upload" endpoint.
	maxParts = 25

	// signedS3UploadEndpoint is the name of the signeds3upload endpoint.
	signedS3UploadEndpoint = "signeds3upload"

	// minutesExpiration is the expiration period of the signed upload URLs.
	// Autodesk default is 2 minutes (1 to 60 minutes).
	minutesExpiration = 60
)

var (
	// defaultChunkSize is the default size of download/upload chunks.
	//  NOTE:
	//  The minimum size seems to be 5 MB.
	//  Using a value < 5 MB causes errors when completing the upload [400, TooSmall].
	defaultChunkSize = int64(100 * megaByte)
)

func newUploadJob(api *OssAPI, bucketKey, objectName, fileToUpload string) (job uploadJob, err error) {

	fileInfo, err := os.Stat(fileToUpload)
	if err != nil {
		return
	}

	// Determine the required number of parts
	// - In the examples, typically a chunk size of 5 or 10 MB is used.
	// - In the old API, the boundary for multipart uploads was 100 MB.
	//   => See const defaultChunkSize
	totalParts := ceilingOfIntDivision(int(fileInfo.Size()), int(defaultChunkSize))
	numberOfBatches := ceilingOfIntDivision(totalParts, maxParts)

	job = uploadJob{
		api:               api,
		bucketKey:         bucketKey,
		objectKey:         objectName,
		fileToUpload:      fileToUpload,
		minutesExpiration: minutesExpiration,
		fileSize:          fileInfo.Size(),
		totalParts:        totalParts,
		numberOfBatches:   numberOfBatches,
		uploadKey:         "",
	}

	log.Println("New upload job:")
	log.Println("- bucketKey:", bucketKey)
	log.Println("- objectName:", objectName)
	log.Println("- fileToUpload:", fileToUpload)
	log.Println("- fileSize:", job.fileSize)
	log.Println("- totalParts:", job.totalParts)
	log.Println("- numberOfBatches:", job.numberOfBatches)

	return
}

func ceilingOfIntDivision(x, y int) int {
	// https://stackoverflow.com/a/54006084
	// https://stackoverflow.com/a/2745086
	return 1 + (x-1)/y
}

func (job *uploadJob) uploadFile() (result UploadResult, err error) {

	file, err := os.Open(job.fileToUpload)
	if err != nil {
		return
	}
	defer file.Close()

	log.Println("Start uploading file...")

	partsCounter := 0
	for i := 0; i < job.numberOfBatches; i++ {

		firstPart := (i * maxParts) + 1

		parts := job.getParts(partsCounter)

		// generate signed S3 upload url(s)
		log.Println("- getting signed URLs...")
		uploadUrls, err := job.getSignedUploadUrlsWithRetries(firstPart, parts)
		if err != nil {
			err = fmt.Errorf("error getting signed URLs for parts %v-%v :\n%w", firstPart, parts, err)
			return result, err
		}

		log.Println("- UploadKey: ", uploadUrls.UploadKey)
		log.Println("- number of signed URLs: ", len(uploadUrls.Urls))

		if i == 0 {
			// remember the uploadKey when requesting signed URLs for the first time
			job.uploadKey = uploadUrls.UploadKey
		}

		// upload the file in chunks to the signed url(s)
		for _, url := range uploadUrls.Urls {

			// read a chunk of the file
			bytesSlice := make([]byte, defaultChunkSize)

			bytesRead, err := file.Read(bytesSlice)
			if err != nil {
				if err != io.EOF {
					err = fmt.Errorf("error reading the file to upload:\n%w", err)
					return result, err
				}
				// EOF reached
			}

			// upload the chunk to the signed URL
			if bytesRead > 0 {
				buffer := bytes.NewBuffer(bytesSlice[:bytesRead])
				err = uploadChunkWithRetries(url, buffer)
				if err != nil {
					err = fmt.Errorf("error uploading a chunk to URL:\n- %v\n%w", url, err)
					return result, err
				}
				log.Println("- number of bytes sent: ", bytesRead)
			}
		}

		partsCounter += parts
	}

	// complete the upload
	log.Println("- completing upload...")
	result, err = job.completeUploadWithRetries()
	if err != nil {
		err = fmt.Errorf("error completing the upload:\n%w", err)
		return result, err
	}
	log.Println("Finished uploading the file:")
	log.Println("- ObjectId: ", result.ObjectId)
	log.Println("- Location: ", result.Location)
	log.Println("- Size: ", result.Size)

	return result, err
}

// getParts gets the number of parts that must be processed in this batch.
func (job *uploadJob) getParts(partsCounter int) int {

	parts := maxParts

	if job.totalParts < (partsCounter + maxParts) {
		// Say totalParts = 20:  part[0]=20, firstPart[0]=1
		// Say totalParts = 30:  part[0]=25, firstPart[0]=1, part[1]= 5, firstPart[1]=26
		// Say totalParts = 40:  part[0]=25, firstPart[0]=1, part[1]=15, firstPart[1]=26
		// Say totalParts = 50:  part[0]=25, firstPart[0]=1, part[1]=25, firstPart[1]=26
		parts = job.totalParts - partsCounter
	}

	return parts
}

// getSignedUploadUrlsWithRetries calls the signedS3UploadEndpoint to get parts signed URLs, retrying max 3 times.
// https://forge.autodesk.com/en/docs/data/v2/reference/http/buckets-:bucketKey-objects-:objectKey-signeds3upload-GET/
func (job *uploadJob) getSignedUploadUrlsWithRetries(firstPart, parts int) (result signedUploadUrls, err error) {

	var statusCode int

	for i := 0; i < 3; i++ {

		statusCode, result, err = job.getSignedUploadUrls(firstPart, parts)

		// 429 - RATE-LIMIT EXCEEDED
		// The maximum number of API calls that a Forge application can make _PER MINUTE_ was exceeded.
		// 500 - INTERNAL SERVER ERROR
		if statusCode == 429 || statusCode == 500 {
			// retry in 1 minute
			time.Sleep(1 * time.Minute)
		} else {
			// done
			break
		}
	}

	return result, err
}

// getSignedUploadUrls calls the signedS3UploadEndpoint to get parts signed URLs.
func (job *uploadJob) getSignedUploadUrls(firstPart, parts int) (statusCode int, result signedUploadUrls, err error) {

	// - https://forge.autodesk.com/en/docs/data/v2/reference/http/buckets-:bucketKey-objects-:objectKey-signeds3upload-GET/

	// - https://forge.autodesk.com/en/docs/data/v2/tutorials/upload-file/#step-4-generate-a-signed-s3-url
	// - https://forge.autodesk.com/en/docs/data/v2/tutorials/app-managed-bucket/#step-2-initiate-a-direct-to-s3-multipart-upload

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
	if job.uploadKey != "" {
		q.Add("uploadKey", job.uploadKey)
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

	statusCode = response.StatusCode

	if response.StatusCode != http.StatusOK {
		content, _ := io.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	err = json.NewDecoder(response.Body).Decode(&result)

	return
}

// uploadChunkWithRetries uploads a chunk of bytes to a given signedUrl, retrying max 3 times.
func uploadChunkWithRetries(signedUrl string, buffer *bytes.Buffer) (err error) {

	var (
		statusCode int

		// A backoff schedule for when and how often to retry failed HTTP
		// requests. The first element is the time to wait after the
		// first failure, the second the time to wait after the second
		// failure, etc. After reaching the last element, retries stop
		// and the request is considered failed.
		// https://brandur.org/fragments/go-http-retry
		backoffSchedule = []time.Duration{
			1 * time.Second,
			3 * time.Second,
			10 * time.Second,
		}
	)

	for _, backoff := range backoffSchedule {

		statusCode, err = uploadChunk(signedUrl, buffer)

		// Consider retrying (for example, with an exponential backoff) individual uploads when the
		// response code is 100-199, 429, or 500-599
		if (statusCode >= 100 && statusCode <= 199) ||
			statusCode == 429 ||
			(statusCode >= 500 && statusCode <= 599) {
			// retry
			time.Sleep(backoff)
		} else {
			// done
			break
		}
	}

	return err
}

// uploadChunk uploads a chunk of bytes to a given signedUrl.
func uploadChunk(signedUrl string, buffer *bytes.Buffer) (statusCode int, err error) {

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

	statusCode = response.StatusCode

	if response.StatusCode == http.StatusOK {
		return
	}

	content, _ := io.ReadAll(response.Body)
	err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))

	return
}

// completeUploadWithRetries instructs OSS to complete the object creation process after the bytes have been uploaded directly to S3, retrying max 3 times.
// https://forge.autodesk.com/en/docs/data/v2/reference/http/buckets-:bucketKey-objects-:objectKey-signeds3upload-POST/
func (job *uploadJob) completeUploadWithRetries() (result UploadResult, err error) {

	var statusCode int

	for i := 0; i < 3; i++ {

		statusCode, result, err = job.completeUpload()

		// 429 - RATE-LIMIT EXCEEDED
		// The maximum number of API calls that a Forge application can make _PER MINUTE_ was exceeded.
		// 500 - INTERNAL SERVER ERROR
		if statusCode == 429 || statusCode == 500 {
			// retry in 1 minute
			time.Sleep(1 * time.Minute)
		} else {
			// done
			break
		}
	}

	return result, err
}

// completeUpload instructs OSS to complete the object creation process after the bytes have been uploaded directly to S3.
// An object will not be accessible until this endpoint is called.
// This endpoint must be called within 24 hours of the upload beginning, otherwise the object will be discarded,
// and the upload must begin again from scratch.
func (job *uploadJob) completeUpload() (statusCode int, result UploadResult, err error) {

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
		UploadKey: job.uploadKey,
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

	statusCode = response.StatusCode

	if response.StatusCode != http.StatusOK {
		content, _ := io.ReadAll(response.Body)
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
	return job.api.Authenticator.GetHostPath() + path.Join(
		job.api.BucketApiPath, job.bucketKey, "objects", job.objectKey, signedS3UploadEndpoint,
	)
}

func (job *uploadJob) authenticate() (accessToken string, err error) {
	bearer, err := job.api.Authenticator.GetToken("data:write data:read")
	if err != nil {
		return
	}
	accessToken = bearer.AccessToken
	return
}
