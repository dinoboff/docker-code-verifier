package verifier

import (
	"github.com/samalba/dockerclient"
	"io"
	"time"
)

// interface for a timeout func
type Timeout func() <-chan time.Time

type StopWatcher interface {
	WatchStop(containerId string, hasStopped chan bool)
}

// partial dockerclient.Client interface, used by Watcher
type MonitorStarter interface {
	StartMonitorEvents(cb dockerclient.Callback, args ...interface{})
}

// partial dockerclient.Client interface, used by Container.Run and verifier.Stop
type ContainerStarter interface {
	CreateContainer(config *dockerclient.ContainerConfig, name string) (string, error)
	StartContainer(id string, config *dockerclient.HostConfig) error
	StopContainer(id string, timeout int) error
}

// partial dockerclient.Client interface, used by Container.GetResult
type ContainerLogger interface {
	ContainerLogs(id string, options *dockerclient.LogOptions) (io.ReadCloser, error)
}
