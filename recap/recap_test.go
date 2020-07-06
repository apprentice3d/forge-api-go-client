package recap_test

import (
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/apprentice3d/forge-api-go-client/recap"
)

func TestReCapAPIWorkflowUsingRemoteLinks(t *testing.T) {

	linkSamples := []string{
		"https://s3.amazonaws.com/adsk-recap-public/forge/lion/DSC_1158.JPG",
		"https://s3.amazonaws.com/adsk-recap-public/forge/lion/DSC_1159.JPG",
		"https://s3.amazonaws.com/adsk-recap-public/forge/lion/DSC_1160.JPG",
		"https://s3.amazonaws.com/adsk-recap-public/forge/lion/DSC_1162.JPG",
		"https://s3.amazonaws.com/adsk-recap-public/forge/lion/DSC_1163.JPG",
		"https://s3.amazonaws.com/adsk-recap-public/forge/lion/DSC_1164.JPG",
		"https://s3.amazonaws.com/adsk-recap-public/forge/lion/DSC_1165.JPG",
	}

	var scene recap.PhotoScene

	testingFormat := "obj"

	// prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		t.Skipf("No Forge credentials present; skipping test")
	}

	recapAPI := recap.NewAPIWithCredentials(clientID, clientSecret)

	t.Run("Creating a new photoScene", func(t *testing.T) {
		var err error
		scene, err = recapAPI.CreatePhotoScene("example", []string{testingFormat}, "object")
		if err != nil {
			t.Fatal(err.Error())
		}
	})

	t.Run("Uploading sample images using links", func(t *testing.T) {
		for _, link := range linkSamples {
			_, err := recapAPI.AddFileToSceneUsingLink(scene.ID, link)
			if err != nil {
				t.Fatal(err.Error())
			}
		}
	})

	t.Run("Starting photoScene processing", func(t *testing.T) {
		if _, err := recapAPI.StartSceneProcessing(scene.ID); err != nil {
			t.Error(err.Error())
		}
	})

	t.Run("Checking scene status each 30 sec", func(t *testing.T) {
		var progressResult recap.SceneProgressReply
		var err error
		for {
			if progressResult, err = recapAPI.GetSceneProgress(scene.ID); err != nil {
				t.Errorf("Failed to get the PhotoScene progress: %s\n", err.Error())
			}

			ratio, err := strconv.ParseFloat(progressResult.PhotoScene.Progress, 64)

			if err != nil {
				t.Fatalf("Failed to parse progress results: %s\n", err.Error())
			}

			if ratio == float64(100.0) {
				break
			}
			time.Sleep(5 * time.Second)
		}
	})

	t.Run("Get the available result", func(t *testing.T) {
		result, err := recapAPI.GetSceneResults(scene.ID, testingFormat)
		if err != nil {
			t.Error(err.Error())
		}
		if len(result.PhotoScene.SceneLink) == 0 {
			t.Error("The received link is empty")
		}
	})

	t.Run("Check the result file size for normal size", func(t *testing.T) {
		response, err := recapAPI.GetSceneResults(scene.ID, testingFormat)
		if err != nil {
			t.Error(err.Error())
		}
		if len(response.PhotoScene.SceneLink) == 0 {
			t.Error("The received link is empty")
		}

		filename := "temp.zip"

		resp, err := http.Get(response.PhotoScene.SceneLink)

		if err != nil {
			return
		}
		defer resp.Body.Close()
		result, err := os.Create(filename)
		if err != nil {
			return
		}
		defer result.Close()

		tempFile, err := os.Stat(filename)

		if tempFile.Size() <= 22 {
			t.Error("The scene was processed, but the result file is abnormally small: ", tempFile.Size())
		}

		err = os.Remove(filename)
		if err != nil {
			t.Error(err.Error())
		}

		return

	})

	t.Run("Delete the scene", func(t *testing.T) {
		_, err := recapAPI.DeleteScene(scene.ID)
		if err != nil {
			t.Fatal(err.Error())
		}
	})

}

func TestReCapAPIWorkflowUsingLocalFiles(t *testing.T) {

	// these files are remotely located, to make them available on remote test servers,
	// so we will have to download them locally, to test the file uploading part
	linkSamples := []string{
		"https://s3.amazonaws.com/adsk-recap-public/forge/lion/DSC_1158.JPG",
		"https://s3.amazonaws.com/adsk-recap-public/forge/lion/DSC_1159.JPG",
		"https://s3.amazonaws.com/adsk-recap-public/forge/lion/DSC_1160.JPG",
		"https://s3.amazonaws.com/adsk-recap-public/forge/lion/DSC_1162.JPG",
		"https://s3.amazonaws.com/adsk-recap-public/forge/lion/DSC_1163.JPG",
		"https://s3.amazonaws.com/adsk-recap-public/forge/lion/DSC_1164.JPG",
		"https://s3.amazonaws.com/adsk-recap-public/forge/lion/DSC_1165.JPG",
	}

	var scene recap.PhotoScene

	testingFormat := "obj"

	// prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		t.Skipf("No Forge credentials present; skipping test")
	}

	recapAPI := recap.NewAPIWithCredentials(clientID, clientSecret)

	t.Run("Creating a new photoScene", func(t *testing.T) {
		var err error
		scene, err = recapAPI.CreatePhotoScene("example", []string{testingFormat}, "object")
		if err != nil {
			t.Fatal(err.Error())
		}
		t.Logf("Created a photoscene with ID=%s", scene.ID)
	})

	t.Run("Uploading sample images using data", func(t *testing.T) {

		//download each link locally and then upload the data
		for _, link := range linkSamples {
			response, err := http.Get(link)
			if err != nil {
				t.Fatal(err.Error())
			}

			data, err := ioutil.ReadAll(response.Body)
			response.Body.Close()
			if err != nil {
				t.Fatal(err.Error())
			}

			_, err = recapAPI.AddFileToSceneUsingData(scene.ID, data)
			if err != nil {
				t.Fatal(err.Error())
			}

		}
	})

	t.Run("Starting photoScene processing", func(t *testing.T) {
		if _, err := recapAPI.StartSceneProcessing(scene.ID); err != nil {
			t.Error(err.Error())
		}
	})

	t.Run("Checking scene status each 30 sec", func(t *testing.T) {
		var progressResult recap.SceneProgressReply
		var err error
		for {
			if progressResult, err = recapAPI.GetSceneProgress(scene.ID); err != nil {
				t.Errorf("Failed to get the PhotoScene progress: %s\n", err.Error())
			}

			ratio, err := strconv.ParseFloat(progressResult.PhotoScene.Progress, 64)

			if err != nil {
				t.Fatalf("Failed to parse progress results: %s\n", err.Error())
			}

			if ratio == float64(100.0) {
				break
			}
			time.Sleep(5 * time.Second)
		}
	})

	t.Run("Get the available result", func(t *testing.T) {
		result, err := recapAPI.GetSceneResults(scene.ID, testingFormat)
		if err != nil {
			t.Error(err.Error())
		}
		if len(result.PhotoScene.SceneLink) == 0 {
			t.Error("The received link is empty")
		}
	})

	t.Run("Check the result file size for normal size", func(t *testing.T) {
		response, err := recapAPI.GetSceneResults(scene.ID, testingFormat)
		if err != nil {
			t.Error(err.Error())
		}
		if len(response.PhotoScene.SceneLink) == 0 {
			t.Error("The received link is empty")
		}

		filename := "temp.zip"

		resp, err := http.Get(response.PhotoScene.SceneLink)

		if err != nil {
			return
		}
		defer resp.Body.Close()
		result, err := os.Create(filename)
		if err != nil {
			return
		}
		defer result.Close()

		tempFile, err := os.Stat(filename)

		if tempFile.Size() <= 22 {
			t.Error("The scene was processed, but the result file is abnormally small: ", tempFile.Size())
		}

		err = os.Remove(filename)
		if err != nil {
			t.Error(err.Error())
		}

		return

	})

	t.Run("Delete the scene", func(t *testing.T) {
		_, err := recapAPI.DeleteScene(scene.ID)
		if err != nil {
			t.Fatal(err.Error())
		}
	})

}

func TestCreatePhotoScene(t *testing.T) {

	// prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		t.Skipf("No Forge credentials present; skipping test")
	}
	recapAPI := recap.NewAPIWithCredentials(clientID, clientSecret)
	var sceneID string

	t.Run("Create a scene", func(t *testing.T) {
		response, err := recapAPI.CreatePhotoScene("testare", nil, "object")

		if err != nil {
			t.Fatalf("Failed to create a photoscene: %s\n", err.Error())
		}

		sceneID = response.ID
	})

	t.Run("Delete the test scene", func(t *testing.T) {
		_, err := recapAPI.DeleteScene(sceneID)

		if err != nil {
			t.Fatalf("Failed to delete the photoscene: %s\n", err.Error())
		}
	})

	t.Run("Check fail on create a scene with empty name", func(t *testing.T) {
		_, err := recapAPI.CreatePhotoScene("", nil, "object")

		if err == nil {
			t.Fatalf("Should fail creating a scene with empty name\n")
		}
	})

}

func ExampleAPI_CreatePhotoScene() {

	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")
	recap := recap.NewAPIWithCredentials(clientID, clientSecret)

	photoScene, err := recap.CreatePhotoScene("test_scene", nil, "object")
	if err != nil {
		// handle error
	}

	if len(photoScene.ID) != 0 {
	}
}
