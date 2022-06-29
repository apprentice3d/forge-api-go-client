package md

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
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

// XHeaders is used when specifying the translation jobs
type XHeaders struct {
	// Format => x-ads-derivative-format header, "latest" (Default) or "fallback"
	Format DerivativeFormat
	// Overwrite => x-ads-force header: false (default) or true
	Overwrite bool
}

// DefaultXHeaders gets XHeaders with default values
func DefaultXHeaders() XHeaders {
	xHeaders := XHeaders{}
	xHeaders.Format = Latest
	xHeaders.Overwrite = false
	return xHeaders
}

// NewXHeaders gets XHeaders with the given values
func NewXHeaders(format DerivativeFormat, overwrite bool) XHeaders {
	xHeaders := XHeaders{}
	xHeaders.Format = format
	xHeaders.Overwrite = overwrite
	return xHeaders
}

// Indicates the value for the xAdsHeaders.Format
type DerivativeFormat string

const (
	Latest   DerivativeFormat = "latest"
	FallBack DerivativeFormat = "fallback"
)

// TranslationResult reflects data received upon successful creation of translation job
type TranslationResult struct {
	Result       string `json:"result"`
	URN          string `json:"urn"`
	AcceptedJobs struct {
		Output OutputSpec `json:"output"`
	}
}

// AdvancedSpec

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

func translate(path string, params TranslationParams, xHeaders XHeaders, token string) (result TranslationResult, err error) {

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
	req.Header.Add("x-ads-derivative-format", string(xHeaders.Format))
	req.Header.Add("x-ads-force", strconv.FormatBool(xHeaders.Overwrite))

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
