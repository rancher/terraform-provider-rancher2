package client

import (
	"context"
)

type Request interface {
	DoRequest(ctx context.Context, client Client) ([]byte, error)
}

type Client interface {
	Create(ctx context.Context, request Request) error
	Read(ctx context.Context, request Request) ([]byte, error)
	Update(ctx context.Context, request Request) error
	Delete(ctx context.Context, request Request) error
}
