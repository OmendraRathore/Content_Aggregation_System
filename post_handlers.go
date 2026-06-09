package main

import (
	"net/http"
	"strconv"

	"github.com/OmendraRathore/Content_Aggregation_System/internal/database"
)

func (app *appConfig) listPosts(w http.ResponseWriter, r *http.Request, user database.User) {
	limitText := r.URL.Query().Get("limit")
	limit := 10
	if requestedLimit, err := strconv.Atoi(limitText); err == nil {
		limit = requestedLimit
	}

	posts, err := app.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get posts for user")
		return
	}

	respondWithJSON(w, http.StatusOK, databasePostsToPosts(posts))
}
