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
	pp "github.com/rancher/terraform-provider-rancher2/internal/provider/pretty_print"
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

	// 1. Setup the Base Transport with TLS config
	baseTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: insecure,
			RootCAs:            rootCAs,
		},
	}

	// 2. Wrap it with our thread-safe AuthTransport
	authTransport := &AuthTransport{
		Base:       baseTransport,
		TokenStore: tokenStore,
	}

	retryClient := retryablehttp.NewClient()
	
	// Turn off the internal logger so it doesn't duplicate your Do() logs
	// and to avoid the context.Background() trap.
	retryClient.Logger = nil 

	retryClient.HTTPClient = &http.Client{
		Timeout:   timeout,
		Transport: authTransport, // Use the wrapped transport!
		
		// 3. Move CheckRedirect HERE so it is completely thread-safe.
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= int(maxRedirects) {
				return fmt.Errorf("stopped after %d redirects", maxRedirects)
			}
			
			// Safely grab the token and inject it on redirects if needed
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
		return fmt.Errorf("doing request: URL is nil")
	}

	var reqBody io.Reader
	if req.Body != nil {
		bodyBytes, err := json.Marshal(req.Body)
		if err != nil {
			return fmt.Errorf("doing request: error marshalling body: %v", err)
		}
		reqBody = bytes.NewBuffer(bodyBytes)
	}

	request, err := retryablehttp.NewRequestWithContext(ctx, req.Method, req.Endpoint, reqBody)
	if err != nil {
		return fmt.Errorf("doing request: %v", err)
	}

	if req.Method == "POST" {
		request.Header.Add("Content-Type", "application/json")
	}

	for key, value := range req.Headers {
		for _, v := range value {
			request.Header.Add(key, v)
		}
	}

	// EXECUTE REQUEST
	res, err := c.client.Do(request)
	if err != nil {
		return fmt.Errorf("doing request: %v", err)
	}
	defer res.Body.Close()

	resp.StatusCode = res.StatusCode
	resp.Headers = res.Header

	// Read Response
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %v", err)
	}
	resp.Body = responseBody

	// --- LOGGING ---
	// Note: Because you are building JSON payloads via bytes.NewBuffer in Do(), 
	// this manual body reading for logs is safe from OOM crashes (unlike large file streams).
	requestReader, err := request.GetBody()
	if err != nil {
		return fmt.Errorf("getting request body: %v", err)
	}
	defer requestReader.Close()
	requestBody, err := io.ReadAll(requestReader)
	if err != nil {
		return fmt.Errorf("reading request body: %v", err)
	}

	var prettyRequest bytes.Buffer
	var rp string
	err = json.Indent(&prettyRequest, requestBody, "", "  ")
	if err != nil {
		rp = string(requestBody)
	} else {
		rp = prettyRequest.String()
	}
	tflog.Debug(ctx, fmt.Sprintf("\nRequest: %s\n%s\nHeaders:\n%s\nBody:\n%s\n", req.Method, req.Endpoint, pp.PrettyPrint(request.Header), rp))

	var prettyResponse bytes.Buffer
	err = json.Indent(&prettyResponse, responseBody, "", "  ")
	if err != nil {
		rp = string(responseBody)
	} else {
		rp = prettyResponse.String()
	}
	tflog.Debug(ctx, fmt.Sprintf("\nResponse: %d\nHeaders:\n%s\nBody:\n%s\n", resp.StatusCode, pp.PrettyPrint(res.Header), rp))

	// ERROR HANDLING
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		tflog.Debug(ctx, fmt.Sprintf("\nSuccessful response! (%d)\n", resp.StatusCode))
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

