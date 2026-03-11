package main

import(
	"encoding/json"
	"net/http"

	"chirpy/internal/auth"
	"chirpy/internal/database"
)


func (cfg *apiConfig) handlerUser (w http.ResponseWriter, r *http.Request) {
	
	type parameters struct {
		Password string `json:"password"`
		Email string    `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not hash password", err)
		return
	}


	tempUser, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email          : params.Email,
		HashedPassword : hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}
	user := User{
		ID        : tempUser.ID,
		CreatedAt : tempUser.CreatedAt,
		UpdatedAt : tempUser.UpdatedAt,
		Email     : tempUser.Email,
	}
	respondWithJSON(w, http.StatusCreated, user)

}