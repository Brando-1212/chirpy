package main

import (
	"encoding/json"
	"net/http"
	"strings"

)


func handlerValidateChirp (w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanBody string `json:"cleaned_body"`
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

	respondWithJSON(w, http.StatusOK, returnVals{
		CleanBody: clean,
	})

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