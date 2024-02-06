package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Mopcho/Golang-api/internal/database"
	"github.com/google/uuid"
)

func (apiCnfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	decoder.Decode(&params)

	feedFollow, err := apiCnfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		UserID: user.ID,
		FeedID: params.FeedID,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error following feed: %v", err))
		return
	}

	respondWithJson(w, 201, databaseFeedFollowToFeedFollow(feedFollow))
}
