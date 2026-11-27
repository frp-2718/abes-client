package abes

import (
	"errors"
	"fmt"
)

// APIError is returned when the ABES server answers with a non-2xx status
// that does not mean "not found" (see NotFoundError).
type APIError struct {
	Op         string // calling operation, e.g. "Multiwhere.GetMultiLocations"
	URL        string
	StatusCode int
	Body       string // response body, possibly truncated
}

func (e *APIError) Error() string {
	return fmt.Sprintf("abes: %s: GET %s: unexpected HTTP status %d: %s", e.Op, e.URL, e.StatusCode, e.Body)
}

// NotFoundError occurs when the requested resource does not exist. ABES
// services report "not found" in different ways: isbn2ppn and the
// <ppn>.xml endpoint answer HTTP 404, while multiwhere answers HTTP 500
// with an <error> element in the body.
type NotFoundError struct {
	Op      string
	URL     string
	Message string // content of the server's <error> element, if any
}

func (e *NotFoundError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("abes: %s: GET %s: not found: %s", e.Op, e.URL, e.Message)
	}
	return fmt.Sprintf("abes: %s: GET %s: not found", e.Op, e.URL)
}

// NetworkError wraps a transport-level failure (DNS, TLS, connection reset...).
type NetworkError struct {
	Op  string
	URL string
	Err error
}

func (e *NetworkError) Error() string { return fmt.Sprintf("abes: %s: GET %s: %v", e.Op, e.URL, e.Err) }
func (e *NetworkError) Unwrap() error { return e.Err }

// SystemError reports a local failure: unreadable body or invalid XML.
type SystemError struct {
	Op  string
	URL string
	Err error
}

func (e *SystemError) Error() string { return fmt.Sprintf("abes: %s: GET %s: %v", e.Op, e.URL, e.Err) }
func (e *SystemError) Unwrap() error { return e.Err }

// NotFound reports whether err is or wraps a *NotFoundError.
func NotFound(err error) bool {
	var nf *NotFoundError
	return errors.As(err, &nf)
}
