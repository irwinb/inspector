package models

type ErrorMessage struct {
	TransactionId int    `json:"transaction_id"`
	MessageType   string `json:"message_type"`
	ErrorMessage  string `json:"error_message"`
	ErrorCode     string `json:"error_code"`
}

type RequestMessage struct {
	TransactionId int      `json:"transaction_id"`
	MessageType   string   `json:"message_type"`
	Project       *Project `json:"project"`
	Request       *Request `json:"request"`
}

type ResponseMessage struct {
	TransactionId int       `json:"transaction_id"`
	MessageType   string    `json:"message_type"`
	Project       *Project  `json:"project"`
	Response      *Response `json:"Response"`
}
