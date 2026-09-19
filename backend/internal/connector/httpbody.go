package connector

import (
	"fmt"
	"io"
)

// MaxResponseBytes caps how much of an upstream HTTP response body is read.
const MaxResponseBytes = 10 << 20 // 10 MiB

// ReadBody reads r up to MaxResponseBytes and errors if the body is larger.
func ReadBody(r io.Reader) ([]byte, error) {
	data, err := io.ReadAll(io.LimitReader(r, MaxResponseBytes+1))
	if err != nil {
		return nil, err
	}
	if len(data) > MaxResponseBytes {
		return nil, fmt.Errorf("upstream response exceeds %d bytes", MaxResponseBytes)
	}
	return data, nil
}

// LimitedBody wraps r so reads stop after MaxResponseBytes (for streaming decoders).
func LimitedBody(r io.Reader) io.Reader { return io.LimitReader(r, MaxResponseBytes) }
