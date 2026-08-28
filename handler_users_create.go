package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (a *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Email string `json:"email"`
	}
	userEmail := requestBody{}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&userEmail); err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong", err)
		return
	}

	dbUser, err := a.db.CreateUser(r.Context(), userEmail.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	parsedUser := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	respondWithJSON(w, http.StatusCreated, parsedUser)
}
