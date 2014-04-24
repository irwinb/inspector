package api

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const (
	ProxyEndpoint = "/rproxy"
	port          = ":8000"
)

type InspectorError struct {
	Error error
	Code  int
}

type ApiHandler func(w http.ResponseWriter, r *http.Request) *InspectorError

func InitAndListen() error {
	log.Println("Initializing HTTP handlers.")

	r := mux.NewRouter()
	initProjectApi(r)
	initProxyApi(r)

	http.Handle("/", r)

	return http.ListenAndServe(port, nil)
}

func (fn ApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := fn(w, r)
	if err == nil || err.Code == 200 {
		return
	}

	switch {
	case err.Code >= 300 && err.Code < 400:
		http.Error(w, getErrorString(err), err.Code)
		return
	case err.Code >= 400 && err.Code < 500:
		http.Error(w, getErrorString(err), err.Code)
		return
	case err.Code >= 500 && err.Code < 600:
		http.Error(w, getErrorString(err), err.Code)
		return
	}

	http.Error(w, err.Error.Error(), err.Code)
}

func getErrorString(e *InspectorError) string {
	if e == nil || e.Error == nil {
		return ""
	} else {
		return e.Error.Error()
	}
}
