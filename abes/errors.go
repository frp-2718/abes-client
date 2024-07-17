package abes

// NotFoundError occurs when the server returns a 404 or 500 response.
// Some ABES API use HTTP 500 errors to report a "not found" error.
type NotFoundError struct {
	ErrString string
}

func (e NotFoundError) Error() string {
	return e.ErrString
}

type NetworkError struct {
	err       error
	ErrString string
}

func (e NetworkError) Error() string {
	return e.ErrString
}
