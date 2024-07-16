package db

import "os"

// TODO: long lived connection
func Save(key, value string) error {
	return nil
}

func Get(key string) string {
	return os.Getenv("OKTO_TOKEN")
}