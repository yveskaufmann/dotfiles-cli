package github

import (
	"encoding/json"
	"net/http"
)

type GetLatestReleaseInfo struct {
	Name    string `json:"name"`
	TagName string `json:"tag_name"`
	URL     string `json:"url"`
}

func GetLatestRelease(repo string) (string, error) {
	release := GetLatestReleaseInfo{}
	err := getJSON("https://api.github.com/repos/"+repo+"/releases/latest", &release)
	if err != nil {
		return "", err
	}

	return release.TagName, nil
}

// TODO: create helper html methods

func getJSON(url string, target interface{}) error {
	client := http.Client{}
	response, err := client.Get(url)
	if err != nil || response.StatusCode != 200 {
		return err
	}
	defer response.Body.Close()

	return json.NewDecoder(response.Body).Decode(target)
}
