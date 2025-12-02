package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerChirpsGetAll(w http.ResponseWriter, r *http.Request) {
	// type response struct {
	// 	Chirps []Chirp
	// }

	chirpsFromDB, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
		return
	}

	chirps := make([]Chirp, 0, len(chirpsFromDB))
	for _, chirpDB := range chirpsFromDB {
		chirps = append(chirps, Chirp{
			ID:        chirpDB.ID,
			CreatedAt: chirpDB.CreatedAt,
			UpdatedAt: chirpDB.UpdatedAt,
			Body:      replaceInvalidWords(chirpDB.Body),
			UserID:    chirpDB.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}
