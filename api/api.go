package api

import (
	"github.com/irwinb/inspector/store"
	"github.com/irwinb/inspector/store/mem"
	"net/http"
)

type InspectorError struct {
	Error error
	Code  int
}

var AnonStore = mem.NewMemStore()

type InspectorHandler func(w http.ResponseWriter, r *http.Request) *InspectorError

func (fn InspectorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
