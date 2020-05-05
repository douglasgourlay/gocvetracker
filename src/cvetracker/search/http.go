package search

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"cvetracker/cve"
	"cvetracker/mongo"

	"go.uber.org/zap"
)

// Server ...
type Server struct {
	mongo  *mongo.Client
	s      *http.Server
	config *Config
}

// NewServer ...
func NewServer(config *Config) (*Server, error) {

	zap.L().Debug("Starting Search Server")

	s := &Server{}

	s.config = config

	var err error

	s.mongo, err = mongo.NewClient(config.Mongo)
	if err != nil {
		return nil, err
	}

	go func() {
		// TODO Port must be part of config
		s.s = &http.Server{Addr: ":8080", Handler: s}
		s.s.ListenAndServe()
	}()
	return s, nil
}

const (
	// Success ...
	Success = 0
	// UnknownFailure ...
	UnknownFailure = 1
	// UnexpectedOperation ...
	UnexpectedOperation = 2
)

// Result ...
type Result struct {
	Status int           `json:"status,omitempty"`
	Error  string        `json:"error,omitempty"`
	Filter cve.DougCVE   `json:"filter,omitempty"`
	CVES   []cve.DougCVE `json:"cves,omitempty"`
}

func (c *Result) String() string {
	pjson, err := json.Marshal(c)
	if err != nil {
		// This is not expected
		return "{\"status\":" + string(UnknownFailure) + ",\"error\":" + err.Error() + "\" + }"
	}
	return string(pjson)
}

// ServeHTTP ...
func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	logConnInfo(nil, req)

	w.Header().Set("Content-Type", "application/json")

	result := &Result{}

	switch req.URL.Path {

	case "/search":
		s.handleSearch(result, w, req)
		break

	default:
		result.Status = UnexpectedOperation
		result.Error = "Endpoint " + req.URL.Path + " not mapped"
		logConnInfo(result, req)
		fmt.Fprintf(w, result.String())
	}

}

func (s *Server) handleSearch(result *Result, w http.ResponseWriter, req *http.Request) {

	switch req.Method {

	case "GET":
		s.handleSearchGET(result, w, req)
		break

	case "POST":
		s.handleSearchPOST(result, w, req)
		break

	default:
		result.Status = UnexpectedOperation
		result.Error = "The method " + req.Method + " is not mapped"
		logConnInfo(result, req)
		fmt.Fprintf(w, result.String())
	}

}

func (s *Server) handleSearchPOST(result *Result, w http.ResponseWriter, req *http.Request) {

	// We expect the post method to have the filter in the body.

	err := json.NewDecoder(req.Body).Decode(&result.Filter)
	if err != nil {
		result.Status = UnknownFailure
		result.Error = fmt.Sprintf("error %s while building filter", err)
		logConnInfo(result, req)
		fmt.Fprintf(w, result.String())
		return
	}

	// Now we hand it off to the GET function. Hence the params can override the filter
	s.handleSearchGET(result, w, req)
}

func (s *Server) handleSearchGET(result *Result, w http.ResponseWriter, req *http.Request) {

	m := getMap(req)
	result.Filter.SetFromMap(m)

	cves, err := s.mongo.Search(&result.Filter)

	if err != nil {
		result.Status = UnknownFailure
		result.Error = fmt.Sprintf("error %s while searching database", err)
		logConnInfo(result, req)
		fmt.Fprintf(w, result.String())
		return
	}

	result.CVES = cves
	result.Status = Success
	result.Error = ""
	fmt.Fprintf(w, result.String())

	logConnInfo(result, req)
}

// Shutdown ...
func (s *Server) Shutdown() {
	zap.L().Debug("Shutting down Search Server")
	s.mongo.Shutdown()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.s.Shutdown(ctx)
}

func getMap(req *http.Request) map[string]string {
	m := make(map[string]string)

	for k, v := range req.URL.Query() {
		//fmt.Printf("%s: %s\n", k, v)
		if len(v) > 0 {
			m[k] = v[0]
		}
	}
	return m
}

func logConnInfo(result *Result, req *http.Request) {

	if result == nil {
		zap.L().Debug(fmt.Sprintf("Request IP:%s, Method:%s(%s) Status: %s", getRemoteIP(req), req.URL.Path, req.Method, "Processing"))
		return
	}

	if result.Status == Success {
		zap.L().Debug(fmt.Sprintf("Request IP:%s, Method:%s(%s) Status: %s", getRemoteIP(req), req.URL.Path, req.Method, "Success"))
		return
	}

	zap.L().Debug(fmt.Sprintf("Request IP:%s, Method:%s(%s) Status: %s Error:%s", getRemoteIP(req), req.URL.Path, req.Method, "Error", result.Error))
}

func getRemoteIP(req *http.Request) string {
	forwarded := req.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return req.RemoteAddr
}
