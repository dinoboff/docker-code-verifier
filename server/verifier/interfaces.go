package verifier

type Client interface {
	Logs(containerId string) (*LogStreams, error)
	Create(image string, cmd []string) (string, error)
	Remove(containerId string) error
	Start(containerId string) error
	Stop(containerId string) error
	Wait(containerId string) (int, error)
}
