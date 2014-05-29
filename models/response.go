package models

import (
	"bytes"
	"net/http"
)

type Response struct {
	Proto            string      `json:"protocol"`
	Header           http.Header `json:"headers"`
	Trailer          http.Header `json:"trailer"`
	Body             string      `json:"body"`
	ContentLength    int64       `json:"content_length"`
	TransferEncoding []string    `json:"transfer_encoding"`
	Status           string      `json:"status"`
	StatusCode       int         `json:"status_code"`
	Timestamp        int64       `json:"timestamp"`
}

func NewResponse(resp *http.Response) (*Response, error) {
	var body bytes.Buffer
	_, err := body.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}

	newResp := Response{
		Proto:            resp.Proto,
		Header:           resp.Header,
		Trailer:          resp.Trailer,
		Body:             body.String(),
		ContentLength:    resp.ContentLength,
		TransferEncoding: resp.TransferEncoding,
		Status:           resp.Status,
		StatusCode:       resp.StatusCode}
	return &newResp, nil
}
