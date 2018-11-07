package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Print("Simple app runningzzz...")
	msg := os.Getenv("SIMPLE_MSG")
	if msg == "" {
		msg = ":( SIMPLE_MSG variable not defined"
	}
	fmt.Fprintf(w, "<h1>%s</h1>", msg)
}

func main() {
	go func() {
		time.Sleep(1 * time.Second)
		fmt.Printf("app-is-exiting\n")
		os.Exit(2)
	}()

	flag.Parse()
	log.Print("Simple app server started...")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
