package cas

import "errors"

// ErrInvalidProtocolVersion is an error message
var ErrInvalidProtocolVersion = errors.New("invalid CAS protocol version")

// ErrInvalidServerURL is an error message
var ErrInvalidServerURL = errors.New("invalid CAS server URL")
