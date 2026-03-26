package client

import (
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
	TokenStore     *TokenStore
	responses      map[string]Response
	requests       map[string]Request
}

func NewTestClient(ctx context.Context, apiURL, caCert string, insecure, ignoreSystemCA bool, timeout time.Duration, maxRedirects int64, tokenStore *TokenStore) *TestClient {
	return &TestClient{
		apiURL:         apiURL,
		caCert:         caCert,
		insecure:       insecure,
		ignoreSystemCA: ignoreSystemCA,
		timeout:        timeout,
		maxRedirects:   maxRedirects,
		TokenStore:     tokenStore,
		requests:       map[string]Request{},
		responses:      map[string]Response{},
	}
}

type LogFields struct {
	ID         string
	Token      string
	Url        string
	Method     string
	StatusCode int
	Headers    map[string][]string
	Body       any
}

func (c *TestClient) Do(ctx context.Context, req *Request, resp *Response) error {
	if req.Endpoint == "" {
		return fmt.Errorf("doing request: URL is nil")
	}

	// This allows the resource to overwrite the client token and the client token to set the request token.
	// if req.Token != "" {
	// 	c.token = req.Token
	// }
	// if c.token != "" {
	//   req.Token = c.token
	// }
	var err error
	var reqBody []byte
	if req.Body == nil {
		if req.Method == "PUT" || req.Method == "PATCH" || req.Method == "POST" {
			req.Body = json.RawMessage("{}")
		}
	} else {
		reqBody, err = json.Marshal(req.Body)
		if err != nil {
			return err
		}
	}

	newReq := Request{
		Method:   req.Method,
		Headers:  req.Headers,
		Endpoint: req.Endpoint,
	}

	if len(reqBody) > 0 {
		newReq.Body = reqBody
	}

	if req.Method == "POST" || req.Method == "PUT" {
		if newReq.Headers == nil {
			newReq.Headers = map[string][]string{}
		}
		newReq.Headers["Content-Type"] = []string{"application/json"}
	}

	// Inject the token, threadsafe due to mutex.
	token := c.TokenStore.GetToken()

	if token != "" {
		if newReq.Headers == nil {
			newReq.Headers = map[string][]string{}
		}
		if len(newReq.Headers["Authorization"]) == 0 {
			newReq.Headers["Authorization"] = []string{fmt.Sprintf("Bearer %s", token)}
		} else {
			tflog.Debug(ctx, "Token already set.")
		}
	} else {
		tflog.Debug(ctx, "Token is empty.")
	}

	requestId := fmt.Sprintf("%s:%s:%s", newReq.Endpoint, newReq.Method, token)
	c.requests[requestId] = newReq

	response, err := c.GetResponse(requestId)
	if err != nil {
		return err
	}

	response, ok := c.responses[requestId]
	if !ok {
		return fmt.Errorf("no response found for request %s", requestId)
	}

	resp.Body = response.Body
	resp.Headers = response.Headers
	resp.StatusCode = response.StatusCode

	logFields := LogFields{
		ID:      requestId,
		Token:   token,
		Url:     req.Endpoint,
		Method:  req.Method,
		Headers: req.Headers,
		Body:    req.Body,
	}

	tflog.Debug(ctx, fmt.Sprintf("HTTP Request: %+v", pp.PrettyPrint(logFields)))

	logFields = LogFields{
		ID:         requestId,
		Token:      token,
		Url:        req.Endpoint,
		Method:     req.Method,
		StatusCode: resp.StatusCode,
		Headers:    resp.Headers,
	}
	if len(resp.Body) > 0 {
		logFields.Body = json.RawMessage(resp.Body)
	}

	tflog.Debug(ctx, fmt.Sprintf("HTTP Response: %s", pp.PrettyPrint(logFields)))

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
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

func (c *TestClient) SetResponse(ctx context.Context, id string, response Response) {
	tflog.Debug(ctx, fmt.Sprintf("Setting response for request %s", id))
	c.responses[id] = response
}

func (c *TestClient) GetResponse(id string) (Response, error) {
	response, ok := c.responses[id]
	if !ok {
		return Response{}, fmt.Errorf("no response found for request %s", id)
	}
	return response, nil
}

func (c *TestClient) GetRequest(id string) Request {
	return c.requests[id]
}

func (c *TestClient) GetRequests() map[string]Request {
	return c.requests
}

func (c *TestClient) ClearToken() {
	c.TokenStore.ClearToken()
}

func (c *TestClient) Token() string {
	return c.TokenStore.GetToken()
}

func (c *TestClient) SetToken(token string) {
	c.TokenStore.SetToken(token)
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
