package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/OmendraRathore/Content_Aggregation_System/internal/database"
	"github.com/google/uuid"
)

func (app *appConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Name string
	}

	jsonDecoder := json.NewDecoder(r.Body)
	payload := requestBody{}
	err := jsonDecoder.Decode(&payload)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := app.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      payload.Name,
	})
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}

func (app *appConfig) getUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}
