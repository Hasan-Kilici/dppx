package native

type Config struct {

	/*
		Vector dimension.
	*/
	Dimension int

	/*
		Storage path.
	*/
	Path string

	/*
		HNSW enabled.
	*/
	UseHNSW bool

	/*
		HNSW max connections.
	*/
	HNSWM int

	/*
		HNSW efConstruction parameter.
	*/
	EFConstruction int

	/*
		HNSW efSearch parameter.
	*/
	EFSearch int

	/*
		Maximum in-memory segment size.
	*/
	MaxSegmentSize int

	/*
		Payload enabled.
	*/
	StorePayload bool

	/*
		WAL enabled.
	*/
	EnableWAL bool
}