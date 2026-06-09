package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/OmendraRathore/Content_Aggregation_System/internal/database"
	"github.com/google/uuid"
)

func (app *appConfig) createFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type requestBody struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	jsonDecoder := json.NewDecoder(r.Body)
	payload := requestBody{}
	err := jsonDecoder.Decode(&payload)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	createdFeed, err := app.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		Name:      payload.Name,
		Url:       payload.URL,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create feed")
		return
	}

	createdFollow, err := app.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    createdFeed.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create feed follow")
		return
	}

	respondWithJSON(w, http.StatusOK, struct {
		Feed       Feed
		FeedFollow FeedFollow
	}{
		Feed:       databaseFeedToFeed(createdFeed),
		FeedFollow: databaseFeedFollowToFeedFollow(createdFollow),
	})
}

func (app *appConfig) listFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := app.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get feeds")
		return
	}

	respondWithJSON(w, http.StatusOK, databaseFeedsToFeeds(feeds))
}
