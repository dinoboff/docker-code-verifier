package verifier

import (
	"encoding/json"
	"github.com/ChrisBoesch/docker-code-verifier/server/verifier/mock"
	"github.com/samalba/dockerclient"
	"github.com/stretchr/testify/assert"
	tmock "github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestNewContainer(t *testing.T) {
	req, err := NewRequest("python3", []byte(`{"solution": "foo=1", "tests": ">>> foo\n1"}`))
	assert.Nil(t, err)

	v, err := NewContainer(nil, req)
	assert.Nil(t, err)

	assert.Equal(t, req, v.Request)
	assert.Equal(t, runtimes["python3"], v.Config.Image)
	assert.Equal(t, true, v.Config.NetworkDisabled)
	assert.Equal(t, false, v.Config.AttachStderr)
	assert.Equal(t, false, v.Config.AttachStdin)
	assert.Equal(t, false, v.Config.AttachStdout)
	assert.Equal(t, []string{"-e", "--tests", b64enc(">>> foo\n1"), b64enc("foo=1")}, v.Config.Cmd)
}

func TestNewContainerNoTest(t *testing.T) {
	req, err := NewRequest("python3", []byte(`{"solution": "foo=1"}`))
	assert.Nil(t, err)

	v, err := NewContainer(nil, req)
	assert.Nil(t, err)

	assert.Equal(t, req, v.Request)
	assert.Equal(t, runtimes["python3"], v.Config.Image)
	assert.Equal(t, true, v.Config.NetworkDisabled)
	assert.Equal(t, false, v.Config.AttachStderr)
	assert.Equal(t, false, v.Config.AttachStdin)
	assert.Equal(t, false, v.Config.AttachStdout)
	assert.Equal(t, []string{"-e", b64enc("foo=1")}, v.Config.Cmd)
}

func TestNewContainerNoRuntime(t *testing.T) {
	var req Request
	body := []byte(`{"solution": "System.out.println(\"Hello World!\");"}`)
	err := json.Unmarshal(body, &req)
	assert.Nil(t, err)

	v, err := NewContainer(nil, &req)
	assert.NotNil(t, err)
	assert.Nil(t, v)
}

func TestRunContainer(t *testing.T) {
	req, err := NewRequest("python3", []byte(`{"solution": "foo=1"}`))
	assert.Nil(t, err)

	client := dockerclient.NewMockClient()
	v, err := NewContainer(client, req)
	assert.Nil(t, err)

	var (
		watcher = NewWatcher()
		result  = make(chan bool, 1)
		timeout = make(chan time.Time, 1)

		containerId = "1234"
		stopEvent   = &dockerclient.Event{Id: containerId, Status: "stop"}
	)

	client.Mock.On("CreateContainer", v.Config, "").Return(containerId, nil)
	client.Mock.On("StartContainer", containerId, tmock.Anything).Return(nil)
	v.Run(watcher, result, func() <-chan time.Time {
		return timeout
	})

	client.Mock.AssertExpectations(t)
	assert.Len(t, watcher.containers[containerId], 1)

	watcher.cb(stopEvent)
	select {
	case b := <-result:
		assert.True(t, b)
	case <-time.After(1 * time.Second):
		t.Error("timeout")
	}
}

func TestRunContainerTimeOut(t *testing.T) {
	req, err := NewRequest("python3", []byte(`{"solution": "foo=1"}`))
	assert.Nil(t, err)

	client := dockerclient.NewMockClient()
	v, err := NewContainer(client, req)
	assert.Nil(t, err)

	var (
		watcher = NewWatcher()
		result  = make(chan bool, 1)
		timeout = make(chan time.Time, 1)

		containerId = "1234"
	)

	client.Mock.On("CreateContainer", v.Config, "").Return(containerId, nil)
	client.Mock.On("StartContainer", containerId, tmock.Anything).Return(nil)
	v.Run(watcher, result, func() <-chan time.Time {
		return timeout
	})

	client.Mock.AssertExpectations(t)

	client.Mock.On("StopContainer", containerId, 1).Return(nil)
	timeout <- time.Now()

	select {
	case b := <-result:
		assert.False(t, b)
		client.Mock.AssertExpectations(t)
	case <-time.After(1 * time.Second):
		t.Error("timeout")
	}
}

func TestRunContainerAlmostTimeOut(t *testing.T) {
	req, err := NewRequest("python3", []byte(`{"solution": "foo=1"}`))
	assert.Nil(t, err)

	client := dockerclient.NewMockClient()
	v, err := NewContainer(client, req)
	assert.Nil(t, err)

	var (
		watcher = NewWatcher()
		result  = make(chan bool, 1)
		timeout = make(chan time.Time, 1)

		containerId = "1234"
	)

	client.Mock.On("CreateContainer", v.Config, "").Return(containerId, nil)
	client.Mock.On("StartContainer", containerId, tmock.Anything).Return(nil)
	v.Run(watcher, result, func() <-chan time.Time {
		return timeout
	})

	client.Mock.AssertExpectations(t)

	// Stop will return an error indicated the container is already stopped
	// i.e. the container finished running the code after the timeour but before
	// we send the request to stop it
	client.Mock.On("StopContainer", containerId, 1).Return(dockerclient.Error{StatusCode: 304, Status: "container already stopped"})
	timeout <- time.Now()

	select {
	case b := <-result:
		assert.True(t, b)
		client.Mock.AssertExpectations(t)
	case <-time.After(1 * time.Second):
		t.Error("timeout")
	}
}

func TestGetResult(t *testing.T) {
	req, err := NewRequest("python3", []byte(`{"solution": "foo=1", "tests": ">>> foo\n1"}`))
	assert.Nil(t, err)

	client := dockerclient.NewMockClient()
	v, err := NewContainer(client, req)
	assert.Nil(t, err)

	var (
		rawResp     = mock.NewResponse([]byte(`{"solved": true, "results": [{"received": "1", "call": "foo", "expected": "1", "correct": true}], "printed": ""}`))
		containerId = "1234"
	)

	v.containerID = containerId
	client.Mock.On("ContainerLogs", containerId, v.LogOptions).Return(rawResp, nil)
	resp, err := v.GetResults()

	assert.Nil(t, err)
	assert.Equal(
		t,
		&Response{
			Solved:  true,
			Printed: "",
			Results: []*Call{&Call{Call: "foo", Expected: "1", Received: "1", Correct: true}},
		},
		resp,
	)
}
