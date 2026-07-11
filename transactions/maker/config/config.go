package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PersonURL string
	Addr      string
}

func LoadConfig() *Config {

	_ = godotenv.Load()

	addr := os.Getenv("PORT")
	if addr == "" {
		addr = ":8010"
	}

	return &Config{
		PersonURL: os.Getenv("PERSONURL"),
		Addr:      addr,
	}
}
