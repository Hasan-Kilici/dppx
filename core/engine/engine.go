package engine

// Engine coordinates search operations using a configured similarity and scoring pipeline.
type Engine struct {
	config Config
}

// New initializes an Engine with the provided configuration.
func New(cfg Config) *Engine {
	return &Engine{
		config: cfg,
	}
}
