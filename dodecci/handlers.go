package handlers

import (
	"bufio"
	"github.com/jtakamine/dodecahedronci/config"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	//Eventually, take a look at the header/body to determine which handler to use.  For now assume it's a github request
	err := gitHubHandle(w, r)
	if err != nil {
		log.Panicf("Error handling GitHub request: %v\n", err)
	}
}

func cloneOrUpdateGitRepo(repoId int, repoUrl string) (dir string, err error) {
	dir = strings.TrimSuffix(config.Get("DODEC_HOME"), "/") + "/" + strconv.Itoa(repoId)

	var cmd *exec.Cmd

	if fInfo, err := os.Stat(dir); os.IsNotExist(err) || !fInfo.IsDir() {
		log.Printf("Cloning git repo from %v\n", repoUrl)
		cmd = exec.Command("git", "clone", repoUrl, dir)
	} else if err == nil {
		log.Printf("Pulling git repo from %v\n", repoUrl)
		cmd = exec.Command("git", "pull", repoUrl)
		cmd.Dir = dir
	} else {
		return "", err
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return "", err
	}

	return dir, nil
}

func buildDockerImages(repoDir string) (err error) {
	dockerFiles := []string{}

	walk := func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(info.Name(), "Dockerfile") {
			dockerFiles = append(dockerFiles, path)
		}

		return err
	}

	err = filepath.Walk(repoDir, walk)
	if err != nil {
		return err
	}

	for _, dFile := range dockerFiles {
		log.Printf("Building Docker file: %v\n", dFile)

		imgName, err := getImageNameHint(dFile)
		if err != nil {
			return err
		}

		cmd := exec.Command("docker", "build", "-t", config.Get("DODEC_DOCKER_USER")+"/"+imgName, ".")
		cmd.Dir = filepath.Dir(dFile)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func getImageNameHint(dockerFile string) (hint string, err error) {
	hint = "builtbydodecci" //default image name hint
	hintPrefix := "#imagenamehint:"

	file, err := os.Open(dockerFile)
	if err != nil {
		return "", nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, hintPrefix) {
			hint = strings.TrimPrefix(line, hintPrefix)
			return hint, nil
		}
	}

	err = scanner.Err()
	if err != nil {
		return "", err
	}

	return hint, nil
}
