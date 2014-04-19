package models

type Operation struct {
	Id        uint      `json:"id"`
	ProjectId uint      `json:"project_id"`
	Request   *Request  `json:"request"`
	Response  *Response `json:"response"`
}
