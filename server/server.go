/*
Server process http requests to run code solution against provided tests

It currently only support python3 solution run against doctest tests.

Supported requests:

- POST /<runtime-name>
  The payload should be json formatted. It should be an object with "solution"
  and "tests" entries.

e.g. for:

    POST /python HTTP/1.1
	Host: localhost:5000
	Cache-Control: no-cache

	{
		"solution": "foo=1\nprint(foo)",
		"tests": ">>> foo\n1"
	}

The response should be:

	{
	    "Solved": true,
	    "Printed": "1\n",
	    "Errors": "",
	    "Results": [
	        {
	            "Call": "foo",
	            "Expected": "1",
	            "Received": "1",
	            "Correct": true
	        }
	    ]
	}

*/
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/ChrisBoesch/docker-code-verifier/server/verifier"
	"github.com/samalba/dockerclient"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	root           = "/"
	defaultAddr    = ":5000"
	defaultMaxJobs = 5
	jsonError      = `{"errors": "unexpected error"}`
)

var (
	showHelp       bool
	addr           string
	dockerHost     string
	dockerCertPath string
	maxJobs        int

	ErrPayloadMissing = errors.New("The request doesn't have any body to test.")
)

func init() {
	var (
		defaultDockerHost     = os.Getenv("DOCKER_HOST")
		defaultDockerCertPath = os.Getenv("DOCKER_CERT_PATH")
	)

	flag.BoolVar(&showHelp, "help", false, "Show help")
	flag.IntVar(&maxJobs, "max-job", defaultMaxJobs, "maximun number of concurrent job running")
	flag.StringVar(&addr, "http", defaultAddr, "Address to bind the server to.")
	flag.StringVar(&dockerHost, "docker-host", defaultDockerHost, "URI of the Docker remote API")
	flag.StringVar(&dockerCertPath, "docker-cert-dir", defaultDockerCertPath, "Path to the X509 Key Pair to authenticate with the Docker remote API")
}

type processError struct {
	err  error
	Msg  string `json:"errors"`
	code int
}

func ProcessError(err error, msg string, code int) *processError {
	return &processError{err, msg, code}
}

func (p *processError) toJson() []byte {
	j, err := json.Marshal(p)
	if err != nil {
		return []byte(jsonError)
	}
	return j
}

type Server struct {
	Docker  dockerclient.Client
	Watcher verifier.StopWatcher

	// increment when a job start;
	// decrement when a job stop;
	// block job starting if full.
	jobs chan struct{}
}

func NewServer(docker dockerclient.Client) *Server {
	var (
		watcher = verifier.NewWatcher()
		jobs    = make(chan struct{}, maxJobs)
	)

	watcher.Start(docker)
	return &Server{docker, watcher, jobs}
}

// fill up the server job channel. Should wait if the channel buffer is full
func (s *Server) startJob() {
	s.jobs <- struct{}{}
}

// free up the channel buffer. Should let an other channel start if it was full.
func (s *Server) stopJob() {
	<-s.jobs
}

func (s *Server) logAccess(req *http.Request) {
	log.Printf("%s %s %s", req.RemoteAddr, req.Method, req.URL)
}

func (s *Server) proccessRequest(req *http.Request) (*verifier.Response, *processError) {
	var (
		runtimeName = req.URL.Path[len(root):]
	)

	err := verifier.SupportedRuntime(runtimeName)
	if err != nil {
		return nil, ProcessError(err, "Unsupported runtimes", http.StatusNotFound)
	}

	body, err := ioutil.ReadAll(req.Body)
	req.Body.Close()
	if err != nil || len(body) == 0 {
		return nil, ProcessError(ErrPayloadMissing, "The solution and tests are missing", http.StatusBadRequest)
	}

	payload, err := verifier.NewRequest(runtimeName, body)
	if err != nil {
		return nil, ProcessError(err, "Could not parse request", http.StatusBadRequest)
	}

	if payload.Solution == "" {
		return nil, ProcessError(verifier.ErrSolutionIsMissing, "Could not parse request", http.StatusBadRequest)
	}

	// Run request
	resp, err := verifier.Run(s.Docker, s.Watcher, payload, false)
	if err != nil {
		return nil, ProcessError(err, "Failed to run code", http.StatusBadRequest)
	}

	return resp, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.logAccess(req)

	s.startJob()
	defer s.stopJob()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	resp, pErr := s.proccessRequest(req)
	if pErr != nil {
		log.Printf("Error Procesing Request: %v", pErr.err)
		w.WriteHeader(pErr.code)
		w.Write(pErr.toJson())
		return
	}

	output, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error Procesing Response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(pErr.toJson())
		return
	}

	w.Write(output)
}

func main() {
	flag.Parse()
	if dockerHost == "" {
		fmt.Println("Docker host wasn't provider. You should either set DOCKER_HOST or use the 'docker-host' flag")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if showHelp {
		flag.PrintDefaults()
		return
	}

	log.Println("Starting server...")
	log.Printf("Docker address: %s", dockerHost)
	log.Printf("Docker cert. path: %s", dockerCertPath)

	docker, err := getClient(dockerHost, dockerCertPath)
	if err != nil {
		log.Fatal(err)
	}

	server := NewServer(docker)
	log.Printf("Binding server to: %s", addr)
	log.Fatal(http.ListenAndServe(addr, server))
}
