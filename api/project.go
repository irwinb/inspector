package api

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"inspector/models"
	"inspector/store"
	"log"
	"net/http"
	"strconv"
)

const (
	projectsEndpoint = "/projects"
)

func initProjectApi(r *mux.Router) {
	r.Handle("/projects", ApiHandler(listProjects)).
		Methods("GET")
	r.Handle("/projects/", ApiHandler(listProjects)).
		Methods("GET")

	r.Handle("/projects/{id:[0-9]+}", ApiHandler(getProject)).
		Methods("GET")
	r.Handle("/projects/{id:[0-9]+}/", ApiHandler(getProject)).
		Methods("GET")

	r.Handle("/projects", ApiHandler(postProject)).
		Methods("POST")
	r.Handle("/projects/", ApiHandler(postProject)).
		Methods("POST")
}

func listProjects(w http.ResponseWriter, r *http.Request) *InspectorError {
	if r.Header.Get("pass") != "awesome" {
		return &InspectorError{nil, 404}
	}

	projects := store.AnonStore.ListProjects()
	enc := json.NewEncoder(w)
	enc.Encode(projects)

	return nil
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
	dec := json.NewDecoder(r.Body)
	var proj models.Project

	if err := dec.Decode(&proj); err != nil {
		log.Println(err)
		return &InspectorError{errors.New("Invalid body."), 400}
	}

	newProj, err := store.AnonStore.SaveProject(&proj)
	if err != nil {
		return &InspectorError{err, 400}
	} else {
		enc := json.NewEncoder(w)
		enc.Encode(newProj)
	}

	return nil
}
