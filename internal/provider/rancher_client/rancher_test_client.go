package rancher_client

import (
  "context"
	"fmt"
	"time"

  "github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ RancherClient = &RancherMemoryClient{} // make sure the RancherMemoryClient implements the RancherClient
var _ RancherRequest = &RancherMemoryRequest{} // make sure the RancherMemoryRequest implements the RancherRequest

type RancherMemoryClient struct{
  ApiURL         string
  CACert         string
  IgnoreSystemCA bool
  Insecure       bool
  Token          string
  MaxRedirects   int
  Timeout        time.Duration
}

func NewRancherMemoryClient(apiURL string, caCert string, ignoreSystemCA bool, insecure bool, token string, maxRedirects int, timeout time.Duration) *RancherMemoryClient {
  return &RancherMemoryClient{
    ApiURL:         apiURL,
    CACert:         caCert,
    Insecure:       insecure,
    IgnoreSystemCA: ignoreSystemCA,
    Token:          token,
    MaxRedirects:   maxRedirects,
    Timeout:        timeout,
  }
}

type RancherMemoryRequest struct {
  Method       string
  Endpoint     string
  Body         interface{}
  Headers      map[string]string
}

func (c *RancherMemoryClient) Create(ctx context.Context, r RancherRequest) error {
  _, err := r.DoRequest(ctx, c)
  return err
}

func (c *RancherMemoryClient) Read(ctx context.Context, r RancherRequest) ([]byte, error) {
  return r.DoRequest(ctx, c)
}

func (c *RancherMemoryClient) Update(ctx context.Context, r RancherRequest) error {
  _, err := r.DoRequest(ctx, c)
  return err
}

func (c *RancherMemoryClient) Delete(ctx context.Context, r RancherRequest) error {
  _, err := r.DoRequest(ctx, c)
  return err
}


func (r *RancherMemoryRequest) DoRequest(ctx context.Context, rc RancherClient) ([]byte, error) {

  c, ok := rc.(*RancherMemoryClient)
  if !ok {
    tflog.Error(ctx, "Doing request: invalid rancher client type")
    return nil, fmt.Errorf("Doing request: invalid rancher client type")
  }

  // This implementation should return mock data for testing purposes.
  // This implementation should log the request and client details for debugging.
  tflog.Info(ctx, fmt.Sprintf("Mock Request: Method=%s, Endpoint=%s, Body=%v, Headers=%v", r.Method, r.Endpoint, r.Body, r.Headers))
  tflog.Info(ctx, fmt.Sprintf("Mock Client: CACert=%s, Insecure=%t, IgnoreSystemCA=%t, Token=%s, MaxRedirects=%d, Timeout=%s",
    c.CACert, c.Insecure, c.IgnoreSystemCA, c.Token, c.MaxRedirects, c.Timeout))
  // Return empty response for now.
  return []byte("{}"), nil
}
