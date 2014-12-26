package handlers

import (
	"bufio"
	"github.com/jtakamine/dodecahedronci/config"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	//Eventually, take a look at the header/body to determine which handler to use.  For now assume it's a github request
	gitHubHandle(w, r)
}

func cloneOrUpdateGitRepo(repoId int, repoUrl string) string {
	dir := strings.TrimSuffix(config.Get("DODEC_HOME"), "/") + "/" + strconv.Itoa(repoId)

	var cmd *exec.Cmd

	if fInfo, err := os.Stat(dir); os.IsNotExist(err) || !fInfo.IsDir() {
		log.Printf("Cloning git repo from %v\n", repoUrl)
		cmd = exec.Command("git", "clone", repoUrl, dir)
	} else {
		log.Printf("Pulling git repo from %v\n", repoUrl)
		cmd = exec.Command("git", "pull", repoUrl)
		cmd.Dir = dir
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Panicf("Error running git operation: %v\n", err)
	}

	return dir
}

func buildDockerImages(repoDir string) {
	dockerFiles := []string{}

	walk := func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(info.Name(), "Dockerfile") {
			dockerFiles = append(dockerFiles, path)
		}

		return nil
	}

	err := filepath.Walk(repoDir, walk)
	if err != nil {
		log.Panicf("Error walking the directory \"%v\": %v\n", repoDir, err)
	}

	for _, dFile := range dockerFiles {
		log.Printf("Building Docker file: %v\n", dFile)

		imgName := getImageNameHint(dFile)

		cmd := exec.Command("docker", "build", "-t", config.Get("DODEC_DOCKER_USER")+"/"+imgName, ".")
		cmd.Dir = filepath.Dir(dFile)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			log.Panicf("Error building Dockerfile: %v\n", err)
		}
	}
}

func getImageNameHint(dockerFile string) string {
	hintPrefix := "#imagenamehint:"

	file, err := os.Open(dockerFile)
	if err != nil {
		log.Panicf("Error opening Dockerfile: %v\n", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, hintPrefix) {
			return strings.TrimPrefix(line, hintPrefix)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Panicf("Error reading Dockerfile: %v\n", err)
	}

	//default arbitrarily to "builtbydodecci" if no hint was found
	return "builtbydodecci"
}
