package search

import (
	"context"

	"bff-go-mvp/internal/model"
)

// Service defines the behavior for the search use case.
type Service interface {
	Search(ctx context.Context, page, perPage int, req model.SearchRequest) (model.SearchResponse, error)
}


