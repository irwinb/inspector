package models

type Project struct {
	Id             uint   `json:"id"`
	TargetEndpoint string `json:"target_endpoint"`
	Name           string `json:"name"`
}
