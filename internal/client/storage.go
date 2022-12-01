package client

import "context"

type Repository interface {
	Create(ctx context.Context, c *Client) error
	FindAll(ctx context.Context) (cl []Client, err error)
	FindOne(ctx context.Context, id string) (Client, error)
	Update(ctx context.Context, c Client) error
}
