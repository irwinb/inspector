package api

import (
	"net/http"
)

func Projects(w http.ResponseWriter, r *http.Request) *InspectorError {
	if r.Method == "GET"
}

func PostProject(w http.ResponseWriter, r *http.Request) *InspectorError {

}

func GetProject(w http.ResponseWriter, r *http.Request) *InspectorError {

	AnonStore.ProjectById(id)
}
