package main

import (
	"bytes"
	"github.com/irwinb/inspector/config"
	"github.com/irwinb/inspector/feeder"
	"github.com/irwinb/inspector/models"
	"github.com/irwinb/inspector/store"
	"log"
	"net/http"
	"time"
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

	http.Handle(config.StaticEndpoint, http.FileServer(http.Dir(config.StaticDir)))

	log.Println("Initializing feeder.")
	feeder.InitializeFeeder()

	log.Println("Starting feeder.")
	if err := http.ListenAndServe(config.ServerPort, nil); err != nil {
		log.Println("Staritng feeder failed: ", err)
	}
}

var httpClient = http.DefaultClient

func handleRequest(w http.ResponseWriter, r *http.Request) *InspectorError {
	reqData, err := models.NewRequest(r)
	if err != nil {
		log.Println("Error creating request: ", err)
		return &InspectorError{
			Error: err,
			Code:  400}
	}

	project, err := store.DefaultStore.ProjectByName(reqData.Project)
	if err != nil {
		return &InspectorError{
			Error: err}
	}

	feeder.Feed(reqData)

	req, err := http.NewRequest(reqData.Method, reqData.GetUrl(project),
		bytes.NewReader(reqData.Body))
	req.Header = reqData.Header

	log.Println("Requesting ", req)

	if err != nil {
		return &InspectorError{
			Error: err,
			Code:  500,
		}
	}

	time.Sleep(1000 * time.Millisecond)

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println("Request failed: ", err)
		return &InspectorError{
			Error: err,
			Code:  400}
	}

	newResponse, err := models.NewResponse(project.Name, resp)
	if err != nil {
		return &InspectorError{
			Error: err,
		}
	}

	for key, vals := range newResponse.Header {
		for _, val := range vals {
			w.Header().Add(key, val)
		}
	}
	w.Write(newResponse.Body)
	resp.Trailer.Write(w)
	return nil
}
