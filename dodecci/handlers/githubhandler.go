package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type gitHubReq struct {
	Repository struct {
		Id        int
		Ssh_url   string
		Clone_url string
	}
}

func gitHubHandle(w http.ResponseWriter, r *http.Request) (err error) {
	req := &gitHubReq{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(req)
	if err != nil {
		return err
	}

	repoDir, err := cloneOrUpdateGitRepo(req.Repository.Id, req.Repository.Clone_url)
	if err != nil {
		return err
	}

	err = buildDockerImages(repoDir)
	if err != nil {
		return err
	}

	fmt.Fprint(w, "build successful\n")

	return nil
}
