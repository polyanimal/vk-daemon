package main

import (
	"bytes"
	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func startUpLog(output *os.File) error {
	lsusb := exec.Command("lsusb")
	res, err := lsusb.Output()
	if err != nil {
		return err
	}

	for _, l := range bytes.SplitAfter(res, []byte("\n")) {
		output.WriteString(string(l))
	}

	return nil
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

	err := startUpLog(f)
	if err != nil {
		log.Println("error:", err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add("/var/log/syslog")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	return
}
