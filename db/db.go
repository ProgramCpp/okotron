package db

import (
	"os"
	"strings"
)

// TODO: long lived connection
func Save(key, value string) error {
	return nil
}

func Get(key string) string {
	if strings.Contains(key, "message"){
		return "/pin"
	} else if strings.Contains(key, "okto_token"){
		return os.Getenv("OKTO_TOKEN")
	} else if strings.Contains(key, "okto_auth_token"){
		return os.Getenv("OKTO_AUTH_TOKEN")
	}
		
	return ""
}