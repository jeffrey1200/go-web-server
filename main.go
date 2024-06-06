package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	mux := http.NewServeMux()
	srv := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	fileHandler := http.FileServer(http.Dir("./"))
	mux.Handle("/", fileHandler)
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
