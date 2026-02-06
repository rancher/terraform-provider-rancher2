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
	ctx            context.Context
	apiURL         string
	caCert         string
	ignoreSystemCA bool
	insecure       bool
	maxRedirects   int64
	timeout        time.Duration
	token          string
	response       Response
	request        Request
}

func NewTestClient(ctx context.Context, apiURL, caCert string, insecure, ignoreSystemCA bool, timeout time.Duration, maxRedirects int64, token string) *TestClient {
	return &TestClient{
		ctx:            ctx,
		apiURL:         apiURL,
		caCert:         caCert,
		insecure:       insecure,
		ignoreSystemCA: ignoreSystemCA,
		timeout:        timeout,
		maxRedirects:   maxRedirects,
		token:          token,
	}
}

func (c *TestClient) Do(req *Request, resp *Response) error {
	start := time.Now()
	ctx := c.ctx

	if req.Endpoint == "" {
		return fmt.Errorf("doing request: URL is nil")
	}

	tflog.Debug(ctx, fmt.Sprintf("Request Object: %v", pp.PrettyPrint(req)))

	c.request.Method = req.Method
	c.request.Headers = req.Headers
	c.request.Endpoint = req.Endpoint
	reqBody, err := json.Marshal(req.Body)
	if err != nil {
		return err
	}
	c.request.Body = reqBody

	resp.Body = c.response.Body
	resp.Headers = c.response.Headers
	resp.StatusCode = c.response.StatusCode

	tflog.Debug(ctx, fmt.Sprintf("Response Time: %f ms", float64((time.Since(start))/time.Millisecond)))
	tflog.Debug(ctx, fmt.Sprintf("Response Status Code: %d", resp.StatusCode))
	tflog.Debug(ctx, fmt.Sprintf("Response Headers: %#v", resp.Headers))
	tflog.Debug(ctx, fmt.Sprintf("Response Body: %s", string(resp.Body)))

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
	c.ctx = client.(*TestClient).ctx
	c.apiURL = client.(*TestClient).apiURL
	c.caCert = client.(*TestClient).caCert
	c.ignoreSystemCA = client.(*TestClient).ignoreSystemCA
	c.insecure = client.(*TestClient).insecure
	c.maxRedirects = client.(*TestClient).maxRedirects
	c.timeout = client.(*TestClient).timeout
	c.response = client.(*TestClient).response
	c.request = client.(*TestClient).request
	return c, nil
}

func (c *TestClient) GetApiUrl() string {
	return c.apiURL
}

func (c *TestClient) SetResponse(response Response) {
	c.response = response
}

func (c *TestClient) GetLastRequest() Request {
	return c.request
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
