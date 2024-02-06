package main

import (
	"fmt"
	"net/http"

	"github.com/Mopcho/Golang-api/internal/auth"
	"github.com/Mopcho/Golang-api/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCnfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)

		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Auth error: %v", err))
			return
		}

		user, err := apiCnfg.DB.GetUserByApiKey(r.Context(), apiKey)

		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Error retrieving user: %v", err))
			return
		}

		handler(w, r, user)
	}
}
