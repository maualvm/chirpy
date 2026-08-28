package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/maualvm/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
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

func validateChirp(body string) (string, error) {
	if len(body) > 140 {
		return "", errors.New("Chirp is too long")
	}

	return replaceProfaneWords(body), nil
}

func (a *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type createChirpRequest struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	createRequest := createChirpRequest{}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&createRequest); err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not decode parameters", err)
		return
	}

	cleanedBody, err := validateChirp(createRequest.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	dbChirp, err := a.db.CreateChrip(r.Context(), database.CreateChripParams{
		Body:   cleanedBody,
		UserID: createRequest.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        dbChirp.ID,
		Body:      dbChirp.Body,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		UserID:    dbChirp.UserID,
	})
}
