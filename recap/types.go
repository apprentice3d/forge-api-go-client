package recap

// PhotoScene holds data encountered in replies like creation of photoScene
type PhotoScene struct {
	ID      string   `json:"photosceneid"`
	Name    string   `json:"name,omitempty"`
	Files   []string `json:",omitempty"`
	Formats []string `json:",omitempty"`
	Metadata []struct {
		Name   string
		Values string
	} `json:",omitempty"`
}

// SceneCreationReply reflects the response content upon scene creation
type SceneCreationReply struct {
	Usage      string     `json:",omitempty"`
	Resource   string     `json:",omitempty"`
	PhotoScene PhotoScene `json:"Photoscene,omitempty"`
	Error      *Error     `json:"Error,omitempty"`
}

// SceneDeletionReply reflects the response content upon scene deletion
type SceneDeletionReply struct {
	Usage    string `json:",omitempty"`
	Resource string `json:",omitempty"`
	Message  string `json:"msg"`
	Error    *Error `json:"Error,omitempty"`
}

// SceneCancelReply reflects the response content upon scene cancel processing
type SceneCancelReply struct {
	Usage    string `json:",omitempty"`
	Resource string `json:",omitempty"`
	Message  string `json:"msg"`
	Error    *Error `json:"Error,omitempty"`
}

// FileUploadingReply reflects the response content upon uploading a file,
// be it a link or a local one
type FileUploadingReply struct {
	Usage    string `json:",omitempty"`
	Resource string `json:",omitempty"`
	Files *struct {
		File struct {
			FileName string `json:"filename"`
			FileID   string `json:"fileid"`
			FileSize string `json:"filesize"`
			Message  string `json:"msg"`
		} `json:"file"`
	} `json:"Files"`

	Error *Error `json:"Error,omitempty"`
}

//type LinksUploadingReply struct {
//	Usage    string `json:",omitempty"`
//	Resource string `json:",omitempty"`
//	Files    *struct {
//		File []struct {
//			FileName string `json:"filename"`
//			FileID   string `json:"fileid"`
//			FileSize string `json:"filesize"`
//			Message  string `json:"msg"`
//		} `json:"file"`
//	} `json:"Files"`
//	Error    *struct {
//		Code    string `json:"code"`
//		Message string `json:"msg"`
//	} `json:"Error"`
//}

// SceneStartProcessingReply reflects the response content upon starting scene processing
type SceneStartProcessingReply struct {
	Message    string     `json:"msg"`
	PhotoScene PhotoScene `json:"Photoscene"`
	Error      *Error     `json:"Error,omitempty"`
}

// SceneProgressReply reflects the response content upon polling for scene status
type SceneProgressReply struct {
	Usage    string `json:",omitempty"`
	Resource string `json:",omitempty"`
	PhotoScene struct {
		ID       string `json:"photosceneid"`
		Message  string `json:"progressmsg"`
		Progress string `json:"progress"`
	} `json:"Photoscene"`

	Error *Error `json:"Error,omitempty"`
}



// SceneResultReply reflects the response content upon requesting the scene results in a certain format
type SceneResultReply struct {
	PhotoScene struct {
		ID        string `json:"photosceneid"`
		Message   string `json:"progressmsg"`
		Progress  string `json:"progress"`
		SceneLink string `json:"scenelink"`
		FileSize  string `json:"filesize"`
	} `json:"Photoscene"`

	Error *Error `json:"Error,omitempty"`
}

// BUG(apprentice3d): SceneResultReply has a slightly different schema when getting results from a successfully
// processed scene and one that failed. In this situation, error of type [JSON DECODING ERROR] should be considered
// as [SCENE FAILED TO PROCESS]

// ErrorMessage represents a struct corresponding to successfully received task, but failed due to some reasons.
// Check the bug section of this documentation for more info.
type ErrorMessage struct {
	Usage    string `json:",omitempty"`
	Resource string `json:",omitempty"`
	Error    *Error `json:"Error"`
}

// BUG(apprentice3d) Frequently the operation succeeded with returning code 200, meaning that the task was
// received successfully, but failed to execute due to reasons specified in message
// (g.e. uploading a file by specifying an wrong link: POST request is successful,
// but internally it failed to download the file because of the wrongly provided link)

// Error is inner struct encountered in cases when the server reported status OK, but still contains details
// on encountered errors. Check the bug section of this documentation for more info.
// 	This bug was reported to the engineering team
type Error struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
}
