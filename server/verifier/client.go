package verifier

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

const (
	defaultCertFile = "cert.pem"
	defaultKeyFile  = "key.pem"

	defaultDialTimeout = 10 * time.Second

	codeImageMissing     = 404
	codeInternalError    = 500
	codeAlreadyStarted   = 304
	codeAlreadyStopped   = 304
	codeContainerMissing = 404
	codeBadRequest       = 400
)

type DockerClient struct {
	URL        *url.URL
	HTTPClient *http.Client
}

func NewClient(host, cert_path string) (*DockerClient, error) {
	var (
		u          *url.URL
		httpClient *http.Client
		err        error
	)

	u, err = url.Parse(host)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse docker host URL (%s): %s", host, err)
	}

	switch u.Scheme {
	case "unix":
		httpClient = unixHTTPClient(u)
	default:
		httpClient, err = tcpHTTPClient(u, cert_path)
	}

	if err != nil {
		return nil, err
	}

	return &DockerClient{u, httpClient}, nil
}

func defaultDial(proto, addr string) (net.Conn, error) {
	return net.DialTimeout(proto, addr, defaultDialTimeout)
}

func tcpHTTPClient(u *url.URL, cert_path string) (*http.Client, error) {
	if cert_path != "" {
		return tcpHTTPSClient(u, cert_path)
	}

	u.Scheme = "http"
	u.Path = ""

	httpTransport := &http.Transport{
		Dial: defaultDial,
	}

	return &http.Client{Transport: httpTransport}, nil
}

func tcpHTTPSClient(u *url.URL, cert_path string) (*http.Client, error) {
	var (
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS10,
		}
		flCert = filepath.Join(cert_path, defaultCertFile)
		flKey  = filepath.Join(cert_path, defaultKeyFile)
	)

	_, err := os.Stat(flCert)
	if err != nil {
		return nil, fmt.Errorf("Failed to find '%s' in '%s' (%s)", defaultCertFile, cert_path, err)
	}

	_, err = os.Stat(flKey)
	if err != nil {
		return nil, fmt.Errorf("Failed to find '%s' in '%s' (%s)", defaultKeyFile, cert_path, err)
	}

	cert, err := tls.LoadX509KeyPair(flCert, flKey)
	if err != nil {
		return nil, fmt.Errorf("Invalid X509 Key Pair ('%s'/'%s') in %s (%s)", defaultCertFile, defaultKeyFile, cert_path, err)
	}

	u.Scheme = "https"
	u.Path = ""

	tlsConfig.Certificates = []tls.Certificate{cert}

	httpTransport := &http.Transport{
		TLSClientConfig: tlsConfig,
		Dial:            defaultDial,
	}

	return &http.Client{Transport: httpTransport}, nil
}

func unixHTTPClient(u *url.URL) *http.Client {
	address := u.Path
	u.Scheme = "http"
	u.Host = "unix.sock"
	u.Path = ""

	httpTransport := &http.Transport{
		Dial: func(proto, addr string) (net.Conn, error) {
			return net.DialTimeout("unix", address, defaultDialTimeout)
		},
	}

	return &http.Client{Transport: httpTransport}
}

type containerReq struct {
	AttachStdin     bool
	AttachStdout    bool
	AttachStderr    bool
	Tty             bool
	Cmd             []string
	Image           string
	Entrypoint      string
	NetworkDisabled bool
}

func newContainerReq(image string, cmd []string) *containerReq {
	return &containerReq{
		AttachStdin:     false,
		AttachStdout:    true,
		AttachStderr:    false,
		Tty:             false,
		Cmd:             cmd,
		Image:           image,
		NetworkDisabled: true,
	}
}

type containerResponse struct {
	Id       string
	Warnings []string
}

func (d *DockerClient) rpost(path string, query url.Values, req interface{}) (*http.Response, error) {
	u := *d.URL
	u.Path = path
	u.RawQuery = query.Encode()

	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := d.HTTPClient.Post(u.String(), "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (d *DockerClient) post(path string, query url.Values, req interface{}) (*http.Response, []byte, error) {
	resp, err := d.rpost(path, query, req)
	if err != nil {
		return nil, nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	return resp, body, nil
}

func (d *DockerClient) Create(image string, cmd []string) (string, error) {
	resp, body, err := d.post("/containers/create", nil, newContainerReq(image, cmd))
	if err != nil {
		return "", err
	}

	if resp.StatusCode == codeImageMissing {
		return "", fmt.Errorf("Image '%s' not found: %s", image, body)
	}

	if resp.StatusCode == codeInternalError {
		return "", fmt.Errorf("Docker failed to create the container: %s", body)
	}

	var c containerResponse
	err = json.Unmarshal(body, &c)
	if err != nil {
		return "", fmt.Errorf("Failed to parse container creation response (%s): %s", string(body), err)
	}

	if len(c.Warnings) > 0 {
		log.Printf("Container creation warning: %+v", c.Warnings)
	}

	return c.Id, nil
}

func (d *DockerClient) Start(containerId string) error {
	resp, body, err := d.post(fmt.Sprintf("/containers/%s/start", containerId), nil, nil)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case codeAlreadyStarted:
		return fmt.Errorf("Container %s already started", containerId)
	case codeContainerMissing:
		return fmt.Errorf("Container %s is missing", containerId)
	case codeInternalError:
		return fmt.Errorf("Docker failed to start container %s: %s", containerId, body)
	default:
		return nil
	}
}

type waitResp struct {
	StatusCode int
}

func (d *DockerClient) Wait(containerId string) (int, error) {
	resp, body, err := d.post(fmt.Sprintf("/containers/%s/wait", containerId), nil, nil)
	if err != nil {
		return 0, err
	}

	switch resp.StatusCode {
	case codeContainerMissing:
		return 0, fmt.Errorf("Container %s is missing", containerId)
	case codeInternalError:
		return 0, fmt.Errorf("Docker failed to wait for container %s: %s", containerId, body)
	default:
	}

	var r waitResp
	err = json.Unmarshal(body, &r)
	if err != nil {
		return 0, err
	}

	return r.StatusCode, nil
}

func (d *DockerClient) Logs(containerId string) (*LogStreams, error) {
	query := url.Values{}
	query.Set("logs", "1")
	query.Set("stdout", "1")

	resp, err := d.rpost(fmt.Sprintf("/containers/%s/attach", containerId), query, nil)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case codeBadRequest:
		return nil, fmt.Errorf("Unexpected error while attaching to container %s stream", containerId)
	case codeContainerMissing:
		return nil, fmt.Errorf("Container %s is missing", containerId)
	case codeInternalError:
		return nil, fmt.Errorf("Docker failed to attach container %s", containerId)
	default:
	}

	logs, err := NewLogStreams(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Failed to parse container %s stream: %s", containerId, err)
	}

	return logs, nil
}

func (d *DockerClient) Stop(containerId string) error {
	resp, body, err := d.post(fmt.Sprintf("/containers/%s/stop", containerId), nil, nil)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case codeAlreadyStopped:
		return fmt.Errorf("Container %s already started", containerId)
	case codeContainerMissing:
		return fmt.Errorf("Container %s is missing", containerId)
	case codeInternalError:
		return fmt.Errorf("Docker failed to start container %s: %s", containerId, body)
	default:
		return nil
	}
}
func (d *DockerClient) Remove(containerId string) error {
	u := *d.URL
	q := url.Values{}
	q.Set("force", "1")
	q.Set("v", "1")

	u.Path = fmt.Sprintf("/containers/%s", containerId)
	u.RawQuery = q.Encode()

	req := http.Request{
		URL:    &u,
		Method: "DELETE",
	}
	resp, err := d.HTTPClient.Do(&req)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case codeBadRequest:
		return fmt.Errorf("Unexpected error while removing container %s", containerId)
	case codeContainerMissing:
		return fmt.Errorf("Container %s is missing", containerId)
	case codeInternalError:
		return fmt.Errorf("Docker failed to start container %s", containerId)
	default:
		return nil
	}
}
