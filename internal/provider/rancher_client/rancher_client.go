package rancher_client

import (
  "context"
)

type RancherRequest interface {
  DoRequest(ctx context.Context, client RancherClient) ([]byte, error)
}

type RancherClient interface {
  Create(ctx context.Context, request RancherRequest) error
  Read(ctx context.Context, request RancherRequest) ([]byte, error)
  Update(ctx context.Context, request RancherRequest) error
  Delete(ctx context.Context, request RancherRequest) error
}
