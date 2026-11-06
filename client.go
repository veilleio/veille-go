package veille

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"mime"
	"net/http"
	"net/url"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	// DefaultBaseURL is the production API origin.
	DefaultBaseURL = "https://api.veille.io"
	// DefaultTimeout bounds every request unless ClientOptions.Timeout says otherwise.
	DefaultTimeout = 120 * time.Second
	// APIKeyEnvironmentVariable is read when NewClient receives an empty API key.
	APIKeyEnvironmentVariable = "VEILLE_API_KEY"
	// DefaultConcurrency is the number of tasks Batch runs at the same time.
	DefaultConcurrency = 10
	// MaxRetries is the default number of retries after a network error or a 408, 429, 500, 502, 503 or 504.
	MaxRetries = 3
	// RetryBackoffFactor multiplies the delay before each new retry (a Retry-After header wins).
	RetryBackoffFactor = 1.5
	// DefaultRetryDelay is the delay before the first retry.
	DefaultRetryDelay = time.Second
	// MaxRetryDelay caps the delay between two attempts.
	MaxRetryDelay = 30 * time.Second
)

// Version is the SDK version sent in the User-Agent header.
const Version = "0.1.0"

// ClientOptions configures a Client. The zero value is valid.
type ClientOptions struct {
	// BaseURL overrides DefaultBaseURL.
	BaseURL string
	// Timeout bounds each request. Zero uses DefaultTimeout; a negative value disables it.
	Timeout time.Duration
	// HTTPClient is copied and used for requests. Redirects are never followed.
	HTTPClient *http.Client
	// Headers are sent with every request.
	Headers map[string]string
	// MaxRetries overrides the number of retries. Zero uses MaxRetries; a negative value disables retries.
	MaxRetries int
	// RetryDelay is the delay before the first retry. Zero uses DefaultRetryDelay.
	RetryDelay time.Duration
}

// Client is a reusable, concurrency-safe Veille API client.
type Client struct {
	apiKey     string
	baseURL    string
	timeout    time.Duration
	headers    map[string]string
	httpClient *http.Client
	maxRetries int
	retryDelay time.Duration
}

// NewClient creates a client. An empty apiKey falls back to the VEILLE_API_KEY environment variable.
func NewClient(apiKey string, options *ClientOptions) (*Client, error) {
	if apiKey == "" {
		apiKey = os.Getenv(APIKeyEnvironmentVariable)
	}
	if strings.TrimSpace(apiKey) == "" {
		return nil, &ConfigurationError{Message: "A non-blank Veille API key is required. Pass it to NewClient or set " + APIKeyEnvironmentVariable + "."}
	}
	if strings.ContainsAny(apiKey, "\r\n") {
		return nil, &ConfigurationError{Message: "The Veille API key cannot be used in an HTTP header."}
	}
	configuration := ClientOptions{}
	if options != nil {
		configuration = *options
	}
	baseURL := strings.TrimRight(configuration.BaseURL, "/")
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	parsed, err := url.Parse(baseURL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return nil, &ConfigurationError{Message: "BaseURL must be an absolute http or https URL."}
	}
	timeout := configuration.Timeout
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	httpClient := &http.Client{}
	if configuration.HTTPClient != nil {
		clone := *configuration.HTTPClient
		httpClient = &clone
	}
	httpClient.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
		return http.ErrUseLastResponse
	}
	headers := make(map[string]string, len(configuration.Headers))
	for name, value := range configuration.Headers {
		headers[name] = value
	}
	maxRetries := configuration.MaxRetries
	if maxRetries == 0 {
		maxRetries = MaxRetries
	}
	retryDelay := configuration.RetryDelay
	if retryDelay <= 0 {
		retryDelay = DefaultRetryDelay
	}
	return &Client{
		apiKey: apiKey, baseURL: baseURL, timeout: timeout, headers: headers, httpClient: httpClient,
		maxRetries: max(maxRetries, 0), retryDelay: retryDelay,
	}, nil
}

// RequestOption customizes one request.
type RequestOption func(*requestOptions)

type requestOptions struct {
	method     string
	headers    map[string]string
	maxRetries *int
}

// WithMaxRetries overrides the number of retries for one request (0 disables them).
func WithMaxRetries(n int) RequestOption {
	return func(o *requestOptions) { o.maxRetries = &n }
}

// WithMethod selects GET or POST. Operations default to GET when they support it.
func WithMethod(method string) RequestOption {
	return func(o *requestOptions) { o.method = strings.ToUpper(method) }
}

// WithHeader adds a header to one request.
func WithHeader(name, value string) RequestOption {
	return func(o *requestOptions) {
		if o.headers == nil {
			o.headers = map[string]string{}
		}
		o.headers[name] = value
	}
}

type param struct {
	name  string
	value any
}

type endpoint struct {
	path     string
	methods  []string
	json     bool
	required []string
}

// Request sends any request, including endpoints released after this SDK version.
// GET parameters travel in the query string; POST parameters travel in a JSON body.
func (c *Client) Request(ctx context.Context, method, path string, params map[string]any, options ...RequestOption) (*Response, error) {
	names := make([]string, 0, len(params))
	for name := range params {
		names = append(names, name)
	}
	sort.Strings(names)
	ordered := make([]param, 0, len(params))
	for _, name := range names {
		ordered = append(ordered, param{name, params[name]})
	}
	opts := collect(options)
	if opts.method == "" {
		opts.method = strings.ToUpper(method)
	}
	return c.send(ctx, opts, path, ordered, false)
}

func (c *Client) call(ctx context.Context, ep endpoint, params []param, options []RequestOption) (*Response, error) {
	opts := collect(options)
	if opts.method == "" {
		opts.method = ep.methods[0]
	}
	if !slices.Contains(ep.methods, opts.method) {
		return nil, &ConfigurationError{Message: fmt.Sprintf("%s supports %s, not %s.", ep.path, strings.Join(ep.methods, ", "), opts.method)}
	}
	for _, name := range ep.required {
		if !slices.ContainsFunc(params, func(p param) bool { return p.name == name }) {
			return nil, &SerializationError{Param: name, Message: fmt.Sprintf("The %q parameter is required.", name)}
		}
	}
	return c.send(ctx, opts, ep.path, params, ep.json)
}

func collect(options []RequestOption) requestOptions {
	var opts requestOptions
	for _, option := range options {
		if option != nil {
			option(&opts)
		}
	}
	return opts
}

func (c *Client) send(ctx context.Context, opts requestOptions, pathTemplate string, params []param, expectJSON bool) (*Response, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if opts.method != http.MethodGet && opts.method != http.MethodPost {
		return nil, &ConfigurationError{Message: fmt.Sprintf("Unsupported HTTP method: %q.", opts.method)}
	}
	if !strings.HasPrefix(pathTemplate, "/") {
		return nil, &ConfigurationError{Message: `The request path must start with "/".`}
	}
	path, rest, err := fillPath(pathTemplate, params)
	if err != nil {
		return nil, err
	}
	target := c.baseURL + path
	var payload []byte
	if opts.method == http.MethodGet {
		query, err := encodeQuery(rest)
		if err != nil {
			return nil, err
		}
		if query != "" {
			target += "?" + query
		}
	} else {
		// The API reads the global parameters from the query string, whatever the method.
		var queryParams []param
		values := make(map[string]any, len(rest))
		for _, p := range rest {
			if p.name == "fields" || p.name == "is_async" || p.name == "soft_fail" {
				queryParams = append(queryParams, p)
				continue
			}
			values[p.name] = p.value
		}
		query, err := encodeQuery(queryParams)
		if err != nil {
			return nil, err
		}
		if query != "" {
			target += "?" + query
		}
		encoded, err := json.Marshal(values)
		if err != nil {
			return nil, &SerializationError{Message: "The request body could not be serialized as JSON.", Err: err}
		}
		payload = encoded
	}
	maxRetries := c.maxRetries
	if opts.maxRetries != nil {
		maxRetries = max(*opts.maxRetries, 0)
	}
	for attempt := 0; ; attempt++ {
		result, retryAfter, err := c.attempt(ctx, opts, target, payload)
		if err == nil {
			if result.StatusCode >= 200 && result.StatusCode <= 299 {
				return c.checkSuccess(result, expectJSON)
			}
			if attempt >= maxRetries || !isRetryable(result) {
				return nil, newAPIError(result)
			}
		} else {
			var transportError *TransportError
			if attempt >= maxRetries || !errors.As(err, &transportError) || ctx.Err() != nil ||
				errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
				return nil, err
			}
		}
		timer := time.NewTimer(retryDelay(c.retryDelay, attempt, retryAfter))
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil, &TransportError{Message: "Veille request was cancelled while waiting to retry", Err: ctx.Err()}
		case <-timer.C:
		}
	}
}

// attempt sends one HTTP request and buffers the response.
func (c *Client) attempt(ctx context.Context, opts requestOptions, target string, payload []byte) (*Response, string, error) {
	if c.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}
	var body io.Reader
	if payload != nil {
		body = bytes.NewReader(payload)
	}
	request, err := http.NewRequestWithContext(ctx, opts.method, target, body)
	if err != nil {
		return nil, "", &ConfigurationError{Message: "The request could not be created.", Err: err}
	}
	request.Header.Set("Accept", "application/json, */*;q=0.8")
	request.Header.Set("User-Agent", "veille-go/"+Version)
	for name, value := range c.headers {
		request.Header.Set(name, value)
	}
	for name, value := range opts.headers {
		request.Header.Set(name, value)
	}
	request.Header.Set("x-api-key", c.apiKey)
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, "", &TransportError{Message: "Veille request failed before an HTTP response was received.", Err: err}
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, "", &TransportError{Message: "Veille response body could not be read.", Err: err}
	}
	result := &Response{
		StatusCode:  response.StatusCode,
		Header:      response.Header.Clone(),
		ContentType: response.Header.Get("Content-Type"),
		Body:        data,
	}
	result.IsJSON = isJSONMediaType(result.ContentType)
	return result, response.Header.Get("Retry-After"), nil
}

func (c *Client) checkSuccess(result *Response, expectJSON bool) (*Response, error) {
	if expectJSON && !result.IsJSON {
		return nil, &ResponseDecodingError{Message: "Veille returned a successful non-JSON response for a JSON operation.", Response: result}
	}
	if result.IsJSON && !json.Valid(result.Body) {
		return nil, &ResponseDecodingError{Message: "Veille returned a successful JSON response that could not be decoded.", Response: result}
	}
	return result, nil
}

var retryableStatuses = map[int]bool{408: true, 500: true, 502: true, 503: true, 504: true}

// isRetryable accepts server errors and plan rate limits: the API reports them as 401 (or 429) with a
// "Rate limit exceeded" message. Per-key quotas ("...: quota daily") do not reset within seconds.
func isRetryable(response *Response) bool {
	if retryableStatuses[response.StatusCode] {
		return true
	}
	if response.StatusCode != http.StatusUnauthorized && response.StatusCode != http.StatusTooManyRequests {
		return false
	}
	message := strings.ToLower(newAPIError(response).Message)
	if strings.Contains(message, "quota") {
		return false
	}
	return strings.Contains(message, "rate limit") || (response.StatusCode == http.StatusTooManyRequests && message == "")
}

func retryDelay(base time.Duration, attempt int, retryAfter string) time.Duration {
	if seconds, err := strconv.ParseFloat(strings.TrimSpace(retryAfter), 64); err == nil && seconds >= 0 {
		return min(time.Duration(seconds*float64(time.Second)), time.Minute)
	}
	delay := float64(base) * math.Pow(RetryBackoffFactor, float64(attempt))
	return min(time.Duration(delay), MaxRetryDelay)
}

func isJSONMediaType(contentType string) bool {
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return false
	}
	return mediaType == "application/json" || strings.HasSuffix(mediaType, "+json")
}

func newAPIError(response *Response) *APIError {
	apiError := &APIError{StatusCode: response.StatusCode, Response: response}
	var object map[string]any
	if json.Unmarshal(response.Body, &object) == nil {
		for _, key := range []string{"error", "message", "error_message", "detail"} {
			if value, ok := object[key].(string); ok && value != "" {
				apiError.Message = value
				break
			}
		}
	}
	return apiError
}

// IsAPIError reports whether err is an *APIError and returns it.
func IsAPIError(err error) (*APIError, bool) {
	var apiError *APIError
	ok := errors.As(err, &apiError)
	return apiError, ok
}
