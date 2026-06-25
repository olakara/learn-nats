package main

import (
	"fmt"
	"time"
)

func schduledJob() {
	fmt.Println("Schduled job running ..")
}

func main() {

	fmt.Println("Application started..")
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		schduledJob()
	}
}
