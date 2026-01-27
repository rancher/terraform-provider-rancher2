package client

import (
	"context"
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
	response       Response
	request        Request
}

func NewTestClient(ctx context.Context, apiURL, caCert string, insecure, ignoreSystemCA bool, timeout time.Duration, maxRedirects int64) *TestClient {
	return &TestClient{
		ctx:            ctx,
		apiURL:         apiURL,
		caCert:         caCert,
		insecure:       insecure,
		ignoreSystemCA: ignoreSystemCA,
		timeout:        timeout,
		maxRedirects:   maxRedirects,
	}
}

func (c *TestClient) Do(req *Request, resp *Response) error {
	start := time.Now()
	ctx := c.ctx

	if req.Endpoint == "" {
		return fmt.Errorf("doing request: URL is nil")
	}

	tflog.Debug(ctx, fmt.Sprintf("Request Object: %v", pp.PrettyPrint(req)))
	c.request = *req

	resp.Body = c.response.Body
	resp.Headers = c.response.Headers
	resp.StatusCode = c.response.StatusCode

	res := *resp

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

	tflog.Debug(ctx, fmt.Sprintf("Response Object: %v", pp.PrettyPrint(res)))
	tflog.Debug(ctx, fmt.Sprintf("Response Time: %f ms", float64((time.Since(start))/time.Millisecond)))

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
