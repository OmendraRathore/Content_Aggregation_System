package main

import (
	"net/http"

	"github.com/OmendraRathore/Content_Aggregation_System/internal/auth"
	"github.com/OmendraRathore/Content_Aggregation_System/internal/database"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (app *appConfig) withAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Couldn't find api key")
			return
		}

		user, err := app.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Couldn't get user")
			return
		}

		handler(w, r, user)
	}
}
