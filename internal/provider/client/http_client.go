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

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ Client = &HttpClient{}

type HttpClient struct {
	apiURL         string
	caCert         string
	ignoreSystemCA bool
	insecure       bool
	maxRedirects   int64
	timeout        time.Duration
	TokenStore     *TokenStore // Replaced the string token with our thread-safe store
	client         *retryablehttp.Client
}

// Pass the initialized TokenStore into the client constructor
func NewHttpClient(ctx context.Context, apiURL, caCert string, insecure, ignoreSystemCA bool, timeout time.Duration, maxRedirects int64, tokenStore *TokenStore) *HttpClient {
	var rootCAs *x509.CertPool
	if ignoreSystemCA {
		rootCAs = x509.NewCertPool()
	} else {
		rootCAs, _ = x509.SystemCertPool()
		if rootCAs == nil {
			rootCAs = x509.NewCertPool()
		}
	}
	if caCert != "" {
		if ok := rootCAs.AppendCertsFromPEM([]byte(caCert)); !ok {
			tflog.Warn(ctx, "no certs appended, using system certs only")
		}
	}

	baseTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: insecure,
			RootCAs:            rootCAs,
		},
	}

	authTransport := &AuthTransport{
		Base:       baseTransport,
		TokenStore: tokenStore,
	}

	retryClient := retryablehttp.NewClient()

	// Turn off the internal logger so it doesn't duplicate logs and to avoid unnecessary contexts.
	retryClient.Logger = nil

	retryClient.HTTPClient = &http.Client{
		Timeout:   timeout,
		Transport: authTransport, // Use the wrapped transport!
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= int(maxRedirects) {
				return fmt.Errorf("stopped after %d redirects", maxRedirects)
			}
			token := tokenStore.GetToken()
			if token != "" {
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
			}
			return nil
		},
	}

	return &HttpClient{
		apiURL:         apiURL,
		caCert:         caCert,
		insecure:       insecure,
		ignoreSystemCA: ignoreSystemCA,
		timeout:        timeout,
		maxRedirects:   maxRedirects,
		TokenStore:     tokenStore,
		client:         retryClient,
	}
}

func (c *HttpClient) Do(ctx context.Context, req *Request, resp *Response) error {
	if req.Endpoint == "" {
		return fmt.Errorf("error doing request: URL is nil")
	}

	var reqBody io.Reader
	if req.Body != nil {
		bodyBytes, err := json.Marshal(req.Body)
		if err != nil {
			return fmt.Errorf("error doing request: error marshalling request body: %v", err)
		}
		reqBody = bytes.NewBuffer(bodyBytes)
	}

	request, err := retryablehttp.NewRequestWithContext(ctx, req.Method, req.Endpoint, reqBody)
	if err != nil {
		return fmt.Errorf("error doing request: %v", err)
	}

	if req.Method == "POST" {
		request.Header.Add("Content-Type", "application/json")
	}

	for key, value := range req.Headers {
		for _, v := range value {
			request.Header.Add(key, v)
		}
	}

	// logging and auth headers handled by round_tripper.go
	res, err := c.client.Do(request)
	if err != nil {
		return fmt.Errorf("error doing request: %v", err)
	}
	defer res.Body.Close()

	resp.StatusCode = res.StatusCode
	resp.Headers = res.Header

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %v", err)
	}
	resp.Body = responseBody

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
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

func (c *HttpClient) GetApiUrl() string {
	return c.apiURL
}

func (c *HttpClient) Set(client Client) (Client, error) {
	httpClient, ok := client.(*HttpClient)
	if !ok {
		return nil, fmt.Errorf("invalid client type, expected: '*HttpClient', got: '%T'", client)
	}
	c.apiURL = httpClient.apiURL
	c.caCert = httpClient.caCert
	c.ignoreSystemCA = httpClient.ignoreSystemCA
	c.insecure = httpClient.insecure
	c.maxRedirects = httpClient.maxRedirects
	c.timeout = httpClient.timeout
	c.TokenStore = httpClient.TokenStore
	c.client = httpClient.client
	return c, nil
}

func (c *HttpClient) SetToken(token string) {
	c.TokenStore.SetToken(token)
}

func (c *HttpClient) Token() string {
	return c.TokenStore.GetToken()
}

func (c *HttpClient) ClearToken() {
	c.TokenStore.ClearToken()
}
