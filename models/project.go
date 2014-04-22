package models

import (
	"errors"
	"time"
)

type Project struct {
	Id          uint       `json:"id"`
	Endpoint    *Endpoint  `json:"endpoint"`
	Name        *string    `json:"name"`
	LastUpdated *time.Time `json:"last_updated"`
}

func (p *Project) Validate() error {
	if len(*p.Name) == 0 {
		return errors.New("The project's name cannot be empty.")
	}
	if len(*p.Name) > 100 {
		return errors.New("The project's name is too long.")
	}
	return nil
}
