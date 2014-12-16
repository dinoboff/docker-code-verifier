/*
Package verifier implement a simple library to launch docker container
to run user code. The container are used as a sandox to run untrusted
code and extract the result.
*/
package verifier

import (
	"errors"
	"time"
)

var (
	// errors
	ErrNotImplemented = errors.New("Unknown runtime")

	// supported runtime
	runtimes = map[string]string{
		"python":  "singpath/verifier-python3",
		"python3": "singpath/verifier-python3",
	}
)

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

func Run(client Client, req *Request, debug bool) (*Response, error) {
	container, err := NewContainer(client, req)
	if err != nil {
		return nil, err
	}

	result := make(chan *Response)
	err = container.Run(result)
	defer container.Remove()
	if err != nil {
		return nil, err
	}

	select {
	case resp := <-result:
		return resp, nil
	case <-time.Tick(6 * time.Second):
		container.Stop()
		return &Response{Errors: "Code failed to run in time"}, nil
	}
}
