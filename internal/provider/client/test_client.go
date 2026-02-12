package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	pp "github.com/rancher/terraform-provider-rancher2/internal/provider/pretty_print"
)

var _ Client = &TestClient{}

type TestClient struct {
	apiURL         string
	caCert         string
	ignoreSystemCA bool
	insecure       bool
	maxRedirects   int64
	timeout        time.Duration
	token          string
	response       Response
	requests       []Request
}

func NewTestClient(ctx context.Context, apiURL, caCert string, insecure, ignoreSystemCA bool, timeout time.Duration, maxRedirects int64, token string) *TestClient {
	return &TestClient{
		apiURL:         apiURL,
		caCert:         caCert,
		insecure:       insecure,
		ignoreSystemCA: ignoreSystemCA,
		timeout:        timeout,
		maxRedirects:   maxRedirects,
		token:          token,
		requests:       make([]Request, 0),
	}
}

func (c *TestClient) Do(ctx context.Context, req *Request, resp *Response) error {
	start := time.Now()

	if req.Endpoint == "" {
		return fmt.Errorf("doing request: URL is nil")
	}

	tflog.Debug(ctx, fmt.Sprintf("Request Object: %v", pp.PrettyPrint(req)))

	reqBody, err := json.Marshal(req.Body)
	if err != nil {
		return err
	}

	newReq := Request{
		Method:   req.Method,
		Headers:  req.Headers,
		Endpoint: req.Endpoint,
		Body:     reqBody,
	}
	c.requests = append(c.requests, newReq)

	resp.Body = c.response.Body
	resp.Headers = c.response.Headers
	resp.StatusCode = c.response.StatusCode

	tflog.Debug(ctx, fmt.Sprintf("Response Time: %f ms", float64((time.Since(start))/time.Millisecond)))
	tflog.Debug(ctx, fmt.Sprintf("Response Status Code: %d", resp.StatusCode))
	tflog.Debug(ctx, fmt.Sprintf("Response Headers: %s", pp.PrettyPrint(resp.Headers)))
	var prettyBody bytes.Buffer
	if err := json.Indent(&prettyBody, resp.Body, "", "  "); err == nil {
		tflog.Debug(ctx, fmt.Sprintf("Response Body (Pretty): \n%s", prettyBody.String()))
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		tflog.Debug(ctx, fmt.Sprintf("Successful response! (%d)", resp.StatusCode))
		return nil
	}

	if resp.StatusCode >= 400 {
		return &ApiError{
			StatusCode: resp.StatusCode,
			Message:    string(resp.Body),
		}
	}

	return nil
}

func (c *TestClient) Set(client Client) (Client, error) {
	testClient, ok := client.(*TestClient)
	if !ok {
		return nil, fmt.Errorf("invalid client type, expected: '*TestClient', got: '%T'", client)
	}
	c.apiURL = testClient.apiURL
	c.caCert = testClient.caCert
	c.ignoreSystemCA = testClient.ignoreSystemCA
	c.insecure = testClient.insecure
	c.maxRedirects = testClient.maxRedirects
	c.timeout = testClient.timeout
	c.response = testClient.response
	c.requests = testClient.requests
	return c, nil
}

func (c *TestClient) GetApiUrl() string {
	return c.apiURL
}

func (c *TestClient) SetResponse(response Response) {
	c.response = response
}

func (c *TestClient) GetLastRequest() Request {
	if len(c.requests) > 0 {
		return c.requests[len(c.requests)-1]
	}
	return Request{}
}

// GetLastRequests returns the last n requests.
func (c *TestClient) GetLastRequests(n int) []Request {
	if len(c.requests) < n {
		return c.requests
	}
	return c.requests[len(c.requests)-n:]
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
