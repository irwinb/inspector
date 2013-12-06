package models

type Transaction struct {
	Id        uint
	ProjectId uint
	Request   *Request
	Response  *Response
}
