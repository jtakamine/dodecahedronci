package main

import (
	"strconv"
	"strings"
	"time"
)

//in-memory log storage--eventually this should be persistent
var buildAddr string
var deployAddr string

type Application struct {
	Name        string
	Description string
}

type Task struct {
	UUID    string
	AppName string
}

type TaskDetails struct {
	Task
	Started   time.Time
	Completed time.Time
	Success   bool
}

type Build struct {
	Task
	Version string
}

type BuildDetails struct {
	TaskDetails
	Version  string
	Artifact string
}

type Deploy struct {
	Task
	BuildUUID string
}

type DeployDetails struct {
	TaskDetails
	BuildUUID string
}

type Log struct {
	TaskUUID string
	Message  string
	Severity int
	Created  time.Time
}

//Log Format: "[src][taskUUID][logType] time_RFC3339\t| msg"
func parseLog(log string) (l Log) {
	whitespace := "\t "

	l = Log{
		TaskUUID: "default",
		Message:  log,
		Severity: 0,
		Created:  time.Now(),
	}

	log = strings.TrimLeft(log, whitespace)

	if !strings.HasPrefix(log, "[") {
		return l
	}

	i := strings.Index(log, "|")
	if i < 0 {
		return l
	}

	pre := strings.Split(log, "|")[0]
	if strings.Count(pre, "[") != 3 || strings.Count(pre, "]") != 3 {
		return l
	}

	pre = strings.TrimRight(pre, whitespace)
	parts := strings.Split(pre, "][")
	if len(parts) != 3 {
		return l
	}

	lastParts := strings.Split(parts[2], "]")
	if len(lastParts) != 2 {
		return l
	}

	lastParts[1] = strings.Trim(lastParts[1], whitespace)
	if _, err := time.Parse(time.RFC3339, lastParts[1]); err != nil {
		return l
	}

	if _, err := strconv.Atoi(lastParts[0]); err != nil {
		return l
	}

	uuid := parts[1]
	msg := strings.Join(strings.Split(log, "|")[1:], "|")
	severity, _ := strconv.Atoi(lastParts[0])
	created, _ := time.Parse(time.RFC3339, lastParts[1])
	//src = strings.TrimPrefix(parts[0], "[")

	l = Log{
		TaskUUID: uuid,
		Message:  msg,
		Severity: severity,
		Created:  created,
	}

	return l
}
