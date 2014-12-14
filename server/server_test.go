package main

import (
	"bytes"
	"encoding/json"
	"github.com/ChrisBoesch/docker-code-verifier/server/verifier"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const (
	solvePayload = `{"solution": "foo=1", "tests": ">>> foo\n1"}`
)

func TestE2eServer(t *testing.T) {
	var (
		dockerCertPath = os.Getenv("DOCKER_CERT_PATH")
		dockerHost     = os.Getenv("DOCKER_HOST")
		resp           verifier.Response
	)

	if dockerHost == "" {
		t.Skip("Docker host is not set.")
	}

	docker, err := getClient(dockerHost, dockerCertPath)
	if err != nil {
		t.Skip("Failed to create a docker client.")
	}

	server := NewServer(docker, 1)
	req, err := http.NewRequest("POST", "http://example.com/python", bytes.NewBuffer([]byte(solvePayload)))
	if err != nil {
		t.Fatalf("Failed to create request: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	server.ServeHTTP(recorder, req)

	if recorder.Code != 200 {
		t.Fatalf("Unexpected code %d", recorder.Code)
	}

	results := recorder.Body.Bytes()
	err = json.Unmarshal(results, &resp)
	if err != nil {
		t.Fatalf("Failed to parse response: %s", results)
	}

	if !resp.Solved {
		t.Fatalf("Unexpect response: %+v", resp)
	}
}

func benchmarkE2eServer(max int, b *testing.B) {
	b.StopTimer()

	var (
		dockerCertPath = os.Getenv("DOCKER_CERT_PATH")
		dockerHost     = os.Getenv("DOCKER_HOST")
	)

	if dockerHost == "" {
		b.Skip("Docker host is not set.")
	}

	docker, err := getClient(dockerHost, dockerCertPath)
	if err != nil {
		b.Skip("Failed to create a docker client.")
	}

	server := NewServer(docker, max)
	results := make(chan []byte, b.N)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		go func() {
			req, err := http.NewRequest("POST", "http://example.com/python", bytes.NewBuffer([]byte(solvePayload)))
			if err != nil {
				b.Fatalf("Failed to create request: %s", err)
			}

			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, req)

			if recorder.Code != 200 {
				b.Fatalf("Unexpected code %d", recorder.Code)
			}

			results <- recorder.Body.Bytes()
		}()
	}

	for i := 0; i < b.N; i++ {
		var (
			resp verifier.Response
		)

		body := <-results
		err = json.Unmarshal(body, &resp)
		if err != nil {
			b.Fatalf("Failed to parse response: %s", body)
		}

		if !resp.Solved {
			b.Fatalf("Unexpect response: %+v", resp)
		}
	}
}

func BenchmarkServer2Job(b *testing.B) { benchmarkE2eServer(2, b) }
func BenchmarkServer3Job(b *testing.B) { benchmarkE2eServer(3, b) }
func BenchmarkServer4Job(b *testing.B) { benchmarkE2eServer(4, b) }
