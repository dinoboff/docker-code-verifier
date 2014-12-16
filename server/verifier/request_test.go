package verifier

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRequest(t *testing.T) {
	var req Request
	body := []byte(`{"solution": "foo=1", "tests": ">>> foo\n1"}`)
	err := json.Unmarshal(body, &req)

	assert.Nil(t, err)
	assert.Equal(t, "foo=1", req.Solution, "Solution should hold the 'solution' property")
	assert.Equal(t, ">>> foo\n1", req.Tests, "Tests should hold the 'tests' property")
}

func TestNewRequest(t *testing.T) {
	req, err := NewRequest("python3", []byte(`{"solution": "foo=1", "tests": ">>> foo\n1"}`))

	assert.Nil(t, err)
	assert.Equal(t, "python3", req.Runtime, "Solution should hold the runtime property")
	assert.Equal(t, "foo=1", req.Solution, "Solution should hold the 'solution' property")
	assert.Equal(t, ">>> foo\n1", req.Tests, "Tests should hold the 'tests' property")
}

func TestToCmd(t *testing.T) {
	req, err := NewRequest("python3", []byte(`{"solution": "foo=1", "tests": ">>> foo\n1"}`))
	assert.Nil(t, err)

	cmd := req.toCmd()
	assert.Len(t, cmd, 3)
	assert.Equal(t, "-e", cmd[0])
	assert.Equal(t, "--tests=Pj4+IGZvbwox", cmd[1])
	assert.Equal(t, "Zm9vPTE=", cmd[2])
}
