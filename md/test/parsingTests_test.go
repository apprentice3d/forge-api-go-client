package md_test

/*
package md_test provides "blackbox" tests for the md package.
These tests are meant to test the public API of the md package.

TODO: add tests to check Region, ProgressReport and Status
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
			manifestJson, err := os.ReadFile("../assets/pending_manifest.json")
			if err != nil {
				t.Fatal(err.Error())
			}

			var manifest md.Manifest
			err = json.Unmarshal(manifestJson, &manifest)
			if err != nil {
				t.Error(err.Error())
			}

			if !manifest.Status.IsPending() {
				t.Error("Status should be pending")
			}

			if manifest.Status.IsTimeout() {
				t.Error("Status should not be timeout")
			}

			if len(manifest.Derivatives) != 0 {
				t.Error("There should not be derivatives")
			}

			pr := manifest.GetProgressReportOfChild("abc", "def")
			if !pr.IsEmpty() {
				t.Error("Progress report should be empty")
			}
		},
	)

	t.Run(
		"Parse in progress manifest", func(t *testing.T) {
			manifestJson, err := os.ReadFile("../assets/in_progress_manifest.json")
			if err != nil {
				t.Fatal(err.Error())
			}

			var manifest md.Manifest
			err = json.Unmarshal(manifestJson, &manifest)
			if err != nil {
				t.Error(err.Error())
			}

			if !manifest.Status.IsInProgress() {
				t.Error("Status should be in progress")
			}

			if len(manifest.Derivatives) != 1 {
				t.Error("Failed to parse derivatives")
			}

			if len(manifest.Derivatives[0].Children) != 1 {
				t.Error("Failed to parse children derivatives")
			}

			if len(manifest.Derivatives[0].Children[0].Children) != 4 {
				t.Error("Failed to parse children of derivative's children [funny]")
			}

			if manifest.Derivatives[0].Children[0].Children[0].URN != "" {
				child := manifest.Derivatives[0].Children[0].Children[0]
				t.Errorf("URN should be empty: %s => %s", child.Name, child.URN)
			}
		},
	)

	t.Run(
		"Parse complete failed manifest", func(t *testing.T) {
			manifestJson, err := os.ReadFile("../assets/failed_manifest.json")
			if err != nil {
				t.Fatal(err.Error())
			}

			var manifest md.Manifest
			err = json.Unmarshal(manifestJson, &manifest)
			if err != nil {
				t.Error(err.Error())
			}

			if !manifest.Status.IsFailed() {
				t.Error("Status should be failed")
			}

			if !manifest.Region.IsUS() {
				t.Error("Region should be US")
			}

			if manifest.Region.IsEMEA() {
				t.Error("Region should not be EMEA")
			}

			if len(manifest.Derivatives) != 1 {
				t.Error("Failed to parse derivatives")
			}

			if len(manifest.Derivatives[0].Children) != 1 {
				t.Error("Failed to parse children derivatives")
			}

			if len(manifest.Derivatives[0].Children[0].Children) != 1 {
				t.Error("Failed to parse children of derivative's children [funny]")
			}

			if manifest.Derivatives[0].Children[0].Children[0].URN == "" {
				t.Error("URN should not be empty")
			}

			if manifest.Derivatives[0].Messages[0].Type != "warning" {
				t.Error("Should contain a warning message")
			}

			if len(manifest.Derivatives[0].Children[0].Messages) != 3 {
				t.Error("Derivative child should contain 3 error message")
			}

			if manifest.Derivatives[0].Children[0].Messages[0].Type != "warning" {
				t.Error("Derivative child message should be a warning message")
			}
			if manifest.Derivatives[0].Children[0].Messages[2].Type != "error" {
				t.Error("Derivative child message should be an error message")
			}

			// use type assertion to check if the message is an array and assign it to a variable
			if messages, ok := manifest.Derivatives[0].Children[0].Messages[2].Message.([]interface{}); !ok {
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

			if manifest.Derivatives[0].Children[0].Children[0].Role != "graphics" {
				t.Error("Failed to parse children of derivative's children [funny]")
			}
		},
	)

	t.Run(
		"Parse complete success manifest", func(t *testing.T) {
			manifestJson, err := os.ReadFile("../assets/success_manifest.json")
			if err != nil {
				t.Fatal(err.Error())
			}

			var manifest md.Manifest
			err = json.Unmarshal(manifestJson, &manifest)
			if err != nil {
				t.Error(err.Error())
			}

			if !manifest.Status.IsSuccess() {
				t.Error("Status should be success")
			}

			if !manifest.Region.IsUS() {
				t.Error("Region should be US")
			}

			if manifest.Region.IsEMEA() {
				t.Error("Region should not be EMEA")
			}

			if len(manifest.Derivatives) != 4 {
				t.Error("Failed to parse derivatives")
			}

			if len(manifest.Derivatives[0].Children) == 1 {
				t.Errorf(
					"Failed to parse childern derivatives, expecting 1, got %d",
					len(manifest.Derivatives[0].Children),
				)
			}

			if len(manifest.Derivatives[0].Children[0].Children) != 4 {
				t.Errorf(
					"Failed to parse childern of derivative's children [funny], expecting 4, got %d",
					len(manifest.Derivatives[0].Children[0].Children),
				)
			}

			if manifest.Derivatives[0].Children[0].Children[0].URN == "" {
				t.Error("URN should not be empty")
			}

			if len(manifest.Derivatives[0].Messages) != 0 {
				t.Error("Derivative should not contain any error messages")
			}

			expectedOutputTypes := []string{"svf", "step", "thumbnail", "obj"}

			for idx := range manifest.Derivatives {
				if manifest.Derivatives[idx].OutputType != expectedOutputTypes[idx] {
					t.Errorf(
						"Wrong derivative type parsing: expectd %s, got %s",
						manifest.Derivatives[idx].OutputType,
						expectedOutputTypes[idx],
					)
				}
			}

			objPr := manifest.GetProgressReportOfChild("obj", "4f981e94-8241-4eaf-b08b-cd337c6b8b1f")
			if objPr.Status != md.StatusSuccess {
				t.Error("Wrong status")
			}
			if objPr.Progress != "" {
				t.Error("Wrong progress")
			}
		},
	)

	t.Run(
		"Parse Revit manifest", func(t *testing.T) {
			manifestJson, err := os.ReadFile("../assets/revit_manifest.json")
			if err != nil {
				t.Fatal(err.Error())
			}

			manifest := md.Manifest{}

			buffer := bytes.NewBuffer(manifestJson)
			decoder := json.NewDecoder(buffer)
			err = decoder.Decode(&manifest)
			if err != nil {
				t.Fatal(err.Error())
			}

			revitFileName := manifest.GetSourceFileName()
			if revitFileName != "20170724_Airport Model.rvt" {
				t.Error("Wrong source file name")
			}

			sp := manifest.GetProgressReport()
			if sp.Status != md.StatusSuccess {
				t.Error("Wrong status")
			}
			if sp.Progress != "complete" {
				t.Error("Wrong progress")
			}

			propDbUrn := manifest.GetPropertiesDatabaseUrn()
			if propDbUrn != "urn:adsk.viewing:fs.file:dXJuOmFkc2sub2JqZWN0czpvcy5vYmplY3Q6dGVzdC1maWxlcy8yMDE3MDcyNF9BaXJwb3J0JTIwTW9kZWwucnZ0/output/Resource/model.sdb" {
				t.Error("Wrong properties database urn")
			}
		},
	)
}
