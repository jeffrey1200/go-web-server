package main

import (
	"log"
	"net/http"
)

type apiConfig struct {
	fileServerHits int
}

func main() {
	const filePathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{
		fileServerHits: 0,
	}

	mux := http.NewServeMux()
	fileHandler := http.FileServer(http.Dir(filePathRoot))
	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileHandler)))
	mux.HandleFunc("GET /api/healthz", handleReadiness)
	mux.HandleFunc("GET /api/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/validate_chirp", handlerChirpsValidate)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	srv := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
