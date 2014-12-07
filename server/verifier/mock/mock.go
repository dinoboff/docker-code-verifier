/*
Usulful mocks for github.com/ChrisBoesch/docker-code-verifier/server/verifier testing
*/

package mock

import (
	"bytes"
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

func NewResponse(jsonResp []byte) *Response {
	return &Response{bytes.NewBuffer(jsonResp)}
}

func (r *Response) Close() error {
	return nil
}
