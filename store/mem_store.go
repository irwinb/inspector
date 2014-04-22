package store

import (
	"github.com/irwinb/inspector/models"
	"github.com/irwinb/inspector/store/mem"
	"net/http"
	"time"
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
type memStore struct {
	projectPq *mem.ProjectMap
	projCount uint
}

func newMemStore() *memStore {
	memStore := &memStore{}
	memStore.projectPq = mem.NewProjectMap()

	return memStore
}

func (ms *memStore) ProjectById(id uint) (*models.Project, error) {
	proj := ms.projectPq.Search(id)
	if proj != nil {
		return proj, nil
	} else {
		return nil, NotFound
	}
}

func (ms *memStore) SaveProject(proj models.Project) error {
	projInStore, err := ms.ProjectById(proj.Id)
	if err != nil {
		return err
	}

	projInStore.Endpoint = proj.Endpoint
	projInStore.Name = proj.Name
	now := time.Now()
	projInStore.LastUpdated = &now

	ms.projectPq.Set(projInStore)
	return nil
}

func (ms *memStore) CreateProject(proj models.Project) error {
	if err := proj.Validate(); err != nil {
		return err
	}

	proj.Id = ms.projCount
	ms.projectPq.Push(&proj)
	ms.projCount++

	return nil
}

func (ms *memStore) NewOperation(httpReq *http.Request) *models.Operation {
	return nil
}
