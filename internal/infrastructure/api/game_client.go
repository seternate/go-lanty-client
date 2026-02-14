package api

import (
	"context"
	"errors"
	"image"
	"maps"
	"sync"

	gameapp "github.com/seternate/go-lanty-client/internal/application/game"
	gameappsrv "github.com/seternate/go-lanty-client/internal/application/game/service"
	"github.com/seternate/go-lanty-client/internal/domain/game"
	gamecontroller "github.com/seternate/go-lanty-client/internal/ui/controller/game"
	"github.com/seternate/go-lanty/pkg/api"
)

var _ gameappsrv.GameCatalogSource = (*GameAPIClient)(nil)
var _ gameappsrv.BlobDownloader = (*GameAPIClient)(nil)
var _ gameappsrv.GameInstallationProgressFetcher = (*GameAPIClient)(nil)
var _ gamecontroller.CatalogIconFetcher = (*GameAPIClient)(nil)

type GameAPIClient struct {
	apiclient *api.Client

	mu       sync.RWMutex
	progress map[string]gameapp.InstallationProgress
}

func NewGameAPIClient(apiclient *api.Client) *GameAPIClient {
	return &GameAPIClient{
		apiclient: apiclient,
	}
}

func (client *GameAPIClient) FetchCatalog() ([]game.GameCatalogItem, error) {
	panic("not implemented")
}

func (client *GameAPIClient) GetIcon(slug string) (image.Image, error) {
	panic("not implemented")
}

func (client *GameAPIClient) Download(ctx context.Context, slug string) (filePath string, err error) {
	panic("not implemented")
}

func (client *GameAPIClient) GetAllProgress() (map[string]gameapp.InstallationProgress, error) {
	client.mu.RLock()
	defer client.mu.RUnlock()

	return maps.Clone(client.progress), nil
}

func (client *GameAPIClient) GetProgress(slug string) (gameapp.InstallationProgress, error) {
	client.mu.RLock()
	defer client.mu.RUnlock()

	progress, ok := client.progress[slug]
	if !ok {
		return gameapp.InstallationProgress{}, errors.New("progress not found")
	}

	return progress, nil
}
