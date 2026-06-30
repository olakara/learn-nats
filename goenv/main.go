package main

import (
	"fmt"
	"goenv/config"
)

func main() {
	config := config.LoadConfig()

	fmt.Printf("Server Port: %d\n", config.Port)
	fmt.Printf("Database URL: %s\n", config.DatabaseURL)
	fmt.Printf("Redis URL: %s\n", config.RedisURL)
}
