package models

import (
	"time"
)

type Project struct {
	Id          uint       `json:"id"`
	Endpoints   []Endpoint `json:"endpoints"`
	Name        string     `json:"name"`
	LastUpdated time.Time  `json:"last_updated"`
}
