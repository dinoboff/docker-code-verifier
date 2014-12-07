package main

import (
	"crypto/tls"
	"errors"
	"github.com/samalba/dockerclient"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultKeyFile  = "key.pem"
	defaultCertFile = "cert.pem"
)

var (
	ErrCertNotFound       = errors.New("cert.pem was not found. Check DOCKER_CERT_PATH is correctly set.")
	ErrKeyNotFound        = errors.New("key.pem was not found. Check DOCKER_CERT_PATH is correctly set.")
	ErrInvalidCertKeyPair = errors.New("Invalid X509 Key Pair. Check DOCKER_CERT_PATH is correctly set.")
)

func getClient(host, dockerCertPath string) (*dockerclient.DockerClient, error) {
	if strings.HasPrefix(host, "tcp:") == false || dockerCertPath == "" {
		return dockerclient.NewDockerClient(host, nil)
	}

	var (
		tlsConfig = tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS10,
		}
		flCert = filepath.Join(dockerCertPath, defaultCertFile)
		flKey  = filepath.Join(dockerCertPath, defaultKeyFile)
	)

	_, err := os.Stat(flCert)
	if err != nil {
		return nil, ErrCertNotFound
	}

	_, err = os.Stat(flKey)
	if err != nil {
		return nil, ErrKeyNotFound
	}

	cert, err := tls.LoadX509KeyPair(flCert, flKey)
	if err != nil {
		return nil, ErrInvalidCertKeyPair
	}

	tlsConfig.Certificates = []tls.Certificate{cert}
	return dockerclient.NewDockerClient(host, &tlsConfig)
}
