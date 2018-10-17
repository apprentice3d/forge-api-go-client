package md

import (
	"log"
	"net/http"
	"bytes"
	"io/ioutil"
	"errors"
	"encoding/json"
	"strconv"
)

//TranslationParams is used when specifying the translation jobs
type TranslationParams struct {
	Input struct {
		URN           string  `json:"urn"`
		CompressedURN *bool   `json:"compressedUrn,omitempty"`
		RootFileName  *string `json:"rootFileName,omitempty"`
	} `json:"input"`
	Output OutputSpec `json:"output"`
}

// TranslationResult reflects data received upon successful creation of translation job
type TranslationResult struct {
	Result string `json:"result"`
	URN    string `json:"urn"`
	AcceptedJobs struct {
		Output OutputSpec `json:"output"`
	}
}

// OutputSpec reflects data found upon creation translation job and receiving translation job status
type OutputSpec struct {
	Destination DestSpec     `json:"destination,omitempty"`
	Formats     []FormatSpec `json:"formats"`
}

// DestSpec is used within OutputSpecs and is useful when specifying the region for translation results
type DestSpec struct {
	Region string `json:"region"`
}

// FormatSpec is used within OutputSpecs and should be used when specifying the expected format and views (2d or/and 3d)
type FormatSpec struct {
	Type  string   `json:"type"`
	Views []string `json:"views"`
}


func translate(path string, params TranslationParams, token string) (result TranslationResult, err error) {

	byteParams, err := json.Marshal(params)
	if err != nil {
		log.Println("Could not marshal the translation parameters")
		return
	}

	req, err := http.NewRequest("POST",
		path+"/job",
		bytes.NewBuffer(byteParams))

	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusOK {
		content, _ := ioutil.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	decoder := json.NewDecoder(response.Body)

	err = decoder.Decode(&result)

	return
}
