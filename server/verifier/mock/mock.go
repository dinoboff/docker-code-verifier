/*
Usulful mocks for github.com/ChrisBoesch/docker-code-verifier/server/verifier testing
*/

package mock

import (
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
)

type MockDockerServer struct {
	Resp []byte
	Code int

	URLs    []*url.URL
	Methods []string
	Bodies  [][]byte

	TS   *httptest.Server
	Host string
}

func NewMockDockerServer(resp []byte, code int) (*MockDockerServer, error) {
	s := &MockDockerServer{
		Resp:    resp,
		Code:    code,
		URLs:    make([]*url.URL, 0, 1),
		Methods: make([]string, 0, 1),
		Bodies:  make([][]byte, 0, 1),
	}
	s.TS = httptest.NewServer(http.Handler(s))

	u, err := url.Parse(s.TS.URL)
	if err != nil {
		s.Close()
		return nil, err
	}

	u.Scheme = "tcp"
	s.Host = u.String()

	return s, nil
}

func (m *MockDockerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("mock.NewMockDockerServer.ServeHTTP: failed to read body: %s", err)
	}
	r.Body.Close()

	m.URLs = append(m.URLs, r.URL)
	m.Methods = append(m.Methods, r.Method)
	m.Bodies = append(m.Bodies, body)

	w.WriteHeader(m.Code)
	if m.Code != 204 {
		w.Write(m.Resp)
	}
}

func (m *MockDockerServer) Close() {
	m.TS.CloseClientConnections()
	m.TS.Close()
	m.TS = nil
}

type Header struct {
	Type uint8
	_    uint8
	_    uint8
	_    uint8
	Size uint32
}

type LogResponse struct {
	io.Reader
}

// See https://docs.docker.com/reference/api/docker_remote_api_v1.15/#attach-to-a-container
func NewResponse(stdout []byte) (io.ReadCloser, error) {
	var (
		header = Header{
			Type: 1,
			Size: uint32(len(stdout)),
		}
		resp = new(bytes.Buffer)
	)

	err := binary.Write(resp, binary.BigEndian, header)
	if err != nil {
		return nil, err
	}

	_, err = resp.Write(stdout)
	if err != nil {
		return nil, err
	}

	return &LogResponse{resp}, nil
}

func (r *LogResponse) Close() error {
	return nil
}
