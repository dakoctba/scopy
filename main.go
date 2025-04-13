package main

import (
	"log"
	"os"

	"github.com/dakoctba/scopy/cmd"
	"github.com/joho/godotenv"
)

func init() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// It's okay if the .env file doesn't exist in production/deployment
		if _, err := os.Stat(".env"); !os.IsNotExist(err) {
			log.Printf("Warning: could not load .env file: %v", err)
		}
	}
}

func main() {
	cmd.Execute()
}
