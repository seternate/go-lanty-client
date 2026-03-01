package game

import (
	"context"
	"maps"
	"slices"
	"sync"
)

var _ InstallationRepository = (*InMemoryInstallationRepository)(nil)

type InMemoryInstallationRepository struct {
	installations map[string]Installation
	mu            sync.RWMutex
}

func NewInMemoryInstallationRepository() *InMemoryInstallationRepository {
	return &InMemoryInstallationRepository{
		installations: make(map[string]Installation),
	}
}

func (repo *InMemoryInstallationRepository) GetAll(ctx context.Context) ([]Installation, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	return slices.Collect(maps.Values(repo.installations)), nil
}

func (repo *InMemoryInstallationRepository) GetBySlug(ctx context.Context, slug string) (Installation, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	installation, found := repo.installations[slug]
	if !found {
		return *NewInstallation(slug), nil
	}

	return installation, nil
}

func (repo *InMemoryInstallationRepository) Store(ctx context.Context, installation Installation) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.installations[installation.Slug] = installation

	return nil
}

func (repo *InMemoryInstallationRepository) Remove(ctx context.Context, slug string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	delete(repo.installations, slug)

	return nil
}
