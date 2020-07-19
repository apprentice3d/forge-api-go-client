// Package recap contains the Go wrappers for calls to Forge Reality Capture API
// https://developer.autodesk.com/api/reality-capture-cover-page/
//
// 	The workflow is the following:
// 		- create a photoScene
//		- upload images to photoScene
//		- start photoScene processing
//		- get the result
package recap

import (
	"github.com/apprentice3d/forge-api-go-client/oauth"
)

// API struct holds all paths necessary to access ReCap API
type ReCapAPI struct {
	Authenticator oauth.ForgeAuthenticator
	ReCapPath string
}

// NewAPI returns a ReCap API client with default configurations
func NewAPI(authenticator oauth.ForgeAuthenticator) ReCapAPI {
	return ReCapAPI{
		authenticator,
		"/photo-to-3d/v1",
	}
}

// CreatePhotoScene prepares a scene with a given name, expected output formats and sceneType
// 	name - should not be empty
// 	formats - should be of type rcm, rcs, obj, ortho or report
// 	sceneType - should be either "aerial" or "object"
func (api ReCapAPI) CreatePhotoScene(name string, formats []string, sceneType string) (scene PhotoScene, err error) {

	bearer, err := api.Authenticator.GetToken("data:write")
	if err != nil {
		return
	}
	path := api.Authenticator.GetHostPath() + api.ReCapPath
	scene, err = createPhotoScene(path, name, formats, sceneType, bearer.AccessToken)

	return
}

// AddFileToSceneUsingLink can be used when the needed images are already available remotely
// and can be uploaded just by providing the remote link
func (api ReCapAPI) AddFileToSceneUsingLink(sceneID string, link string) (uploads FileUploadingReply, err error) {

	bearer, err := api.Authenticator.GetToken("data:write")
	if err != nil {
		return
	}
	path := api.Authenticator.GetHostPath() + api.ReCapPath

	uploads, err = addFileToSceneUsingLink(path, sceneID, link, bearer.AccessToken)
	return
}

// AddFileToSceneUsingData can be used when the image is already available as a byte slice,
// be it read from a local file or as a result/body of a POST request
func (api ReCapAPI) AddFileToSceneUsingData(sceneID string, data []byte) (uploads FileUploadingReply, err error) {

	bearer, err := api.Authenticator.GetToken("data:write")
	if err != nil {
		return
	}
	path := api.Authenticator.GetHostPath() + api.ReCapPath

	uploads, err = addFileToSceneUsingFileData(path, sceneID, data, bearer.AccessToken)

	return
}

// StartSceneProcessing will trigger the processing of a specified scene that can be canceled any time
func (api ReCapAPI) StartSceneProcessing(sceneID string) (result SceneStartProcessingReply, err error) {
	bearer, err := api.Authenticator.GetToken("data:write")
	if err != nil {
		return
	}
	path := api.Authenticator.GetHostPath() + api.ReCapPath
	result, err = startSceneProcessing(path, sceneID, bearer.AccessToken)
	return
}

// GetSceneProgress polls the scene processing status and progress
//	Note: instead of polling, consider using the callback parameter that can be specified upon scene creation
func (api ReCapAPI) GetSceneProgress(sceneID string) (progress SceneProgressReply, err error) {
	bearer, err := api.Authenticator.GetToken("data:read")
	if err != nil {
		return
	}
	path := api.Authenticator.GetHostPath() + api.ReCapPath
	progress, err = getSceneProgress(path, sceneID, bearer.AccessToken)
	return
}

// GetSceneResults requests result in a specified format
//	Note: The link specified in SceneResultReplies will be available for the time specified in reply,
//	even if the scene is deleted
func (api ReCapAPI) GetSceneResults(sceneID string, format string) (result SceneResultReply, err error) {
	bearer, err := api.Authenticator.GetToken("data:read")
	if err != nil {
		return
	}
	path := api.Authenticator.GetHostPath() + api.ReCapPath
	result, err = getSceneResult(path, sceneID, bearer.AccessToken, format)
	return
}

// CancelSceneProcessing stops the scene processing, without affecting the already uploaded resources
func (api ReCapAPI) CancelSceneProcessing(sceneID string) (ID string, err error) {
	bearer, err := api.Authenticator.GetToken("data:write")
	if err != nil {
		return
	}
	path := api.Authenticator.GetHostPath() + api.ReCapPath
	_, err = cancelSceneProcessing(path, sceneID, bearer.AccessToken)

	return sceneID, err
}

// DeleteScene removes all the resources associated with given scene.
func (api ReCapAPI) DeleteScene(sceneID string) (ID string, err error) {
	bearer, err := api.Authenticator.GetToken("data:write")
	if err != nil {
		return
	}
	path := api.Authenticator.GetHostPath() + api.ReCapPath
	_, err = deleteScene(path, sceneID, bearer.AccessToken)
	ID = sceneID
	return
}
