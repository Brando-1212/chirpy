package main

import(
	"encoding/json"
	"net/http"
	"time"

	"chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin (w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "coundn't decode parameters", err)
		return
	}


	dbUser, err := cfg.db.GetUserByEmail(r.Context(), params.Email) 
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't find matching email", err )
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, dbUser.HashedPassword) 
	if err != nil || !match {
		respondWithError(w, http.StatusUnauthorized, "password didn't match", err)
		return
	}

	expirationTime := time.Hour
	if params.ExpiresInSeconds > 0 && params.ExpiresInSeconds < 3600 {
		expirationTime = time.Duration(params.ExpiresInSeconds) * time.Second
	}

	accessToken, err := auth.MakeJWT(
		dbUser.ID,
		cfg.jwtSecret,
		expirationTime,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT", err)
		return
	}


	user := User{
		ID        : dbUser.ID,
		CreatedAt : dbUser.CreatedAt,
		UpdatedAt : dbUser.UpdatedAt,
		Email     : dbUser.Email,
	}
	respondWithJSON(w, http.StatusOK, response{
		User: user,
		Token: accessToken,
	})
}