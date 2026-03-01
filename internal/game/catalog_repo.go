package game

import "context"

type CatalogRepository interface {
	GetAll(ctx context.Context) ([]CatalogItem, error)
	GetBySlug(ctx context.Context, slug string) (CatalogItem, error)
	ReplaceAll(ctx context.Context, items ...CatalogItem) error
}
