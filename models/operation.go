package models

import (
	"errors"
)

type Operation struct {
	Id        uint      `json:"id"`
	ProjectId *uint     `json:"project_id"`
	Request   *Request  `json:"request"`
	Response  *Response `json:"response"`
}

func (o *Operation) Validate() error {
	if o.ProjectId == nil {
		return errors.New("An operation's project_id cannot be empty.")
	}

	return nil
}
