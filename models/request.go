package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/irwinb/inspector/config"
	"log"
	"net/http"
	"strings"
)

type Request struct {
	Project          string      `json:"project"`
	Path             string      `json:"string"`
	Method           string      `json:"method"`
	Proto            string      `json:"protocol"`
	Header           http.Header `json:"header"`
	Body             []byte      `json:"body"`
	ContentLength    int64       `json:"content_length"`
	TransferEncoding []string    `json:"transfer_encoding"`
	Host             string      `json:"host"`
	RemoteAddr       string      `json:"remote_addr"`
	RequestURI       string      `json:"request_uri"`
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

	withoutProxy := request.RequestURI[len(config.ProxyEndpoint):]
	tokens := strings.SplitN(withoutProxy, "/", 2)
	if len(tokens) < 2 {
		return nil, errors.New("Invalid url.  No project name found.")
	}
	log.Println("Request tokens: ", tokens)

	req := Request{
		Project:          tokens[0],
		Path:             tokens[1],
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

func (r *Request) GetUrl(proj *Project) string {
	// Remove userinfo@host/project to create a new URL
	buff := bytes.NewBufferString("http://")
	buff.WriteString(proj.TargetEndpoint)
	buff.WriteString("/")
	buff.WriteString(r.Path)
	return buff.String()
}
