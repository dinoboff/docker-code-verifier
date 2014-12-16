package verifier

// The test result from the container.
type Response struct {
	Solved  bool    `json:"solved"`
	Printed string  `json:"printed,omitempty"`
	Errors  string  `json:"errors,omitempty"`
	Results []*Call `json:"results,omitempty"`
}

// One call result in a test.
type Call struct {
	Call     string `json:"call"`
	Expected string `json:"expected"`
	Received string `json:"received"`
	Correct  bool   `json:"correct"`
}
