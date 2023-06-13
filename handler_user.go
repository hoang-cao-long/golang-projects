package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Kutelata/rssagg/internal/database"
	"github.com/google/uuid"
)

func (apiConfig *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decode := json.NewDecoder(r.Body)

	params := parameters{}

	err := decode.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	user, err := apiConfig.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create user: %v", err))
		return
	}

	respondWithJSON(w, 201, databaseUserToUser(user))
}

func (apiConfig *apiConfig) handleGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, 200, databaseUserToUser(user))
}

func (apiConfig *apiConfig) handleGetPostsForUser(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := apiConfig.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  10,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get post: %v", err))
		return
	}

	respondWithJSON(w, 200, databasePostsToPosts(posts))
}
