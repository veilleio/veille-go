// Code generated from the Veille OpenAPI document. DO NOT EDIT.

package veille

// IntelligenceDomainParams holds the parameters of IntelligenceDomain (GET /v1/intelligence/domain).
// Required fields are plain values; optional fields are pointers (use veille.String, veille.Float, veille.Bool) or nil slices and maps.
type IntelligenceDomainParams struct {
	// The domain to validate. Accepts raw domains (tempmail.org), full URLs (https://www.tempmail.org/page), or domains with www. prefix, all are normalized automatically.
	Query string `json:"query" api:"query,required"`
	// Enable Certificate Transparency lookup. Returns SSL/TLS certificate history, current certificate, issuers, and all Subject Alternative Names (SANs) ever associated with the domain. Adds 1-5 seconds of latency. Default: false.
	CertificateTransparency *bool `json:"certificate_transparency" api:"certificate_transparency"`
	// Comma-separated list of fields to keep in the response, with dot notation for nested fields (e.g. name,location.city). Does not change the credit cost.
	Fields *string `json:"fields" api:"fields"`
	// Process the request asynchronously: the response is a job_id to poll on /v1/jobs/{id} (see waitForJob). Default: false.
	IsAsync *bool `json:"is_async" api:"is_async"`
	// Return HTTP 200 even when an error occurs; the error is then in the response body. Default: false.
	SoftFail *bool `json:"soft_fail" api:"soft_fail"`
}

func (p *IntelligenceDomainParams) params() []param {
	if p == nil {
		return nil
	}
	values := make([]param, 0, 5)
	if p.Query != "" {
		values = append(values, param{"query", p.Query})
	}
	if p.CertificateTransparency != nil {
		values = append(values, param{"certificate_transparency", *p.CertificateTransparency})
	}
	if p.Fields != nil {
		values = append(values, param{"fields", *p.Fields})
	}
	if p.IsAsync != nil {
		values = append(values, param{"is_async", *p.IsAsync})
	}
	if p.SoftFail != nil {
		values = append(values, param{"soft_fail", *p.SoftFail})
	}
	return values
}

// IntelligenceEmailParams holds the parameters of IntelligenceEmail (GET /v1/intelligence/email).
// Required fields are plain values; optional fields are pointers (use veille.String, veille.Float, veille.Bool) or nil slices and maps.
type IntelligenceEmailParams struct {
	// The email address to validate (e.g. john@example.com). Automatically normalized to lowercase.
	Query string `json:"query" api:"query,required"`
	// Comma-separated list of fields to keep in the response, with dot notation for nested fields (e.g. name,location.city). Does not change the credit cost.
	Fields *string `json:"fields" api:"fields"`
	// Process the request asynchronously: the response is a job_id to poll on /v1/jobs/{id} (see waitForJob). Default: false.
	IsAsync *bool `json:"is_async" api:"is_async"`
	// Return HTTP 200 even when an error occurs; the error is then in the response body. Default: false.
	SoftFail *bool `json:"soft_fail" api:"soft_fail"`
}

func (p *IntelligenceEmailParams) params() []param {
	if p == nil {
		return nil
	}
	values := make([]param, 0, 4)
	if p.Query != "" {
		values = append(values, param{"query", p.Query})
	}
	if p.Fields != nil {
		values = append(values, param{"fields", *p.Fields})
	}
	if p.IsAsync != nil {
		values = append(values, param{"is_async", *p.IsAsync})
	}
	if p.SoftFail != nil {
		values = append(values, param{"soft_fail", *p.SoftFail})
	}
	return values
}

// IntelligenceEmailqParams holds the parameters of IntelligenceEmailq (GET /v1/intelligence/emailq).
// Required fields are plain values; optional fields are pointers (use veille.String, veille.Float, veille.Bool) or nil slices and maps.
type IntelligenceEmailqParams struct {
	// The email address to qualify (e.g. john@example.com). Automatically normalized to lowercase.
	Query string `json:"query" api:"query,required"`
	// Comma-separated list of fields to keep in the response, with dot notation for nested fields (e.g. name,location.city). Does not change the credit cost.
	Fields *string `json:"fields" api:"fields"`
	// Process the request asynchronously: the response is a job_id to poll on /v1/jobs/{id} (see waitForJob). Default: false.
	IsAsync *bool `json:"is_async" api:"is_async"`
	// Return HTTP 200 even when an error occurs; the error is then in the response body. Default: false.
	SoftFail *bool `json:"soft_fail" api:"soft_fail"`
}

func (p *IntelligenceEmailqParams) params() []param {
	if p == nil {
		return nil
	}
	values := make([]param, 0, 4)
	if p.Query != "" {
		values = append(values, param{"query", p.Query})
	}
	if p.Fields != nil {
		values = append(values, param{"fields", *p.Fields})
	}
	if p.IsAsync != nil {
		values = append(values, param{"is_async", *p.IsAsync})
	}
	if p.SoftFail != nil {
		values = append(values, param{"soft_fail", *p.SoftFail})
	}
	return values
}

// IntelligenceIPParams holds the parameters of IntelligenceIP (GET /v1/intelligence/ip).
// Required fields are plain values; optional fields are pointers (use veille.String, veille.Float, veille.Bool) or nil slices and maps.
type IntelligenceIPParams struct {
	// The IP address to analyze (IPv4 or IPv6).
	Query string `json:"query" api:"query,required"`
	// Comma-separated list of fields to keep in the response, with dot notation for nested fields (e.g. name,location.city). Does not change the credit cost.
	Fields *string `json:"fields" api:"fields"`
	// Process the request asynchronously: the response is a job_id to poll on /v1/jobs/{id} (see waitForJob). Default: false.
	IsAsync *bool `json:"is_async" api:"is_async"`
	// Return HTTP 200 even when an error occurs; the error is then in the response body. Default: false.
	SoftFail *bool `json:"soft_fail" api:"soft_fail"`
}

func (p *IntelligenceIPParams) params() []param {
	if p == nil {
		return nil
	}
	values := make([]param, 0, 4)
	if p.Query != "" {
		values = append(values, param{"query", p.Query})
	}
	if p.Fields != nil {
		values = append(values, param{"fields", *p.Fields})
	}
	if p.IsAsync != nil {
		values = append(values, param{"is_async", *p.IsAsync})
	}
	if p.SoftFail != nil {
		values = append(values, param{"soft_fail", *p.SoftFail})
	}
	return values
}

// IntelligencePhoneParams holds the parameters of IntelligencePhone (GET /v1/intelligence/phone).
// Required fields are plain values; optional fields are pointers (use veille.String, veille.Float, veille.Bool) or nil slices and maps.
type IntelligencePhoneParams struct {
	// The phone number to validate in E.164 format (must start with +, e.g. +33612345678).
	Query string `json:"query" api:"query,required"`
	// Comma-separated list of fields to keep in the response, with dot notation for nested fields (e.g. name,location.city). Does not change the credit cost.
	Fields *string `json:"fields" api:"fields"`
	// Process the request asynchronously: the response is a job_id to poll on /v1/jobs/{id} (see waitForJob). Default: false.
	IsAsync *bool `json:"is_async" api:"is_async"`
	// Return HTTP 200 even when an error occurs; the error is then in the response body. Default: false.
	SoftFail *bool `json:"soft_fail" api:"soft_fail"`
}

func (p *IntelligencePhoneParams) params() []param {
	if p == nil {
		return nil
	}
	values := make([]param, 0, 4)
	if p.Query != "" {
		values = append(values, param{"query", p.Query})
	}
	if p.Fields != nil {
		values = append(values, param{"fields", *p.Fields})
	}
	if p.IsAsync != nil {
		values = append(values, param{"is_async", *p.IsAsync})
	}
	if p.SoftFail != nil {
		values = append(values, param{"soft_fail", *p.SoftFail})
	}
	return values
}

// JobsParams holds the parameters of Jobs (GET /v1/jobs/{id}).
// Required fields are plain values; optional fields are pointers (use veille.String, veille.Float, veille.Bool) or nil slices and maps.
type JobsParams struct {
	// The job_id returned by an endpoint called with is_async=true.
	ID string `json:"id" api:"id,required"`
}

func (p *JobsParams) params() []param {
	if p == nil {
		return nil
	}
	values := make([]param, 0, 1)
	if p.ID != "" {
		values = append(values, param{"id", p.ID})
	}
	return values
}

// StatsParams holds the parameters of Stats (GET /v1/stats).
// Required fields are plain values; optional fields are pointers (use veille.String, veille.Float, veille.Bool) or nil slices and maps.
type StatsParams struct {
	// API family to check, for example vat. Provide at least engine or path.
	Engine *string `json:"engine" api:"engine"`
	// Specific endpoint to check, for example /v1/vat. If you also pass engine, use an endpoint from that API.
	Path *string `json:"path" api:"path"`
}

func (p *StatsParams) params() []param {
	if p == nil {
		return nil
	}
	values := make([]param, 0, 2)
	if p.Engine != nil {
		values = append(values, param{"engine", *p.Engine})
	}
	if p.Path != nil {
		values = append(values, param{"path", *p.Path})
	}
	return values
}

// VATParams holds the parameters of VAT (GET /v1/vat).
// Required fields are plain values; optional fields are pointers (use veille.String, veille.Float, veille.Bool) or nil slices and maps.
type VATParams struct {
	// The VAT number to validate (e.g. FR40303265045).
	Query string `json:"query" api:"query,required"`
	// Comma-separated list of fields to keep in the response, with dot notation for nested fields (e.g. name,location.city). Does not change the credit cost.
	Fields *string `json:"fields" api:"fields"`
	// Process the request asynchronously: the response is a job_id to poll on /v1/jobs/{id} (see waitForJob). Default: false.
	IsAsync *bool `json:"is_async" api:"is_async"`
	// Return HTTP 200 even when an error occurs; the error is then in the response body. Default: false.
	SoftFail *bool `json:"soft_fail" api:"soft_fail"`
}

func (p *VATParams) params() []param {
	if p == nil {
		return nil
	}
	values := make([]param, 0, 4)
	if p.Query != "" {
		values = append(values, param{"query", p.Query})
	}
	if p.Fields != nil {
		values = append(values, param{"fields", *p.Fields})
	}
	if p.IsAsync != nil {
		values = append(values, param{"is_async", *p.IsAsync})
	}
	if p.SoftFail != nil {
		values = append(values, param{"soft_fail", *p.SoftFail})
	}
	return values
}

// VATIBANParams holds the parameters of VATIBAN (GET /v1/vat/iban).
// Required fields are plain values; optional fields are pointers (use veille.String, veille.Float, veille.Bool) or nil slices and maps.
type VATIBANParams struct {
	// The IBAN to validate (e.g. FR7630006000011234567890189).
	Query string `json:"query" api:"query,required"`
	// Comma-separated list of fields to keep in the response, with dot notation for nested fields (e.g. name,location.city). Does not change the credit cost.
	Fields *string `json:"fields" api:"fields"`
	// Process the request asynchronously: the response is a job_id to poll on /v1/jobs/{id} (see waitForJob). Default: false.
	IsAsync *bool `json:"is_async" api:"is_async"`
	// Return HTTP 200 even when an error occurs; the error is then in the response body. Default: false.
	SoftFail *bool `json:"soft_fail" api:"soft_fail"`
}

func (p *VATIBANParams) params() []param {
	if p == nil {
		return nil
	}
	values := make([]param, 0, 4)
	if p.Query != "" {
		values = append(values, param{"query", p.Query})
	}
	if p.Fields != nil {
		values = append(values, param{"fields", *p.Fields})
	}
	if p.IsAsync != nil {
		values = append(values, param{"is_async", *p.IsAsync})
	}
	if p.SoftFail != nil {
		values = append(values, param{"soft_fail", *p.SoftFail})
	}
	return values
}

// VATRatesParams holds the parameters of VATRates (GET /v1/vat/rates).
// Required fields are plain values; optional fields are pointers (use veille.String, veille.Float, veille.Bool) or nil slices and maps.
type VATRatesParams struct {
	// Base currency ISO 4217 code (defaults to EUR). Default: "EUR".
	Base *string `json:"base" api:"base"`
	// Comma-separated target currency codes to filter (e.g. USD,GBP,JPY).
	Symbols *string `json:"symbols" api:"symbols"`
	// Comma-separated list of fields to keep in the response, with dot notation for nested fields (e.g. name,location.city). Does not change the credit cost.
	Fields *string `json:"fields" api:"fields"`
	// Process the request asynchronously: the response is a job_id to poll on /v1/jobs/{id} (see waitForJob). Default: false.
	IsAsync *bool `json:"is_async" api:"is_async"`
	// Return HTTP 200 even when an error occurs; the error is then in the response body. Default: false.
	SoftFail *bool `json:"soft_fail" api:"soft_fail"`
}

func (p *VATRatesParams) params() []param {
	if p == nil {
		return nil
	}
	values := make([]param, 0, 5)
	if p.Base != nil {
		values = append(values, param{"base", *p.Base})
	}
	if p.Symbols != nil {
		values = append(values, param{"symbols", *p.Symbols})
	}
	if p.Fields != nil {
		values = append(values, param{"fields", *p.Fields})
	}
	if p.IsAsync != nil {
		values = append(values, param{"is_async", *p.IsAsync})
	}
	if p.SoftFail != nil {
		values = append(values, param{"soft_fail", *p.SoftFail})
	}
	return values
}
