package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        int
	DatabaseURL string
	RedisURL    string
}

func LoadConfig() *Config {

	_ = godotenv.Load()
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	databaseURL := os.Getenv("DATABASE_URL")
	redisURL := os.Getenv("REDIS_URL")

	return &Config{
		Port:        port,
		DatabaseURL: databaseURL,
		RedisURL:    redisURL,
	}
}
