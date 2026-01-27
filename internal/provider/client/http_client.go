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

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ Client = &HttpClient{}

type HttpClient struct {
	ctx            context.Context
	apiURL         string
	caCert         string
	ignoreSystemCA bool
	insecure       bool
	maxRedirects   int64
	timeout        time.Duration
}

func NewHttpClient(ctx context.Context, apiURL, caCert string, insecure, ignoreSystemCA bool, timeout time.Duration, maxRedirects int64) *HttpClient {
	return &HttpClient{
		ctx:            ctx,
		apiURL:         apiURL,
		caCert:         caCert,
		insecure:       insecure,
		ignoreSystemCA: ignoreSystemCA,
		timeout:        timeout,
		maxRedirects:   maxRedirects,
	}
}

func (c *HttpClient) Do(req *Request, resp *Response) error {
	start := time.Now()
	ctx := c.ctx

	if req.Endpoint == "" {
		return fmt.Errorf("doing request: URL is nil")
	}

	tflog.Debug(ctx, fmt.Sprintf("Request Object: %#v", req))

	MaxRedirectCheckFunction := func(r *http.Request, via []*http.Request) error {
		if len(via) >= int(c.maxRedirects) {
			return fmt.Errorf("stopped after %d redirects", c.maxRedirects)
		}
		return nil
	}

	var rootCAs *x509.CertPool
	if c.ignoreSystemCA {
		rootCAs = x509.NewCertPool()
	} else {
		// Get the SystemCertPool, continue with an empty pool on error
		rootCAs, _ = x509.SystemCertPool()
		if rootCAs == nil {
			rootCAs = x509.NewCertPool()
		}
	}

	if c.caCert != "" {
		// Append our cert to the cert pool
		if ok := rootCAs.AppendCertsFromPEM([]byte(c.caCert)); !ok {
			tflog.Warn(ctx, "no certs appended, using system certs only")
		}
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: c.insecure,
		RootCAs:            rootCAs,
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
		Proxy:           http.ProxyFromEnvironment,
	}

	client := &http.Client{
		Timeout:       c.timeout,
		CheckRedirect: MaxRedirectCheckFunction,
		Transport:     transport,
	}

	var reqBody *bytes.Buffer
	if req.Body != nil {
		bodyBytes, err := json.Marshal(req.Body)
		if err != nil {
			return fmt.Errorf("doing request: error marshalling body: %v", err)
		}
		reqBody = bytes.NewBuffer(bodyBytes)
	} else {
		reqBody = &bytes.Buffer{}
	}

	url := fmt.Sprintf("%s", req.Endpoint)
	request, err := http.NewRequest(req.Method, url, reqBody)
	if err != nil {
		return fmt.Errorf("doing request: %v", err)
	}

	for key, value := range req.Headers {
		request.Header.Add(key, value)
	}

	res, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("doing request: %v", err)
	}
	defer res.Body.Close()

	// Timings recorded as part of internal metrics
	tflog.Debug(ctx, fmt.Sprintf("Response Time: %f ms", float64((time.Since(start))/time.Millisecond)))
	tflog.Debug(ctx, fmt.Sprintf("Response Status: %s", res.Status))
	tflog.Debug(ctx, fmt.Sprintf("Response Headers: %#v", res.Header))

	switch resp.StatusCode {
	case 200:
		tflog.Debug(ctx, "Successful response! (200)")
	case 202:
		tflog.Debug(ctx, "Accepted! (202)")
	case 400:
		return fmt.Errorf("Bad request! (400)")
	case 401:
		return fmt.Errorf("Unauthorized! (401)")
	case 403:
		return fmt.Errorf("Forbidden! (403)")
	case 404:
		return fmt.Errorf("Not found! (404)")
	case 500:
		return fmt.Errorf("Internal server error! (500)")
	default:
		return fmt.Errorf("Unknown status code! (%d)", resp.StatusCode)
	}

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %v", err)
	}
	tflog.Debug(ctx, fmt.Sprintf("Response Body: %v", string(responseBody)))

	return nil
}

func (c *HttpClient) Set(client Client) (Client, error) {
	c.ctx = client.(*HttpClient).ctx
	c.apiURL = client.(*HttpClient).apiURL
	c.caCert = client.(*HttpClient).caCert
	c.ignoreSystemCA = client.(*HttpClient).ignoreSystemCA
	c.insecure = client.(*HttpClient).insecure
	c.maxRedirects = client.(*HttpClient).maxRedirects
	c.timeout = client.(*HttpClient).timeout
	return c, nil
}

func (c *HttpClient) GetApiUrl() string {
	return c.apiURL
}
