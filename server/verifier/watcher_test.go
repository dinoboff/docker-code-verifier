package verifier

import (
	"github.com/ChrisBoesch/docker-code-verifier/server/verifier/mock"
	"github.com/samalba/dockerclient"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewWatcher(t *testing.T) {
	w := NewWatcher()
	assert.NotNil(t, w.containers)
}

func TestStartWatcher(t *testing.T) {
	watcher := NewWatcher()
	client := &mock.MonitorClient{}
	watcher.Start(client)

	assert.ObjectsAreEqual(watcher.cb, client.CBs[0])
	assert.Len(t, client.Args[0], 0)
}

func TestStartWatcherWIthTestifyMock(t *testing.T) {
	t.Skip("Fix mocking...")

	w := NewWatcher()
	c := dockerclient.NewMockClient()
	c.Mock.On("StartMonitorEvents", w.cb).Return()

	w.Start(c)
	c.Mock.AssertExpectations(t)
}

func TestWatchStopped(t *testing.T) {
	w := NewWatcher()
	assert.NotNil(t, w.containers)

	stopped := make(chan bool, 1)
	stopEvent := &dockerclient.Event{
		Id:     "1234",
		Status: "stop",
		From:   "debian:stable",
		Time:   1234567890,
	}

	w.WatchStop("1234", stopped)
	assert.Equal(t, 1, len(w.containers["1234"]))

	w.cb(stopEvent)
	select {
	case <-stopped:
	default:
		t.Error("the watcher failed to notice the stop event")
	}

	assert.Equal(t, 0, len(w.containers["1234"]))
}
