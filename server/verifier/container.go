package verifier

import (
	"encoding/json"
	"errors"
	"log"
)

// A verifier docker container.
type Container struct {
	docker      Client
	request     *Request
	image       string
	containerID string
}

var (
	ErrVerifierNotCreated = errors.New("The container for the verifier is not created.")
)

// Create a new verifier container
func NewContainer(client Client, request *Request) (*Container, error) {
	image, ok := runtimes[request.Runtime]
	if !ok {
		return nil, ErrNotImplemented
	}

	return &Container{client, request, image, ""}, nil
}

func (v *Container) Run(result chan<- *Response) error {
	id, err := v.docker.Create(v.image, v.request.toCmd())
	v.containerID = id
	if err != nil {
		return err
	}

	err = v.docker.Start(id)
	if err != nil {
		return err
	}

	go func() {
		// TODO: merge the 2 request. Must use a lower level api
		// to start the container once the hijack connection is started.
		_, err = v.docker.Wait(id)
		if err != nil {
			log.Print(err)
			result <- nil
			return
		}

		stream, err := v.docker.Logs(id)
		if err != nil {
			log.Print(err)
			result <- nil
			return
		}

		var r Response
		err = json.Unmarshal(stream.Stdout.Bytes(), &r)
		if err != nil {
			log.Printf("Failed to parse response from %s (%+u): %s", id, stream, err)
			result <- nil
			return
		}

		result <- &r
	}()

	return nil
}

func (v *Container) Stop() error {
	if v.containerID == "" {
		return ErrVerifierNotCreated
	}

	return v.docker.Stop(v.containerID)
}

func (v *Container) Remove() error {
	if v.containerID == "" {
		return ErrVerifierNotCreated
	}

	return v.docker.Remove(v.containerID)
}
