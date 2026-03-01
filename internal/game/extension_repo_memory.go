package game

import (
	"context"
	"slices"
	"sync"
)

var _ ExtensionRepository = (*InMemoryExtensionRepository)(nil)

type InMemoryExtensionRepository struct {
	bySlug map[string][]Extension
	mu     sync.RWMutex
}

func NewInMemoryExtensionRepository() *InMemoryExtensionRepository {
	return &InMemoryExtensionRepository{
		bySlug: make(map[string][]Extension),
	}
}

func (repo *InMemoryExtensionRepository) ListBySlug(ctx context.Context, slug string) ([]Extension, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	list := repo.bySlug[slug]
	if list == nil {
		return nil, nil
	}

	return slices.Clone(list), nil
}

func (repo *InMemoryExtensionRepository) Store(ctx context.Context, ext Extension) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	list := repo.bySlug[ext.Slug]
	if list == nil {
		list = []Extension{}
	}

	for i, e := range list {
		if e.Filename == ext.Filename {
			list[i] = ext
			repo.bySlug[ext.Slug] = list
			return nil
		}
	}

	repo.bySlug[ext.Slug] = append(list, ext)

	return nil
}
