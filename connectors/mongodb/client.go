package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	raw *mongo.Client
}

func NewClient(
	ctx context.Context,
	cfg Config,
) (*Client, error) {

	raw, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(
			cfg.URI,
		),
	)

	if err != nil {
		return nil, err
	}

	return &Client{
		raw: raw,
	}, nil
}

/*
	Raw exposes the original
	MongoDB client.

	Use this when you need:
	- manual collections
	- custom aggregations
	- transactions
	- Atlas Search
	- advanced Mongo features
*/
func (c *Client) Raw() *mongo.Client {
	return c.raw
}