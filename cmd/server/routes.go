package server

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	uuid "github.com/satori/go.uuid"

	json "github.com/json-iterator/go"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
)

type health struct {
	Status string `json:"status"`
}

func (a *App) Health(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	respondWithJson(w, http.StatusOK, health{Status: "Pets is up and available"})
}

func (a *App) CreateCat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithJson(w, http.StatusInternalServerError, Response{
			Result: nil,
			Error: Error{
				Status:  http.StatusInternalServerError,
				Message: "error reading body"},
		})
		log.Info().Msgf("error reading body: %v", err)
		return
	}

	if !json.Valid(body) {
		respondWithJson(w, http.StatusBadRequest,
			Response{
				Error: Error{
					Status:  http.StatusBadRequest,
					Message: "invalid json in request body"},
			})
		log.Warn().Msgf("received invalid json in request body: %v", err)
		return
	}

	if id, err := a.Storage.Insert(body); err != nil {
		log.Info().Msgf("error storing cat %s: %v", string(body), err)
		respondWithJson(w, http.StatusInternalServerError, Response{
			Error: Error{
				Status:  http.StatusInternalServerError,
				Message: "error storing cat"},
		})
		return
	} else {
		respondWithJson(w, http.StatusCreated,
			Response{
				Result: id,
			},
		)
		return
	}
}

func (a *App) GetCat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cat, err := a.Storage.Select(ps.ByName("id"))
	switch err {
	case nil:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(cat)
	case sql.ErrNoRows:
		respondWithJson(w, http.StatusNotFound, Response{Result: "cat not found"})
	default:
		log.Debug().Msgf("error getting cat: %v", err)
		respondWithJson(w, http.StatusInternalServerError, Response{Result: "error getting cat"})
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
	case nil:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(all)
	case sql.ErrNoRows:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[]`))
	default:
		respondWithJson(w, http.StatusInternalServerError, Response{
			Error: Error{
				Status:  http.StatusInternalServerError,
				Message: err.Error()},
		})
	}
}

func (a *App) UpdateCat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithJson(w, http.StatusInternalServerError, Response{
			Error: Error{
				Status:  http.StatusInternalServerError,
				Message: "error reading body",
			}})
		log.Info().Msgf("error reading body: %v", err)
		return
	}

	if !json.Valid(body) {
		respondWithJson(w, http.StatusBadRequest,
			Response{
				Error: Error{
					Status:  http.StatusBadRequest,
					Message: "invalid json in request body"},
			})
		log.Warn().Msgf("received invalid json in request body: %v", err)
		return
	}

	id := ps.ByName("id")
	err = a.Storage.Update(id, body)
	switch err {
	case nil:
		respondWithJson(w, http.StatusOK, Response{
			Result: "success",
		})
	case sql.ErrNoRows:
		respondWithJson(w, http.StatusNotFound, Response{
			Error: Error{
				Status:  http.StatusNotFound,
				Message: fmt.Sprintf("cat id %s not found", id)},
		})
	default:
		respondWithJson(w, http.StatusInternalServerError, Response{Result: "error storing object"})

	}
}

func (a *App) DeleteCat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	if uuid.FromStringOrNil(id) == uuid.Nil {
		respondWithJson(w, http.StatusBadRequest, Response{
			Error: Error{
				Status:  http.StatusBadRequest,
				Message: fmt.Sprintf("invalid cat id: %s", id)},
		})
		return
	}
	err := a.Storage.Delete(id)
	switch err {
	case nil:
		respondWithJson(w, http.StatusOK, Response{Result: "success"})
	case sql.ErrNoRows:
		respondWithJson(w, http.StatusNotFound, Response{
			Error: Error{
				Status:  http.StatusNotFound,
				Message: fmt.Sprintf("cat id %s not found", id)},
		})
		log.Debug().Msgf("error getting cat: %v", err)
	default:
		respondWithJson(w, http.StatusInternalServerError, Response{Result: err.Error()})
		log.Debug().Msgf("error getting cat: %v", err)
	}
}

//func (a *App) MassCreateCat(w http.ResponseWriter, r *http.Request) {
//	respondWithJson(w, http.StatusOK, Response{Result: "not implemented"})
//}
