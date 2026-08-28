package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type successResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}

	chirp := chirp{}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&chirp); err != nil {
		log.Printf("Error decoding JSON request: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if len(chirp.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	responseBody := successResponse{
		CleanedBody: replaceProfaneWords(chirp.Body),
	}

	respondWithJSON(w, http.StatusOK, responseBody)
}

func replaceProfaneWords(chirpBody string) string {
	profaneWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	chirpWords := strings.Split(chirpBody, " ")
	for i, chirpWord := range chirpWords {
		lowerChirpWord := strings.ToLower(chirpWord)
		if _, ok := profaneWords[lowerChirpWord]; ok {
			chirpWords[i] = "****"
		}
	}

	return strings.Join(chirpWords, " ")
}
