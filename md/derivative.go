package md

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type derivativeDownloadUrl struct {
	Etag        string `json:"etag"`
	Size        int    `json:"size"`
	Url         string `json:"url"`
	ContentType string `json:"content-type"`
	Expiration  int64  `json:"expiration"`
}

func getDerivative(path, urn, derivativeUrn, token string) (result []byte, err error) {

	task := http.Client{}

	req, err := http.NewRequest("GET", path+"/"+urn+"/manifest/"+derivativeUrn+"/signedcookies", nil)

	if err != nil {
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	response, err := task.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		content, _ := io.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	var getUrlResult derivativeDownloadUrl

	// deserialize the response
	err = json.NewDecoder(response.Body).Decode(&getUrlResult)
	if err != nil {
		return
	}

	signedCookieValue := strings.Join(response.Header.Values("Set-Cookie"), ";")

	return downloadDerivative(getUrlResult, signedCookieValue)
}

func downloadDerivative(downloadUrl derivativeDownloadUrl, cookieValue string) (result []byte, err error) {

	task := http.Client{}

	req, err := http.NewRequest("GET", downloadUrl.Url, nil)
	if err != nil {
		return
	}

	req.Header.Set("Cookie", cookieValue)
	req.Header.Set("Content-Type", downloadUrl.ContentType)
	req.Header.Set("Content-Length", strconv.Itoa(downloadUrl.Size))
	response, err := task.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		content, _ := io.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	result, err = io.ReadAll(response.Body)
	if err != nil {
		return
	}
	if len(result) != downloadUrl.Size {
		err = errors.New("downloaded file size is different than the expected size")
		return
	}

	return result, nil
}
