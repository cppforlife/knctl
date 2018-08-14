package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Print("Simple app runningzzz...")
	msg := os.Getenv("SIMPLE_MSG_WITHOUT_DOCKERFILE")
	if msg == "" {
		msg = ":( SIMPLE_MSG_WITHOUT_DOCKERFILE variable not defined"
	}
	fmt.Fprintf(w, "<h1>%s</h1>", msg)
}

func main() {
	flag.Parse()
	log.Print("Simple app server started...")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
