package github

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// HTTPClient interface
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	Client HTTPClient
)

func init() {
	Client = &http.Client{}
}

// GetRepos takes a username and retreives
func GetRepos(username string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos?sort=created&direction=desc", username)

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	response, err := Client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	m := []map[string]interface{}{}
	err = json.NewDecoder(response.Body).Decode(&m)
	if err != nil {
		return nil, err
	}

	return m, nil
}
