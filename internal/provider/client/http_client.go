package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ Client = &HttpClient{}
var _ retryablehttp.LeveledLogger = &retryableHTTPLogger{}

type retryableHTTPLogger struct {
	ctx context.Context
}

func (l *retryableHTTPLogger) toMap(keysAndValues ...any) map[string]any {
	if len(keysAndValues) == 0 {
		return nil
	}
	if len(keysAndValues)%2 != 0 {
		keysAndValues = append(keysAndValues, "MISSING")
	}
	fields := make(map[string]any, len(keysAndValues)/2)
	for i := 0; i < len(keysAndValues); i += 2 {
		k, ok := keysAndValues[i].(string)
		if !ok {
			k = fmt.Sprintf("%v", keysAndValues[i])
		}
		fields[k] = keysAndValues[i+1]
	}
	return fields
}

func (l *retryableHTTPLogger) Error(msg string, keysAndValues ...any) {
	if fields := l.toMap(keysAndValues...); fields != nil {
		tflog.Error(l.ctx, msg, fields)
	} else {
		tflog.Error(l.ctx, msg)
	}
}

func (l *retryableHTTPLogger) Info(msg string, keysAndValues ...any) {
	if fields := l.toMap(keysAndValues...); fields != nil {
		tflog.Info(l.ctx, msg, fields)
	} else {
		tflog.Info(l.ctx, msg)
	}
}

func (l *retryableHTTPLogger) Debug(msg string, keysAndValues ...any) {
	if fields := l.toMap(keysAndValues...); fields != nil {
		tflog.Debug(l.ctx, msg, fields)
	} else {
		tflog.Debug(l.ctx, msg)
	}
}

func (l *retryableHTTPLogger) Warn(msg string, keysAndValues ...any) {
	if fields := l.toMap(keysAndValues...); fields != nil {
		tflog.Warn(l.ctx, msg, fields)
	} else {
		tflog.Warn(l.ctx, msg)
	}
}

type HttpClient struct {
	ctx            context.Context
	apiURL         string
	caCert         string
	ignoreSystemCA bool
	insecure       bool
	maxRedirects   int64
	timeout        time.Duration
	token          string
	client         *retryablehttp.Client
}

func NewHttpClient(ctx context.Context, apiURL, caCert string, insecure, ignoreSystemCA bool, timeout time.Duration, maxRedirects int64, token string) *HttpClient {
	var rootCAs *x509.CertPool
	if ignoreSystemCA {
		rootCAs = x509.NewCertPool()
	} else {
		// Get the SystemCertPool, continue with an empty pool on error
		rootCAs, _ = x509.SystemCertPool()
		if rootCAs == nil {
			rootCAs = x509.NewCertPool()
		}
	}
	if caCert != "" {
		// Append our cert to the cert pool
		if ok := rootCAs.AppendCertsFromPEM([]byte(caCert)); !ok {
			tflog.Warn(ctx, "no certs appended, using system certs only")
		}
	}
	retryClient := retryablehttp.NewClient()
	retryClient.Logger = &retryableHTTPLogger{ctx: ctx}
	retryClient.HTTPClient = &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: insecure,
				RootCAs:            rootCAs,
			},
		},
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			if len(via) >= int(maxRedirects) {
				return fmt.Errorf("stopped after %d redirects", maxRedirects)
			}
			return nil
		},
	}

	return &HttpClient{
		ctx:            ctx,
		apiURL:         apiURL,
		caCert:         caCert,
		insecure:       insecure,
		ignoreSystemCA: ignoreSystemCA,
		timeout:        timeout,
		maxRedirects:   maxRedirects,
		client:         retryClient,
	}
}

func (c *HttpClient) Do(req *Request, resp *Response) error {
	start := time.Now()
	ctx := c.ctx

	if req.Endpoint == "" {
		return fmt.Errorf("doing request: URL is nil")
	}

	tflog.Debug(ctx, fmt.Sprintf("Request Object: %#v", req))

	// This allows the resource to overwrite the token.
	if req.Token != "" {
		c.token = req.Token
	}

	var reqBody io.Reader
	if req.Body != nil {
		bodyBytes, err := json.Marshal(req.Body)
		if err != nil {
			return fmt.Errorf("doing request: error marshalling body: %v", err)
		}
		reqBody = bytes.NewBuffer(bodyBytes)
	}

	request, err := retryablehttp.NewRequest(req.Method, req.Endpoint, reqBody)
	if err != nil {
		return fmt.Errorf("doing request: %v", err)
	}
	request = request.WithContext(ctx)

	for key, value := range req.Headers {
		request.Header.Add(key, value)
	}

	res, err := c.client.Do(request)
	if err != nil {
		return fmt.Errorf("doing request: %v", err)
	}
	defer res.Body.Close()

	resp.StatusCode = res.StatusCode
	resp.Headers = res.Header

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %v", err)
	}
	resp.Body = responseBody

	// Timings recorded as part of internal metrics
	tflog.Debug(ctx, fmt.Sprintf("Response Time: %f ms", float64((time.Since(start))/time.Millisecond)))
	tflog.Debug(ctx, fmt.Sprintf("Response Status: %s", res.Status))
	tflog.Debug(ctx, fmt.Sprintf("Response Headers: %#v", res.Header))
	tflog.Debug(ctx, fmt.Sprintf("Response Body: %s", string(responseBody)))

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		tflog.Debug(ctx, fmt.Sprintf("Successful response! (%d)", resp.StatusCode))
		return nil
	}

	if resp.StatusCode >= 400 {
		return &ApiError{
			StatusCode: resp.StatusCode,
			Message:    string(responseBody),
		}
	}

	return nil
}

func (c *HttpClient) Set(client Client) (Client, error) {
	httpClient, ok := client.(*HttpClient)
	if !ok {
		return nil, fmt.Errorf("invalid client type, expected: '*HttpClient', got: '%T'", client)
	}
	c.ctx = httpClient.ctx
	c.apiURL = httpClient.apiURL
	c.caCert = httpClient.caCert
	c.ignoreSystemCA = httpClient.ignoreSystemCA
	c.insecure = httpClient.insecure
	c.maxRedirects = httpClient.maxRedirects
	c.timeout = httpClient.timeout
	return c, nil
}

func (c *HttpClient) GetApiUrl() string {
	return c.apiURL
}
