package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func PostBuild(app string, version string, build Build, targetUrl string) (err error) {
	reqObj := struct {
		App               string
		Version           string
		Artifact          string
		DockerRegistryUrl string
	}{
		App:               app,
		Version:           version,
		Artifact:          build.Artifact,
		DockerRegistryUrl: build.DockerRegistryUrl,
	}

	reqData, err := json.Marshal(reqObj)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", targetUrl, strings.NewReader(string(reqData)))
	if err != nil {
		return err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Printf("%v", string(respBody))
	return nil
}

func GetBuild(app string, version string, targetUrl string) (build Build, err error) {
	vals := url.Values{"app": {app}, "version": {version}}
	url := targetUrl + "?" + vals.Encode()

	resp, err := http.Get(url)
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Build{}, err
	}

	log.Printf("get build response: %v\n", string(respBody))

	build = Build{}
	err = json.Unmarshal(respBody, &build)
	if err != nil {
		return Build{}, err
	}

	return build, nil
}
