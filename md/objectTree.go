package md

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
)

type ObjectTree struct {
	Data struct {
		Type    string           `json:"type"`
		Objects []ObjectTreeNode `json:"objects"`
	} `json:"data"`
}

type ObjectTreeNode struct {
	ObjectId int              `json:"objectid"`
	Name     string           `json:"name"`
	Objects  []ObjectTreeNode `json:"objects"`
}

func getObjectTree(path, urn, guid, token string, forceGet bool, xHeaders XAdsHeaders) (
	result ObjectTree, err error,
) {
	task := http.Client{}

	url := path + "/" + urn + "/metadata/" + guid + "?forceget=" + strconv.FormatBool(forceGet)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Add("x-ads-force", strconv.FormatBool(xHeaders.Overwrite))
	req.Header.Add("x-ads-derivative-format", string(xHeaders.Format))

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

	err = json.NewDecoder(response.Body).Decode(&result)

	return
}
