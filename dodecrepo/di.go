package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	//"github.com/lib/pq"
)

var getConnStr = func() (connStr string, err error) {
	missing := []string{}

	userEnv := "PGUSER"
	passEnv := "PGPASSWORD"
	dbAddrEnv := "DODEC_REPODBADDR"

	user := os.Getenv(userEnv)
	if user == "" {
		missing = append(missing, userEnv)
	}

	pass := os.Getenv(passEnv)
	if pass == "" {
		missing = append(missing, passEnv)
	}

	dbAddr := os.Getenv(dbAddrEnv)
	if dbAddr == "" {
		missing = append(missing, dbAddrEnv)
	}

	if len(missing) > 0 {
		msg := fmt.Sprintf("Missing environment variables: %v", missing)
		return "", errors.New(msg)
	}

	dbAddrParts := strings.Split(dbAddr, ":")
	if len(dbAddrParts) != 2 {
		msg := fmt.Sprintf("Could not parse address: %s", dbAddr)
		return "", errors.New(msg)
	}

	host := dbAddrParts[0]
	port := dbAddrParts[1]

	return fmt.Sprintf("user=%s password=%s host=%s port=%s sslmode=disable", user, pass, host, port), nil
}
