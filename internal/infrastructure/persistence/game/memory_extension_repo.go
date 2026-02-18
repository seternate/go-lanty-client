package game

import (
	"slices"
	"sync"

	"github.com/seternate/go-lanty-client/internal/domain/game"
)

var _ game.GameExtensionRepository = (*InMemoryExtensionRepository)(nil)

type InMemoryExtensionRepository struct {
	bySlug map[string][]game.Extension
	mu     sync.RWMutex
}

func NewInMemoryExtensionRepository() *InMemoryExtensionRepository {
	return &InMemoryExtensionRepository{
		bySlug: make(map[string][]game.Extension),
	}
}

func (r *InMemoryExtensionRepository) ListByGameSlug(slug string) ([]game.Extension, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := r.bySlug[slug]
	if list == nil {
		return nil, nil
	}
	return slices.Clone(list), nil
}

func (r *InMemoryExtensionRepository) Store(ext game.Extension) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	list := r.bySlug[ext.GameSlug]
	if list == nil {
		list = []game.Extension{}
	}
	for i, e := range list {
		if e.Filename == ext.Filename {
			list[i] = ext
			r.bySlug[ext.GameSlug] = list
			return nil
		}
	}
	r.bySlug[ext.GameSlug] = append(list, ext)
	return nil
}
