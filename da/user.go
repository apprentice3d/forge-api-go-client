package da

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func getUserID(path string, token string) (nickname string, err error) {
	task := http.Client{}
	req, err := http.NewRequest("GET",
		path+"/forgeapps/me",
		nil,
	)

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

	data, err := io.ReadAll(response.Body)

	if err != nil {
		return
	}

	//TODO: Review why the data has quotes in its content and find a more elegant way to remove them
	nickname = strings.Replace(string(data), "\"", "", -1)

	return
}
