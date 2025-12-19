package client

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ Client = &TestClient{}   // make sure the TestClient implements the Client
var _ Request = &TestRequest{} // make sure the TestRequest implements the Request

type TestClient struct {
	ApiURL         string
	CACert         string
	IgnoreSystemCA bool
	Insecure       bool
	AccessKey      string
	SecretKey      string
	Token          string
	MaxRedirects   int
	Timeout        time.Duration
}

func NewTestClient(apiURL string, caCert string, ignoreSystemCA bool, insecure bool, accessKey string, secretKey string, token string, maxRedirects int, timeout time.Duration) *TestClient {
	return &TestClient{
		ApiURL:         apiURL,
		CACert:         caCert,
		Insecure:       insecure,
		IgnoreSystemCA: ignoreSystemCA,
		AccessKey:      accessKey,
		SecretKey:      secretKey,
		Token:          token,
		MaxRedirects:   maxRedirects,
		Timeout:        timeout,
	}
}

type TestRequest struct {
	Method   string
	Endpoint string
	Body     interface{}
	Headers  map[string]string
}

func (c *TestClient) Create(ctx context.Context, r Request) error {
	_, err := r.DoRequest(ctx, c)
	return err
}

func (c *TestClient) Read(ctx context.Context, r Request) ([]byte, error) {
	return r.DoRequest(ctx, c)
}

func (c *TestClient) Update(ctx context.Context, r Request) error {
	_, err := r.DoRequest(ctx, c)
	return err
}

func (c *TestClient) Delete(ctx context.Context, r Request) error {
	_, err := r.DoRequest(ctx, c)
	return err
}

func (r *TestRequest) DoRequest(ctx context.Context, rc Client) ([]byte, error) {

	c, ok := rc.(*TestClient)
	if !ok {
		tflog.Error(ctx, "Doing request: invalid rancher client type")
		return nil, fmt.Errorf("Doing request: invalid rancher client type")
	}

	// This implementation should return mock data for testing purposes.
	// This implementation should log the request and client details for debugging.
	tflog.Info(ctx, fmt.Sprintf("Mock Request: Method=%s, Endpoint=%s, Body=%v, Headers=%v", r.Method, r.Endpoint, r.Body, r.Headers))
	tflog.Info(ctx, fmt.Sprintf(
		"Mock Client: CACert=%s, Insecure=%t, IgnoreSystemCA=%t, Token=%s, MaxRedirects=%d, Timeout=%s",
		c.CACert,
		c.Insecure,
		c.IgnoreSystemCA,
		c.Token,
		c.MaxRedirects,
		c.Timeout,
	))
	// Return empty response for now.
	return []byte("{}"), nil
}
