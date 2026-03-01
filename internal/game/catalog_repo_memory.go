package game

import (
	"context"
	"maps"
	"slices"
	"sync"
)

var _ CatalogRepository = (*InMemoryCatalogRepository)(nil)

type InMemoryCatalogRepository struct {
	catalog map[string]CatalogItem
	mu      sync.RWMutex
}

func NewInMemoryCatalogRepository() *InMemoryCatalogRepository {
	return &InMemoryCatalogRepository{
		catalog: make(map[string]CatalogItem),
	}
}

func (repo *InMemoryCatalogRepository) GetAll(ctx context.Context) ([]CatalogItem, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	return slices.Collect(maps.Values(repo.catalog)), nil
}

func (repo *InMemoryCatalogRepository) GetBySlug(ctx context.Context, slug string) (CatalogItem, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	return repo.catalog[slug], nil
}

func (repo *InMemoryCatalogRepository) ReplaceAll(ctx context.Context, items ...CatalogItem) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.catalog = make(map[string]CatalogItem, len(items))
	for _, item := range items {
		repo.catalog[item.Slug] = item
	}

	return nil
}
