package main

import (
	"net/http"
	
)


func (cfg *apiConfig) handlerReset (w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")

	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Can't delete users outside of dev mode", nil)
		return
	}
	cfg.fileserverHits.Store(0)
	err := cfg.db.DeleteUsers(r.Context())
	if err != nil {
		respondWithError(w,http.StatusInternalServerError, "Couldn't delete users", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("file server hits reset to 0"))
	
}