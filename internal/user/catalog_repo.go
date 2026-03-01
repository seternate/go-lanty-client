package user

import "context"

type CatalogRepository interface {
	GetAll(ctx context.Context) ([]CatalogItem, error)
	GetByIP(ctx context.Context, ip string) (CatalogItem, error)
	ReplaceAll(ctx context.Context, items ...CatalogItem) error
}
