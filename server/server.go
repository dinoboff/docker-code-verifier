/*
Server process http requests to run code solution against provided tests

It currently only support python3 solution run against doctest tests.

*/
package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/ChrisBoesch/docker-code-verifier/server/verifier"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	version = "0.1.0-dev"
	root    = "/"

	defaultAddr       = ":5000"
	defaultDockerHost = "unix:///var/run/docker.sock"
	defaultMaxJobs    = 5

	jsonError = `{"errors": "unexpected error"}`
)

var (
	showHelp       bool
	showVersion    bool
	addr           string
	dockerCertPath string
	dockerHost     string
	maxJobs        int

	ErrMethodNotSupported = errors.New("Method not supported.")
	ErrPayloadMissing     = errors.New("The request doesn't have any payload to test.")
)

func init() {
	var (
		flagDefaultDockerCertPath = os.Getenv("DOCKER_CERT_PATH")
		flagDefaultDockerHost     = os.Getenv("DOCKER_HOST")
	)

	if flagDefaultDockerHost == "" {
		flagDefaultDockerHost = defaultDockerHost
	}

	flag.BoolVar(&showHelp, "help", false, "Show help")
	flag.BoolVar(&showVersion, "version", false, "Show version")
	flag.IntVar(&maxJobs, "max-job", defaultMaxJobs, "maximun number of concurrent job running")
	flag.StringVar(&addr, "http", defaultAddr, "Address to bind the server to.")
	flag.StringVar(&dockerHost, "docker-host", flagDefaultDockerHost, "URI of the Docker remote API")
	flag.StringVar(
		&dockerCertPath,
		"docker-cert-dir",
		flagDefaultDockerCertPath,
		"Path to the X509 Key Pair to authenticate with the Docker remote API",
	)
}

func main() {
	flag.Parse()

	if showVersion {
		fmt.Println(version)
		return
	}

	if showHelp {
		flag.PrintDefaults()
		return
	}

	if dockerHost == "" {
		fmt.Println("Docker host wasn't provider. You should either set DOCKER_HOST or use the 'docker-host' flag")
		flag.PrintDefaults()
		os.Exit(1)
	}

	log.Println("Starting server...")
	log.Printf("Docker address: %s", dockerHost)
	log.Printf("Docker cert. path: %s", dockerCertPath)

	docker, err := verifier.NewClient(dockerHost, dockerCertPath)
	if err != nil {
		log.Fatal(err)
	}

	server := NewServer(docker, maxJobs)
	http.Handle(root, server)
	log.Printf("Binding server to: %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func urlb64Decode(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

type processError struct {
	err  error
	Msg  string `json:"errors"`
	code int
}

// Web app Helper to return an error (for system use), a user friendly message
// and the http error code
func ProcessError(err error, msg string, code int) *processError {
	return &processError{err, msg, code}
}

// The schema of the index return value.
//
// Should hold all the supported runtimes.
type Index struct {
	Runtimes []string `json:"runtimes"`
}

// The main server handler.
type Server struct {
	Docker verifier.Client

	// increment when a job start;
	// decrement when a job stop;
	// block starting jobs if full.
	jobs chan struct{}
}

// Create a new server.
//
// `max` limit the number of concurent jobs.
func NewServer(docker verifier.Client, max int) *Server {
	if max < 1 {
		max = 1
	}

	return &Server{docker, make(chan struct{}, max)}
}

// Main server handler. Run user provided code and return the results.
//
// The code to run is sent via a json encoded object with a `solution`
// and a `tests` fields.
//
// The server will look for the payload in:
// - the first `jsonrequest` variable query string of the URL (for GET request).
//   json can eith be plain or URL safe base64 encoded.
// - the body of POST request with a `application/json` content type.
// - the first `jsonrequest` variable of the x-www-form-urlencoded post request
//
// The runtime name use to run the code is extracted from the URL:
//
// - /python -> will be run by the default python runtime
// - /python3 -> will be run by the default python3 runtine
//
// The response will either be json or jsonp encoded (for GET request with a
// `vcallback` variable in the URL query string).
//
// If no runtime is given, a list of supported runtime will be returned.
//
func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.logAccess(req)

	s.startJob()
	defer s.stopJob()

	resp, pErr := s.proccessRequest(req)
	if pErr != nil {
		log.Printf("Error Procesing Request: %v", pErr.err)
		s.processResponse(w, req, pErr, pErr.code)
		return
	}

	s.processResponse(w, req, resp, 200)
}

func (s *Server) logAccess(req *http.Request) {
	log.Printf("%s %s %s", req.RemoteAddr, req.Method, req.URL)
}

// fill up the server job channel. Should wait if the channel buffer is full
func (s *Server) startJob() {
	s.jobs <- struct{}{}
}

// free up the channel buffer. Should let an other channel start if it was full.
func (s *Server) stopJob() {
	<-s.jobs
}

// TODO: handle some errors as internal errors
func (s *Server) proccessRequest(req *http.Request) (interface{}, *processError) {
	var (
		runtimeName = req.URL.Path[len(root):]
	)

	if runtimeName == "" {
		return s.getIndex(), nil
	}

	err := verifier.SupportedRuntime(runtimeName)
	if err != nil {
		return nil, ProcessError(err, "Unsupported runtimes", http.StatusNotFound)
	}

	body, pErr := s.getPayload(req)
	if pErr != nil {
		return nil, pErr
	}

	payload, err := verifier.NewRequest(runtimeName, body)
	if err == verifier.ErrSolutionIsMissing {
		return nil, ProcessError(verifier.ErrSolutionIsMissing, "A solution is required", http.StatusBadRequest)
	}
	if err != nil {
		return nil, ProcessError(err, "Could not parse the json request", http.StatusBadRequest)
	}

	resp, err := verifier.Run(s.Docker, payload, false)
	if err != nil {
		return nil, ProcessError(err, "Failed to run code", http.StatusBadRequest)
	}

	return resp, nil
}

func (s *Server) getIndex() *Index {
	return &Index{verifier.SupportedtRuntimeList()}
}

func (s *Server) getPayload(req *http.Request) ([]byte, *processError) {
	switch req.Method {
	case "GET":
		return s.getPayloadFromQuery(req)
	case "POST":
		return s.getPayloadFromBody(req)
	default:
		return nil, ProcessError(
			ErrMethodNotSupported,
			"Unsupported method",
			http.StatusMethodNotAllowed,
		)
	}
}

func (s *Server) getPayloadFromQuery(req *http.Request) ([]byte, *processError) {
	jsonRequest := req.URL.Query().Get("jsonrequest")
	if jsonRequest == "" {
		return nil, ProcessError(
			ErrPayloadMissing,
			"jsonrequest missing from the query string",
			http.StatusBadRequest,
		)
	}

	jsonRequest = strings.TrimSpace(jsonRequest)
	if jsonRequest[0:1] == "{" {
		return []byte(jsonRequest), nil
	}

	payload, err := urlb64Decode(jsonRequest)
	if err != nil {
		return nil, ProcessError(
			fmt.Errorf("Failed to parse payload: %s", jsonRequest),
			"Failed to parse jsonrequest query string (expecting it plain, or alternate base64 encoded)",
			http.StatusBadRequest,
		)
	}
	return payload, nil
}

func (s *Server) getPayloadFromBody(req *http.Request) ([]byte, *processError) {
	var (
		payload []byte
		err     error
	)

	switch req.Header.Get("Content-Type") {
	case "application/json":
		payload, err = ioutil.ReadAll(req.Body)
		req.Body.Close()
		if err != nil {
			return nil, ProcessError(err, "Failed to parse payload POST body", http.StatusBadRequest)
		}
	default:
		payload = []byte(req.FormValue("jsonrequest"))
	}

	if len(payload) == 0 {
		return nil, ProcessError(ErrPayloadMissing, "jsonrequest is missing", http.StatusBadRequest)
	}

	return payload, nil
}

func (s *Server) processResponse(w http.ResponseWriter, req *http.Request, resp interface{}, code int) {
	output, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error Procesing Response (%+v): %v", resp, err)
		code = http.StatusInternalServerError
		output = []byte(jsonError)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	if cb, m := req.URL.Query().Get("vcallback"), req.Method; m == "GET" && cb != "" {
		w.Header().Set("Content-Type", "application/javascript")
		w.WriteHeader(code)
		fmt.Fprintf(w, "%s(%s)", cb, output)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(output)
	}
}
