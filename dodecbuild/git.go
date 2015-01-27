package main

import (
	"github.com/jtakamine/dodecahedronci/configutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func cloneOrUpdateGitRepo(repoUrl string) (dir string, err error) {
	dir = strings.TrimSuffix(configutil.Get("DODEC_HOME"), "/") + "/" + repoUrlToDir(repoUrl)

	var cmd *exec.Cmd

	if fInfo, err := os.Stat(dir); os.IsNotExist(err) || !fInfo.IsDir() {
		log.Printf("Cloning git repo from %v\n", repoUrl)
		cmd = exec.Command("git", "clone", repoUrl, dir)
	} else if err == nil {
		log.Printf("Pulling git repo from %v\n", repoUrl)
		cmd = exec.Command("git", "pull", "--rebase", repoUrl)
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

func repoUrlToDir(url string) string {
	dirRunes := make([]rune, 0, len(url))
	validRunes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789_"

	for _, r := range url {
		if strings.ContainsRune(validRunes, r) {
			dirRunes = append(dirRunes, r)
		}
	}

	return string(dirRunes)
}
