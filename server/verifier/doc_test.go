package verifier

import (
	"github.com/samalba/dockerclient"
	"log"
	"os"
)

func Example() {
	// Create a docker client.
	docker, err := dockerclient.NewDockerClient(os.Getenv("DOCKER_HOST"), nil)
	if err != nil {
		log.Fatalf("Failed to create the docker client: %s", err)
	}

	// Create a monitor to know when a container has finished running a test
	watcher := NewWatcher()
	watcher.Start(docker)

	// Create a test request
	request, err := NewRequest("python3", []byte(`{"solution": "foo=1", "tests": ">>> foo\n1"}`))
	if err != nil {
		log.Fatalf("Failed to create the test request: %s", err)
	}

	// Run the tests
	resp, err := Run(docker, watcher, request, false)
	if err != nil {
		log.Fatalf("Failed to run a verifier container: %s", err)
	}

	log.Println("Test solved?: %s", resp.Solved)
	log.Println("%+u", resp)
}
