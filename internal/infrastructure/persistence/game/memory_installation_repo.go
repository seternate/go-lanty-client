package game

import (
	"maps"
	"slices"
	"sync"

	"github.com/seternate/go-lanty-client/internal/domain/game"
	gamecontroller "github.com/seternate/go-lanty-client/internal/ui/controller/game"
)

var _ game.GameInstallationRepository = (*InMemoryInstallationRepository)(nil)
var _ gamecontroller.InstallationReader = (*InMemoryInstallationRepository)(nil)

type InMemoryInstallationRepository struct {
	installations map[string]game.GameInstallation
	mu            sync.RWMutex
}

func NewInMemoryInstallationRepository() *InMemoryInstallationRepository {
	return &InMemoryInstallationRepository{
		installations: make(map[string]game.GameInstallation),
	}
}

func (r *InMemoryInstallationRepository) GetAll() ([]game.GameInstallation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return slices.Collect(maps.Values(r.installations)), nil
}

func (r *InMemoryInstallationRepository) GetBySlug(slug string) (game.GameInstallation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	installation, ok := r.installations[slug]
	if !ok {
		return game.GameInstallation{
			Slug:   slug,
			Status: game.InstallationStatusNotInstalled,
		}, nil
	}

	return installation, nil
}

func (r *InMemoryInstallationRepository) Store(installation game.GameInstallation) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.installations[installation.Slug] = installation
	return nil
}

func (r *InMemoryInstallationRepository) Remove(slug string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.installations, slug)

	return nil
}
