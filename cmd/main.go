package main

import (
	"log"

	"github.com/flpnascto/rate-limiter-go/internal/infra/database"
	"github.com/flpnascto/rate-limiter-go/internal/infra/web"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile("./cmd/config.json")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error loading config file: %s", err)
	}
	viper.AutomaticEnv()
}

func main() {

	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	db := database.NewRedisClient(client)
	web.NewWebServer(db)
}
