/*
Package verifier implement a simple library to launch docker container
to run user code. The container are used as a sandox to run untrusted
code and extract the result.
*/
package verifier

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/samalba/dockerclient"
	"io"
	"io/ioutil"
	"sync"
	"time"
)

var (
	// errors
	ErrNotImplemented    = errors.New("Unknown runtime")
	ErrSolutionIsMissing = errors.New("The user solution was not specified")
	ErrNotJson           = errors.New("Failed to read the container response")

	// supported runtime
	runtimes = map[string]string{
		"python":  "singpath/verifier-python3",
		"python3": "singpath/verifier-python3",
	}

	defaultTimeoutDuration = 5 * time.Second

	containerEnded = map[string]bool{
		"die":  true,
		"kill": true,
		"stop": true,
	}
)

const (
	containerAlreadyStopped = 304
	openBracket             = 123
	logHeaderSize           = 8
)

//Default timeout func.
//
//Set as function instead of a chanel directly, because we need the timeout
//to start after the container start running, before it's initialisation.
//
func DefaultTimeOut() <-chan time.Time {
	return time.After(defaultTimeoutDuration)
}

// String to standard base64 encoding bytes.
func b64enc(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// Test if the runtime is supported
func SupportedRuntime(name string) error {
	_, ok := runtimes[name]
	if ok {
		return nil
	}

	return ErrNotImplemented
}

func SupportedtRuntimeList() []string {
	result := make([]string, 0, len(runtimes))
	for key := range runtimes {
		result = append(result, key)
	}
	return result
}

func Run(client dockerclient.Client, watcher StopWatcher, req *Request, debug bool) (*Response, error) {
	container, err := NewContainer(client, req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if debug {
			return
		}

		cErr := container.Remove()
		if err == nil && cErr != nil {
			err = cErr
		}
	}()

	stopped := make(chan bool, 1)
	err = container.Run(watcher, stopped, DefaultTimeOut)
	if err != nil {
		return nil, err
	}

	<-stopped
	resp, err := container.GetResults()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// A User solution tests request.
type Request struct {
	Runtime  string
	Solution string
	Tests    string
}

func NewRequest(runtimeName string, body []byte) (*Request, error) {
	var req Request

	if runtimes[runtimeName] == "" {
		return &req, ErrNotImplemented
	}

	req.Runtime = runtimeName
	err := json.Unmarshal(body, &req)
	if err != nil {
		return nil, err
	}

	if req.Solution == "" {
		return nil, ErrSolutionIsMissing
	}

	return &req, err
}

// The test result from the container.
type Response struct {
	Solved  bool    `json:"solved"`
	Printed string  `json:"printed,omitempty"`
	Errors  string  `json:"errors,omitempty"`
	Results []*Call `json:"results,omitempty"`
}

// One call result in a test.
type Call struct {
	Call     string `json:"call"`
	Expected string `json:"expected"`
	Received string `json:"received"`
	Correct  bool   `json:"correct"`
}

// Watch docker event (only support the stop event).
//
// Usage to notify when a verifier container has stopped running the tests
//
type Watcher struct {
	sync.Mutex
	containers map[string][]chan bool
}

// Create a new watcher
func NewWatcher() *Watcher {
	return &Watcher{containers: make(map[string][]chan bool)}
}

// Start the watcher
func (w *Watcher) Start(client MonitorStarter) {
	client.StartMonitorEvents(w.cb)
}

// Set the the watcher to report stop event for a container ID.
func (w *Watcher) WatchStop(containerId string, hasStopped chan bool) {
	w.Lock()
	w.containers[containerId] = append(w.containers[containerId], hasStopped)
	w.Unlock()
}

func (w *Watcher) cb(event *dockerclient.Event, args ...interface{}) {
	if !containerEnded[event.Status] {
		return
	}

	w.Lock()

	stoppedList, ok := w.containers[event.Id]
	if !ok || len(stoppedList) == 0 {
		w.Unlock()
		return
	}

	delete(w.containers, event.Id)
	w.Unlock()

	for _, stopped := range stoppedList {
		select {
		case stopped <- true:
			continue // fine
		default:
			continue // too bad, but it won't block an other channel
		}
	}
}

// A verifier docker container.
type Container struct {
	Docker      dockerclient.Client
	Request     *Request
	Config      *dockerclient.ContainerConfig
	LogOptions  *dockerclient.LogOptions
	containerID string
}

// Create a new verifier container
func NewContainer(client dockerclient.Client, request *Request) (*Container, error) {
	var cmd []string

	runtimeImage, ok := runtimes[request.Runtime]
	if !ok {
		return nil, ErrNotImplemented
	}

	if request.Solution == "" {
		return nil, ErrSolutionIsMissing
	}

	if request.Tests == "" {
		cmd = []string{"-e", b64enc(request.Solution)}
	} else {
		cmd = []string{"-e", "--tests", b64enc(request.Tests), b64enc(request.Solution)}
	}

	config := &dockerclient.ContainerConfig{
		Image:           runtimeImage,
		Cmd:             cmd,
		NetworkDisabled: true,
	}
	logOptions := &dockerclient.LogOptions{
		Stdout: true,
	}
	return &Container{client, request, config, logOptions, ""}, nil
}

func (v *Container) Run(watcher StopWatcher, result chan<- bool, timeout Timeout) error {
	var err error

	v.containerID, err = v.Docker.CreateContainer(v.Config, "")
	if err != nil {
		return err
	}

	stopped := make(chan bool, 1)
	watcher.WatchStop(v.containerID, stopped)

	err = v.Docker.StartContainer(v.containerID, nil)
	if err != nil {
		return err
	}

	go func() {
		select {
		case <-stopped:
			result <- true
		case <-timeout():
			success := false
			defer func() {
				result <- success
			}()

			err = v.Stop(1)
			if err == nil {
				return
			}

			dockerErr, ok := err.(dockerclient.Error)
			if !ok {
				return
			}

			if dockerErr.StatusCode == containerAlreadyStopped {
				success = true
			}
		}
	}()
	return nil
}

func (v *Container) Stop(timeout int) error {
	return v.Docker.StopContainer(v.containerID, timeout)
}

// Query the the stdout logs of the container an parse the first line.
func (v *Container) GetResults() (*Response, error) {
	var resp Response

	logReader, err := v.Docker.ContainerLogs(v.containerID, v.LogOptions)
	if err != nil {
		return &resp, err
	}
	defer logReader.Close()

	logStreams, err := NewLogStreams(logReader)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(logStreams.Stdout)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse response (%s): %s", body, err)
	}

	return &resp, err
}

func (v *Container) Remove() error {
	return v.Docker.RemoveContainer(v.containerID, true)
}

type LogStreams struct {
	Stdout io.Reader
	Stderr io.Reader
}

// Parse container logs
// See https://docs.docker.com/reference/api/docker_remote_api_v1.15/#attach-to-a-container
func NewLogStreams(logs io.Reader) (*LogStreams, error) {
	var (
		sdtout = bytes.NewBuffer(nil)
		sdterr = bytes.NewBuffer(nil)

		target *bytes.Buffer
		header = make([]byte, logHeaderSize)
		frame  []byte
		size   uint32
		err    error
		i      int
	)

	for {
		i, err = logs.Read(header)
		if err != nil {
			break
		}
		if i != logHeaderSize {
			err = fmt.Errorf("wrong size header %v (%d)", header, i)
			break
		}

		switch header[0] {
		case 0, 1:
			target = sdtout
		case 2:
			target = sdterr
		default:
			err = fmt.Errorf("wrong stream type: %v", header[0])
			break
		}

		size = binary.BigEndian.Uint32(header[4:])
		frame = make([]byte, size)
		i, err = logs.Read(frame)
		if err != nil && err != io.EOF {
			break
		}
		if uint32(i) < size {
			err = fmt.Errorf("frame too short %s (%d)", frame, i)
			break
		}

		_, err = target.Write(frame)
		if err != nil {
			break
		}
	}

	if err != nil && err != io.EOF {
		return nil, err
	}

	return &LogStreams{sdtout, sdterr}, nil
}
