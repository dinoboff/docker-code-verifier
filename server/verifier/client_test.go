package verifier

import (
	"encoding/json"
	"github.com/ChrisBoesch/docker-code-verifier/server/verifier/mock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

const key = `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEA6FoApIVTJtOeXZoPlOJnkQ+9zbC6fpssDe+o4X8lUYa7jTZn
/6b9L8OiBnVkCDJMTx9kDrP6tM0lg89sDCh/qjbheN1BOJrzesyNsAq76J2VJdjU
9gcLceDAPsjzPa1KsKANydymCSdJBrwea7yhDKr8Ux3QAUQivAKud3c6n2VXWDMv
9qaqL/tvNsYMXkw3KDcUvpQZBD7lVhDcQodXznJ+oA7ytfnYNg5qVfldOO6RJls9
gpotTc3JmexVQ7U+ii5EEY4GyraiXuD5Xhen2Zidj3LHSSv4BAVV2OY0NPJXTDce
434osgo6his2bQcPz1yuMDHEuSjfT8TRCuSKCwIDAQABAoIBAHkVliomlLuyImBz
DdWv2vr8shQEGlwqL14f2+mPofoHdavUg4P2GRGQKNqmyHeBBsVg/XqwDmG0Wu2C
5bK8VDN3IC8lVnzSOzpuaRQps904aeZsRibkavFwh57wh9pHeZnr/uOIijpQ77yD
wnKwvVjlrlL+QUy3nkZOO29kgfdYEWKtEA8z1XRdlUUBPGQxXwVMX3mY9nj6oT3K
6NOaw9wILOm/U0GTI455PhvgbNnNp4aNHd82+otdHfIkrvtbeb5BKSqUwxTpMFQP
72H4UdElqm59+xANUBvGMf21oyiq2YRwAydNqtTYfa05tQOHun1lHd970zwzyjTK
B16Vc+ECgYEA8Rvao2q/tdQCQez7c9KJQ05owQl+ogwrgthUp08VjLl4oEHHjPZ5
jAoZBjoCfEVJP/ZLu+tRtRugZZ6okBMxGpifIk9DWSAQxvJXF8LsF8LW0arazsNs
jBwg8unqUcVGmNNokd3HcINX/93ej4JU5bPs2rWcUqwYPte3yfcxATECgYEA9rOw
U063OX94/Q5/K97O1s65chpB3HIt+k42IH+22PcocJCLC8NdsGJmtNY6YR5h4i7C
ThXsehE+ax2NgR0J+ZwpWjQPy3q37WtSg8xjmDuzGb/YwmGCw+eLiM9XBoaQB/cD
4NJY4Ujtva0y32FFpjwXSPvwYp4CPVOkgbWOj/sCgYEA3C7MrpnIszsGSMArLa1h
jqanQUnza/bjMV1viU7OZjHmN6t6mX9opnt+ONJ2/JelehTpOpZ+in7NLqACXXg5
SomAIavy3AxNZfFfmaJ3Soey98wof9O1aAo0CMGXK8+VVfESOMso29YGYfJy0el6
sD5smZpqRJFGnvUOsRDdnrECgYEA3guvedwQqCzmzgX9SpQ5YTghy+R8QRl37qH0
r92jyrby7BX9QLIwInD+9mcXlpBNE9J4SuYKuXfJ0YmA8qQbdVIsGidfzAqBf60o
UL5nKf8Z3eRCCfrQQtmmSpYsQxBclP6su+831lXYve8lKc+Ya94MK0GwBGMpqt8c
4y5xyX8CgYEArWV0H7Mrj/5YSYImfH5HA/auu+F86sHcgLMBEhD2NtYpNg3z+gj4
xxUIULihu81ymEj/yIOZPBeBWl3uLWTx/oSyUXG+AadjwqoTz6VLfZJHlylUpdmb
oAAihiyaHLxHLQGOCPpvJ+trx3d5L18UGPs8M739yCajFQt7e9rm/d8=
-----END RSA PRIVATE KEY-----`

const cert = `-----BEGIN CERTIFICATE-----
MIIC5zCCAdGgAwIBAgIQWlczqbSIKznD4oqpNaVJazALBgkqhkiG9w0BAQswFjEU
MBIGA1UEChMLQm9vdDJEb2NrZXIwHhcNMTQxMTI0MTgxNDI2WhcNMTcxMTA4MTgx
NDI2WjAWMRQwEgYDVQQKEwtCb290MkRvY2tlcjCCASIwDQYJKoZIhvcNAQEBBQAD
ggEPADCCAQoCggEBAOhaAKSFUybTnl2aD5TiZ5EPvc2wun6bLA3vqOF/JVGGu402
Z/+m/S/DogZ1ZAgyTE8fZA6z+rTNJYPPbAwof6o24XjdQTia83rMjbAKu+idlSXY
1PYHC3HgwD7I8z2tSrCgDcncpgknSQa8Hmu8oQyq/FMd0AFEIrwCrnd3Op9lV1gz
L/amqi/7bzbGDF5MNyg3FL6UGQQ+5VYQ3EKHV85yfqAO8rX52DYOalX5XTjukSZb
PYKaLU3NyZnsVUO1PoouRBGOBsq2ol7g+V4Xp9mYnY9yx0kr+AQFVdjmNDTyV0w3
HuN+KLIKOoYrNm0HD89crjAxxLko30/E0QrkigsCAwEAAaM1MDMwDgYDVR0PAQH/
BAQDAgCAMBMGA1UdJQQMMAoGCCsGAQUFBwMCMAwGA1UdEwEB/wQCMAAwCwYJKoZI
hvcNAQELA4IBAQBp+EjPOGghQa+NixO3FR43GJ2uzvtSChlPCXn41VW7vRBPDm3Y
DpHEwASOhBrAEWaDx1yNDFJ4txk3aKwk6R2N8TycJ5IkmbdJzcWO0i8+Sg4zdFrC
AWofpGNs2f03aNsdgLObgWnQpGamFm2W9h3lQAqhdxPuiMUAhd3WaP34Z1P050sF
GM+/PNkeF/BTvY+ni+K9IP+UsdqZTTR71hobbr9m0nssJ/K/INOsRSMPKyDdYLN6
OyXK/ciq5Mv7n9XZAVcvW1elVSqdlqqysOT2uYCblH2FV9INwfHDeEY/px0o+i/+
qYdFP1/ejZb50EmbYf1fcHZCrtWUAqNW1DhY
-----END CERTIFICATE-----`

func TestNewTCPDockerClient(t *testing.T) {
	client, err := NewClient("tcp://localhost:5000", "")
	assert.Nil(t, err)
	assert.Equal(t, client.URL.String(), "http://localhost:5000")

	transport, ok := client.HTTPClient.Transport.(*http.Transport)
	assert.True(t, ok)
	assert.Nil(t, transport.TLSClientConfig)
}

func TestNewSecureTCPDockerCLient(t *testing.T) {
	temp, err := ioutil.TempDir("", "golangTest")
	assert.Nil(t, err)
	defer os.Remove(temp)

	ioutil.WriteFile(filepath.Join(temp, defaultKeyFile), []byte(key), 0755)
	ioutil.WriteFile(filepath.Join(temp, defaultCertFile), []byte(cert), 0755)

	client, err := NewClient("tcp://localhost:5000", temp)
	assert.Nil(t, err)
	assert.Equal(t, client.URL.String(), "https://localhost:5000")

	transport, ok := client.HTTPClient.Transport.(*http.Transport)
	assert.True(t, ok)
	assert.NotNil(t, transport.TLSClientConfig)
}

func TestCreateContainer(t *testing.T) {
	var (
		resp = []byte(`{"Id": "1234", "Warnings": null}`)
	)

	server, err := mock.NewMockDockerServer(resp, 200)
	assert.Nil(t, err)
	defer server.Close()

	client, err := NewClient(server.Host, "")
	assert.Nil(t, err)

	id, err := client.Create("singpath/verifier-python3", []string{"print('foo')"})
	assert.Nil(t, err)
	assert.Equal(t, "1234", id)

	assert.Len(t, server.URLs, 1)
	assert.Equal(t, "/containers/create", server.URLs[0].Path)
	assert.Equal(t, "POST", server.Methods[0])

	assert.Len(t, server.Bodies, 1)

	var req containerReq
	err = json.Unmarshal(server.Bodies[0], &req)
	assert.Nil(t, err)
	assert.True(t, req.NetworkDisabled)
	assert.False(t, req.Tty)
	assert.False(t, req.AttachStdin)
	assert.False(t, req.AttachStderr)
	assert.True(t, req.AttachStdout)
	assert.Equal(t, []string{"print('foo')"}, req.Cmd)
	assert.Equal(t, "singpath/verifier-python3", req.Image)
	assert.Equal(t, "", req.Entrypoint)
}

func TestStartContainer(t *testing.T) {
	server, err := mock.NewMockDockerServer([]byte{}, 204)
	assert.Nil(t, err)
	defer server.Close()

	client, err := NewClient(server.Host, "")
	assert.Nil(t, err)

	err = client.Start("1234")
	assert.Nil(t, err)
	assert.Equal(t, "/containers/1234/start", server.URLs[0].Path)
	assert.Equal(t, "POST", server.Methods[0])
}

func TestContianerLogs(t *testing.T) {
	r, _ := mock.NewResponse([]byte(`{"ok": true}`))
	resp, _ := ioutil.ReadAll(r)
	r.Close()

	server, err := mock.NewMockDockerServer(resp, 200)
	assert.Nil(t, err)
	defer server.Close()

	client, err := NewClient(server.Host, "")
	assert.Nil(t, err)

	logs, err := client.Logs("1234")
	assert.Nil(t, err)
	assert.Len(t, server.URLs, 1)

	assert.Equal(t, "/containers/1234/attach", server.URLs[0].Path)
	assert.Equal(t, "1", server.URLs[0].Query().Get("stdout"))
	assert.Equal(t, "", server.URLs[0].Query().Get("stream"))
	assert.Equal(t, "1", server.URLs[0].Query().Get("logs"))

	parsedResp, err := ioutil.ReadAll(logs.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(`{"ok": true}`), parsedResp)
}

func TestStopContainer(t *testing.T) {
	server, err := mock.NewMockDockerServer([]byte{}, 204)
	assert.Nil(t, err)
	defer server.Close()

	client, err := NewClient(server.Host, "")
	assert.Nil(t, err)

	err = client.Stop("1234")
	assert.Nil(t, err)
	assert.Equal(t, "/containers/1234/stop", server.URLs[0].Path)
	assert.Equal(t, "POST", server.Methods[0])
}

func TestRemoveContainer(t *testing.T) {
	server, err := mock.NewMockDockerServer([]byte{}, 204)
	assert.Nil(t, err)
	defer server.Close()

	client, err := NewClient(server.Host, "")
	assert.Nil(t, err)

	err = client.Remove("1234")
	assert.Nil(t, err)
	assert.Equal(t, "/containers/1234", server.URLs[0].Path)
	assert.Equal(t, "1", server.URLs[0].Query().Get("force"))
	assert.Equal(t, "1", server.URLs[0].Query().Get("v"))
	assert.Equal(t, "DELETE", server.Methods[0])
}
