package models

import (
	"bytes"
	"net/http"
)

type Response struct {
	Project          string      `json:"prjoect"`
	Proto            string      `json:"protocol"`
	Header           http.Header `json:"header"`
	Trailer          http.Header `json:"trailer"`
	Body             []byte      `json:"body"`
	ContentLength    int64       `json:"content_length"`
	TransferEncoding []string    `json:"transfer_encoding"`
	Status           string      `json:"status"`
	StatusCode       int         `json:"status_code"`
}

func NewResponse(proj string, resp *http.Response) (*Response, error) {
	var body bytes.Buffer
	_, err := body.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}

	newResp := Response{
		Project:          proj,
		Proto:            resp.Proto,
		Header:           resp.Header,
		Trailer:          resp.Trailer,
		Body:             body.Bytes(),
		ContentLength:    resp.ContentLength,
		TransferEncoding: resp.TransferEncoding,
		Status:           resp.Status,
		StatusCode:       resp.StatusCode}

	return &newResp, nil
}
