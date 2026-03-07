package main

import(

	"encoding/json"
	"net/http"
)


func (cfg *apiConfig) handlerUser (w http.ResponseWriter, r *http.Request)  {
	
	type parameters struct {
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	
	tempUser, err := cfg.db.CreateUser(r.Context(),params.Email)
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