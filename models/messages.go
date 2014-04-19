package models

type ErrorMessage struct {
	Id           int    `json:"id"`
	MessageType  string `json:"message_type"`
	ErrorMessage string `json:"error_message"`
	ErrorCode    string `json:"error_code"`
}
