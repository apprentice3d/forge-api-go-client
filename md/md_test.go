package md_test

import (
	"testing"
	"os"
	"github.com/apprentice3d/forge-api-go-client/dm"
	"github.com/apprentice3d/forge-api-go-client/md"
	"io/ioutil"
	"encoding/json"
	"encoding/base64"
	"bytes"
)

func TestAPI_TranslateToSVF(t *testing.T) {
	// prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")
	bucketAPI := dm.NewBucketAPIWithCredentials(clientID, clientSecret)
	mdAPI := md.NewAPIWithCredentials(clientID, clientSecret)

	tempBucketName := "go_testing_md_bucket"
	testFilePath := "../assets/HelloWorld.rvt"

	var testObject dm.ObjectDetails

	t.Run("Create a temporary bucket", func(t *testing.T) {
		_, err := bucketAPI.CreateBucket(tempBucketName, "transient")

		if err != nil {
			t.Errorf("Failed to create a bucket: %s\n", err.Error())
		}
	})

	t.Run("Get bucket details", func(t *testing.T) {
		_, err := bucketAPI.GetBucketDetails(tempBucketName)

		if err != nil {
			t.Fatalf("Failed to get bucket details: %s\n", err.Error())
		}
	})

	t.Run("Upload an object into temp bucket", func(t *testing.T) {
		file, err := os.Open(testFilePath)
		if err != nil {
			t.Fatal("Cannot open testfile for reading")
		}
		defer file.Close()
		data, err := ioutil.ReadAll(file)
		if err != nil {
			t.Fatal("Cannot read the testfile")
		}

		testObject, err = bucketAPI.UploadObject(tempBucketName, "temp_file.rvt", data)

		if err != nil {
			t.Fatal("Could not upload the test object, got: ", err.Error())
		}

		if testObject.Size == 0 {
			t.Fatal("The test object was uploaded but it is zero-sized")
		}
	})

	t.Run("Translate object into SVF", func(t *testing.T) {

		result, err := mdAPI.TranslateToSVF(testObject.ObjectID)

		if err != nil {
			t.Error("Could not translate the test object, got: ", err.Error())
		}

		if result.Result != "created" {
			t.Error("The test object was uploaded, but failed to create the translation job")
		}
	})

	t.Run("Delete the temporary bucket", func(t *testing.T) {
		err := bucketAPI.DeleteBucket(tempBucketName)

		if err != nil {
			t.Fatalf("Failed to delete bucket: %s\n", err.Error())
		}
	})
}

func TestAPI_TranslateToSVF2_JSON_Creation(t *testing.T) {

	params := md.TranslationSVFPreset
	params.Input.URN = base64.RawStdEncoding.EncodeToString([]byte("just a test urn"))

	output, err := json.Marshal(&params)
	if err != nil {
		t.Fatal("Could not marshal the preset into JSON: ", err.Error())
	}

	referenceExample := `
{
        "input": {
          "urn": "anVzdCBhIHRlc3QgdXJu"
        },
        "output": {
			"destination": {
        		"region": "us"
      		},
          	"formats": [
            {
              "type": "svf",
              "views": [
                "2d",
                "3d"
              ]
            }
          ]
        }
      }
`

	var example md.TranslationParams
	err = json.Unmarshal([]byte(referenceExample), &example)
	if err != nil {
		t.Fatal("Could not unmarshal the reference example: ", err.Error())
	}

	expected, err := json.Marshal(example)
	if err != nil {
		t.Fatal("Could not marshal the reference example into JSON: ", err.Error())
	}

	if bytes.Compare(expected, output) != 0 {
		t.Fatalf("The translation params are not correct:\nexpected: %s\n created: %s",
			string(expected),
			string(output))

	}

}

func TestParseManifest(t *testing.T) {
	t.Run("Parse pending manifest", func(t *testing.T) {
		manifest := `
			{
			  "type": "manifest",
			  "hasThumbnail": "false",
			  "status": "pending",
			  "progress": "0% complete",
			  "region": "US",
			  "urn": "dXJuOmFkc2sub2JqZWN0czpvcy5vYmplY3Q6bW9kZWxkZXJpdmF0aXZlL0E1LnppcA",
			  "derivatives": [
			  ]
			}
			`
		var decodedManifest md.Manifest
		err := json.Unmarshal([]byte(manifest), &decodedManifest)
		if err != nil {
			t.Error(err.Error())
		}

		if len(decodedManifest.Derivatives) != 0 {
			t.Error("There should not be derivatives")
		}

	})

	t.Run("Parse in progress manifest", func(t *testing.T) {
		manifest := `
			{
				  "type": "manifest",
				  "hasThumbnail": "true",
				  "status": "inprogress",
				  "progress": "99% complete",
				  "region": "US",
				  "urn": "dXJuOmFkc2sub2JqZWN0czpvcy5vYmplY3Q6bW9kZWxkZXJpdmF0aXZlL0E1LnppcA",
				  "derivatives": [
					{
					  "name": "A5.iam",
					  "hasThumbnail": "true",
					  "status": "success",
					  "progress": "99% complete",
					  "outputType": "svf",
					  "children": [
						{
						  "guid": "d998268f-eeb4-da87-0db4-c5dbbc4926d0",
						  "type": "geometry",
						  "role": "3d",
						  "name": "Scene",
						  "status": "success",
						  "progress": "99% complete",
						  "hasThumbnail": "true",
						  "children": [
							{
							  "guid": "4f981e94-8241-4eaf-b08b-cd337c6b8b1f",
							  "type": "resource",
							  "progress": "99% complete",
							  "role": "graphics",
							  "mime": "application/autodesk-svf"
							},
							{
							  "guid": "d718eb7e-fa8a-42f9-8b32-e323c0fbea0c",
							  "type": "resource",
							  "urn": "urn:adsk.viewing:fs.file:dXJuOmFkc2sub2JqZWN0czpvcy5vYmplY3Q6bW9kZWxkZXJpdmF0aXZlL0E1LnppcA/output/1/A5.svf.png01_thumb_400x400.png",
							  "resolution": [
								400.0,
								400.0
							  ],
							  "mime": "image/png",
							  "role": "thumbnail"
							},
							{
							  "guid": "34dc340b-835f-47f7-9da5-b8219aefe741",
							  "type": "resource",
							  "urn": "urn:adsk.viewing:fs.file:dXJuOmFkc2sub2JqZWN0czpvcy5vYmplY3Q6bW9kZWxkZXJpdmF0aXZlL0E1LnppcA/output/1/A5.svf.png01_thumb_200x200.png",
							  "resolution": [
								200.0,
								200.0
							  ],
							  "mime": "image/png",
							  "role": "thumbnail"
							},
							{
							  "guid": "299c6ba6-650e-423e-bbd6-3aaff44ee104",
							  "type": "resource",
							  "urn": "urn:adsk.viewing:fs.file:dXJuOmFkc2sub2JqZWN0czpvcy5vYmplY3Q6bW9kZWxkZXJpdmF0aXZlL0E1LnppcA/output/1/A5.svf.png01_thumb_100x100.png",
							  "resolution": [
								100.0,
								100.0
							  ],
							  "mime": "image/png",
							  "role": "thumbnail"
							}
						  ]
						}
					  ]
					}
				  ]
			}
			`
		var decodedManifest md.Manifest
		err := json.Unmarshal([]byte(manifest), &decodedManifest)
		if err != nil {
			t.Error(err.Error())
		}

		if len(decodedManifest.Derivatives) != 1 {
			t.Error("Failed to parse derivatives")
		}

		if len(decodedManifest.Derivatives[0].Children) != 1 {
			t.Error("Failed to parse childern derivatives")
		}

		if len(decodedManifest.Derivatives[0].Children[0].Children) != 4 {
			t.Error("Failed to parse childern of derivative's children [funny]")
		}

		if decodedManifest.Derivatives[0].Children[0].Children[0].URN != "" {
			child := decodedManifest.Derivatives[0].Children[0].Children[0]
			t.Errorf("URN should be empty: %s => %s", child.Name, child.URN)
		}

	})

	t.Run("Parse complete failed manifest", func(t *testing.T) {
		manifest := `
			{
			  "type": "manifest",
			  "hasThumbnail": "false",
			  "status": "failed",
			  "progress": "complete",
			  "region": "US",
			  "urn": "dXJuOmFkc2sub2JqZWN0czpvcy5vYmplY3Q6bW9kZWxkZXJpdmF0aXZlL0E1LnppcA",
			  "derivatives": [
				{
				  "name": "A5.iam",
				  "hasThumbnail": "false",
				  "status": "failed",
				  "progress": "complete",
				  "messages": [
					{
					  "type": "warning",
					  "message": "The drawing's thumbnails were not properly created.",
					  "code": "TranslationWorker-ThumbnailGenerationFailed"
					}
				  ],
				  "outputType": "svf",
				  "children": [
					{
					  "guid": "d998268f-eeb4-da87-0db4-c5dbbc4926d0",
					  "type": "geometry",
					  "role": "3d",
					  "name": "Scene",
					  "status": "success",
					  "messages": [
						{
						  "type": "warning",
						  "code": "ATF-1023",
						  "message": [
							"The file: {0} does not exist.",
							"C:\\Users\\ADSK\\Documents\\A5\\Top.ipt"
						  ]
						},
						{
						  "type": "warning",
						  "code": "ATF-1023",
						  "message": [
							"The file: {0} does not exist.",
							"C:\\Users\\ADSK\\Documents\\A5\\Bottom.ipt"
						  ]
						},
						{
						  "type": "error",
						  "code": "ATF-1026",
						  "message": [
							"The file: {0} is empty.",
							"C:/worker/viewing-inventor-lmv/tmp/job-1/5/output/1/A5.svf"
						  ]
						}
					  ],
					  "progress": "complete",
					  "hasThumbnail": "false",
					  "children": [
						{
						  "guid": "4f981e94-8241-4eaf-b08b-cd337c6b8b1f",
						  "type": "resource",
						  "urn": "urn:adsk.viewing:fs.file:dXJuOmFkc2sub2JqZWN0czpvcy5vYmplY3Q6bW9kZWxkZXJpdmF0aXZlL0E1LnppcA/output/1/A5.svf",
						  "role": "graphics",
						  "mime": "application/autodesk-svf"
						}
					  ]
					}
				  ]
				}
			  ]
			}
			`
		var decodedManifest md.Manifest
		err := json.Unmarshal([]byte(manifest), &decodedManifest)
		if err != nil {
			t.Error(err.Error())
		}

		if len(decodedManifest.Derivatives) != 1 {
			t.Error("Failed to parse derivatives")
		}

		if len(decodedManifest.Derivatives[0].Children) != 1 {
			t.Error("Failed to parse childern derivatives")
		}

		if len(decodedManifest.Derivatives[0].Children[0].Children) != 1 {
			t.Error("Failed to parse childern of derivative's children [funny]")
		}

		if decodedManifest.Derivatives[0].Children[0].Children[0].URN == "" {
			t.Error("URN should not be empty")
		}

		if decodedManifest.Derivatives[0].Messages[0].Type != "warning" {
			t.Error("Chould contain a warning message")
		}

		if len(decodedManifest.Derivatives[0].Children[0].Messages) != 3 {
			t.Error("Derivative child should contain 3 error message")
		}

		if decodedManifest.Derivatives[0].Children[0].Messages[0].Type != "warning" {
			t.Error("Derivative child message should be a warning message")
		}
		if decodedManifest.Derivatives[0].Children[0].Messages[2].Type != "error" {
			t.Error("Derivative child message should be an error message")
		}

		if len(decodedManifest.Derivatives[0].Children[0].Messages[2].Message) != 2 {
			t.Error("Derivative child message should contain 2 message descriptions")
		}

		if decodedManifest.Derivatives[0].Children[0].Children[0].Role != "graphics" {
			t.Error("Failed to parse childern of derivative's children [funny]")
		}

	})

	t.Run("Parse complete success manifest", func(t *testing.T) {
		manifest := `
			{
			  "type": "manifest",
			  "hasThumbnail": "true",
			  "status": "success",
			  "progress": "complete",
			  "region": "US",
			  "urn": "dXJuOmFkc2sub2JqZWN0czpvcy5vYmplY3Q6bW9kZWxkZXJpdmF0aXZlL0E1LnppcA",
			  "derivatives": [
				{
				  "name": "A5.iam",
				  "hasThumbnail": "true",
				  "status": "success",
				  "progress": "complete",
				  "outputType": "svf",
				  "children": [
					{
					  "guid": "d998268f-eeb4-da87-0db4-c5dbbc4926d0",
					  "type": "geometry",
					  "role": "3d",
					  "name": "Scene",
					  "status": "success",
					  "progress": "complete",
					  "hasThumbnail": "true",
					  "children": [
						{
						  "guid": "4f981e94-8241-4eaf-b08b-cd337c6b8b1f",
						  "type": "resource",
						  "urn": "urn:adsk.viewing:fs.file:dXJuOmFkc2sub2JqZWN0czpvcy5vYmplY3Q6bW9kZWxkZXJpdmF0aXZlL0E1LnppcA/output/1/A5.svf",
						  "role": "graphics",
						  "mime": "application/autodesk-svf"
						},
						{
						  "guid": "d718eb7e-fa8a-42f9-8b32-e323c0fbea0c",
						  "type": "resource",
						  "urn": "urn:adsk.viewing:fs.file:dXJuOmFkc2sub2JqZWN0czpvcy5vYmplY3Q6bW9kZWxkZXJpdmF0aXZlL0E1LnppcA/output/1/A5.svf.png01_thumb_400x400.png",
						  "resolution": [
							400.0,
							400.0
						  ],
						  "mime": "image/png",
						  "role": "thumbnail"
						},
						{
						  "guid": "34dc340b-835f-47f7-9da5-b8219aefe741",
						  "type": "resource",
						  "urn": "urn:adsk.viewing:fs.file:dXJuOmFkc2sub2JqZWN0czpvcy5vYmplY3Q6bW9kZWxkZXJpdmF0aXZlL0E1LnppcA/output/1/A5.svf.png01_thumb_200x200.png",
						  "resolution": [
							200.0,
							200.0
						  ],
						  "mime": "image/png",
						  "role": "thumbnail"
						},
						{
						  "guid": "299c6ba6-650e-423e-bbd6-3aaff44ee104",
						  "type": "resource",
						  "urn": "urn:adsk.viewing:fs.file:dXJuOmFkc2sub2JqZWN0czpvcy5vYmplY3Q6bW9kZWxkZXJpdmF0aXZlL0E1LnppcA/output/1/A5.svf.png01_thumb_100x100.png",
						  "resolution": [
							100.0,
							100.0
						  ],
						  "mime": "image/png",
						  "role": "thumbnail"
						}
					  ]
					},
					{
					  "guid": "b86dcf4d-dd4e-561a-1b52-50ee01f7af4f",
					  "hasThumbnail": "true",
					  "progress": "complete",
					  "role": "2d",
					  "status": "success",
					  "type": "geometry",
					  "children": [
						{
						  "guid": "cfe81eb4-fbc6-17c0-beba-3ab845d228f0",
						  "mime": "image/png",
						  "resolution": [
							100,
							100
						  ],
						  "role": "thumbnail",
						  "status": "success",
						  "type": "resource",
						  "urn": "urn:adsk.viewing:fs.file:dXJuOmFkc2sub2JqZWN0czpvcy5vYmplY3Q6bW9kZWxkZXJpdmF0aXZlL0E1LnppcA/output/661c6096-056d-e58c-6c87-38769662932f_f2d/02___Floor1.png"
						},
						{
						  "guid": "03c34714-36c7-b2bf-eb19-245f26c15e50",
						  "mime": "image/png",
						  "resolution": [
							200,
							200
						  ],
						  "role": "thumbnail",
						  "status": "success",
						  "type": "resource",
						  "urn": "urn:adsk.viewing:fs.file:dXJuOmFkc2sub2JqZWN0czpvcy5vYmplY3Q6bW9kZWxkZXJpdmF0aXZlL0E1LnppcA/output/661c6096-056d-e58c-6c87-38769662932f_f2d/02___Floor2.png"
						},
						{
						  "guid": "b680b9ec-5240-6858-b7ef-7e9adafd9d9a",
						  "mime": "image/png",
						  "resolution": [
							400,
							400
						  ],
						  "role": "thumbnail",
						  "status": "success",
						  "type": "resource",
						  "urn": "urn:adsk.viewing:fs.file:dXJuOmFkc2sub2JqZWN0czpvcy5vYmplY3Q6bW9kZWxkZXJpdmF0aXZlL0E1LnppcA/output/661c6096-056d-e58c-6c87-38769662932f_f2d/02___Floor4.png"
						},
						{
						  "guid": "a81433d1-e3e7-97f8-17f2-e85c1bbc1f66",
						  "mime": "application/autodesk-f2d",
						  "role": "graphics",
						  "status": "success",
						  "type": "resource",
						  "urn": "urn:adsk.viewing:fs.file:dXJuOmFkc2sub2JqZWN0czpvcy5vYmplY3Q6bW9kZWxkZXJpdmF0aXZlL0E1LnppcA/output/661c6096-056d-e58c-6c87-38769662932f_f2d/primaryGraphics.f2d"
						},
						{
						  "guid": "5d2d63c3-943e-4111-b0fe-75abfeb85cb8",
						  "name": "Floor Plan: 02 - Floor",
						  "role": "2d",
						  "type": "view",
						  "viewbox": [
							0,
							0,
							279.4,
							215.9
						  ]
						}
					  ]
					}
				  ]
				},
				{
				  "status": "success",
				  "progress": "complete",
				  "outputType": "step",
				  "children": [
					{
					  "guid": "a6128518-dcf0-967b-31a1-3439a375daeb",
					  "role": "STEP",
					  "mime": "application/octet-stream",
					  "urn": "urn:adsk.viewing:fs.file:dXJuOmFkc2sub2JqZWN0czpvcy5vYmplY3Q6bW9kZWxkZXJpdmF0aXZlL0E1LnppcA/output/1/A5.stp",
					  "status": "success",
					  "type": "resource"
					}
				  ]
				},
				{
				  "name": "A5.iam",
				  "hasThumbnail": "true",
				  "status": "success",
				  "progress": "complete",
				  "outputType": "thumbnail",
				  "children": [
					{
					  "guid": "63c50197-c285-411b-bcfd-b3f19b1d37ef",
					  "mime": "image/png",
					  "resolution": [
						256,
						256
					  ],
					  "role": "thumbnail",
					  "type": "resource",
					  "urn": "urn:adsk.viewing:fs.file:dXJuOmFkc2sub2JqZWN0czpvcy5vYmplY3Q6bW9kZWxkZXJpdmF0aXZlL0E1LnppcA/output/256x256.png"
					}
				  ]
				},
				{
				  "status": "success",
				  "progress": "complete",
				  "outputType": "obj",
				  "children": [
					{
					  "guid": "1122e136-ea24-31ee-a7ef-ad065fafad42",
					  "type": "resource",
					  "role": "obj",
					  "modelGUID": "4f981e94-8241-4eaf-b08b-cd337c6b8b1f",
					  "objectIds": [
						2,
						3,
						4
					  ],
					  "status": "success",
					  "urn": "urn:adsk.viewing:fs.file:dXJuOmFkc2sub2JqZWN0czpvcy5vYmplY3Q6bW9kZWxkZXJpdmF0aXZlL0E1LnppcA/output/geometry/bc3339b2-73cd-4fba-9cb3-15363703a354.obj"
					},
					{
					  "guid": "29c1c0d4-7a35-350a-b3e5-fb221b054e29",
					  "type": "resource",
					  "role": "obj",
					  "modelGUID": "4f981e94-8241-4eaf-b08b-cd337c6b8b1f",
					  "objectIds": [
						2,
						3,
						4
					  ],
					  "status": "success",
					  "urn": "urn:adsk.viewing:fs.file:dXJuOmFkc2sub2JqZWN0czpvcy5vYmplY3Q6bW9kZWxkZXJpdmF0aXZlL0E1LnppcA/output/geometry/bc3339b2-73cd-4fba-9cb3-15363703a354.mtl"
					},
					{
					  "guid": "3e9752f1-5989-38b1-bff1-1f2d81841c8a",
					  "type": "resource",
					  "role": "obj",
					  "modelGUID": "4f981e94-8241-4eaf-b08b-cd337c6b8b1f",
					  "objectIds": [
						2,
						3,
						4
					  ],
					  "status": "success",
					  "urn": "urn:adsk.viewing:fs.file:dXJuOmFkc2sub2JqZWN0czpvcy5vYmplY3Q6bW9kZWxkZXJpdmF0aXZlL0E1LnppcA/output/geometry/bc3339b2-73cd-4fba-9cb3-15363703a354.zip"
					}
				  ]
				}
			  ]
			}
			`
		var decodedManifest md.Manifest
		err := json.Unmarshal([]byte(manifest), &decodedManifest)
		if err != nil {
			t.Error(err.Error())
		}

		if len(decodedManifest.Derivatives) != 4 {
			t.Error("Failed to parse derivatives")
		}

		if len(decodedManifest.Derivatives[0].Children) == 1 {
			t.Errorf("Failed to parse childern derivatives, expecting 1, got %d",
				len(decodedManifest.Derivatives[0].Children))
		}

		if len(decodedManifest.Derivatives[0].Children[0].Children) != 4 {
			t.Errorf("Failed to parse childern of derivative's children [funny], expecting 4, got %d",
				len(decodedManifest.Derivatives[0].Children[0].Children))
		}

		if decodedManifest.Derivatives[0].Children[0].Children[0].URN == "" {
			t.Error("URN should not be empty")
		}

		if len(decodedManifest.Derivatives[0].Messages) != 0 {
			t.Error("Derivative should not contain any error messages")
		}

		expectedOutputTypes := []string{"svf", "step", "thumbnail", "obj"}

		for idx := range decodedManifest.Derivatives {
			if decodedManifest.Derivatives[idx].OutputType != expectedOutputTypes[idx] {
				t.Error("Wrong derivative type parsing: expectd %s, got %s",
					decodedManifest.Derivatives[idx].OutputType,
					expectedOutputTypes[idx],
				)
			}
		}

	})

}


func TestModelDerivativeAPI_GetManifest(t *testing.T) {
	// prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")
	bucketAPI := dm.NewBucketAPIWithCredentials(clientID, clientSecret)
	mdAPI := md.NewAPIWithCredentials(clientID, clientSecret)

	tempBucketName := "go_testing_md_bucket"
	testFilePath := "../assets/HelloWorld.rvt"

	var testObject dm.ObjectDetails
	var translationResult md.TranslationResult

	t.Run("Create a temporary bucket", func(t *testing.T) {
		_, err := bucketAPI.CreateBucket(tempBucketName, "transient")

		if err != nil {
			t.Errorf("Failed to create a bucket: %s\n", err.Error())
		}
	})

	t.Run("Get bucket details", func(t *testing.T) {
		_, err := bucketAPI.GetBucketDetails(tempBucketName)

		if err != nil {
			t.Fatalf("Failed to get bucket details: %s\n", err.Error())
		}
	})

	t.Run("Upload an object into temp bucket", func(t *testing.T) {
		file, err := os.Open(testFilePath)
		if err != nil {
			t.Fatal("Cannot open testfile for reading")
		}
		defer file.Close()
		data, err := ioutil.ReadAll(file)
		if err != nil {
			t.Fatal("Cannot read the testfile")
		}

		testObject, err = bucketAPI.UploadObject(tempBucketName, "temp_file.rvt", data)

		if err != nil {
			t.Fatal("Could not upload the test object, got: ", err.Error())
		}

		if testObject.Size == 0 {
			t.Fatal("The test object was uploaded but it is zero-sized")
		}
	})

	t.Run("Translate object into SVF", func(t *testing.T) {
		var err error
		translationResult, err = mdAPI.TranslateToSVF(testObject.ObjectID)

		if err != nil {
			t.Error("Could not translate the test object, got: ", err.Error())
		}

		if translationResult.Result != "created" {
			t.Error("The test object was uploaded, but failed to create the translation job")
		}
	})

	t.Run("Get manifest of the object", func(t *testing.T) {
		manifest, err := mdAPI.GetManifest(translationResult.URN)
		if err != nil {
			t.Errorf("Problems getting the manifest for %s: %s", translationResult.URN, err.Error())
		}

		if manifest.Type != "manifest" {
			t.Error("Expecting 'manifest' type, got ", manifest.Type)
		}

		if manifest.URN != translationResult.URN {
			t.Errorf("URN not matching: translation=%s\tmanifest=%s", translationResult.URN, manifest.URN)
		}

		status := manifest.Status
		if status != "failed" && status != "success" && status != "inprogress" && status != "pending" {
			t.Errorf("Got unexpected status: %s", status)
		}


		if status == "success" && len(manifest.Derivatives) != 2 {
			t.Errorf("Expecting to have 2 derivative, got %d", len(manifest.Derivatives))
		}

		outputType := manifest.Derivatives[0].OutputType
		if status == "success" && outputType != "svf" {
			t.Errorf("Expecting first derivative to be 'svf', got %s", outputType)
		}

	})

	t.Run("Delete the temporary bucket", func(t *testing.T) {
		err := bucketAPI.DeleteBucket(tempBucketName)

		if err != nil {
			t.Fatalf("Failed to delete bucket: %s\n", err.Error())
		}
	})



}