package main

import (
	"strconv"
	"strings"
	"time"
)

//in-memory log storage--eventually this should be persistent
var logs = make(map[string]map[string][]LogEntry)

type LogEntry struct {
	Type int
	Msg  string
}

//Log Format: "[src][taskID][logType] time_RFC3339\t| msg"
func parseLog(log string) (msg string, t time.Time, logType int, taskID string, src string) {
	msg = log
	t = time.Now()
	logType = 0
	taskID = "default"
	src = "default"
	whitespace := "\t "

	msg = strings.TrimLeft(msg, whitespace)

	if !strings.HasPrefix(msg, "[") {
		return msg, t, logType, taskID, src
	}

	i := strings.Index(msg, "|")
	if i < 0 {
		return msg, t, logType, taskID, src
	}

	pre := strings.Split(msg, "|")[0]
	if strings.Count(pre, "[") != 3 || strings.Count(pre, "]") != 3 {
		return msg, t, logType, taskID, src
	}

	pre = strings.TrimRight(pre, whitespace)
	parts := strings.Split(pre, "][")
	if len(parts) != 3 {
		return msg, t, logType, taskID, src
	}

	lastParts := strings.Split(parts[2], "]")
	if len(lastParts) != 2 {
		return msg, t, logType, taskID, src
	}

	lastParts[1] = strings.Trim(lastParts[1], whitespace)
	if _, err := time.Parse(time.RFC3339, lastParts[1]); err != nil {
		return msg, t, logType, taskID, src
	}

	if _, err := strconv.Atoi(lastParts[0]); err != nil {
		return msg, t, logType, taskID, src
	}

	msg = strings.Join(strings.Split(msg, "|")[1:], "|")
	t, _ = time.Parse(time.RFC3339, lastParts[1])
	logType, _ = strconv.Atoi(lastParts[0])
	taskID = parts[1]
	src = strings.TrimPrefix(parts[0], "[")

	return msg, t, logType, taskID, src
}
