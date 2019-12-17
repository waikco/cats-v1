package app

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"

	"github.com/waikco/cats-v1/model"

	"github.com/julienschmidt/httprouter"
)

type health struct {
	Status string `json:"status"`
}

type catResponse struct {
	Result string `json:"result, omitempty"`
	ID     string `json:"id, omitempty"`
}

func (a *App) Health(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	respondWithJson(w, http.StatusOK, health{Status: "Pets is up and available"})
}

func (a *App) CreateCat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var cat model.Cat
	if err := json.NewDecoder(r.Body).Decode(&cat); err != nil {
		respondWithJson(w, http.StatusInternalServerError, catResponse{Result: "error marshalling json"})
		return
	}

	if id, err := a.Storage.Insert(cat); err != nil {
		log.Info().Msgf("error storing cat %+v: %v", cat, err)
		respondWithJson(w, http.StatusInternalServerError, catResponse{Result: "error storing object"})
		return
	} else {
		respondWithJson(w, http.StatusCreated, catResponse{ID: id})
		return
	}
}

func (a *App) GetCat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cat, err := a.Storage.Select(ps.ByName("id"))
	switch err {
	case sql.ErrNoRows:
		respondWithJson(w, http.StatusNotFound, catResponse{Result: "cat not found"})
	case nil:
		respondWithJson(w, http.StatusOK, cat)
	default:
		log.Debug().Msgf("error getting cat: %v", err)
		respondWithJson(w, http.StatusInternalServerError, catResponse{Result: err.Error()})
	}
}

func (a *App) GetCats(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}

	if start < 0 {
		start = 0
	}

	all, err := a.Storage.SelectAll(count, start)
	switch err {
	case sql.ErrNoRows:
		respondWithJson(w, http.StatusOK, `[]`)
	case nil:
		respondWithJson(w, http.StatusOK, all)
	default:
		respondWithJson(w, http.StatusInternalServerError, catResponse{Result: err.Error()})
	}
}

func (a *App) UpdateCat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var cat model.Cat
	if err := json.NewDecoder(r.Body).Decode(&cat); err != nil {
		respondWithJson(w, http.StatusInternalServerError, catResponse{Result: "error marshalling json"})
		return
	}

	if err := a.Storage.Update(ps.ByName("id"), cat); err != nil {
		respondWithJson(w, http.StatusInternalServerError, catResponse{Result: "error storing object"})
		return
	} else {
		respondWithJson(w, http.StatusOK, catResponse{
			Result: "success",
		})
		return
	}
}

func (a *App) DeleteCat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if err := a.Storage.Delete(ps.ByName("id")); err != nil {
		respondWithJson(w, http.StatusInternalServerError, catResponse{Result: err.Error()})
	} else {
		respondWithJson(w, http.StatusOK, catResponse{Result: "success"})
	}
}

func (a *App) MassCreateCat(w http.ResponseWriter, r *http.Request) {
	respondWithJson(w, http.StatusOK, catResponse{Result: "not implemented"})
}
