package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"time"
)

type Application struct {
	Name        string
	Description string
}

type Build struct {
	UUID    string
	AppName string
	Version string
}

type BuildDetails struct {
	UUID      string
	AppName   string
	Started   time.Time
	Completed time.Time
	Success   bool
	Version   string
	Artifact  string
}

type Deploy struct {
	UUID      string
	AppName   string
	BuildUUID string
}

type DeployDetails struct {
	UUID      string
	AppName   string
	Started   time.Time
	Completed time.Time
	Success   bool
	BuildUUID string
}

type Log struct {
	TaskUUID string
	Message  string
	Severity int
	Created  time.Time
}

func initApp() (app *cli.App) {
	app = cli.NewApp()
	app.Name = "dodec-cli"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "endpoint, e",
			Usage:  "CLI target endpoint",
			EnvVar: "DODEC_ENDPOINT",
		},
	}
	app.Action = func(c *cli.Context) {
		fmt.Println("Display usage info here...")
	}

	app.Commands = []cli.Command{
		{
			Name:      "listbuilds",
			ShortName: "lb",
			Usage:     "List all builds",
			Action:    newAction(listBuilds),
		},
		{
			Name:      "getbuild",
			ShortName: "gb",
			Usage:     "Get a build by its UUID",
			Action:    newAction(getBuild),
		},
		{
			Name:      "listdeploys",
			ShortName: "ld",
			Usage:     "List all deploys",
			Action:    newAction(listDeploys),
		},
		{
			Name:      "getdeploy",
			ShortName: "gd",
			Usage:     "Get a deploy by its UUID",
			Action:    newAction(getDeploy),
		},
		{
			Name:      "getlogs",
			ShortName: "gl",
			Usage:     "Get build or deploy logs for the specified UUID",
			Action:    newAction(getLogs),
		},
	}

	return app
}

func listBuilds(endpt string, c *cli.Context) {
	addr := endpt + "builds"

	var builds []Build
	err := req(addr, "GET", &builds)
	if err != nil {
		panic(err)
	}

	printRows(builds)
}

func getBuild(endpt string, c *cli.Context) {
	uuid := requireArg(c.Args(), "getbuild", "Build UUID", "de08deb8e1b0ce5a")
	addr := endpt + "builds/" + uuid

	var b BuildDetails
	err := req(addr, "GET", &b)
	if err != nil {
		panic(err)
	}

	if b.UUID != "" {
		printObj(b)
	} else {
		fmt.Printf("No build found with UUID = \"%s\"\n", uuid)
	}

}

func listDeploys(endpt string, c *cli.Context) {
	addr := endpt + "deploys"

	var deploys []Deploy
	err := req(addr, "GET", &deploys)
	if err != nil {
		panic(err)
	}

	printRows(deploys)
}

func getDeploy(endpt string, c *cli.Context) {
	uuid := requireArg(c.Args(), "getdeploy", "Deploy UUID", "de08deb8e1b0ce5a")
	addr := endpt + "deploys/" + uuid

	var d DeployDetails
	err := req(addr, "GET", &d)
	if err != nil {
		panic(err)
	}

	if d.UUID != "" {
		printObj(d)
	} else {
		fmt.Printf("No deploy found with UUID = \"%s\"\n", uuid)
	}
}

func getLogs(endpt string, c *cli.Context) {
	uuid := requireArg(c.Args(), "getlogs", "Build/Deploy UUID", "de08deb8e1b0ce5a")
	addr := endpt + "task/" + uuid + "/logs"

	var logs []Log
	err := req(addr, "GET", &logs)
	if err != nil {
		panic(err)
	}

	printLogs(logs)
}
