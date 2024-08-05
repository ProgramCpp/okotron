package db

import (
	"os"
	"strings"

	"github.com/redis/go-redis/v9"
)

// TODO: move this to main function with dependency injection

var rdb *redis.Client

func init(){
	rdb = redis.NewClient(&redis.Options{
		Addr:           os.Getenv("REDIS_ADDR"),
		Password:        "",
		DB:              0,
		DisableIndentity: true, // Disable set-info on connect
	})
}

func RedisClient() *redis.Client {
	return rdb
}

// TODO: long lived connection. 
// TODO: move all save calls to redis client
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
