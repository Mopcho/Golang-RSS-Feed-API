package main

import "net/http"

func handle_readiness(w http.ResponseWriter, r *http.Request) {
	respondWithJson(w, 200, struct{}{})
}
