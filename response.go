package veille

import (
	"encoding/json"
	"net/http"
)

// Response is a buffered, successful API response.
type Response struct {
	StatusCode  int
	Header      http.Header
	ContentType string
	// Body holds the raw bytes: JSON text, or a file such as an image or PDF.
	Body []byte
	// IsJSON reports whether the body is JSON (Content-Type application/json or +json).
	IsJSON bool
}

// Decode unmarshals a JSON body into v.
func (r *Response) Decode(v any) error {
	if !r.IsJSON {
		return &ResponseDecodingError{Message: "The response is not JSON (" + r.ContentType + ").", Response: r}
	}
	if err := json.Unmarshal(r.Body, v); err != nil {
		return &ResponseDecodingError{Message: "The JSON response could not be decoded.", Response: r, Err: err}
	}
	return nil
}

// Map decodes a JSON object body.
func (r *Response) Map() (map[string]any, error) {
	var value map[string]any
	err := r.Decode(&value)
	return value, err
}
