package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}


func main() {
	filename, ok := os.LookupEnv("FILENAME")
	if !ok {
		log.Fatalf("Failed to export file name from env")
	}

	fmt.Println(filename)
}