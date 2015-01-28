package main

import (
	"encoding/json"
	dodecpubsub_API "github.com/jtakamine/dodecahedronci/dodecpubsub/api"
	dodecregistry_API "github.com/jtakamine/dodecahedronci/dodecregistry/api"
	"gopkg.in/yaml.v2"
	"strconv"
	"time"
)

var saveBuild = func(app string, version string, fFile figFile, dockerRegistryUrl string) (err error) {
	data, err := yaml.Marshal(fFile.Config)
	if err != nil {
		return err
	}
	artifact := string(data)

	build := dodecregistry_API.Build{Artifact: artifact, DockerRegistryUrl: dockerRegistryUrl}

	err = dodecregistry_API.PostBuild(app, version, build, "http://dodecregistry:8000/")
	if err != nil {
		return err
	}

	return nil
}

var log = func(msg string, lType logType) (err error) {
	msgObj := &struct {
		Msg       string
		TimeStamp time.Time
	}{
		msg,
		time.Now(),
	}

	msgData, err := json.Marshal(msgObj)
	if err != nil {
		panic(err.Error())
	}
	msgJson := string(msgData)

	for i := 0; i <= int(lType); i++ {
		err = dodecpubsub_API.Publish(msgJson, strconv.Itoa(i), "dodecpubsub:6379")
		if err != nil {
			panic(err.Error())
		}
	}

	return nil
}
