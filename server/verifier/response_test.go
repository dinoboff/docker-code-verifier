package verifier

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	solved_response = []byte(`{"solved": true, "results": [{"received": "1", "call": "foo", "expected": "1", "correct": true}], "printed": ""}`)
)

func TestResponse(t *testing.T) {
	var resp Response
	err := json.Unmarshal(solved_response, &resp)

	assert.Nil(t, err)
	assert.Equal(t, true, resp.Solved)
	assert.Equal(t, 1, len(resp.Results))
	assert.Equal(t, "", resp.Printed)
	assert.Equal(t, "", resp.Errors)
}
