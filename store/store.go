package store

import (
	"inspector/models"
)

var AnonStore = newMemStore()

var NotFound = NewStoreError("Not found.")
var ProjectInvalid = NewStoreError("Invalid project format.")
var NumProjectsExceeded = NewStoreError("Number of Projects exceeded.  Try again later.")
var NumEndpointsExceeded = NewStoreError("Number of endpoints limit exceeded.")

// ALl save operations hould have the behaviour that if it exists,
// the object at hand should be overwritten.  If not, create it.
type Store interface {
	ProjectById(id uint) *models.Project
	SaveProject(p *models.Project) (*models.Project, error)
	ListProjects() []models.Project
	SaveEndpoint(p *models.Project, ep *models.Endpoint) error
	SaveOperation(p *models.Project, op *models.Operation) error
}

type StoreError struct {
	Reason  string
	Details string
}

func (se *StoreError) Error() string {
	return "Store error: " + se.Reason + ".  " + se.Details
}

func NewStoreError(reason string) *StoreError {
	return &StoreError{Reason: reason}
}
