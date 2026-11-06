// Code generated from the Veille OpenAPI document. DO NOT EDIT.

package veille

import "context"

// IntelligenceDomain — Domain Validation. Check any domain instantly. Detects disposable email providers, public webmail, relay services, and education domains. Returns DNS health with MX provider detection, domain history, and optional SSL/TLS certificate transparency. Credit cost: 2 credits per call.
//
// GET /v1/intelligence/domain
func (c *Client) IntelligenceDomain(ctx context.Context, params *IntelligenceDomainParams, options ...RequestOption) (*Response, error) {
	return c.call(ctx, endpoints["intelligence_domain"], params.params(), options)
}

// IntelligenceEmail — Email Validation. Validate any email address in a single call. Detects disposable and temporary inboxes, identifies role accounts, public webmail, relay services, and returns a 0-100 risk score. For deep deliverability verification (SMTP check, catch-all de… Credit cost: 3 credits per call.
//
// GET /v1/intelligence/email
func (c *Client) IntelligenceEmail(ctx context.Context, params *IntelligenceEmailParams, options ...RequestOption) (*Response, error) {
	return c.call(ctx, endpoints["intelligence_email"], params.params(), options)
}

// IntelligenceEmailq — Email Qualification. Deep email deliverability verification. Performs full Email Validation plus SMTP inbox check, catch-all detection, and returns a qualification verdict (valid, invalid, risky, disposable, unknown). Credit cost: 5 credits per call.
//
// GET /v1/intelligence/emailq
func (c *Client) IntelligenceEmailq(ctx context.Context, params *IntelligenceEmailqParams, options ...RequestOption) (*Response, error) {
	return c.call(ctx, endpoints["intelligence_emailq"], params.params(), options)
}

// IntelligenceIP — IP Reputation. Assess the reputation of any IP address. Detects VPNs, proxies, Tor exit nodes, datacenter hosting, abusive traffic, and bots. Returns geolocation, network owner, and a 0-100 threat score. Credit cost: 2 credits per call.
//
// GET /v1/intelligence/ip
func (c *Client) IntelligenceIP(ctx context.Context, params *IntelligenceIPParams, options ...RequestOption) (*Response, error) {
	return c.call(ctx, endpoints["intelligence_ip"], params.params(), options)
}

// IntelligencePhone — Phone Validation. Validate any phone number in E.164 format. Detects disposable and virtual numbers, identifies carrier and line type (mobile/fixed/VoIP), returns country and timezone, and computes a 0-100 risk score. Credit cost: 1 credit per call.
//
// GET /v1/intelligence/phone
func (c *Client) IntelligencePhone(ctx context.Context, params *IntelligencePhoneParams, options ...RequestOption) (*Response, error) {
	return c.call(ctx, endpoints["intelligence_phone"], params.params(), options)
}

// Jobs — Jobs. Poll the status, progress (0–100) and result of a job started with is_async=true. Polling is free. Poll until status is terminal: CompletedJob, Error, NotFound or BadRequest.
//
// GET /v1/jobs/{id}
func (c *Client) Jobs(ctx context.Context, params *JobsParams, options ...RequestOption) (*Response, error) {
	return c.call(ctx, endpoints["jobs"], params.params(), options)
}

// Stats — Stats. Infrastructure health for an engine or endpoint over the last hour. Not scoped to your account or tenant. Enabled per account: HTTP 403 means contact support to activate it.
//
// GET /v1/stats
func (c *Client) Stats(ctx context.Context, params *StatsParams, options ...RequestOption) (*Response, error) {
	return c.call(ctx, endpoints["stats"], params.params(), options)
}

// Usage — Usage. Retrieve real-time usage statistics and credit balance for your Veille account.
//
// GET /v1/usage
func (c *Client) Usage(ctx context.Context, options ...RequestOption) (*Response, error) {
	return c.call(ctx, endpoints["usage"], nil, options)
}

// VAT — VAT Validation. Validate a European VAT number against the VIES database. Returns company name, registered address, and validity status for any EU member state. Credit cost: 1 credit per call.
//
// GET /v1/vat
func (c *Client) VAT(ctx context.Context, params *VATParams, options ...RequestOption) (*Response, error) {
	return c.call(ctx, endpoints["vat"], params.params(), options)
}

// VATIBAN — IBAN Validation. Validate any IBAN and retrieve associated bank details. Returns the BIC/SWIFT code, bank name, country, BBAN structural breakdown, and SEPA zone membership. Credit cost: 1 credit per call.
//
// GET /v1/vat/iban
func (c *Client) VATIBAN(ctx context.Context, params *VATIBANParams, options ...RequestOption) (*Response, error) {
	return c.call(ctx, endpoints["vat_iban"], params.params(), options)
}

// VATRates — Exchange Rates. Get current ECB exchange rates for any supported base currency. Optionally filter to specific target currencies. Updated daily from the European Central Bank. Credit cost: 1 credit per call.
//
// GET /v1/vat/rates
func (c *Client) VATRates(ctx context.Context, params *VATRatesParams, options ...RequestOption) (*Response, error) {
	return c.call(ctx, endpoints["vat_rates"], params.params(), options)
}

// endpoints lists every documented operation by its snake_case name.
var endpoints = map[string]endpoint{
	"intelligence_domain": {path: "/v1/intelligence/domain", methods: []string{"GET"}, json: true, required: []string{"query"}},
	"intelligence_email":  {path: "/v1/intelligence/email", methods: []string{"GET"}, json: true, required: []string{"query"}},
	"intelligence_emailq": {path: "/v1/intelligence/emailq", methods: []string{"GET"}, json: true, required: []string{"query"}},
	"intelligence_ip":     {path: "/v1/intelligence/ip", methods: []string{"GET"}, json: true, required: []string{"query"}},
	"intelligence_phone":  {path: "/v1/intelligence/phone", methods: []string{"GET"}, json: true, required: []string{"query"}},
	"jobs":                {path: "/v1/jobs/{id}", methods: []string{"GET"}, json: true, required: []string{"id"}},
	"stats":               {path: "/v1/stats", methods: []string{"GET"}, json: true, required: []string{}},
	"usage":               {path: "/v1/usage", methods: []string{"GET"}, json: true, required: []string{}},
	"vat":                 {path: "/v1/vat", methods: []string{"GET"}, json: true, required: []string{"query"}},
	"vat_iban":            {path: "/v1/vat/iban", methods: []string{"GET"}, json: true, required: []string{"query"}},
	"vat_rates":           {path: "/v1/vat/rates", methods: []string{"GET"}, json: true, required: []string{}},
}
