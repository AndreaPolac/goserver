package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

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

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
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

	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: replaceInvalidWords(params.Body),
	})
}
