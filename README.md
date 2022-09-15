[![GoDoc](https://godoc.org/github.com/woweh/forge-api-go-client?status.svg)](https://godoc.org/github.com/woweh/forge-api-go-client)
[![Go Report Card](https://goreportcard.com/badge/github.com/woweh/forge-api-go-client)](https://goreportcard.com/report/github.com/woweh/forge-api-go-client)

# forge-api-go-client


**Forge API:**  
[![oAuth2](https://img.shields.io/badge/oAuth2-v2-green.svg)](http://developer-autodesk.github.io/)
[![Data-Management](https://img.shields.io/badge/Data%20Management-v1-green.svg)](http://autodesk-forge.github.io/)
[![OSS](https://img.shields.io/badge/OSS-v2-green.svg)](http://autodesk-forge.github.io/)
[![Model-Derivative](https://img.shields.io/badge/Model%20Derivative-v2-green.svg)](http://autodesk-forge.github.io/)
[![Reality-Capture](https://img.shields.io/badge/Reality%20Capture-v1-green.svg)](http://developer-autodesk.github.io/)


Golang SDK for building Forge based applications.

This is a fork of the [Golang SDK by Denis Grigor](https://github.com/apprentice3d/forge-api-go-client).

The SDK has been extended with the following features:

### Model Derivative API:
- Add support for x-ads-headers and advanced translation options.  
  See: https://forge.autodesk.com/en/docs/model-derivative/v2/reference/http/jobs/job-POST/

### Data Management API:
- Update the upload object and download object to use the new direct-to-s3 approach.  
  Note that UploadObject method has a breaking change!  
  See:
  - https://forge.autodesk.com/blog/data-management-oss-object-storage-service-migrating-direct-s3-approach  
  - https://forge.autodesk.com/en/docs/data/v2/reference/http/buckets-:bucketKey-objects-:objectKey-signeds3upload-GET/
  - https://forge.autodesk.com/en/docs/data/v2/reference/http/buckets-:bucketKey-objects-:objectKey-signeds3download-GET/

TODO:
- Update ListBuckets to list all buckets (support paging).