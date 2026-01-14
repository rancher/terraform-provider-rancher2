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

var _ Client = &HttpClient{}   // make sure the HttpClient implements the Client
var _ Request = &HttpRequest{} // make sure the HttpRequest implements the Request

type HttpClient struct {
	ApiURL         string
	CACert         string
	IgnoreSystemCA bool
	Insecure       bool
	MaxRedirects   int64
	Timeout        time.Duration
}

func NewHttpClient(ctx context.Context, apiURL string, caCert string, ignoreSystemCA bool, insecure bool, maxRedirects int64, timeout string) *HttpClient {
	to, err := time.ParseDuration(timeout)
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("error parsing timeout: %v", err))
		return nil
	}
	return &HttpClient{
		ApiURL:         apiURL,
		CACert:         caCert,
		IgnoreSystemCA: ignoreSystemCA,
		Insecure:       insecure,
		MaxRedirects:   maxRedirects,
		Timeout:        to,
	}
}

type HttpRequest struct {
	Method   string
	Endpoint string
	Body     any
	Headers  map[string]string
}

func (c *HttpClient) Create(ctx context.Context, r Request) error {
	_, err := r.DoRequest(ctx, c)
	return err
}

func (c *HttpClient) Read(ctx context.Context, r Request) ([]byte, error) {
	return r.DoRequest(ctx, c)
}

func (c *HttpClient) Update(ctx context.Context, r Request) error {
	_, err := r.DoRequest(ctx, c)
	return err
}

func (c *HttpClient) Delete(ctx context.Context, r Request) error {
	_, err := r.DoRequest(ctx, c)
	return err
}

func (r *HttpRequest) DoRequest(ctx context.Context, rc Client) ([]byte, error) {
	start := time.Now()

	c, ok := rc.(*HttpClient)
	if !ok {
		tflog.Error(ctx, "doing request: invalid rancher client type")
		return nil, fmt.Errorf("doing request: invalid rancher client type")
	}

	if r.Endpoint == "" {
		tflog.Error(ctx, "doing request: URL is nil")
		return nil, fmt.Errorf("doing request: URL is nil")
	}

	tflog.Debug(ctx, fmt.Sprintf("Request Object: %#v", r))

	MaxRedirectCheckFunction := func(req *http.Request, via []*http.Request) error {
		if len(via) >= int(c.MaxRedirects) {
			tflog.Error(ctx, fmt.Sprintf("stopped after %d redirects", c.MaxRedirects))
			return fmt.Errorf("stopped after %d redirects", c.MaxRedirects)
		}
		// if len(c.Token) > 0 {
		// 	// make sure the auth token is added to redirected requests
		// 	req.Header.Add("Authorization", "Bearer "+c.Token)
		// }
		return nil
	}

	var rootCAs *x509.CertPool
	if c.IgnoreSystemCA {
		rootCAs = x509.NewCertPool()
	} else {
		// Get the SystemCertPool, continue with an empty pool on error
		rootCAs, _ = x509.SystemCertPool()
		if rootCAs == nil {
			rootCAs = x509.NewCertPool()
		}
	}

	if c.CACert != "" {
		// Append our cert to the cert pool
		if ok := rootCAs.AppendCertsFromPEM([]byte(c.CACert)); !ok {
			tflog.Warn(ctx, "no certs appended, using system certs only")
		}
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: c.Insecure,
		RootCAs:            rootCAs,
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
		Proxy:           http.ProxyFromEnvironment,
	}

	client := &http.Client{
		Timeout:       c.Timeout,
		CheckRedirect: MaxRedirectCheckFunction,
		Transport:     transport,
	}

	var reqBody *bytes.Buffer
	if r.Body != nil {
		bodyBytes, err := json.Marshal(r.Body)
		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("doing request: error marshalling body: %v", err))
			return nil, fmt.Errorf("doing request: error marshalling body: %v", err)
		}
		reqBody = bytes.NewBuffer(bodyBytes)
	} else {
		reqBody = &bytes.Buffer{}
	}

	request, err := http.NewRequest(r.Method, r.Endpoint, reqBody)
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("doing request: %v", err))
		return nil, fmt.Errorf("doing request: %v", err)
	}

	for key, value := range r.Headers {
		request.Header.Add(key, value)
	}

	// if len(c.Token) > 0 {
	// 	request.Header.Add("Authorization", "Bearer "+c.Token)
	// }

	resp, err := client.Do(request)
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("doing request: %v", err))
		return nil, fmt.Errorf("doing request: %v", err)
	}
	defer resp.Body.Close()

	// Timings recorded as part of internal metrics
	tflog.Debug(ctx, fmt.Sprintf("Response Time: %f ms", float64((time.Since(start))/time.Millisecond)))
	tflog.Debug(ctx, fmt.Sprintf("Response Status: %s", resp.Status))
	tflog.Debug(ctx, fmt.Sprintf("Response Headers: %#v", resp.Header))
	tflog.Debug(ctx, fmt.Sprintf("Response Body: %v", resp.Body))

	return io.ReadAll(resp.Body)
}
