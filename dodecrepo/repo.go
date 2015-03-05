package main

import (
	"time"
)

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

type TaskCompletionInfo struct {
	UUID    string
	Success bool
}

type Artifact struct {
	Artifact  string
	Type      string
	BuildUUID string
}

type LogQuery struct {
	TaskUUID string
	Severity int
}

type Log struct {
	TaskUUID string
	Message  string
	Severity int
	Created  time.Time
}
