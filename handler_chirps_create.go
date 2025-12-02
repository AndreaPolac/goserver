package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/AndreaPolac/goserver/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	type response struct {
		Chirp
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

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   params.Body,
		UserID: params.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		Chirp: Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      replaceInvalidWords(chirp.Body),
			UserID:    chirp.UserID,
		},
	})
}

func replaceInvalidWords(body string) string {
	invalidWords := []string{"kerfuffle", "sharbert", "fornax"}
	replacement := "****"
	for _, word := range invalidWords {
		body = replaceIgnoreCase(body, word, replacement)
	}
	return body
}

func replaceIgnoreCase(s, old, new string) string {
	lowerS := strings.ToLower(s)
	lowerOld := strings.ToLower(old)

	var result strings.Builder
	start := 0
	for {
		index := strings.Index(lowerS[start:], lowerOld)
		if index == -1 {
			result.WriteString(s[start:])
			break
		}
		result.WriteString(s[start : start+index])
		result.WriteString(new)
		start += index + len(old)
	}
	return result.String()
}
