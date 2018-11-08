package da

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
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
		content, _ := ioutil.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	data, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return
	}

	nickname = string(data)

	return
}
