package main

import (
	"io"
	"log"
	"net/http"
)

func main() {

	r, w := io.Pipe()
	log.Default().SetOutput(w)

	go func() {
		log.Println("hello world 1")
		// log.Println("this is your brand new instantlogs")
		// log.Println("check the console https://instantlogs.io/")
		w.Close()
	}()

	req, _ := http.NewRequest(http.MethodPost, "https://instantlogs.io/v1/loggers/logger-d8da571e-6056-47ec-8c98-1034f1865443/ingest", r)
	req.Header.Set("Api-Key", "96bb6b95-acca-4a59-8a65-effe8b1cb00d")
	req.Header.Set("Api-Secret", "26a96ac3-4ae7-4983-9df8-cc53a28fbc3d")
	http.DefaultClient.Do(req)

}
