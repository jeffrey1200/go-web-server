package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type paramaters struct {
		Body string `json:"body"`
	}
	// type returnVals struct {
	// 	CleanedBody string `json:"cleaned_body"`
	// }
	// type cleanedReturnVals struct {
	// Cleaned_body string `json:"cleaned_body"`
	// }
	decoder := json.NewDecoder(r.Body)
	params := paramaters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode paramaters")
		return
	}
	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	// fmt.Printf("does the cleaned body exist? %s\n", cleaned)
	chirp, err := cfg.DB.CreateChirp(cleaned)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}
	// fmt.Printf("does the created chirp in the struct exist? %v\n", chirp)
	respondWithJSON(w, http.StatusCreated, Chirp{ID: chirp.ID, Body: chirp.Body})
	// const maxChirpLength = 140
	// if len(params.Body) > maxChirpLength {
	// 	respondWithError(w, http.StatusBadRequest, "Chirp is too long")
	// 	return
	// }
	// badWords := map[string]struct{}{
	// 	"kerfuffle": {},
	// 	"sharbert":  {},
	// 	"fornax":    {},
	// }
	// cleanedString := returnVals{CleanedBody: removeProfaneWords(params.Body, badWords)}
	// respondWithJSON(w, http.StatusOK, cleanedString)

}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("chirp is too long")
	}

	badwords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleanedBody := removeProfaneWords(body, badwords)
	return cleanedBody, nil
}

func removeProfaneWords(words string, badWords map[string]struct{}) string {

	separatedWords := strings.Split(words, " ")
	// log.Print(words)
	for i, word := range separatedWords {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			separatedWords[i] = "****"
		}
		// separatedWords[i] = strings.ToLower(word)
		// if strings.HasSuffix(separatedWords[i], "!") {
		// 	continue
		// 	}
		// log.Printf("the separated words are: %s", separatedWords[i])
		// if strings.ToLower(separatedWords[i]) == "kerfuffle" || strings.ToLower(separatedWords[i]) == "sharbert" || strings.ToLower(separatedWords[i]) == "fornax" {

		// separatedWords[i] = strings.ReplaceAll(separatedWords[i], "kerfuffle", "****")
		// separatedWords[i] = strings.ReplaceAll(separatedWords[i], "sharbert", "****")
		// separatedWords[i] = strings.ReplaceAll(separatedWords[i], "Fornax", "****")
		// }
	}
	normalizedString := strings.Join(separatedWords, " ")

	return normalizedString
}
