package models

import (
	"errors"
)

type Endpoint struct {
	Id         uint        `json:"id"`
	Target     *string     `json:"target"`
	Name       *string     `json:"name"`
	Operations []Operation `json:"operations"`
}

func (e *Endpoint) Validate() error {
	if e.Target == nil || len(*e.Target) == 0 {
		return errors.New("The endpoint's target cannot be empty.")
	}
	if len(*e.Target) > 100 {
		return errors.New("The endpoint's target is too long.  < 100 plz.")
	}

	if e.Name == nil || len(*e.Name) == 0 {
		return errors.New("The endpoint's name cannot be empty.")
	}
	if len(*e.Name) > 100 {
		return errors.New("The endpoint's name is too long.  < 100 plz.")
	}

	return nil
}
