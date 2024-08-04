package db

import (
	"os"
	"strings"
)

// TODO: long lived connection. which db to use?
func Save(key, value string) error {
	return nil
}

func Get(key string) string {
	if strings.Contains(key, "message") {
		return "/swap/source-token" // "/setup-profile"
	} else if strings.Contains(key, "okto_token") {
		return os.Getenv("OKTO_TOKEN")
	} else if strings.Contains(key, "okto_auth_token") {
		return os.Getenv("OKTO_AUTH_TOKEN")
	} else if strings.Contains(key, "google_id_token") {
		return os.Getenv("GOOGLE_TOKEN")
	}

	return ""
}
