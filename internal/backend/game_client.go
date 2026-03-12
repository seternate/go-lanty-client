package backend

import (
	"context"
	"errors"
	"image"
	"maps"
	"sync"

	"github.com/seternate/go-lanty-client/internal/game"
	"github.com/seternate/go-lanty/pkg/api"
)

type GameAPIClient struct {
	apiclient *api.Client

	mu       sync.RWMutex
	progress map[string]game.InstallationProgress
}

func NewGameAPIClient(apiclient *api.Client) *GameAPIClient {
	return &GameAPIClient{
		apiclient: apiclient,
	}
}

func (client *GameAPIClient) FetchCatalog() ([]game.CatalogItem, error) {
	panic("not implemented")
}

func (client *GameAPIClient) GetIcon(slug string) (image.Image, error) {
	panic("not implemented")
}

func (client *GameAPIClient) Download(ctx context.Context, slug string) (filePath string, err error) {
	panic("not implemented")
}

func (client *GameAPIClient) GetAllProgress() (map[string]game.InstallationProgress, error) {
	client.mu.RLock()
	defer client.mu.RUnlock()

	return maps.Clone(client.progress), nil
}

func (client *GameAPIClient) GetProgress(slug string) (game.InstallationProgress, error) {
	client.mu.RLock()
	defer client.mu.RUnlock()

	progress, ok := client.progress[slug]
	if !ok {
		return game.InstallationProgress{}, errors.New("progress not found")
	}

	return progress, nil
}
