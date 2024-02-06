package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Mopcho/Golang-api/internal/database"
	"github.com/go-chi/chi"
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
		ID:     uuid.New(),
		UserID: user.ID,
		FeedID: params.FeedID,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error following feed: %v", err))
		return
	}

	respondWithJson(w, 201, databaseFeedFollowToFeedFollow(feedFollow))
}

func (apiCnfg *apiConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := apiCnfg.DB.GetFeedFollows(r.Context(), user.ID)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error gettings feed follows: %v", err))
		return
	}

	respondWithJson(w, 201, databaseFeedFollowsToFeedFollows(feedFollows))
}

func (apiCnfg *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIDStr := chi.URLParam(r, "feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDStr)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error deleting feed follows: %v", err))
		return
	}

	err = apiCnfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		ID:     feedFollowID,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error deleting feed follows: %v", err))
		return
	}

	respondWithJson(w, 200, struct{}{})
}
