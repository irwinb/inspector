package store

import (
	"github.com/irwinb/inspector/models"
)

var DefaultStore = NewDefaultStore()

var NotFound = NewStoreError("Not found.")
var ProjectInvalid = NewStoreError("Invalid project format.")
var ProjectNameExists = NewStoreError("Invalid project format.")

type Store interface {
	ProjectById(id uint) *models.Project
	ProjectByName(name string) *models.Project
	SaveProject(proj *models.Project) error
	NewTransaction(proj *models.Project, req *models.Request) *models.Transaction
}

type StoreError struct {
	Reason  string
	Details string
}

func NewDefaultStore() *MemStore {
	memStore := &MemStore{}
	memStore.projectsByName = make(map[string]models.Project)
	memStore.projectsById = make(map[uint]models.Project)

	memStore.CreateProject(models.Project{
		Name:           "google",
		TargetEndpoint: "google.ca"})
	memStore.CreateProject(models.Project{
		Name:           "facebook_graph",
		TargetEndpoint: "graph.facebook.com"})
	return memStore
}

func (se *StoreError) Error() string {
	return "Store error: " + se.Reason + ".  " + se.Details
}

func NewStoreError(reason string) *StoreError {
	return &StoreError{Reason: reason}
}
