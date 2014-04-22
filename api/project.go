package api

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/irwinb/inspector/store"
	"net/http"
	"strconv"
)

const (
	projectsEndpoint = "/projects"
)

func initProjectApi(r *mux.Router) {
	projectR := r.Path(projectsEndpoint).Subrouter()

	projectR.Handle("/{id:[0-9]+}/", ApiHandler(getProject)).
		Methods("GET")
	projectR.Handle("/", ApiHandler(postProject)).
		Methods("POST")
}

func getProject(w http.ResponseWriter, r *http.Request) *InspectorError {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseUint(idStr, 10, strconv.IntSize)
	if err != nil {
		return &InspectorError{errors.New("Invalid project id."), 404}
	}

	proj, err := store.AnonStore.ProjectById(uint(id))
	if err != nil {
		return &InspectorError{err, 500}
	}
	if proj == nil {
		return &InspectorError{errors.New("Project does not exist."), 404}
	}
	enc := json.NewEncoder(w)
	enc.Encode(proj)

	return nil
}

func postProject(w http.ResponseWriter, r *http.Request) *InspectorError {
	return nil
}
