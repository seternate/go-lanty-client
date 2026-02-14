package user

import (
	"maps"
	"slices"
	"sync"

	"github.com/seternate/go-lanty-client/internal/domain/user"
	"github.com/seternate/go-lanty-client/internal/ui/adapter"
	usercontroller "github.com/seternate/go-lanty-client/internal/ui/controller/user"
)

var _ user.UserCatalogRepository = (*InMemoryCatalogRepository)(nil)
var _ usercontroller.CatalogReader = (*InMemoryCatalogRepository)(nil)
var _ adapter.UserProvider = (*InMemoryCatalogRepository)(nil)

type InMemoryCatalogRepository struct {
	catalog map[string]user.UserCatalogItem
	mu      sync.RWMutex
}

func NewInMemoryCatalogRepository() *InMemoryCatalogRepository {
	return &InMemoryCatalogRepository{
		catalog: make(map[string]user.UserCatalogItem),
	}
}

func (r *InMemoryCatalogRepository) GetAll() ([]user.UserCatalogItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return slices.Collect(maps.Values(r.catalog)), nil
}

func (r *InMemoryCatalogRepository) ListUsers() ([]user.UserCatalogItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return slices.Collect(maps.Values(r.catalog)), nil
}

func (r *InMemoryCatalogRepository) GetByIP(ip string) (user.UserCatalogItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.catalog[ip], nil
}

func (r *InMemoryCatalogRepository) ReplaceAll(items ...user.UserCatalogItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.catalog = make(map[string]user.UserCatalogItem, len(items))
	for _, item := range items {
		r.catalog[item.IP] = item
	}

	return nil
}
