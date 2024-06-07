package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// type jsonRequest struct {
// 	Body string `json:"body"`
// }

// type jsonResponse struct {
// 	respError string
// 	valid     bool
// }

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type paramaters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		Valid bool `json:"valid"`
	}
	decoder := json.NewDecoder(r.Body)
	params := paramaters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode paramaters")
		return
	}
	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{Valid: true})

}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5xx error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{Error: msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(data)
}

// func validateChirpHandler(w http.ResponseWriter, r *http.Request) {

// 	decoder := json.NewDecoder(r.Body)
// 	params := jsonRequest{}
// 	err := decoder.Decode(&params)
// 	if err != nil {
// 		log.Printf("Error decoding parameters: %s", err)
// 		w.WriteHeader(500)
// 		return
// 	}
// 	// responseBody := map[string]string{}
// 	t := jsonResponse{}
// 	if len(params.Body) > 140 {
// 		t.respError = "Chirp is too large"
// 		// responseBody["error"] = "Chirp is too large"
// 	}
// 	t.valid = true
// 	data, err := json.Marshal(t)
// 	if err != nil {

// 		log.Printf("Error marshalling json: %s", err)
// 		w.WriteHeader(500)
// 		return
// 	}
// 	log.Print(t)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(200)
// 	// w.Write([]byte("the request was succesful"))
// 	w.Write(data)

// }
