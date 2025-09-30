package main

import (
	"encoding/json"
	"fmt"
	db "github/gojogourav/RSSAggregator/db/sqlc"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateFeedFollows(w http.ResponseWriter, r *http.Request, user db.User) {
	type parameters struct {
		FeedId uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON : %v", err))
		return
	}

	feed, err := apiCfg.DB.CreateFeedFollow(r.Context(), db.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedId,
	},
	)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create feed follow : %v", err))
		return
	}

	respondWithJson(w, 201, databaseFeedFollowToFeedFollow(feed))

}
