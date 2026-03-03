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
	responses      map[string]Response
	requests       map[string]Request
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
		requests:       map[string]Request{},
    responses:      map[string]Response{},
	}
}

func (c *TestClient) Do(ctx context.Context, req *Request, resp *Response) error {
	if req.Endpoint == "" {
		return fmt.Errorf("doing request: URL is nil")
	}

  // This allows the resource to overwrite the client token and the client token to set the request token.
	if req.Token != "" {
		c.token = req.Token
	}
  if c.token != "" {
    req.Token = c.token
  }

  reqBody, err := json.Marshal(req.Body)
	if err != nil {
		return err
	}

	newReq := Request{
		Method:   req.Method,
		Headers:  req.Headers,
		Endpoint: req.Endpoint,
		Body:     reqBody,
    Token:    req.Token,
	}

  if req.Method == "POST" {
    if newReq.Headers == nil {
      newReq.Headers = map[string][]string{}
    }
    newReq.Headers["Content-Type"] = []string{"application/json"}
  }

  if c.token != "" {
    if newReq.Headers == nil {
      newReq.Headers = map[string][]string{}
    }
    if len(newReq.Headers["Authorization"]) == 0 {
      newReq.Headers["Authorization"] = []string{fmt.Sprintf("Bearer %s", c.token)}
    } else {
      tflog.Debug(ctx, "Token already set.")
    }
	} else {
    tflog.Debug(ctx, "Token is empty.")
  }

  requestId := fmt.Sprintf("%s:%s:%s", newReq.Endpoint, newReq.Method, c.token)
	c.requests[requestId] = newReq

  response := c.responses[requestId]

	resp.Body = response.Body
	resp.Headers = response.Headers
	resp.StatusCode = response.StatusCode

  tflog.Debug(ctx, fmt.Sprintf("\nRequest: %s\n%s\nHeaders:\n%s\nBody:\n%s\n", newReq.Method, newReq.Endpoint, pp.PrettyPrint(newReq.Headers), newReq.Body))

	var prettyResponse bytes.Buffer
  var rp string
  err = json.Indent(&prettyResponse, resp.Body, "", "  ")
  if err != nil {
    rp = string(resp.Body)
  } else {
    rp = prettyResponse.String()
  }
  tflog.Debug(ctx, fmt.Sprintf("\nResponse: %d\nHeaders:\n%s\nBody:\n%s\n", resp.StatusCode, pp.PrettyPrint(resp.Headers), rp))

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		tflog.Debug(ctx, fmt.Sprintf("\nSuccessful response! (%d)\n", resp.StatusCode))
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
	c.responses = testClient.responses
	c.requests = testClient.requests
	return c, nil
}

func (c *TestClient) GetApiUrl() string {
	return c.apiURL
}

func (c *TestClient) SetResponse(id string, response Response) {
	c.responses[id] = response
}

func (c *TestClient) GetRequest(id string) Request {
	return c.requests[id]
}

func (c *TestClient) ClearToken() {
	c.token = ""
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
