package main

import (
	"github.com/jtakamine/dodecahedronci/configutil"
	"github.com/jtakamine/dodecahedronci/logutil"
	"os"
	"os/exec"
	"strings"
)

func cloneOrUpdateGitRepo(repoUrl string, writer *logutil.Writer) (dir string, err error) {
	w := writer.WriteType

	dir = strings.TrimSuffix(configutil.Get("DODEC_HOME"), "/") + "/" + repoUrlToDir(repoUrl)

	var cmd *exec.Cmd

	if fInfo, err := os.Stat(dir); os.IsNotExist(err) || !fInfo.IsDir() {
		w("Cloning git repo from "+repoUrl+"...", logutil.Info)
		cmd = exec.Command("git", "clone", repoUrl, dir)
	} else if err == nil {
		w("Pulling git repo from "+repoUrl+"...", logutil.Info)
		cmd = exec.Command("git", "pull", "--rebase", repoUrl)
		cmd.Dir = dir
	} else {
		return "", err
	}

	wc := writer.CreateChild()
	cmd.Stdout = wc.CreateWriter(logutil.Verbose)
	cmd.Stderr = wc.CreateWriter(logutil.Error)

	err = cmd.Run()
	if err != nil {
		return "", err
	}

	w("Done cloning/pulling git repo.", logutil.Info)

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
