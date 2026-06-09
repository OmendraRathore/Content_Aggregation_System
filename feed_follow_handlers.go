package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/OmendraRathore/Content_Aggregation_System/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (app *appConfig) listFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := app.DB.GetFeedFollowsForUser(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create feed follow")
		return
	}

	respondWithJSON(w, http.StatusOK, databaseFeedFollowsToFeedFollows(feedFollows))
}

func (app *appConfig) createFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type requestBody struct {
		FeedID uuid.UUID
	}

	jsonDecoder := json.NewDecoder(r.Body)
	payload := requestBody{}
	err := jsonDecoder.Decode(&payload)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	createdFollow, err := app.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    payload.FeedID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create feed follow")
		return
	}

	respondWithJSON(w, http.StatusOK, databaseFeedFollowToFeedFollow(createdFollow))
}

func (app *appConfig) deleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIDText := chi.URLParam(r, "feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDText)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid feed follow ID")
		return
	}

	err = app.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		ID:     feedFollowID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create feed follow")
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}
