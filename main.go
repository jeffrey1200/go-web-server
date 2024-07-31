package main

import (
	"log"
	"net/http"

	"github.com/jeffrey1200/go-web-server/internal/database"
)

type apiConfig struct {
	fileServerHits int
	DB             *database.DB
}

func main() {
	const filePathRoot = "."
	const port = "8080"

	db, err := database.NewDB("/home/jeffrey/Documents/go-web-server/database.json")
	if err != nil {
		log.Fatal(err)
	}
	apiCfg := apiConfig{
		fileServerHits: 0,
		DB:             db,
	}

	mux := http.NewServeMux()
	fileHandler := http.FileServer(http.Dir(filePathRoot))
	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileHandler)))
	mux.HandleFunc("GET /api/healthz", handleReadiness)
	mux.HandleFunc("GET /api/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsRetrieve)
	mux.HandleFunc("GET /api/chirps/{id}", apiCfg.handlerRetrieveChirpById)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	srv := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
