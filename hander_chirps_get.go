package main

import (
	"log"
	"net/http"
	"sort"
	"strconv"
)

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}

	chirps := []Chirp{}
	for _, dbchrip := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:   dbchrip.ID,
			Body: dbchrip.Body,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerRetrieveChirpById(w http.ResponseWriter, r *http.Request) {
	pathIdString := r.PathValue("id")
	pathIdInt, _ := strconv.Atoi(pathIdString)
	log.Printf("Am I getting the id from path? %s", pathIdString)

	dbChirp, err := cfg.DB.GetChirp(pathIdInt)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not retrieve chirp")
	}
	if dbChirp.Body == "" && dbChirp.ID == 0 {
		respondWithError(w, http.StatusNotFound, "404 page not found")
	} else {

		respondWithJSON(w, http.StatusOK, dbChirp)
	}
	// if strings == 0 {
	// 	respondWithError(w,http.StatusNotFound,"This chirp doesn't exist")
	// }

}
