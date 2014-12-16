package verifier

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type Header struct {
	Type uint8
	_    uint8
	_    uint8
	_    uint8
	Size uint32
}

func (h *Header) target(s *LogStreams) (*bytes.Buffer, error) {
	switch h.Type {
	case 0, 1:
		return s.Stdout, nil
	case 2:
		return s.Stderr, nil
	default:
		return nil, fmt.Errorf("wrong stream type: %+v", h)
	}
}

type LogStreams struct {
	Stdout *bytes.Buffer
	Stderr *bytes.Buffer
}

// Parse container logs
// See https://docs.docker.com/reference/api/docker_remote_api_v1.15/#attach-to-a-container
func NewLogStreams(logs io.Reader) (*LogStreams, error) {
	var (
		streams = &LogStreams{new(bytes.Buffer), new(bytes.Buffer)}

		header Header
		err    error
		frame  []byte
		target *bytes.Buffer
	)

	for {
		err = binary.Read(logs, binary.BigEndian, &header)
		if err == io.ErrUnexpectedEOF {
			err = fmt.Errorf("Failed parsing frame header (frame so far %+v.", streams)
			break
		}
		if err != nil {
			break
		}

		target, err = header.target(streams)
		if err != nil {
			break
		}

		frame = make([]byte, header.Size)
		_, err = io.ReadFull(logs, frame)
		if err != nil {
			err = fmt.Errorf("Failed parsing frame header (frame so far %v): %s.", frame, err)
			break
		}

		_, err = target.Write(frame)
		if err != nil {
			err = fmt.Errorf("Failed to write frame (%s): %s", frame, err)
			break
		}
	}

	if err != nil && err != io.EOF {
		return nil, err
	}

	return streams, nil
}
