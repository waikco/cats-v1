package server

import (
	"net/http"

	json "github.com/json-iterator/go"
)

//respondWithJson wraps a message into json and returns it in the Response,along with a header and Response code
func respondWithJson(w http.ResponseWriter, code int, message interface{}) {
	response, _ := json.Marshal(message)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}

type Response struct {
	Result interface{} `json:"result,omitempty"`
	Error  Error       `json:"error,omitempty"`
}

type Error struct {
	Status  int         `json:"status,omitempty"`
	Message interface{} `json:"message,omitempty"`
}
