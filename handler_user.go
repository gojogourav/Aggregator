package main

import (
	"encoding/json"
	"fmt"
	db "github/gojogourav/RSSAggregator/db/sqlc"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlreCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON : %v", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(
		r.Context(), db.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name:      params.Name,
		},
	)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Failed to create new User : %v", err))
		return
	}

	respondWithJson(w, 200, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user db.User) {

	respondWithJson(w, 200, databaseUserToUser(user))

}
