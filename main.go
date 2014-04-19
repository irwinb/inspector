package main

import (
	"bytes"
	"errors"
	"github.com/irwinb/inspector/config"
	"github.com/irwinb/inspector/feeder"
	"github.com/irwinb/inspector/models"
	"github.com/irwinb/inspector/store"
	"github.com/mreiferson/go-httpclient"
	"log"
	"net/http"
	"strings"
	"sync"
)

type InspectorError struct {
	Error error
	Code  int
}

type inspectorHandler func(w http.ResponseWriter, r *http.Request) *InspectorError

func (fn inspectorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := fn(w, r)
	if err == nil || err.Code == 200 {
		return
	}

	if err.Error == store.NotFound {
		err.Code = 404
	}

	switch {
	case err.Code >= 300 && err.Code < 400:
		http.Error(w, err.Error.Error(), err.Code)
		return
	case err.Code >= 400 && err.Code < 500:
		http.Error(w, err.Error.Error(), err.Code)
		return
	case err.Code >= 500 && err.Code < 600:
		http.Error(w, err.Error.Error(), err.Code)
		return
	}

	http.Error(w, err.Error.Error(), err.Code)
}

func main() {
	log.Println("Starting server.")

	log.Println("Initializing HTTP handlers.")

	http.Handle(config.ProxyEndpoint, http.StripPrefix(config.ProxyEndpoint,
		inspectorHandler(handleRequest)))

	log.Println("Initializing feeder.")
	feeder.InitializeFeeder()

	log.Println("Starting feeder.")
	if err := http.ListenAndServe(config.ServerPort, nil); err != nil {
		log.Println("Staritng feeder failed: ", err)
	}
}

var httpTransport = &httpclient.Transport{
	ResponseHeaderTimeout: config.RequestTimeout,
}
var httpClient = http.Client{Transport: httpTransport}

func createTargetUrl(path string, proj *models.Project) string {
	buff := bytes.NewBufferString("http://")
	buff.WriteString(proj.TargetEndpoint)
	buff.WriteString("/")
	buff.WriteString(path)
	return buff.String()
}

var operationId uint = 0
var idMutex sync.Mutex

func handleRequest(w http.ResponseWriter, r *http.Request) *InspectorError {
	// Url will be /[proxy_endpoint]/[project]/[path]
	requestURI := strings.Trim(r.RequestURI, "/")
	tokens := strings.SplitN(requestURI, "/", 3)
	if len(tokens) < 2 {
		return &InspectorError{
			Error: errors.New("Project not found."),
			Code:  404}
	}

	if len(tokens) < 3 {
		tokens = append(tokens, "")
	}

	reqInbound, err := models.NewRequest(r)
	if err != nil {
		log.Println("Error creating request: ", err)
		return &InspectorError{
			Error: err,
			Code:  400}
	}

	project, err := store.DefaultStore.ProjectByName(tokens[1])
	if err != nil {
		return &InspectorError{
			Error: err}
	}

	idMutex.Lock()
	operationId += 1
	id := operationId
	idMutex.Unlock()

	operation := models.Operation{
		Id:        id,
		ProjectId: project.Id,
		Request:   reqInbound}

	log.Println("Requesting project", project)

	req, err := http.NewRequest(reqInbound.Method,
		createTargetUrl(tokens[2], project),
		bytes.NewReader(reqInbound.Body))
	req.Header = reqInbound.Header

	if err != nil {
		log.Println("Error creating outbound request.", err)
		return &InspectorError{
			Error: err,
			Code:  500}
	}

	feeder.FeedOperation(&operation)

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println("Request failed: ", err)
		return &InspectorError{
			Error: err,
			Code:  400}
	}

	respOutbound, err := models.NewResponse(resp)
	if err != nil {
		return &InspectorError{
			Error: err,
		}
	}

	for key, vals := range respOutbound.Header {
		for _, val := range vals {
			w.Header().Add(key, val)
		}
	}
	w.Write(respOutbound.Body)
	resp.Trailer.Write(w)

	operation.Response = respOutbound
	feeder.FeedOperation(&operation)

	return nil
}
