package main

import(
	"encoding/json"
	"net/http"
	"strings"

	"chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirps (w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	clean := cleanBody(params.Body, badWords)

	tempChirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body   : clean,
		UserID : params.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}
	chirp := Chirp{
		ID : tempChirp.ID,
		CreatedAt : tempChirp.CreatedAt,
		UpdatedAt : tempChirp.UpdatedAt,
		Body : tempChirp.Body,
		UserID : tempChirp.UserID,
	}
	respondWithJSON(w, http.StatusCreated, chirp)
}



func cleanBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		lowWord := strings.ToLower(word)
		if _, ok := badWords[lowWord]; ok {
			words[i] = "****"
		}
	}
	clean := strings.Join(words, " ")
	return clean
}