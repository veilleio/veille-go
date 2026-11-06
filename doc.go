// Package veille is the official Go SDK for the Veille API (https://www.veille.io).
//
// Validate emails, domains, IPs, VAT numbers and IBANs in milliseconds to block fake sign-ups, bots and fraud.
//
// Create one client and reuse it:
//
//	client, err := veille.NewClient("", nil) // reads VEILLE_API_KEY
//	response, err := client.IntelligenceEmail(ctx, ...)
//
// Every documented endpoint has a typed method; Client.Request reaches any other path.
package veille
