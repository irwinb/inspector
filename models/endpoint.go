package models

type Endpoint struct {
	Id         uint        `json:"id"`
	Target     string      `json:"target"`
	Name       string      `json:"name"`
	Operations []Operation `json:"operations"`
}
