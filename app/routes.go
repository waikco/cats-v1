package app

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type health struct {
	Status string `json:"status"`
}

type response struct {
	Message string `json:"message"`
}

func (a *App) Health(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	respondWithJson(w, http.StatusOK, health{Status: "Pets is up and available"})
}

func (a *App) CreateCat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	respondWithJson(w, http.StatusOK, response{Message: "not implemented"})
}

func (a *App) GetCat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	respondWithJson(w, http.StatusOK, response{Message: ps.ByName("id")})
}

func (a *App) GetCats(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	respondWithJson(w, http.StatusOK, response{Message: "not implemented"})
}

func (a *App) UpdateCat(w http.ResponseWriter, r *http.Request) {
	respondWithJson(w, http.StatusOK, response{Message: "not implemented"})
}

func (a *App) DeleteCat(w http.ResponseWriter, r *http.Request) {
	respondWithJson(w, http.StatusOK, response{Message: "not implemented"})
}

func (a *App) MassCreateCat(w http.ResponseWriter, r *http.Request) {
	respondWithJson(w, http.StatusOK, response{Message: "not implemented"})
}
