package main

import (
	"net/http"
	"fmt"
	"github.com/ScriptedPranav/rss-aggregator/internal/database"
	"github.com/ScriptedPranav/rss-aggregator/internal/database/auth"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, http.StatusNetworkAuthenticationRequired, fmt.Sprintf("Error getting API key: %v", err))
			return
		}
	
		user, err := cfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusNotFound, fmt.Sprintf("Error getting user: %v", err))
			return
		} 

		handler(w, r, user)
	}
}