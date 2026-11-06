# Veille Go SDK

The official Go SDK for the [Veille](https://www.veille.io) fraud prevention API.

Validate emails, domains, IPs, VAT numbers and IBANs in milliseconds to block fake sign-ups, bots and fraud.

The SDK exposes **11 typed methods**, one per documented endpoint, plus a generic `Request()`. It has no dependencies outside the standard library.

## Requirements

- Go 1.22 or newer

## Installation

```sh
go get github.com/veilleio/veille-go
```

## Get your API key

[Sign up for Veille](https://app.veille.io/register) or sign in, then copy your API key from the [dashboard](https://app.veille.io). `NewClient("", nil)` reads the `VEILLE_API_KEY` environment variable.

## Quick start

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/veilleio/veille-go"
)

func main() {
	client, err := veille.NewClient("", nil) // reads VEILLE_API_KEY
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	response, err := client.IntelligenceEmail(ctx, &veille.IntelligenceEmailParams{Query: "jane.doe@example.com"})
	if err != nil {
		log.Fatal(err)
	}
	data, err := response.Map()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(data)
}
```

Parameter structs mirror the [API reference](https://app.veille.io/openapi). Required fields are plain values; optional fields are pointers — use `veille.String`, `veille.Float`, `veille.Int` and `veille.Bool`.

Decode into your own types with `response.Decode(&value)`.

## Async jobs

Any call accepts `is_async: true` (useful past the 60-second server limit): the response is then a `job_id`. `WaitForJob` polls `GET /v1/jobs/{id}` until a terminal status (`CompletedJob`, `Error`, `NotFound`, `BadRequest`):

```go
submitted, err := client.IntelligenceEmail(ctx, &veille.IntelligenceEmailParams{Query: "jane.doe@example.com", IsAsync: veille.Bool(true)})
if err != nil {
	log.Fatal(err)
}
var accepted struct{ JobID string `json:"job_id"` }
if err := submitted.Decode(&accepted); err != nil {
	log.Fatal(err)
}
job, err := client.WaitForJob(ctx, accepted.JobID, &veille.WaitForJobOptions{
	OnProgress: func(job map[string]any) { fmt.Println(job["status"], job["progress"]) },
})
```

## GET or POST

Methods use GET when the endpoint supports it (parameters in the query string); `WithMethod("POST")` sends them as a JSON body.

## Any endpoint

```go
response, err := client.Request(ctx, "GET", "/v1/usage", nil)
```

## Errors

```go
response, err := client.IntelligenceEmail(ctx, &veille.IntelligenceEmailParams{Query: "jane.doe@example.com"})
if apiError, ok := veille.IsAPIError(err); ok {
	log.Printf("HTTP %d: %s", apiError.StatusCode, apiError.Message)
}
```

| Error | When |
| --- | --- |
| `*ConfigurationError` | Missing API key, invalid options or unsupported method. |
| `*SerializationError` | A required parameter is missing or a value can't be sent. |
| `*TransportError` | Network failure, timeout or cancelled context (`errors.Is(err, context.DeadlineExceeded)`). |
| `*APIError` | Non-2xx status (`StatusCode`, `Message`, `Response`). |
| `*ResponseDecodingError` | A success response could not be decoded. |

## Retries

Requests are retried up to 3 times after a network error, a `408`, `500`, `502`, `503` or `504`, or a rate limit — which the API returns as `401` (or `429`) with a `Rate limit exceeded` message. Per-key quota errors (`…: quota daily`) are not retried, nor are `400`, `402`, `404` or client-side timeouts. The delay starts at 1 second and grows 1.5× per retry (capped at 30 seconds); a `Retry-After` header wins.

```go
client, err := veille.NewClient("", &veille.ClientOptions{
	MaxRetries: 5,                      // negative disables retries
	RetryDelay: 500 * time.Millisecond,
})

// or for one request
response, err := client.IntelligenceEmail(ctx, &veille.IntelligenceEmailParams{Query: "jane.doe@example.com"}, veille.WithMaxRetries(0))
```

## Batch requests

`veille.Batch` runs a function per item with at most 10 calls in flight and returns one `BatchResult` per item, in input order. A failing item never stops the others.

```go
inputs := []string{"jane.doe@example.com", "john@mailinator.com", "contact@veille.io"}
results := veille.Batch(ctx, inputs, 5, func(ctx context.Context, value string) (*veille.Response, error) {
	return client.IntelligenceEmail(ctx, &veille.IntelligenceEmailParams{Query: value})
})
for _, result := range results {
	if result.Err != nil {
		log.Printf("%s: %v", inputs[result.Index], result.Err)
		continue
	}
	fmt.Println(string(result.Value.Body))
}
```

| Constant | Default | Description |
| --- | --- | --- |
| `DefaultConcurrency` | `10` | Tasks run at the same time by batch |
| `DefaultTimeout` | `120 * time.Second` | Request timeout |
| `MaxRetries` | `3` | Retry attempts after a network error, a 408/5xx or a rate limit |
| `RetryBackoffFactor` | `1.5` | Exponential backoff multiplier |
| `DefaultRetryDelay` | `time.Second` | Delay before the first retry |

## Timeouts and custom HTTP clients

```go
client, err := veille.NewClient("", &veille.ClientOptions{
	Timeout:    60 * time.Second, // default 120 s; negative disables
	HTTPClient: &http.Client{Transport: myTransport},
})
```

Every method also honours the deadline and cancellation of the `context.Context` you pass.

## Endpoints

11 methods, generated from the [OpenAPI document](https://app.veille.io/openapi).

| Method | HTTP | Description |
| --- | --- | --- |
| `IntelligenceDomain` | GET `/v1/intelligence/domain` | Check any domain instantly. Detects disposable email providers, public webmail, relay services, and education domains. Returns DNS health w… |
| `IntelligenceEmail` | GET `/v1/intelligence/email` | Validate any email address in a single call. Detects disposable and temporary inboxes, identifies role accounts, public webmail, relay serv… |
| `IntelligenceEmailq` | GET `/v1/intelligence/emailq` | Deep email deliverability verification. Performs full Email Validation plus SMTP inbox check, catch-all detection, and returns a qualificat… |
| `IntelligenceIP` | GET `/v1/intelligence/ip` | Assess the reputation of any IP address. Detects VPNs, proxies, Tor exit nodes, datacenter hosting, abusive traffic, and bots. Returns geol… |
| `IntelligencePhone` | GET `/v1/intelligence/phone` | Validate any phone number in E.164 format. Detects disposable and virtual numbers, identifies carrier and line type (mobile/fixed/VoIP), re… |
| `Jobs` | GET `/v1/jobs/{id}` | Poll the status, progress (0–100) and result of a job started with is_async=true. Polling is free. Poll until status is terminal: Completed… |
| `Stats` | GET `/v1/stats` | Infrastructure health for an engine or endpoint over the last hour. Not scoped to your account or tenant. Enabled per account: HTTP 403 mea… |
| `Usage` | GET `/v1/usage` | Retrieve real-time usage statistics and credit balance for your Veille account. |
| `VAT` | GET `/v1/vat` | Validate a European VAT number against the VIES database. Returns company name, registered address, and validity status for any EU member s… |
| `VATIBAN` | GET `/v1/vat/iban` | Validate any IBAN and retrieve associated bank details. Returns the BIC/SWIFT code, bank name, country, BBAN structural breakdown, and SEPA… |
| `VATRates` | GET `/v1/vat/rates` | Get current ECB exchange rates for any supported base currency. Optionally filter to specific target currencies. Updated daily from the Eur… |

## All Veille SDKs

| SDK | Package | Install |
| --- | --- | --- |
| Node.js | [`@veilleio/sdk`](https://www.npmjs.com/package/@veilleio/sdk) | `npm install @veilleio/sdk` |
| Python | [`veille`](https://pypi.org/project/veille/) | `pip install veille` |
| Go | [`veille-go`](https://pkg.go.dev/github.com/veilleio/veille-go) | `go get github.com/veilleio/veille-go` |
| .NET | [`Veille`](https://www.nuget.org/packages/Veille) | `dotnet add package Veille` |
| Ruby | [`veille`](https://rubygems.org/gems/veille) | `gem install veille` |
| PHP | [`veille/sdk`](https://packagist.org/packages/veille/sdk) | `composer require veille/sdk` |
| Java | [`veille-java`](https://github.com/veilleio/veille-java) | build from source, see its README |
| CLI | [`@veilleio/cli`](https://www.npmjs.com/package/@veilleio/cli) | `npm install -g @veilleio/cli` |
| MCP server | [`@veilleio/mcp`](https://www.npmjs.com/package/@veilleio/mcp) | `npx -y @veilleio/mcp` |

## License

MIT — see [LICENSE](LICENSE).
