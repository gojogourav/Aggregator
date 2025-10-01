package main

import (
	"encoding/json"
	"fmt"
	db "github/gojogourav/RSSAggregator/db/sqlc"
	"net/http"
	"time"

	"github.com/go-chi/chi"
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

	FeedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), db.CreateFeedFollowParams{
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

	respondWithJson(w, 201, databaseFeedFollowToFeedFollow(FeedFollow))

}

func (apiCfg *apiConfig) handlerGetFeedFollow(w http.ResponseWriter, r *http.Request, user db.User) {
	feed_follows, err := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get feed follows Error : %v ", err))
		return
	}

	respondWithJson(w, 201, databaseFeedFollowsToFeedFollows(feed_follows))
}

func (apiCfg *apiConfig) handlerDeleteFeedFollows(w http.ResponseWriter, r *http.Request, user db.User) {
	feedFollowIdStr := chi.URLParam(r, "feedFollowId")
	feedFollowID, err := uuid.Parse(feedFollowIdStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't parse feedFollowId"))
		return
	}
	err = apiCfg.DB.DeleteFeedFollows(r.Context(), db.DeleteFeedFollowsParams{
		UserID: user.ID,
		ID:     feedFollowID,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Cloudn't delete feed"))
		return
	}

	respondWithJson(w, 200, struct{}{})
}
