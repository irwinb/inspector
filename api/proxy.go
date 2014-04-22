package api

import (
	"bytes"
	"errors"
	"github.com/gorilla/mux"
	"github.com/irwinb/inspector/config"
	"github.com/irwinb/inspector/feeder"
	"github.com/irwinb/inspector/models"
	"github.com/irwinb/inspector/store"
	"github.com/mreiferson/go-httpclient"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var httpTransport = &httpclient.Transport{
	ResponseHeaderTimeout: config.RequestTimeout,
}
var httpClient = http.Client{Transport: httpTransport}

func initProxyApi(r *mux.Router) {
	proxyR := r.Path(ProxyEndpoint).Subrouter()
	proxyR.Handle("/", ApiHandler(handleProxy))
}

func createTargetUrl(path string, ep *models.Endpoint) string {
	buff := bytes.NewBufferString("http://")
	buff.WriteString(ep.Target)
	buff.WriteString("/")
	buff.WriteString(path)
	return buff.String()
}

var operationId uint = 0
var idMutex sync.Mutex

// Url will be /[project_name]/[path]
func handleProxy(w http.ResponseWriter, r *http.Request) *InspectorError {
	requestURL := strings.Trim(r.URL.Path, "/")
	tokens := strings.SplitN(requestURL, "/", 2)
	if len(tokens) < 2 {
		return &InspectorError{
			Error: errors.New("Project not found."),
			Code:  404}
	}

	if len(tokens) < 2 {
		tokens = append(tokens, "")
	}

	reqInbound, err := models.NewRequest(r)
	if err != nil {
		log.Println("Error creating request: ", err)
		return &InspectorError{
			Error: err,
			Code:  400}
	}
	projId, err := strconv.ParseUint(tokens[1], 10, strconv.IntSize)
	if err != nil {
		return &InspectorError{
			Error: err}
	}

	project, err := store.AnonStore.ProjectById(uint(projId))
	if err != nil {
		return &InspectorError{
			Error: err}
	}

	if project.Endpoint == nil {
		return &InspectorError{
			Error: errors.New("Project has no endpoint."),
			Code:  400}
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
		createTargetUrl(tokens[2], project.Endpoint),
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
