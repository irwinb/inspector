package models

type ErrorMessage struct {
	Id           int    `json:"operation_id"`
	MessageType  string `json:"message_type"`
	ErrorMessage string `json:"error_message"`
	ErrorCode    string `json:"error_code"`
}

type RequestMessage struct {
	Id          int      `json:"operation_id"`
	MessageType string   `json:"message_type"`
	Project     *Project `json:"project"`
	Request     *Request `json:"request"`
}

type ResponseMessage struct {
	Id          int       `json:"operation_id"`
	MessageType string    `json:"message_type"`
	Project     *Project  `json:"project"`
	Response    *Response `json:"Response"`
}

type Operation struct {
	Id        uint      `json:"operation_id"`
	ProjectId uint      `json:"project"`
	Request   *Request  `json:"request"`
	Response  *Response `json:"response"`
}
