package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"strings"
)

func InitDB() (db *sql.DB, err error) {
	connStr, err := GetConnStr()
	if err != nil {
		return err
	}

	*db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

}

func GetConnStr() (connStr string, err error) {
	missing := []string{}

	userEnv := "POSTGRES_USER"
	passEnv := "POSTGRES_PASSWORD"
	dbAddrEnv := "DODEC_POSTGRESADDR"

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
		return "", errrors.New(msg)
	}

	return fmt.Sprintf("user=%s password=%s host=%s port=%s sslmode=disable"), nil
}
