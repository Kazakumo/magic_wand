package embedding

import "context"

// Embedder converts text into a vector representation.
type Embedder interface {
	Embed(ctx context.Context, text string) ([]float32, error)
}
