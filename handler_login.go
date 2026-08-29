package main

import (
	"encoding/json"
	"net/http"

	"github.com/maualvm/chirpy/internal/auth"
)

func (a *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	request := requestBody{}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&request); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode request body", err)
		return
	}

	dbUser, err := a.db.FindUserByEmail(r.Context(), request.Email)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Could not find user with email "+request.Email, err)
		return
	}

	passwordMatched, err := auth.CheckPasswordHash(request.Password, dbUser.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error logging in", err)
		return
	}

	if !passwordMatched {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	})
}
