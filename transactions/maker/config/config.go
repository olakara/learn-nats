package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PersonURL string
}

func LoadConfig() *Config {

	_ = godotenv.Load()
	return &Config{
		PersonURL: os.Getenv("PERSONURL"),
	}
}
