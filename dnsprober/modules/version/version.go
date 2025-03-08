package version

import (
	"encoding/json"
	"net/http"
	"time"
)

type Gitjson struct {
	Version string `json:"tag_name"`
}

func GitVersion() (string, error) {
	request, err := http.NewRequest("GET", "https://api.github.com/repos/RevoltSecurities/Dnsprober/releases/latest", nil)
	if err != nil {
		return "", err
	}
	client := &http.Client{Timeout: 10 * time.Second}

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var gitjson Gitjson
	if err := json.NewDecoder(response.Body).Decode(&gitjson); err != nil {
		return "", err
	}
	return gitjson.Version, nil
}
