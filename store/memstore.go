package store

import (
	"github.com/irwinb/inspector/models"
	"net/http"
)

const (
	MaxProjects   = 10000
	MaxOperations = 10
	MaxEndpoints  = 1
	ProjectMaxAge = 24 * 60 * 60
)

// Where N is the number of projects,
// Running time requirements:
// ProjectById      O(1)
// SaveProject      O(logN)
// SaveEndpoint     O(logN)
// SaveOperation    O(logN)
//
// Memory requirements:
// Projects   O(2N)
// Operations O(10*N)
// Endpoints  O(N)
type MemStore struct {
	projectsById   map[uint]models.Project
	endpointsById  map[uint]*models.Endpoint
	operationsById map[uint]*models.Operation
	projCount      uint
}

func NewMemStore() *MemStore {
	memStore := &MemStore{}
	memStore.projectsByName = make(map[string]models.Project)

	return memStore
}

func (ms *MemStore) ProjectById(id uint) (*models.Project, error) {
	proj, ok := ms.projectsById[id]
	if ok {
		return &proj, nil
	} else {
		return nil, NotFound
	}
}

func (ms *MemStore) SaveProject(proj models.Project) error {
	projInStore, err := ms.ProjectById(proj.Id)
	if err != nil {
		return err
	}

	// Validate project name.
	projOfSameName, err := ms.ProjectByName(proj.Name)
	if err == nil && projOfSameName.Id != proj.Id {
		return ProjectNameExists
	}

	projInStore.Endpoints = proj.Endpoints
	projInStore.Name = proj.Name
	projInStore.Name = proj.Name
	return nil
}

func (ms *MemStore) CreateProject(proj models.Project) error {
	if proj.TargetEndpoint == "" || proj.Name == "" {
		return ProjectInvalid
	}
	_, err := ms.ProjectByName(proj.Name)
	if err == nil {
		return ProjectNameExists
	}

	projInStore := &models.Project{
		Id:             ms.projCount,
		TargetEndpoint: proj.TargetEndpoint,
		Name:           proj.Name}

	ms.projCount++

	ms.projectsById[projInStore.Id] = *projInStore
	ms.projectsByName[projInStore.Name] = *projInStore

	return nil
}

func (ms *MemStore) NewOperation(httpReq *http.Request) *models.Operation {
	return nil
}
