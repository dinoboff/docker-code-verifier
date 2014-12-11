/*
Usulful mocks for github.com/ChrisBoesch/docker-code-verifier/server/verifier testing
*/

package mock

import (
	"bytes"
	"encoding/binary"
	"github.com/samalba/dockerclient"
	"io"
	"sync"
)

type MonitorClient struct {
	sync.Mutex
	CBs  []dockerclient.Callback
	Args [][]interface{}
}

func (m *MonitorClient) StartMonitorEvents(cb dockerclient.Callback, args ...interface{}) {
	m.Lock()
	defer m.Unlock()

	m.CBs = append(m.CBs, cb)
	m.Args = append(m.Args, args)
}

// Mock for the container stdout stream
type Response struct {
	io.Reader
}

// See https://docs.docker.com/reference/api/docker_remote_api_v1.15/#attach-to-a-container
func NewResponse(stdout []byte) *Response {
	header := []byte{1, 0, 0, 0}
	size := make([]byte, 4)
	binary.BigEndian.PutUint32(size, uint32(len(stdout)))
	header = append(header, size...)
	return &Response{bytes.NewBuffer(append(header, stdout...))}
}

func (r *Response) Close() error {
	return nil
}
