package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	pp "github.com/rancher/terraform-provider-rancher2/internal/provider/pretty_print"
)

// AuthTransport injects auth tokens and safely logs requests and responses.
type AuthTransport struct {
	Base       http.RoundTripper
	TokenStore *TokenStore
}

// RoundTrip executes a single HTTP transaction.
func (t *AuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()

	// Inject the token, threadsafe due to mutex.
	token := t.TokenStore.GetToken()
	// Clone is important to prevent inherent race conditions.
	reqCopy := req.Clone(ctx)
	if token != "" {
		reqCopy.Header.Set("Authorization", "Bearer "+token)
	}

	// Log Request, checking the size to prevent OOMs
	// Replicates the stream because reading flushes it out.
	t.logRequest(ctx, reqCopy)

	base := t.Base
	if base == nil {
		base = http.DefaultTransport
	}
	resp, err := base.RoundTrip(reqCopy)

	if err != nil {
		logFields := map[string]any{
			"method": reqCopy.Method,
			"url":    reqCopy.URL.String(),
			"error":  err.Error(),
		}
		tflog.Debug(ctx, fmt.Sprintf("HTTP Request Failed: %s", pp.PrettyPrint(logFields)))
		return resp, err
	}

	// Log Response, checking the size and replicating the stream.
	t.logResponse(ctx, reqCopy.URL.String(), resp)

	return resp, nil
}

// Inspect and log the request.
func (t *AuthTransport) logRequest(ctx context.Context, req *http.Request) {
	logFields := map[string]any{
		"method": req.Method,
		"url":    req.URL.String(),
	}

	if req.Body == nil {
		tflog.Debug(ctx, fmt.Sprintf("HTTP Request (No Body): %s", pp.PrettyPrint(logFields)))
		return
	}

	const maxLogSize = 2 << 20 // 2 MB

	if req.ContentLength > maxLogSize || req.ContentLength == -1 {
		logFields["content_length"] = req.ContentLength
		tflog.Debug(ctx, fmt.Sprintf("HTTP Request (Body too large/chunked to log) %s", pp.PrettyPrint(logFields)))
		return
	}

	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		tflog.Warn(ctx, "Failed to read request body for logging", map[string]any{"error": err.Error()})
		return
	}

	marshalledResponse := MarshalledData{
		Data: bodyBytes,
	}

	logFields["body"] = marshalledResponse
	tflog.Debug(ctx, fmt.Sprintf("HTTP Request: %s", pp.PrettyPrint(logFields)))

	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
}

// Inspect and log the response.
func (t *AuthTransport) logResponse(ctx context.Context, url string, resp *http.Response) {
	logFields := map[string]any{
		"url":         url,
		"status_code": resp.StatusCode,
	}

	if resp == nil || resp.Body == nil {
		tflog.Debug(ctx, fmt.Sprintf("HTTP Response (No Body): %s", pp.PrettyPrint(logFields)))
		return
	}

	const maxLogSize = 2 << 20 // 2 MB

	if resp.ContentLength > maxLogSize || resp.ContentLength == -1 {
		logFields["content_length"] = resp.ContentLength
		tflog.Debug(ctx, fmt.Sprintf("HTTP Response (Body too large/chunked to log): %s", pp.PrettyPrint(logFields)))
		return
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		tflog.Warn(ctx, "Failed to read response body for logging", map[string]any{"error": err.Error()})
		return
	}

	// marshalledResponse := MarshalledData{
	//   Data: bodyBytes,
	// }
	logFields["body"] = bodyBytes //marshalledResponse
	tflog.Debug(ctx, fmt.Sprintf("HTTP Response: %s", pp.PrettyPrint(logFields)))

	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
}
