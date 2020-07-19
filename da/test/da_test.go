package da_test

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"github.com/apprentice3d/forge-api-go-client/da"
	"github.com/apprentice3d/forge-api-go-client/oauth"
	"os"
	"reflect"
	"testing"
)

func TestAPI_UserId(t *testing.T) {
	// prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

	authenticator := oauth.NewTwoLegged(clientID, clientSecret)
	daApi := da.NewAPI(authenticator)

	t.Run("Get the user ID", func(t *testing.T) {
		id, err := daApi.UserId()
		if err != nil {
			t.Fatal(err.Error())
		}
		if len(id) == 0 {
			t.Fatal("ID/nickname is empty")
		}
	})

}

func TestAPI_EngineList(t *testing.T) {

	// prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

	authenticator := oauth.NewTwoLegged(clientID, clientSecret)
	daApi := da.NewAPI(authenticator)

	t.Run("List the available engines", func(t *testing.T) {
		list, err := daApi.EngineList()
		if err != nil {
			t.Fatal(err.Error())
		}
		if len(list.Data) == 0 {
			t.Fatal("No data on available engines")
		}
	})

	t.Run("Get details on an engine", func(t *testing.T) {
		engineId := "Autodesk.3dsMax+2019"
		details, err := daApi.EngineDetails(engineId)
		if err != nil {
			t.Fatal(err.Error())
		}
		if len(details.Description) == 0 {
			t.Fatal("missing engine description")
		}
		if len(details.Id) == 0 {
			t.Fatal("missing engine id")
		}
		if len(details.ProductVersion) == 0 {
			t.Fatal("missing engine product version")
		}
	})

	t.Run("Get details on non-existing engine", func(t *testing.T) {
		engineId := "Autodesk.3dsMax+1995"
		_, err := daApi.EngineDetails(engineId)
		if err == nil {
			t.Fatal(err.Error())
		}

	})

}

func TestAPI_AppBundle(t *testing.T) {
	// prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")


	authenticator := oauth.NewTwoLegged(clientID, clientSecret)
	daApi := da.NewAPI(authenticator)

	testAppName := "GolangSDKTest"
	testEngine := "Autodesk.3dsMax+2019"
	//var testAlias da.Alias
	var app da.AppBundle

	nickname, err := daApi.UserId()
	if err != nil {
		t.Fatal("Could not get the user ID")
	}

	t.Run("Create an app", func(t *testing.T) {
		app, err = daApi.CreateApp(testAppName, testEngine)
		if err != nil {
			t.Fatal(err.Error())
		}

		if app.ID != nickname+"."+testAppName {
			t.Fatalf("The id of created app mismatch: expect '%s', got '%s'",
				nickname+"."+testAppName,
				app.ID)
		}
	})

	t.Run("Check if app available", func(t *testing.T) {
		list, err := daApi.AppList()
		if err != nil {
			t.Fatal(err.Error())
		}

		found := false
		for _, v := range list.Data {
			if v == nickname+"."+testAppName+"+$LATEST" {
				found = true
			}
		}

		if !found {
			t.Fatalf("The previously created app '%s' was not found", testAppName)
		}
	})

	t.Run("Create alias for app", func(t *testing.T) {
		_, err := app.CreateAlias("test", 1)
		if err != nil {
			t.Fatalf("Could not create alias: %s", err.Error())
		}
	})

	t.Run("Get app details", func(t *testing.T) {
		details, err := app.Details("test")
		if err != nil {
			t.Fatal(err.Error())
		}

		if app.Version != details.Version ||
			app.Engine != details.Engine {
			t.Fatalf("Mismatching between app data and details data: %+v", details)
		}

	})

	t.Run("Create an app with same id", func(t *testing.T) {
		_, err := daApi.CreateApp(testAppName, testEngine)
		if err == nil {
			t.Fatal("Creating  app with same id should fail, but it doesn't")
		}
	})

	t.Run("Delete the previously created app", func(t *testing.T) {
		appID := app.ID
		err := app.Delete()
		if err != nil {
			t.Fatal(err.Error())
		}

		list, err := daApi.AppList()

		found := false
		for _, v := range list.Data {
			if v == appID {
				found = true
			}
		}

		if found {
			t.Fatalf("The deleted app '%s' should not be in the list", appID)
		}

		if app.ID != "" {
			t.Fatal("The app was not properly cleared")
		}

	})
}

func TestAPI_AppList(t *testing.T) {

	// prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

	authenticator := oauth.NewTwoLegged(clientID, clientSecret)
	daApi := da.NewAPI(authenticator)

	t.Run("List the available apps", func(t *testing.T) {
		list, err := daApi.AppList()
		if err != nil {
			t.Fatal(err.Error())
		}

		if len(list.Data) == 0 {
			t.Fatal("No data on available apps")
		}

		exampleApp := "3dsMax.UVUnwrap+latest"
		found := false
		for _, v := range list.Data {
			if v == exampleApp {
				found = true
			}
		}

		if !found {
			t.Fatalf("The example app '%s' was not found", exampleApp)
		}
	})

}

func TestAppBundle_Aliases(t *testing.T) {

	// prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

	authenticator := oauth.NewTwoLegged(clientID, clientSecret)
	daApi := da.NewAPI(authenticator)

	testAppName := "GolangSDKTest"
	testEngine := "Autodesk.3dsMax+2019"
	testAliasId := "tester"
	var app da.AppBundle
	defer app.Delete()
	var err error

	t.Run("Create an app", func(t *testing.T) {
		app, err = daApi.CreateApp(testAppName, testEngine)
		if err != nil {
			t.Fatal(err.Error())
		}
	})

	t.Run("Create an alias", func(t *testing.T) {
		alias, err := app.CreateAlias(testAliasId, 1)
		if err != nil {
			t.Fatal(err.Error())
		}

		if alias.ID != testAliasId {
			t.Fatal("Mismatching aliases")
		}

	})

	t.Run("Get alias details", func(t *testing.T) {
		details, err := app.AliasDetail(testAliasId)
		if err != nil {
			t.Fatal(err.Error())
		}

		if details.ID != testAliasId {
			t.Fatalf("Alias details mismatch: %+v", details)
		}

	})

	t.Run("Create another alias", func(t *testing.T) {
		_, err := app.CreateAlias(testAliasId+"_again", 1)
		if err != nil {
			t.Errorf("Could not create another alias: %s", err.Error())
		}
	})

	t.Run("List the aliases for the app", func(t *testing.T) {
		aliases, err := app.Aliases()
		if err != nil {
			t.Fatal(err.Error())
		}

		if len(aliases.Data) == 0 {
			t.Fatal("No data on aliases")
		}

		if len(aliases.Data) != 4 {
			t.Fatalf("Expecting 3 aliases, but got %d: %+v", len(aliases.Data), aliases)
		}
	})

	t.Run("Create a new version", func(t *testing.T) {
		_, err := app.CreateVersion("Autodesk.3dsMax+2018")
		if err != nil {
			t.Fatalf("Could not create new version of app: %s", err.Error())
		}
	})

	t.Run("Modify an alias by switching to v2", func(t *testing.T) {
		alias, err := app.ModifyAlias(testAliasId, 2)
		if err != nil {
			t.Fatal(err.Error())
		}

		if alias.ID != testAliasId ||
			alias.Version != 2 {
			t.Fatalf("Alias was not properly modified: %+v", alias)
		}

	})

	t.Run("Get alias details after modification", func(t *testing.T) {
		details, err := app.AliasDetail(testAliasId)
		if err != nil {
			t.Fatal(err.Error())
		}

		if details.ID != testAliasId ||
			details.Version != 2 {
			t.Fatalf("Alias was not properly modified: %+v", details)
		}

	})

	t.Run("Delete an alias", func(t *testing.T) {
		err = app.DeleteAlias(testAliasId)
		if err != nil {
			t.Fatalf("Could not delete the alias: %s", err.Error())
		}

		aliases, err := app.Aliases()
		if err != nil {
			t.Fatal(err.Error())
		}

		if len(aliases.Data) != 3 {
			t.Fatalf("Expecting 2 aliases, but got %d: %+v", len(aliases.Data), aliases)
		}

	})

	// WARNING: removed this step and switching to defer app.Delete()
	//t.Run("Delete the app", func(t *testing.T) {
	//	err := app.Delete()
	//	if err != nil {
	//		t.Fatalf("Could not delete the app after test:%s",err.Error())
	//	}
	//})

}

func TestAppBundle_Versions(t *testing.T) {

	// prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

	authenticator := oauth.NewTwoLegged(clientID, clientSecret)
	daApi := da.NewAPI(authenticator)

	testAppName := "GolangSDKTest"
	testEngine := "Autodesk.3dsMax+2019"
	newVersionEngine := "Autodesk.3dsMax+2018"
	testAliasId := "tester"
	var app da.AppBundle
	var app2 da.AppBundle
	defer app.Delete()
	var err error

	t.Run("Create an app", func(t *testing.T) {
		app, err = daApi.CreateApp(testAppName, testEngine)
		if err != nil {
			t.Fatal(err.Error())
		}
	})

	t.Run("Create an alias", func(t *testing.T) {
		alias, err := app.CreateAlias(testAliasId, 1)
		if err != nil {
			t.Fatal(err.Error())
		}

		if alias.ID != testAliasId {
			t.Fatal("Mismatching aliases")
		}

	})

	t.Run("List app versions", func(t *testing.T) {
		versions, err := app.Versions()
		if err != nil {
			t.Fatalf("Could not get list of versions: %s", err.Error())
		}

		if len(versions.Data) != 1 {
			t.Fatalf("Expecting 1 versions, but got %d: %+v", len(versions.Data), versions)
		}
	})

	t.Run("Create a new version", func(t *testing.T) {
		app2, err = app.CreateVersion(newVersionEngine)
		if err != nil {
			t.Fatalf("Could not create new version of app: %s", err.Error())
		}

		if app2.Engine != newVersionEngine {
			t.Fatal("Newly created app doesn't have the new engine")
		}

	})

	t.Run("List app versions after addition", func(t *testing.T) {
		versions, err := app.Versions()
		if err != nil {
			t.Fatalf("Could not get list of versions: %s", err.Error())
		}

		if len(versions.Data) != 2 {
			t.Fatalf("Expecting 2 versions, but got %d: %+v", len(versions.Data), versions)
		}

	})

	t.Run("Get version details after modification", func(t *testing.T) {
		details, err := app.VersionDetails(2)
		if err != nil {
			t.Fatal(err.Error())
		}

		if details.ID != testAppName ||
			details.Version != 2 ||
			details.Engine != newVersionEngine {
			t.Fatalf("Version details are not as expected: %+v", details)
		}

	})

	t.Run("Delete a version", func(t *testing.T) {
		err = app.DeleteVersion(2)
		if err != nil {
			t.Fatalf("Could not delete the version: %s", err.Error())
		}

		versions, err := app.Versions()
		if err != nil {
			t.Fatal(err.Error())
		}

		if len(versions.Data) != 1 {
			t.Fatalf("Expecting 1 version, but got %d: %+v", len(versions.Data), versions)
		}

	})

}

func TestAppBundle_Upload(t *testing.T) {
	// prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

	authenticator := oauth.NewTwoLegged(clientID, clientSecret)
	daApi := da.NewAPI(authenticator)

	testAppName := "GolangSDKTest"
	testEngine := "Autodesk.3dsMax+2019"
	var app da.AppBundle
	defer app.Delete()
	var err error

	t.Run("Create an app", func(t *testing.T) {
		app, err = daApi.CreateApp(testAppName, testEngine)
		if err != nil {
			t.Fatal(err.Error())
		}
	})

	t.Run("Upload a test bundle", func(t *testing.T) {

		//data, err := ioutil.ReadFile("ListMyObjectsApplicationPackage.zip")
		//if err != nil {
		//	t.Fatal("Could not read file to upload it")
		//}

		data := []byte("some test load")
		err = app.Upload(data)

		if err != nil {
			t.Fatal(err.Error())
		}
	})
}



func TestAPI_Activity(t *testing.T) {
	// prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

	authenticator := oauth.NewTwoLegged(clientID, clientSecret)
	daApi := da.NewAPI(authenticator)

	testActivityName := "GolangSDKTest"
	testEngine := "Autodesk.3dsMax+2019"
	//testAlias := "golangTest"

	var testActivity da.Activity

	nickname, err := daApi.UserId()
	if err != nil {
		t.Fatal("Could not get the user ID")
	}

	t.Run("Create an activity", func(t *testing.T) {

		config := da.ActivityConfig{
			ID: testActivityName,
			Engine:testEngine,
			CommandLine:[]string{"dir"},

		}

		testActivity, err = daApi.CreateActivity(config)
		if err != nil {
			t.Fatal(err.Error())
		}

		if testActivity.ID != nickname+"."+testActivityName {
			t.Fatalf("The id of created activity mismatch: expect '%s', got '%s'",
				testActivityName,
				testActivity.ID)
		}
	})

	//t.Run("Check if activity is available", func(t *testing.T) {
	//	list, err := daApi.AppList()
	//	if err != nil {
	//		t.Fatal(err.Error())
	//	}
	//
	//	found := false
	//	for _, v := range list.Data {
	//		if v == nickname+"."+testAppName+"+$LATEST" {
	//			found = true
	//		}
	//	}
	//
	//	if !found {
	//		t.Fatalf("The previously created app '%s' was not found", testAppName)
	//	}
	//})

	//t.Run("Create alias for app", func(t *testing.T) {
	//	testAlias, err = app.CreateAlias("test", 1)
	//	if err != nil {
	//		t.Fatalf("Could not create alias: %s", err.Error())
	//	}
	//})
	//
	//t.Run("Get app details", func(t *testing.T) {
	//	details, err := app.Details("test")
	//	if err != nil {
	//		t.Fatal(err.Error())
	//	}
	//
	//	if app.Version != details.Version ||
	//		app.Engine != details.Engine {
	//		t.Fatalf("Mismatching between app data and details data: %+v", details)
	//	}
	//
	//})
	//
	//t.Run("Create an app with same id", func(t *testing.T) {
	//	_, err := daApi.CreateApp(testAppName, testEngine)
	//	if err == nil {
	//		t.Fatal("Creating  app with same id should fail, but it doesn't")
	//	}
	//})
	//
	t.Run("Delete the previously created app", func(t *testing.T) {
		//activityID := testActivity.ID
		//err := testActivity.Delete()
		//if err != nil {
		//	t.Fatal(err.Error())
		//}
		//
		//list, err := daApi.AppList()
		//
		//found := false
		//for _, v := range list.Data {
		//	if v == appID {
		//		found = true
		//	}
		//}
		//
		//if found {
		//	t.Fatalf("The deleted app '%s' should not be in the list", appID)
		//}
		//
		//if app.ID != "" {
		//	t.Fatal("The app was not properly cleared")
		//}

	})
}














func TestAPI_DataParsing(t *testing.T) {
	t.Run("Check EngineList struct parsing", func(t *testing.T) {
		jsonString := `
{
    "paginationToken": null,
    "data": [
        "Autodesk.Revit+28_1",
        "Autodesk.AutoCAD+22",
        "Autodesk.Inventor+23",
        "Autodesk.Revit+29_3",
        "Autodesk.AutoCAD+23",
        "Autodesk.3dsMax+2018",
        "Autodesk.Revit+2018",
        "Autodesk.Test+Latest",
        "Autodesk.Inventor+22",
        "Autodesk.3dsMax+2019",
        "Autodesk.AutoCAD+21",
        "Autodesk.Revit+29_5",
        "Autodesk.Revit+2019",
        "Autodesk.AutoCAD+20_1"
    ]
}`
		result := da.EngineList{}

		buffer := bytes.NewBufferString(jsonString)
		decoder := json.NewDecoder(buffer)
		err := decoder.Decode(&result)

		if err != nil {
			t.Fatal(err.Error())
		}
		if len(result.Data) == 0 {
			t.Fatal("No data on available engines")
		}

	})

	t.Run("Check App creation result parsing", func(t *testing.T) {
		appCreationResult := `
{
    "uploadParameters": {
        "endpointURL": "https://dasprod-store.s3.amazonaws.com",
        "formData": {
            "key": "apps/U4E38tF2P9hpSdth39FfTMsUphf5gHsP/DenixTest/1",
            "content-type": "application/octet-stream",
            "policy": "eyJleHBpcmF0aW9uIjoiMjAxOC0xMS0wNlQxOToyODo1OS43MjI0MzcxWiIsImNvbmRpdGlvbnMiOlt7ImtleSI6ImFwcHMvVTRFMzh0RjJQOWhwU2R0aDM5RmZUTXNVcGhmNWdIc1AvRGVuaXhUZXN0LzEifSx7ImJ1Y2tldCI6ImRhc3Byb2Qtc3RvcmUifSx7InN1Y2Nlc3NfYWN0aW9uX3N0YXR1cyI6IjIwMCJ9LFsic3RhcnRzLXdpdGgiLCIkc3VjY2Vzc19hY3Rpb25fcmVkaXJlY3QiLCIiXSxbInN0YXJ0cy13aXRoIiwiJGNvbnRlbnQtVHlwZSIsImFwcGxpY2F0aW9uL29jdGV0LXN0cmVhbSJdLHsieC1hbXotc2VydmVyLXNpZGUtZW5jcnlwdGlvbiI6IkFFUzI1NiJ9LFsiY29udGVudC1sZW5ndGgtcmFuZ2UiLCIwIiwiMTA0ODU3NjAwIl0seyJ4LWFtei1jcmVkZW50aWFsIjoiQVNJQVRHVkpaS00zRFdaQ01WVzMvMjAxODExMDYvdXMtZWFzdC0xL3MzL2F3czRfcmVxdWVzdC8ifSx7IngtYW16LWFsZ29yaXRobSI6IkFXUzQtSE1BQy1TSEEyNTYifSx7IngtYW16LWRhdGUiOiIyMDE4MTEwNlQxODI4NTlaIn0seyJ4LWFtei1zZWN1cml0eS10b2tlbiI6IkZRb0daWEl2WVhkekVEc2FETGlhT0t0RTI4ZXZKc1d3dnlMOUFReFRGOXY5cmxyVlZraHZzeDVjSytrV2o0dHZNL2FZT2doY2NhWGhWaEQ5UlBuTDVKNnMvaHFFWGh6UGkxNkQrT3p0c1U2L0xLcTFhV2RaT281cXFTYi90WmNoUUhQSFZJWTdCY0RZa0EyeDlSRFlCZDMrWHJhZ0IzeWR4UWM0ekZ3aDZDV3psN3RMeXl1cWhMa0NVL3o1c2hFWE42bis5TExJdCtHWEFmWGl2NjlGSFZobXo5L3Q3ZHlSK1lXWTl2aUdOVWN4MHM2Si83cnVpR1k2bmtmd3lMS1NBR3NXU3hjeC9wbGZIbXRHTXZoTnY3TERweUI1NDdScGwyRFEwVEdaYndkVHRreGJGNEJ0U2N5ZUdacjF2ZWVoZ3MrN1pTeklaek50Z2hsKzhnOXlDdGg1Smdta2VGV2pIOXJjUkI0TTBLM2lqK2RlaDhqRXZxWW8ycWlIM3dVPSJ9XX0=",
            "success_action_status": "200",
            "success_action_redirect": "",
            "x-amz-signature": "0dfb670d1d275feedaef370a44d3f29c0b448993c94304cd5ebb50285015e3cb",
            "x-amz-credential": "ASIATGVJZKM3DWZCMVW3/20181106/us-east-1/s3/aws4_request/",
            "x-amz-algorithm": "AWS4-HMAC-SHA256",
            "x-amz-date": "20181106T182859Z",
            "x-amz-server-side-encryption": "AES256",
            "x-amz-security-token": "FQoGZXIvYXdzEDsaDLiaOKtE28evJsWwvyL9AQxTF9v9rlrVVkhvsx5cK+kWj4tvM/aYOghccaXhVhD9RPnL5J6s/hqEXhzPi16D+OztsU6/LKq1aWdZOo5qqSb/tZchQHPHVIY7BcDYkA2x9RDYBd3+XragB3ydxQc4zFwh6CWzl7tLyyuqhLkCU/z5shEXN6n+9LLIt+GXAfXiv69FHVhmz9/t7dyR+YWY9viGNUcx0s6J/7ruiGY6nkfwyLKSAGsWSxcx/plfHmtGMvhNv7LDpyB547Rpl2DQ0TGZbwdTtkxbF4BtScyeGZr1veehgs+7ZSzIZzNtghl+8g9yCth5JgmkeFWjH9rcRB4M0K3ij+deh8jEvqYo2qiH3wU="
        }
    },
    "engine": "Autodesk.3dsMax+2019",
    "version": 1,
    "id": "U4E38tF2P9hpSdth39FfTMsUphf5gHsP.DenixTest"
}`
		result := da.AppBundle{}

		buffer := bytes.NewBufferString(appCreationResult)
		decoder := json.NewDecoder(buffer)
		err := decoder.Decode(&result)

		if err != nil {
			t.Fatal(err.Error())
		}

		if result.ID != "U4E38tF2P9hpSdth39FfTMsUphf5gHsP.DenixTest" {
			t.Fatalf("Could not properly extract the id: "+
				"expecting 'U4E38tF2P9hpSdth39FfTMsUphf5gHsP.DenixTest', "+
				"got '%s'", result.ID)
		}

	})

	t.Run("Check App uploading error parsing", func(t *testing.T) {
		appUploadingError := `
<?xml version="1.0" encoding="UTF-8"?>
<Error>
    <Code>InvalidArgument</Code>
    <Message>POST requires exactly one file upload per request.</Message>
    <ArgumentName>file</ArgumentName>
    <ArgumentValue>0</ArgumentValue>
    <RequestId>E348E16638E3065E</RequestId>
    <HostId>Wd8GvmCwzBF5EI3QkOlurAs9pxTuixDQIajEGRhH1mC8O8vRMkSKynxpEK5mtKzVRaQyTI7Awyw=</HostId>
</Error>
`
		result := struct {
			Code          string `xml:"Code"`
			Message       string `xml:"Message"`
			Argument      string `xml:"Argument"`
			ArgumentValue string `xml:"ArgumentValue"`
			RequestID     string `xml:"RequestId"`
			HostID        string `xml:"HostId"`
		}{}

		buffer := bytes.NewBufferString(appUploadingError)
		decoder := xml.NewDecoder(buffer)
		err := decoder.Decode(&result)

		if err != nil {
			t.Fatal(err.Error())
		}

	})

	t.Run("Check Activity creation JSON formation", func(t *testing.T) {
		aliasCreationJSON := `
{
	"id": "Denis.Incercare02",
    "commandLine": [
		"$(engine.path)/3dsmaxbatch.exe -sceneFile \"$(args[InputFile].path)\" \"$(settings[script].path)\""
	],
    "description": "Export a single max file to FBX",
    "appbundles": [
    	],
    "engine" : "Autodesk.3dsMax+2019",
    "parameters": {
		"InputFile" : {
		    "zip": false,
			"description": "Input 3ds Max file",
            "ondemand": false,
			"required": true,
            "verb": "get",
            "localName": "input.max"
		},
		"OutputFile": {
		    "zip": false,
            "ondemand": false,
            "verb": "put",
            "description": "Output FBX file",
            "required": true,
            "localName": "output.fbx"
		}
    },
    "settings": {
       "script": "exportFile (sysInfo.currentdir + \"/output.fbx\") #noPrompt using:FBXEXP"
   }
}
`

		result := da.ActivityConfig{}

		buffer := bytes.NewBufferString(aliasCreationJSON)
		decoder := json.NewDecoder(buffer)
		err := decoder.Decode(&result)

		if err != nil {
			t.Fatal(err.Error())
		}

		aliasCreationResponseJSON := `
{
    "commandLine": [
        "$(engine.path)/3dsmaxbatch.exe -sceneFile \"$(args[InputFile].path)\" \"$(settings[script].path)\""
    ],
    "parameters": {
        "InputFile": {
            "verb": "get",
            "description": "Input 3ds Max file",
            "required": true,
            "localName": "input.max"
        },
        "OutputFile": {
            "verb": "put",
            "description": "Output FBX file",
            "required": true,
            "localName": "output.fbx"
        }
    },
    "engine": "Autodesk.3dsMax+2019",
    "appbundles": [],
    "settings": {
        "script": "exportFile (sysInfo.currentdir + \"/output.fbx\") #noPrompt using:FBXEXP"
    },
    "description": "Export a single max file to FBX",
    "version": 1,
    "id": "Denis.Incercare02"
}
`

		response := da.ActivityConfig{}

		buffer2 := bytes.NewBufferString(aliasCreationResponseJSON)
		decoder2 := json.NewDecoder(buffer2)
		err = decoder2.Decode(&response)

		if err != nil {
			t.Fatal(err.Error())
		}


		if !reflect.DeepEqual(result, response) {
			t.Fatal("failed to properly parse the Activity JSON")
		}






	})

}
