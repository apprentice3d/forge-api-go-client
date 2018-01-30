package recap

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestTheEntireWorkflow(t *testing.T) {

	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

	recapAPI := NewReCapAPIWithCredentials(clientID, clientSecret)

	format := "obj"

	t.Log("Creating a scene ...")
	scene, err := recapAPI.CreatePhotoScene("example", []string{format})
	if err != nil {
		t.Fatal(err.Error())
	}

	fileSamples := []string{
		"https://s3.amazonaws.com/adsk-recap-public/forge/lion/DSC_1158.JPG",
		"https://s3.amazonaws.com/adsk-recap-public/forge/lion/DSC_1159.JPG",
		"https://s3.amazonaws.com/adsk-recap-public/forge/lion/DSC_1160.JPG",
		"https://s3.amazonaws.com/adsk-recap-public/forge/lion/DSC_1162.JPG",
		"https://s3.amazonaws.com/adsk-recap-public/forge/lion/DSC_1163.JPG",
		"https://s3.amazonaws.com/adsk-recap-public/forge/lion/DSC_1164.JPG",
		"https://s3.amazonaws.com/adsk-recap-public/forge/lion/DSC_1165.JPG",
	}

	t.Log("Uploading sample images ...")
	uploadResults, err := recapAPI.AddFilesToSceneUsingLinks(&scene, fileSamples)
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Logf("Successfully uploaded: %s\n", uploadResults.Files.File[0].FileName)


	t.Log("Starting scene processing ...")
	if _, err := recapAPI.StartSceneProcessing(scene); err != nil {
		t.Error(err.Error())
	}

	t.Log("Checking scene status ...")
	var progressResult SceneProgressReply
	for {
		if progressResult, err = recapAPI.GetSceneProgress(scene); err != nil {
			t.Errorf("Failed to get the PhotoScene progress: %s\n", err.Error())
		}

		ratio, err := strconv.ParseFloat(progressResult.PhotoScene.Progress, 64)

		if err != nil {
			t.Fatalf("Failed to parse progress results: %s\n", err.Error())
		}

		if ratio == float64(100.0) {
			break
		}
		t.Logf("Scene progress = %.2f%%\n", ratio)
		time.Sleep(5 * time.Second)
	}

	t.Log("Finished processing the scene, now getting the results ...")
	result, err := recapAPI.GetSceneResults(scene, format)
	if err != nil {
		t.Error(err.Error())
	}
	t.Logf("Received the following link %s\n", result.PhotoScene.SceneLink)

	t.Log("Deleting the scene ...")

	_, err = recapAPI.DeleteScene(scene)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("Scene deleted successfully!")
}

func TestCreatePhotoSceneReCap(t *testing.T) {

	// prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")
	recapAPI := NewReCapAPIWithCredentials(clientID, clientSecret)

	t.Run("Create a scene", func(t *testing.T) {
		_, err := recapAPI.CreatePhotoScene("testare", nil)

		if err != nil {
			t.Fatalf("Failed to create a photoscene: %s\n", err.Error())
		}
	})

	t.Run("Check fail on create a scene with empty name", func(t *testing.T) {
		_, err := recapAPI.CreatePhotoScene("", nil)

		if err == nil {
			t.Fatalf("Should fail creating a scene with empty name\n")
		}
	})

}

func ExampleReCapAPI_CreatePhotoScene() {
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")
	recapAPI := NewReCapAPIWithCredentials(clientID, clientSecret)

	photoscene, err := recapAPI.CreatePhotoScene("test_scene", nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	if len(photoscene.ID) != 0 {
		fmt.Println("Scene was successfully created")
	}

	//Output:
	//Scene was successfully created
}
