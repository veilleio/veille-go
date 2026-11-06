package veille

import "fmt"

// ConfigurationError reports a missing or unusable client configuration.
type ConfigurationError struct {
	Message string
	Err     error
}

func (e *ConfigurationError) Error() string { return e.Message }
func (e *ConfigurationError) Unwrap() error { return e.Err }

// SerializationError reports a parameter that is missing or cannot be sent.
type SerializationError struct {
	Param   string
	Message string
	Err     error
}

func (e *SerializationError) Error() string { return e.Message }
func (e *SerializationError) Unwrap() error { return e.Err }

// TransportError reports a failure before a complete HTTP response was received,
// including client-side timeouts and context cancellation (see Unwrap).
type TransportError struct {
	Message string
	Err     error
}

func (e *TransportError) Error() string { return e.Message + ": " + e.Err.Error() }
func (e *TransportError) Unwrap() error { return e.Err }

// APIError is returned for non-2xx responses.
type APIError struct {
	StatusCode int
	// Message is the "error" (or "message") field of the JSON body, when present.
	Message  string
	Response *Response
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("Veille API error %d: %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("Veille API request failed with status %d.", e.StatusCode)
}

// ResponseDecodingError reports a success response that could not be decoded.
type ResponseDecodingError struct {
	Message  string
	Response *Response
	Err      error
}

func (e *ResponseDecodingError) Error() string { return e.Message }
func (e *ResponseDecodingError) Unwrap() error { return e.Err }

// JobTimeoutError reports an async job that did not finish before the deadline.
type JobTimeoutError struct {
	JobID        string
	LastResponse *Response
}

func (e *JobTimeoutError) Error() string {
	return "Job " + e.JobID + " did not finish before the deadline."
}
