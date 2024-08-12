package db

import (
	"os"

	"github.com/redis/go-redis/v9"
)

// all the redis namespaces for keys. equivalent to tables in SQL
const (
	// okto
	OKTO_AUTH_TOKEN_KEY = "okto_auth_token_%d"
	OKTO_ADDRESSES_KEY  = "okto_addresses_%d"

	// google
	GOOGLE_ID_TOKEN_KEY = "google_id_token_%d"

	// telegram
	SUB_COMMAND_KEY = "subcommand_%d"
	REQUEST_KEY     = "request_%d"
)

// TODO: move this to main function with dependency injection
var rdb *redis.Client

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:             os.Getenv("REDIS_ADDR"),
		Password:         "",
		DB:               0,
		DisableIndentity: true, // Disable set-info on connect
	})
}

func RedisClient() *redis.Client {
	return rdb
}
