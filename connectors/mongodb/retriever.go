package mongodb

import (
	"context"

	"github.com/hasan-kilici/dppx/core/retriever"
	"github.com/hasan-kilici/dppx/types"

	"go.mongodb.org/mongo-driver/bson"
)

var _ retriever.Retriever = (*Retriever)(nil)

type Retriever struct {
	client *Client

	database   string
	collection string

	vectorField string
	indexName   string

	/*
		true:
		use Atlas native vector search

		false:
		use local Mongo retrieval
		and let DPPX handle reranking
	*/
	useAtlasVectorSearch bool
}

/*
	Expose raw DPPX Mongo client.

	Useful for:
	- inserts
	- admin operations
	- custom pipelines
	- indexes
*/
func (r *Retriever) Client() *Client {
	return r.client
}

func New(
	ctx context.Context,
	cfg Config,
) (*Retriever, error) {

	client, err := NewClient(
		ctx,
		cfg,
	)

	if err != nil {
		return nil, err
	}

	return &Retriever{
		client: client,

		database: cfg.Database,

		collection: cfg.Collection,

		vectorField: cfg.VectorField,

		indexName: cfg.IndexName,

		useAtlasVectorSearch:
			cfg.UseAtlasVectorSearch,
	}, nil
}

/*
	Search dispatches retrieval strategy.

	Atlas mode:
	- uses MongoDB Atlas Vector Search

	Local mode:
	- uses regular Mongo queries
	- DPPX handles similarity/reranking
*/
func (r *Retriever) Search(
	ctx context.Context,
	query types.Query,
	limit int,
) ([]types.Item, error) {

	if r.useAtlasVectorSearch {

		return r.searchAtlas(
			ctx,
			query,
			limit,
		)
	}

	return r.searchLocal(
		ctx,
		query,
		limit,
	)
}

/*
	searchAtlas uses MongoDB Atlas
	native ANN vector retrieval.
*/
func (r *Retriever) searchAtlas(
	ctx context.Context,
	query types.Query,
	limit int,
) ([]types.Item, error) {

	coll := r.client.Raw().
		Database(r.database).
		Collection(r.collection)

	pipeline := bson.A{
		bson.M{
			"$vectorSearch": bson.M{

				/*
					Atlas Vector Search index.
				*/
				"index": r.indexName,

				/*
					Embedding field path.
				*/
				"path": r.vectorField,

				/*
					Query embedding vector.
				*/
				"queryVector": query.Vector,

				/*
					ANN candidate pool.
				*/
				"numCandidates": limit * 10,

				/*
					Final retrieval limit.
				*/
				"limit": limit,
			},
		},
	}

	cursor, err := coll.Aggregate(
		ctx,
		pipeline,
	)

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var docs []bson.M

	if err := cursor.All(
		ctx,
		&docs,
	); err != nil {
		return nil, err
	}

	items := make(
		[]types.Item,
		0,
		len(docs),
	)

	for _, doc := range docs {

		item := mapDocument(
			doc,
			r.vectorField,
		)

		items = append(
			items,
			item,
		)
	}

	return items, nil
}

/*
	searchLocal is used for:
	- local MongoDB
	- Docker Mongo
	- self-hosted Mongo

	This mode does NOT require
	Atlas Vector Search.

	DPPX handles:
	- cosine similarity
	- reranking
	- diversification
	- scoring
*/
func (r *Retriever) searchLocal(
	ctx context.Context,
	_ types.Query,
	limit int,
) ([]types.Item, error) {

	coll := r.client.Raw().
		Database(r.database).
		Collection(r.collection)

	cursor, err := coll.Find(
		ctx,
		bson.M{},
	)

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var docs []bson.M

	if err := cursor.All(
		ctx,
		&docs,
	); err != nil {
		return nil, err
	}

	items := make(
		[]types.Item,
		0,
		len(docs),
	)

	for _, doc := range docs {

		item := mapDocument(
			doc,
			r.vectorField,
		)

		items = append(
			items,
			item,
		)
	}

	/*
		Optional safety limit.

		DPPX engine will still rerank
		and diversify these candidates.
	*/

	if len(items) > limit {
		items = items[:limit]
	}

	return items, nil
}