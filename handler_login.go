package main

import(
	"encoding/json"
	"net/http"

	"chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin (w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Password string `json:"password"`
		Email string    `json:"email"`
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
	user := User{
		ID        : dbUser.ID,
		CreatedAt : dbUser.CreatedAt,
		UpdatedAt : dbUser.UpdatedAt,
		Email     : dbUser.Email,
	}
	respondWithJSON(w, http.StatusOK, user)
}