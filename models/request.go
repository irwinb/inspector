package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Our abstracted represnetation of a request.
type Request struct {
	Method           string      `json:"method"`
	Proto            string      `json:"protocol"`
	Header           http.Header `json:"headers,omitempty"`
	Body             []byte      `json:"body"`
	ContentLength    int64       `json:"content_length"`
	TransferEncoding []string    `json:"transfer_encoding,omitempty"`
	Host             string      `json:"host"`
	RemoteAddr       string      `json:"remote_addr"`
	RequestURI       string      `json:"request_uri"`
	Timestamp        int64       `json:"timestamp"`
}

func (r *Request) String() string {
	result, err := json.Marshal(r)
	if err != nil {
		return fmt.Sprint("Could not String() request: ", *r)
	}
	return string(result)
}

func NewRequest(request *http.Request) (*Request, error) {
	var body bytes.Buffer
	_, err := body.ReadFrom(request.Body)
	if err != nil {
		return nil, err
	}

	req := Request{
		Method:           request.Method,
		Proto:            request.Proto,
		Header:           request.Header,
		Body:             body.Bytes(),
		ContentLength:    request.ContentLength,
		TransferEncoding: request.TransferEncoding,
		Host:             request.Host,
		RemoteAddr:       request.RemoteAddr,
		RequestURI:       request.RequestURI}

	return &req, nil
}
