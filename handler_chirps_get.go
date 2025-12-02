package main

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID format", err)
		return
	}
	chirpFromDB, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err == sql.ErrNoRows {
		respondWithError(w, http.StatusNotFound, "Chirp not found", nil)
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirp", err)
		return
	}

	chirp := Chirp{
		ID:        chirpFromDB.ID,
		CreatedAt: chirpFromDB.CreatedAt,
		UpdatedAt: chirpFromDB.UpdatedAt,
		Body:      replaceInvalidWords(chirpFromDB.Body),
		UserID:    chirpFromDB.UserID,
	}

	respondWithJSON(w, http.StatusOK, chirp)
}
