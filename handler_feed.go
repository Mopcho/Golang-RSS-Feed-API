package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Mopcho/Golang-api/internal/database"
	"github.com/google/uuid"
)

func (apiCnfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `name`
		Url  string `url`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	decoder.Decode(&params)

	feed, err := apiCnfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:     uuid.New(),
		Name:   params.Name,
		Url:    params.Url,
		UserID: user.ID,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating feed: %v", err))
		return
	}

	respondWithJson(w, 201, databaseFeedToFeed(feed))
}

func (apiCnfg *apiConfig) handleGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCnfg.DB.GetFeeds(r.Context())

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get feeds: %v", err))
		return
	}

	respondWithJson(w, 200, databaseFeedsToFeeds(feeds))
}
