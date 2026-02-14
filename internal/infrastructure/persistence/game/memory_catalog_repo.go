package game

import (
	"maps"
	"slices"
	"sync"

	"github.com/seternate/go-lanty-client/internal/domain/game"
	gamecontroller "github.com/seternate/go-lanty-client/internal/ui/controller/game"
)

var _ game.GameCatalogRepository = (*InMemoryCatalogRepository)(nil)
var _ gamecontroller.CatalogReader = (*InMemoryCatalogRepository)(nil)

type InMemoryCatalogRepository struct {
	catalog map[string]game.GameCatalogItem
	mu      sync.RWMutex
}

func NewInMemoryCatalogRepository() *InMemoryCatalogRepository {
	return &InMemoryCatalogRepository{
		catalog: make(map[string]game.GameCatalogItem),
	}
}

func (r *InMemoryCatalogRepository) GetAll() ([]game.GameCatalogItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return slices.Collect(maps.Values(r.catalog)), nil
}

func (r *InMemoryCatalogRepository) GetBySlug(slug string) (game.GameCatalogItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.catalog[slug], nil
}

func (r *InMemoryCatalogRepository) ReplaceAll(items ...game.GameCatalogItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.catalog = make(map[string]game.GameCatalogItem, len(items))
	for _, item := range items {
		r.catalog[item.Slug] = item
	}

	return nil
}
