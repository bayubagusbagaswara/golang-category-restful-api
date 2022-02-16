package exception

// mengikuti kontrak dari interface error
type NotFoundError struct {
	Error string
}

// bikin method untuk struct NotFoundError
func NewNotFoundError(error string) NotFoundError {
	return NotFoundError{Error: error}
}
