package store

import (
	"github.com/irwinb/inspector/models"
	"net/http"
)

type MemStore struct {
	projectsByName   map[string]models.Project
	projectsById     map[uint]models.Project
	transactionsById map[uint]models.Transaction
	projCount        uint
}

func (ms *MemStore) ProjectById(id uint) (*models.Project, error) {
	proj, ok := ms.projectsById[id]
	if ok {
		return &proj, nil
	} else {
		return nil, NotFound
	}
}

func (ms *MemStore) ProjectByName(name string) (*models.Project, error) {
	proj, ok := ms.projectsByName[name]
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

	projOfSameName, err := ms.ProjectByName(proj.Name)
	if err == nil && projOfSameName.Id != proj.Id {
		return ProjectNameExists
	}

	projInStore.TargetEndpoint = proj.TargetEndpoint
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

func (ms *MemStore) NewTransaction(httpReq *http.Request) *models.Transaction {
	models.NewRequest(httpReq)
	return nil
}
