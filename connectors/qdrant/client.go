package qdrant

import (
	"context"
	"fmt"

	qdrant "github.com/qdrant/go-client/qdrant"
)

/*
	Client wraps the official Qdrant client.

	This abstraction allows DPPX to:
	- expose simplified APIs
	- keep connector logic isolated
	- still provide raw Qdrant access
	  for advanced/custom operations
*/
type Client struct {
	raw *qdrant.Client
}

/*
	NewClient creates a new Qdrant client connection.

	Example:
		client, err := qdrant.NewClient(qdrant.Config{
			Host: "localhost",
			Port: 6334,
		})
*/
func NewClient(
	cfg Config,
) (*Client, error) {

	raw, err := qdrant.NewClient(
		&qdrant.Config{
			Host: cfg.Host,
			Port: cfg.Port,
		},
	)

	if err != nil {
		return nil, err
	}

	return &Client{
		raw: raw,
	}, nil
}

/*
	Raw exposes the original Qdrant client.

	Use this when you need advanced Qdrant features
	that are not abstracted by DPPX.

	Example:
		client.Raw().CreateCollection(...)
		client.Raw().Upsert(...)
		client.Raw().Scroll(...)
*/
func (c *Client) Raw() *qdrant.Client {
	return c.raw
}

/*
	CreateCollection is a helper wrapper around
	the native Qdrant CreateCollection API.

	This is optional convenience functionality.
*/
func (c *Client) CreateCollection(
	ctx context.Context,
	req *qdrant.CreateCollection,
) error {

	return c.raw.CreateCollection(
		ctx,
		req,
	)
}

/*
	Upsert inserts or updates points.

	This wrapper exists to simplify
	common DPPX ingestion workflows.
*/
func (c *Client) Upsert(
	ctx context.Context,
	req *qdrant.UpsertPoints,
) (*qdrant.UpdateResult, error) {

	return c.raw.Upsert(
		ctx,
		req,
	)
}

/*
	Query performs ANN vector search.

	Most users should use the Retriever abstraction
	instead of calling this directly.
*/
func (c *Client) Query(
	ctx context.Context,
	req *qdrant.QueryPoints,
) ([]*qdrant.ScoredPoint, error) {

	return c.raw.Query(
		ctx,
		req,
	)
}

/*
	Close closes the underlying Qdrant connection.
*/
func (c *Client) Close() error {

	if c.raw == nil {
		return nil
	}

	return c.raw.Close()
}

/*
	Ping verifies connectivity to Qdrant.

	Useful for:
	- health checks
	- startup validation
	- monitoring
*/
func (c *Client) Ping(
	ctx context.Context,
) error {

	info, err := c.raw.HealthCheck(
		ctx,
	)

	if err != nil {
		return err
	}

	if info == nil {
		return fmt.Errorf(
			"qdrant healthcheck returned nil response",
		)
	}

	return nil
}