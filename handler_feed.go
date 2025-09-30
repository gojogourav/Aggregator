package main

import (
	"encoding/json"
	"fmt"
	db "github/gojogourav/RSSAggregator/db/sqlc"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user db.User) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), db.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create feed: %v", err))
		return
	}

	respondWithJson(w, 201, databaseFeedToFeed(feed))
}

func (apiConfig *apiConfig) handlerGetFeed(w http.ResponseWriter, r *http.Request) {
	feed, err := apiConfig.DB.GetFeeds(r.Context())

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get feeds : %v", err))
		return
	}

	respondWithJson(w, 201, databaseFeedsToFeeds(feed))

}
