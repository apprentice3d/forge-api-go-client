package md

import (
	"encoding/json"
	"errors"
	"io"
	"log"
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

func getDerivative(path, urn, derivativeUrn, token string, writer io.Writer) (written int64, err error) {

	task := http.Client{}

	req, err := http.NewRequest("GET", path+"/"+urn+"/manifest/"+derivativeUrn+"/signedcookies", nil)
	if err != nil {
		return 0, err
	}

	log.Println("Requesting derivative URN...")
	log.Println("- Base64 encoded design URL: ", urn)
	log.Println("- URL: ", req.URL.String())

	req.Header.Set("Authorization", "Bearer "+token)
	response, err := task.Do(req)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		content, _ := io.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return 0, err
	}

	var getUrlResult derivativeDownloadUrl

	// deserialize the response
	err = json.NewDecoder(response.Body).Decode(&getUrlResult)
	if err != nil {
		return 0, err
	}
	// https://aps.autodesk.com/en/docs/model-derivative/v2/reference/http/urn-manifest-derivativeUrn-signedcookies-GET/#http-headers
	// Signed cookie to use with download URL.
	// There will be three headers in the response named Set-Cookie
	if len(response.Header.Values("Set-Cookie")) != 3 {
		err = errors.New("invalid number of Set-Cookie headers in the response")
		return 0, err
	}

	signedCookieValue := strings.Join(response.Header.Values("Set-Cookie"), ";")

	return downloadDerivative(getUrlResult, signedCookieValue, writer)
}

func downloadDerivative(downloadUrl derivativeDownloadUrl, cookieValue string, writer io.Writer) (
	written int64, err error,
) {

	task := http.Client{}

	req, err := http.NewRequest("GET", downloadUrl.Url, nil)
	if err != nil {
		return
	}

	log.Println("Start downloading derivative...")
	log.Println("- URL: ", downloadUrl.Url)

	req.Header.Set("Cookie", cookieValue)
	req.Header.Set("Content-Type", downloadUrl.ContentType)
	req.Header.Set("Content-Length", strconv.Itoa(downloadUrl.Size))
	response, err := task.Do(req)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		content, _ := io.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return 0, err
	}

	written, err = io.Copy(writer, response.Body)
	if err != nil {
		return 0, err
	}

	if written != int64(downloadUrl.Size) {
		err = errors.New("downloaded file size is different than the expected size")
		return
	}

	log.Println("Finished downloading derivative.")

	return written, nil
}
