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

	var f *os.File
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		f, err = os.Create(filename)
		if err != nil {
			log.Fatalf("Failed to create file")
		}
	} else {
		f, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			log.Fatalf("Failed to open file")
		}
	}

	defer f.Close()

	_, err := f.WriteString("xxx\n")
	if err != nil {
		fmt.Println(err)
	}
}