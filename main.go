package main

import (
	"log"
	"os"
	"path/filepath"
	"weather/cmd"

	"github.com/joho/godotenv"
)

func init() {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fullPath := filepath.Join(path, ".env")
	err = godotenv.Load(fullPath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	cmd.Execute()
}
