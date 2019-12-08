package app

import (
	"encoding/json"
	"net/http"

	log "github.com/rs/zerolog/log"
)

//respondWithJson wraps a message into json and returns it in the response,along with a header and response code
func respondWithJson(w http.ResponseWriter, code int, cargo interface{}) {
	response, err := json.Marshal(cargo)

	if err != nil {
		log.Warn().Msgf("error marshalling json: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
