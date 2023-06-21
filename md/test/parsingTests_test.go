package md_test

/*
package md_test provides "blackbox" tests for the md package.
These tests are meant to test the public API of the md package.
*/

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/woweh/forge-api-go-client"
	"github.com/woweh/forge-api-go-client/md"
)

func TestAPI_DefaultTranslationParams_JSON_Creation(t *testing.T) {

	mdApi := md.NewMdApi(nil, forge.US)
	params := mdApi.DefaultTranslationParams("just a test urn")

	output, err := json.Marshal(&params)
	if err != nil {
		t.Fatal("Could not marshal the preset into JSON: ", err.Error())
	}

	referenceExample := `
{
        "input": {
          "urn": "just a test urn"
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

	if !bytes.Equal(expected, output) {
		t.Fatalf(
			"The translation params are not correct:\nexpected: %s\n created: %s",
			string(expected),
			string(output),
		)

	}
}

func TestParseManifest(t *testing.T) {
	t.Run(
		"Parse pending manifest", func(t *testing.T) {
			manifest, err := os.ReadFile("../assets/pending_manifest.json")
			if err != nil {
				t.Fatal(err.Error())
			}

			var decodedManifest md.Manifest
			err = json.Unmarshal(manifest, &decodedManifest)
			if err != nil {
				t.Error(err.Error())
			}

			if len(decodedManifest.Derivatives) != 0 {
				t.Error("There should not be derivatives")
			}
		},
	)

	t.Run(
		"Parse in progress manifest", func(t *testing.T) {
			manifest, err := os.ReadFile("../assets/in_progress_manifest.json")
			if err != nil {
				t.Fatal(err.Error())
			}

			var decodedManifest md.Manifest
			err = json.Unmarshal(manifest, &decodedManifest)
			if err != nil {
				t.Error(err.Error())
			}

			if len(decodedManifest.Derivatives) != 1 {
				t.Error("Failed to parse derivatives")
			}

			if len(decodedManifest.Derivatives[0].Children) != 1 {
				t.Error("Failed to parse children derivatives")
			}

			if len(decodedManifest.Derivatives[0].Children[0].Children) != 4 {
				t.Error("Failed to parse children of derivative's children [funny]")
			}

			if decodedManifest.Derivatives[0].Children[0].Children[0].URN != "" {
				child := decodedManifest.Derivatives[0].Children[0].Children[0]
				t.Errorf("URN should be empty: %s => %s", child.Name, child.URN)
			}
		},
	)

	t.Run(
		"Parse complete failed manifest", func(t *testing.T) {
			manifest, err := os.ReadFile("../assets/failed_manifest.json")
			if err != nil {
				t.Fatal(err.Error())
			}

			var decodedManifest md.Manifest
			err = json.Unmarshal(manifest, &decodedManifest)
			if err != nil {
				t.Error(err.Error())
			}

			if len(decodedManifest.Derivatives) != 1 {
				t.Error("Failed to parse derivatives")
			}

			if len(decodedManifest.Derivatives[0].Children) != 1 {
				t.Error("Failed to parse children derivatives")
			}

			if len(decodedManifest.Derivatives[0].Children[0].Children) != 1 {
				t.Error("Failed to parse children of derivative's children [funny]")
			}

			if decodedManifest.Derivatives[0].Children[0].Children[0].URN == "" {
				t.Error("URN should not be empty")
			}

			if decodedManifest.Derivatives[0].Messages[0].Type != "warning" {
				t.Error("Should contain a warning message")
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

			// use type assertion to check if the message is an array and assign it to a variable
			if messages, ok := decodedManifest.Derivatives[0].Children[0].Messages[2].Message.([]interface{}); !ok {
				t.Error("Derivative child message should be an array")
			} else {
				// check if the message is an array of strings
				if _, okay := messages[0].(string); !okay {
					t.Error("Derivative child message should be an array of strings")
				}

				if len(messages) != 2 {
					t.Error("Derivative child message should contain 2 message descriptions")
				}
			}

			if decodedManifest.Derivatives[0].Children[0].Children[0].Role != "graphics" {
				t.Error("Failed to parse children of derivative's children [funny]")
			}
		},
	)

	t.Run(
		"Parse complete success manifest", func(t *testing.T) {
			manifest, err := os.ReadFile("../assets/success_manifest.json")
			if err != nil {
				t.Fatal(err.Error())
			}

			var decodedManifest md.Manifest
			err = json.Unmarshal(manifest, &decodedManifest)
			if err != nil {
				t.Error(err.Error())
			}

			if len(decodedManifest.Derivatives) != 4 {
				t.Error("Failed to parse derivatives")
			}

			if len(decodedManifest.Derivatives[0].Children) == 1 {
				t.Errorf(
					"Failed to parse childern derivatives, expecting 1, got %d",
					len(decodedManifest.Derivatives[0].Children),
				)
			}

			if len(decodedManifest.Derivatives[0].Children[0].Children) != 4 {
				t.Errorf(
					"Failed to parse childern of derivative's children [funny], expecting 4, got %d",
					len(decodedManifest.Derivatives[0].Children[0].Children),
				)
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
					t.Errorf(
						"Wrong derivative type parsing: expectd %s, got %s",
						decodedManifest.Derivatives[idx].OutputType,
						expectedOutputTypes[idx],
					)
				}
			}
		},
	)

	t.Run(
		"Parse Revit manifest", func(t *testing.T) {
			manifestJson, err := os.ReadFile("../assets/revit_manifest.json")
			if err != nil {
				t.Fatal(err.Error())
			}

			result := md.Manifest{}

			buffer := bytes.NewBuffer(manifestJson)
			decoder := json.NewDecoder(buffer)
			err = decoder.Decode(&result)
			if err != nil {
				t.Fatal(err.Error())
			}

			revitFileName := result.GetSourceFileName()
			if revitFileName != "20170724_Airport Model.rvt" {
				t.Error("Wrong source file name")
			}

			sp := result.GetProgressReport()
			if sp.Status != md.StatusSuccess {
				t.Error("Wrong status")
			}
			if sp.Progress != "complete" {
				t.Error("Wrong progress")
			}

			propDbUrn := result.GetPropertiesDatabaseUrn()
			if propDbUrn != "urn:adsk.viewing:fs.file:dXJuOmFkc2sub2JqZWN0czpvcy5vYmplY3Q6dGVzdC1maWxlcy8yMDE3MDcyNF9BaXJwb3J0JTIwTW9kZWwucnZ0/output/Resource/model.sdb" {
				t.Error("Wrong properties database urn")
			}

			svfSp := result.GetProgressReportOfChild("svf", "6fac95cb-af5d-3e4f-b943-8a7f55847ff1")
			if svfSp.Status != md.StatusSuccess {
				t.Error("Wrong status")
			}
			if svfSp.Progress != "" {
				t.Error("Wrong progress")
			}

			tnSP := result.GetProgressReportOfChild("thumbnail", "db899ab5-939f-e250-d79d-2d1637ce4565")
			if tnSP.Status != md.StatusSuccess {
				t.Error("Wrong status")
			}
			if tnSP.Progress != "" {
				t.Error("Wrong progress")
			}
		},
	)
}
