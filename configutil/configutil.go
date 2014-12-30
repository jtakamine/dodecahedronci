package configutil

import (
	"encoding/json"
	"errors"
	"os"
)

func Load(configFile string) error {
	file, err_open := os.Open(configFile)
	if err_open != nil {
		return err_open
	}

	decoder := json.NewDecoder(file)

	var cfg map[string]string

	err_decode := decoder.Decode(&cfg)
	if err_decode != nil {
		return err_decode
	}

	for key, val := range cfg {
		Set(key, val)
	}

	return nil
}

func Get(key string) string {
	return os.Getenv(key)
}

func Set(key string, val string) {
	os.Setenv(key, val)
}

func Require(keys []string) error {
	missingKeys := []string{}
	for _, key := range keys {
		if Get(key) == "" {
			missingKeys = append(missingKeys, key)
		}
	}
	if len(missingKeys) > 0 {
		err := "The following configuration values are required, but were not supplied:\n"
		for _, key := range missingKeys {
			err += "\t- " + key + "\n"
		}

		return errors.New(err)
	}

	return nil
}
