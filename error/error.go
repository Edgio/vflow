package error

import (
	"bytes"
	"errors"
	"fmt"
)

// ErrInvalidSetID is returned when the set id is 0.
var ErrInvalidSetID = errors.New("failed to decodeSet / invalid setID")

// ErrNoFieldsDecoded is returned when no field was decoded from the data.
var ErrNoFieldsDecoded = errors.New("no field decoded from data")

var _ error = &InvalidProtocolVersionError{}

// InvalidProtocolVersionError indicates the received version is invalid.
type InvalidProtocolVersionError struct {
	// Expected is the expected version.
	Expected uint16
	// Received is the received version.
	Received uint16
	// Protocol is the protocol name.
	Protocol string
}

// Error implements error.
func (e *InvalidProtocolVersionError) Error() string {
	return fmt.Sprintf("invalid %s version (expected: %d) (received: %d)",
		e.Protocol, e.Expected, e.Received)
}

var _ error = &CombinedErrors{}

// CombinedErrors is a collection of errors.
type CombinedErrors struct {
	Errors []error
}

// Error implements error.
func (e *CombinedErrors) Error() string {
	if len(e.Errors) == 1 {
		return e.Errors[0].Error()
	}
	var errMsg bytes.Buffer
	errMsg.WriteString("Multiple errors:")
	for _, subError := range e.Errors {
		errMsg.WriteString("\n- " + subError.Error())
	}
	return errMsg.String()
}

// CombineErrors returns a CombinedErrors containing the given errors.
func CombineErrors(errorSlice ...error) (err error) {
	if len(errorSlice) == 0 {
		return nil
	}
	return &CombinedErrors{Errors: errorSlice}
}
