package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Mopcho/Golang-api/internal/auth"
	"github.com/Mopcho/Golang-api/internal/database"
	"github.com/google/uuid"
)

func (apiCnfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `name`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing json %v", err))
		return
	}

	user, err := apiCnfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create an user: %v", err))
		return
	}

	respondWithJson(w, 201, databaseUserToUser(user))
}

func (apiCnfg *apiConfig) handleGetUser(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Auth error: %v", err))
	}

	user, err := apiCnfg.DB.GetUserByApiKey(r.Context(), apiKey)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error retrieving user: %v", err))
		return
	}

	respondWithJson(w, 200, databaseUserToUser(user))
}
