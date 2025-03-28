package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client PromocodeServiceClient
}

func NewClient(serverAddr string) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	client := NewPromocodeServiceClient(conn)
	return &Client{
		conn:   conn,
		client: client,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) CreatePromocode(ctx context.Context, req *CreatePromocodeRequest) (*PromocodeResponse, error) {
	return c.client.CreatePromocode(ctx, req)
}

func (c *Client) GetPromocode(ctx context.Context, req *GetPromocodeRequest) (*PromocodeResponse, error) {
	return c.client.GetPromocode(ctx, req)
}

func (c *Client) UpdatePromocode(ctx context.Context, req *UpdatePromocodeRequest) (*PromocodeResponse, error) {
	return c.client.UpdatePromocode(ctx, req)
}

func (c *Client) DeletePromocode(ctx context.Context, req *DeletePromocodeRequest) (*DeletePromocodeResponse, error) {
	return c.client.DeletePromocode(ctx, req)
}

func (c *Client) ListPromocodes(ctx context.Context, req *ListPromocodesRequest) (*ListPromocodesResponse, error) {
	return c.client.ListPromocodes(ctx, req)
}
