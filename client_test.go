package veille

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

type recorded struct {
	request *http.Request
	body    []byte
}

func newTestClient(t *testing.T, respond func(*http.Request) *http.Response) (*Client, *[]recorded) {
	t.Helper()
	var calls []recorded
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		var body []byte
		if r.Body != nil {
			body, _ = io.ReadAll(r.Body)
		}
		calls = append(calls, recorded{request: r, body: body})
		return respond(r), nil
	})
	client, err := NewClient("test-key", &ClientOptions{HTTPClient: &http.Client{Transport: transport}})
	if err != nil {
		t.Fatal(err)
	}
	return client, &calls
}

func jsonResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"application/json; charset=utf-8"}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func okJSON(*http.Request) *http.Response { return jsonResponse(200, `{"ok":true}`) }

func TestNewClientRequiresAPIKey(t *testing.T) {
	t.Setenv(APIKeyEnvironmentVariable, "")
	var configurationError *ConfigurationError
	if _, err := NewClient("", nil); !errors.As(err, &configurationError) {
		t.Fatalf("expected ConfigurationError, got %v", err)
	}
	if _, err := NewClient("   ", nil); !errors.As(err, &configurationError) {
		t.Fatalf("expected ConfigurationError, got %v", err)
	}
	if _, err := NewClient("k", &ClientOptions{BaseURL: "ftp://x"}); !errors.As(err, &configurationError) {
		t.Fatalf("expected ConfigurationError for BaseURL, got %v", err)
	}
}

func TestNewClientReadsEnvironment(t *testing.T) {
	t.Setenv(APIKeyEnvironmentVariable, "env-key")
	client, err := NewClient("", nil)
	if err != nil || client.apiKey != "env-key" {
		t.Fatalf("unexpected client %v %v", client, err)
	}
}

func TestFeaturedOperation(t *testing.T) {
	client, calls := newTestClient(t, func(*http.Request) *http.Response {
		return jsonResponse(200, `{"hello":"world"}`)
	})
	ctx := context.Background()
	response, err := client.IntelligenceEmail(ctx, &IntelligenceEmailParams{Query: "jane.doe@example.com"})
	if err != nil {
		t.Fatal(err)
	}
	data, err := response.Map()
	if err != nil || data["hello"] != "world" {
		t.Fatalf("unexpected data %v %v", data, err)
	}
	request := (*calls)[0].request
	if request.Method != "GET" || request.URL.Scheme+"://"+request.URL.Host+request.URL.Path != "https://api.veille.io/v1/intelligence/email" {
		t.Fatalf("unexpected request %s %s", request.Method, request.URL)
	}
	if request.Header.Get("x-api-key") != "test-key" || !strings.HasPrefix(request.Header.Get("User-Agent"), "veille-go/") {
		t.Fatalf("unexpected headers %v", request.Header)
	}
}

func TestQuerySerialization(t *testing.T) {
	client, calls := newTestClient(t, okJSON)
	_, err := client.Request(context.Background(), "GET", "/v1/anything", map[string]any{
		"text": "a b&c", "flag": false, "count": 0, "ratio": 1.5, "whole": 2.0, "list": []string{"x", "y"}, "nothing": nil,
	})
	if err != nil {
		t.Fatal(err)
	}
	got := (*calls)[0].request.URL.RawQuery
	want := "count=0&flag=false&list=x&list=y&ratio=1.5&text=a%20b%26c&whole=2"
	if got != want {
		t.Fatalf("query = %q, want %q", got, want)
	}
}

func TestPostSendsJSONBody(t *testing.T) {
	client, calls := newTestClient(t, okJSON)
	_, err := client.Request(context.Background(), "POST", "/v1/anything", map[string]any{"query": "x", "page": 2, "gone": nil})
	if err != nil {
		t.Fatal(err)
	}
	call := (*calls)[0]
	var body map[string]any
	if err := json.Unmarshal(call.body, &body); err != nil {
		t.Fatal(err)
	}
	if call.request.Method != "POST" || call.request.Header.Get("Content-Type") != "application/json" || body["query"] != "x" || body["page"] != 2.0 || len(body) != 2 {
		t.Fatalf("unexpected POST %s %v", call.request.Method, body)
	}
}

func TestPostKeepsGlobalParametersInQuery(t *testing.T) {
	client, calls := newTestClient(t, okJSON)
	_, err := client.Request(context.Background(), "POST", "/v1/anything", map[string]any{"query": "x", "is_async": true, "fields": "a"})
	if err != nil {
		t.Fatal(err)
	}
	call := (*calls)[0]
	var body map[string]any
	if err := json.Unmarshal(call.body, &body); err != nil {
		t.Fatal(err)
	}
	query := call.request.URL.Query()
	if len(body) != 1 || body["query"] != "x" || query.Get("is_async") != "true" || query.Get("fields") != "a" {
		t.Fatalf("unexpected POST %v ?%s", body, call.request.URL.RawQuery)
	}
}

func TestPathParameters(t *testing.T) {
	client, calls := newTestClient(t, okJSON)
	if _, err := client.Request(context.Background(), "GET", "/v1/items/{id}", map[string]any{"id": "a/b c", "other": 1}); err != nil {
		t.Fatal(err)
	}
	if got := (*calls)[0].request.URL.RawPath; got != "/v1/items/a%2Fb%20c" {
		t.Fatalf("path = %q", got)
	}
	var serializationError *SerializationError
	if _, err := client.Request(context.Background(), "GET", "/v1/items/{id}", nil); !errors.As(err, &serializationError) {
		t.Fatalf("expected SerializationError, got %v", err)
	}
}

func TestRequiredParameters(t *testing.T) {
	client, calls := newTestClient(t, okJSON)
	ctx := context.Background()
	_, err := client.IntelligenceEmail(ctx, &IntelligenceEmailParams{})
	var serializationError *SerializationError
	if !errors.As(err, &serializationError) || len(*calls) != 0 {
		t.Fatalf("expected SerializationError before sending, got %v", err)
	}
}

func TestUnsupportedMethod(t *testing.T) {
	client, _ := newTestClient(t, okJSON)
	ctx := context.Background()
	_, err := client.IntelligenceDomain(ctx, &IntelligenceDomainParams{Query: "test"}, WithMethod("POST"))
	var configurationError *ConfigurationError
	if !errors.As(err, &configurationError) {
		t.Fatalf("expected ConfigurationError, got %v", err)
	}
}

func TestEveryEndpoint(t *testing.T) {
	client, calls := newTestClient(t, func(*http.Request) *http.Response { return jsonResponse(200, `{}`) })
	ctx := context.Background()
	checks := []struct {
		name string
		call func() (*Response, error)
	}{
		{"intelligence_domain", func() (*Response, error) {
			return client.IntelligenceDomain(ctx, &IntelligenceDomainParams{Query: "test"})
		}},
		{"intelligence_email", func() (*Response, error) {
			return client.IntelligenceEmail(ctx, &IntelligenceEmailParams{Query: "john@example.com"})
		}},
		{"intelligence_emailq", func() (*Response, error) {
			return client.IntelligenceEmailq(ctx, &IntelligenceEmailqParams{Query: "john@example.com"})
		}},
		{"intelligence_ip", func() (*Response, error) { return client.IntelligenceIP(ctx, &IntelligenceIPParams{Query: "test"}) }},
		{"intelligence_phone", func() (*Response, error) {
			return client.IntelligencePhone(ctx, &IntelligencePhoneParams{Query: "+33612345678"})
		}},
		{"jobs", func() (*Response, error) { return client.Jobs(ctx, &JobsParams{ID: "test"}) }},
		{"stats", func() (*Response, error) { return client.Stats(ctx, &StatsParams{}) }},
		{"usage", func() (*Response, error) { return client.Usage(ctx) }},
		{"vat", func() (*Response, error) { return client.VAT(ctx, &VATParams{Query: "test"}) }},
		{"vat_iban", func() (*Response, error) { return client.VATIBAN(ctx, &VATIBANParams{Query: "test"}) }},
		{"vat_rates", func() (*Response, error) { return client.VATRates(ctx, &VATRatesParams{}) }},
	}
	for i, check := range checks {
		if _, err := check.call(); err != nil {
			t.Fatalf("%s: %v", check.name, err)
		}
		ep := endpoints[check.name]
		request := (*calls)[i].request
		if request.URL.Path != strings.ReplaceAll(ep.path, "{id}", "test") || request.Method != ep.methods[0] {
			t.Fatalf("%s: unexpected %s %s", check.name, request.Method, request.URL.Path)
		}
	}
	if len(checks) != len(endpoints) {
		t.Fatalf("%d checks for %d endpoints", len(checks), len(endpoints))
	}
}

func TestBinaryResponse(t *testing.T) {
	client, _ := newTestClient(t, func(*http.Request) *http.Response {
		return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": []string{"image/png"}}, Body: io.NopCloser(strings.NewReader("\x89PNG"))}
	})
	response, err := client.Request(context.Background(), "GET", "/v1/anything", nil)
	if err != nil {
		t.Fatal(err)
	}
	if response.IsJSON || string(response.Body) != "\x89PNG" || response.ContentType != "image/png" {
		t.Fatalf("unexpected response %+v", response)
	}
}

func TestNonJSONSuccessForJSONOperation(t *testing.T) {
	client, _ := newTestClient(t, func(*http.Request) *http.Response {
		return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": []string{"text/html"}}, Body: io.NopCloser(strings.NewReader("<html>"))}
	})
	ctx := context.Background()
	_, err := client.IntelligenceEmail(ctx, &IntelligenceEmailParams{Query: "jane.doe@example.com"})
	var decodingError *ResponseDecodingError
	if !errors.As(err, &decodingError) {
		t.Fatalf("expected ResponseDecodingError, got %v", err)
	}
}

func TestAPIError(t *testing.T) {
	client, _ := newTestClient(t, func(*http.Request) *http.Response {
		return jsonResponse(401, `{"error":"Invalid API Key"}`)
	})
	_, err := client.Request(context.Background(), "GET", "/v1/usage", nil)
	apiError, ok := IsAPIError(err)
	if !ok || apiError.StatusCode != 401 || apiError.Message != "Invalid API Key" {
		t.Fatalf("unexpected error %v", err)
	}
}

func TestTransportError(t *testing.T) {
	client, err := NewClient("k", &ClientOptions{MaxRetries: -1, HTTPClient: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return nil, errors.New("network down")
	})}})
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Request(context.Background(), "GET", "/v1/usage", nil)
	var transportError *TransportError
	if !errors.As(err, &transportError) {
		t.Fatalf("expected TransportError, got %v", err)
	}
}

func TestTimeout(t *testing.T) {
	client, err := NewClient("k", &ClientOptions{Timeout: time.Millisecond, HTTPClient: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		<-r.Context().Done()
		return nil, r.Context().Err()
	})}})
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Request(context.Background(), "GET", "/v1/usage", nil)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected deadline exceeded, got %v", err)
	}
}

func TestRetries(t *testing.T) {
	statuses := []int{503, 429, 200}
	var calls int
	client, err := NewClient("k", &ClientOptions{RetryDelay: time.Millisecond, HTTPClient: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		calls++
		if calls == 1 {
			return nil, errors.New("connection reset")
		}
		status := statuses[0]
		statuses = statuses[1:]
		response := jsonResponse(status, `{}`)
		if status == 429 {
			response.Header.Set("Retry-After", "0")
		}
		return response, nil
	})}})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.Request(context.Background(), "POST", "/v1/usage", map[string]any{"q": "x"}); err != nil || calls != 4 {
		t.Fatalf("expected success after 4 calls, got %v after %d", err, calls)
	}
}

func TestRetryLimitsAndClientErrors(t *testing.T) {
	newClient := func(status int, options ...RequestOption) (int, error) {
		var calls int
		client, _ := NewClient("k", &ClientOptions{MaxRetries: 2, RetryDelay: time.Millisecond, HTTPClient: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			calls++
			return jsonResponse(status, `{"error":"x"}`), nil
		})}})
		_, err := client.Request(context.Background(), "GET", "/v1/usage", nil, options...)
		return calls, err
	}
	if calls, err := newClient(503); calls != 3 || err == nil {
		t.Fatalf("503: %d calls, %v", calls, err)
	}
	if calls, _ := newClient(402); calls != 1 {
		t.Fatalf("402 should not be retried, got %d calls", calls)
	}
	if calls, _ := newClient(500, WithMaxRetries(0)); calls != 1 {
		t.Fatalf("WithMaxRetries(0) should disable retries, got %d calls", calls)
	}
}

func TestRateLimitRetries(t *testing.T) {
	cases := []struct {
		message string
		calls   int
	}{
		{"Rate limit exceeded: You've exceeded the rate limit on your subscription", 2},
		{"Rate limit exceeded for the API key: quota daily", 1},
		{"Invalid API Key", 1},
	}
	for _, c := range cases {
		calls := 0
		client, _ := NewClient("k", &ClientOptions{RetryDelay: time.Millisecond, HTTPClient: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			calls++
			if calls == 1 {
				body, _ := json.Marshal(map[string]string{"error": c.message})
				return jsonResponse(401, string(body)), nil
			}
			return jsonResponse(200, `{}`), nil
		})}})
		_, _ = client.Request(context.Background(), "GET", "/v1/usage", nil)
		if calls != c.calls {
			t.Fatalf("%q: %d calls, want %d", c.message, calls, c.calls)
		}
	}
}

func TestGlobalParameters(t *testing.T) {
	client, calls := newTestClient(t, func(*http.Request) *http.Response { return jsonResponse(200, `{"job_id":"j1"}`) })
	ctx := context.Background()
	if _, err := client.IntelligenceEmail(ctx, &IntelligenceEmailParams{Query: "jane.doe@example.com", Fields: String("a,b.c"), IsAsync: Bool(true), SoftFail: Bool(false)}); err != nil {
		t.Fatal(err)
	}
	query := (*calls)[0].request.URL.Query()
	if query.Get("fields") != "a,b.c" || query.Get("is_async") != "true" || query.Get("soft_fail") != "false" {
		t.Fatalf("unexpected query %v", query)
	}
}

func TestBatch(t *testing.T) {
	var mu sync.Mutex
	running, peak := 0, 0
	results := Batch(context.Background(), []int{1, 2, 3, 4, 5, 6}, 2, func(_ context.Context, n int) (int, error) {
		mu.Lock()
		running++
		peak = max(peak, running)
		mu.Unlock()
		time.Sleep(5 * time.Millisecond)
		mu.Lock()
		running--
		mu.Unlock()
		if n == 4 {
			return 0, errors.New("boom")
		}
		return n * 10, nil
	})
	if peak != 2 {
		t.Fatalf("peak concurrency %d", peak)
	}
	for i, r := range results {
		if r.Index != i || (i == 3) != (r.Err != nil) || (r.Err == nil && r.Value != (i+1)*10) {
			t.Fatalf("unexpected result %d: %+v", i, r)
		}
	}
}

func TestWaitForJob(t *testing.T) {
	statuses := []string{"Queued", "Running", "CompletedJob"}
	client, calls := newTestClient(t, func(*http.Request) *http.Response {
		status := statuses[0]
		statuses = statuses[1:]
		return jsonResponse(200, `{"status":"`+status+`","progress":50}`)
	})
	seen := 0
	response, err := client.WaitForJob(context.Background(), "abc", &WaitForJobOptions{Interval: time.Millisecond, OnProgress: func(map[string]any) { seen++ }})
	if err != nil {
		t.Fatal(err)
	}
	job, _ := response.Map()
	if job["status"] != "CompletedJob" || len(*calls) != 3 || seen != 3 || (*calls)[0].request.URL.Path != "/v1/jobs/abc" {
		t.Fatalf("unexpected job %v after %d calls", job, len(*calls))
	}
}

func TestWaitForJobTimesOut(t *testing.T) {
	client, _ := newTestClient(t, func(*http.Request) *http.Response { return jsonResponse(200, `{"status":"Running"}`) })
	_, err := client.WaitForJob(context.Background(), "abc", &WaitForJobOptions{Interval: 10 * time.Millisecond, Timeout: time.Millisecond})
	var timeoutError *JobTimeoutError
	if !errors.As(err, &timeoutError) {
		t.Fatalf("expected JobTimeoutError, got %v", err)
	}
}
