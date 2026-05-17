package qdrant

import qdrant "github.com/qdrant/go-client/qdrant"

func NewClient(
	cfg Config,
) (*qdrant.Client, error) {

	qcfg := &qdrant.Config{
		Host: cfg.Host,
		Port: cfg.Port,
	}

	// optional api key
	if cfg.APIKey != "" {
		qcfg.APIKey = cfg.APIKey
	}

	return qdrant.NewClient(qcfg)
}