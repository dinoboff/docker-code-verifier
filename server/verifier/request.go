package verifier

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrSolutionIsMissing = errors.New("The user solution was not specified")
)

// String to standard base64 encoding bytes.
func b64enc(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// A User solution tests request.
type Request struct {
	Runtime  string
	Solution string
	Tests    string
}

func NewRequest(runtimeName string, body []byte) (*Request, error) {
	var req Request

	err := SupportedRuntime(runtimeName)
	if err != nil {
		return nil, err
	}

	req.Runtime = runtimeName
	err = json.Unmarshal(body, &req)
	if err != nil {
		return nil, err
	}

	if req.Solution == "" {
		return nil, ErrSolutionIsMissing
	}

	return &req, err
}

func (r *Request) toCmd() []string {
	return []string{
		"-e",
		fmt.Sprintf("--tests=%s", b64enc(r.Tests)),
		b64enc(r.Solution),
	}
}
