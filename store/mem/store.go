package mem

import (
	"github.com/irwinb/inspector/models"
	"github.com/irwinb/inspector/store"
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
type MemStore struct {
	projectPq *projectMap
	projCount uint
}

func NewMemStore() *MemStore {
	memStore := &MemStore{}
	memStore.projectPq = newProjectMap()

	return memStore
}

func (ms *MemStore) ProjectById(id uint) (*models.Project, error) {
	proj := ms.projectPq.Search(id)
	if proj != nil {
		return proj, nil
	} else {
		return nil, store.NotFound
	}
}

func (ms *MemStore) SaveProject(proj models.Project) error {
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

func (ms *MemStore) CreateProject(proj models.Project) error {
	if err := proj.Validate(); err != nil {
		return err
	}

	proj.Id = ms.projCount
	ms.projectPq.Push(&proj)
	ms.projCount++

	return nil
}

func (ms *MemStore) NewOperation(httpReq *http.Request) *models.Operation {
	return nil
}
