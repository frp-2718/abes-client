package abes

// NotFoundError occurs when the server returns a 404 or 500 response.
// Some ABES API use HTTP 500 errors to report a "not found" error.
type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

type NetworkError struct {
	Message string
}

func (e *NetworkError) Error() string {
	return e.Message
}
