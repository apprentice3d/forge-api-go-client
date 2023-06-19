[![GoDoc](https://godoc.org/github.com/woweh/forge-api-go-client?status.svg)](https://godoc.org/github.com/woweh/forge-api-go-client)
[![Go Report Card](https://goreportcard.com/badge/github.com/woweh/forge-api-go-client)](https://goreportcard.com/report/github.com/woweh/forge-api-go-client)

# forge-api-go-client


**Autodesk Platform Services APIs:**  
[![Authentication](https://img.shields.io/badge/Authentication-v2-green.svg)](https://aps.autodesk.com/en/docs/oauth/v2/developers_guide/overview/)
[![Data-Management](https://img.shields.io/badge/Data%20Management-v2-green.svg)](https://aps.autodesk.com/en/docs/data/v2/developers_guide/)
[![Model-Derivative](https://img.shields.io/badge/Model%20Derivative-v2-green.svg)](https://aps.autodesk.com/en/docs/model-derivative/v2/developers_guide/)


Golang client for building APS based applications ([Autodesk Platform Services], formerly *"Forge"*).

This is a fork of the [forge-api-go-client by Denis Grigor](https://github.com/apprentice3d/forge-api-go-client).  
The original client is no longer maintained.  
This forks has been updated and extended with new features.

---
## Supported and maintained APIs
Note that this client only covers a subset of the APS APIs.  

At the time of writing (2023/06/15), only the following APIs are maintained:
1. Authentication (oauth)
2. Data Management (dm)
3. Model Derivative (md)

The following APIs are not maintained:
1. Reality Capture (rc)
2. Design Automation (da)

You are invited to contribute!  
Please fork and add missing APIs.

---
## Updates from the original client

The client has been extended with the following features:


### Authentication (oauth):
- Update to OAuth V2.  
  See: https://aps.autodesk.com/blog/migration-guide-oauth2-v1-v2


### Data Management API (dm):
- Update the upload object and download object to use the direct-to-s3 approach.  
  Note that UploadObject method has a breaking change!  
  See:
  - https://forge.autodesk.com/blog/data-management-oss-object-storage-service-migrating-direct-s3-approach
  - https://forge.autodesk.com/en/docs/data/v2/reference/http/buckets-:bucketKey-objects-:objectKey-signeds3upload-GET/
  - https://forge.autodesk.com/en/docs/data/v2/reference/http/buckets-:bucketKey-objects-:objectKey-signeds3download-GET/
- Add support for regions (US <> EMEA):  
  See: https://aps.autodesk.com/blog/data-management-and-model-derivative-regions


### Model Derivative API (md):
- Add support for x-ads-headers and advanced translation options (input > formats > advanced).  
  See: https://forge.autodesk.com/en/docs/model-derivative/v2/reference/http/jobs/job-POST/
- Update GET derivatives to use the `signedcookies` endpoint.  
  See:
  - https://aps.autodesk.com/blog/data-management-oss-object-storage-service-migrating-direct-s3-approach
  - https://aps.autodesk.com/en/docs/model-derivative/v2/reference/http/urn-manifest-derivativeUrn-signedcookies-GET/
- Add support for regions (US <> EMEA):  
  See: https://aps.autodesk.com/blog/data-management-and-model-derivative-regions
- Add support for downloading all properties.  
  See: https://aps.autodesk.com/en/docs/model-derivative/v2/reference/http/metadata/urn-metadata-guid-properties-GET/ 


---

## To Do:
- (dm) Update ListBuckets to list all buckets (support paging).
- (md) Add support for fetching the object tree.

[Autodesk Platform Services]: https://aps.autodesk.com/