package app

import (
	"context"
	"image"
	"sync"

	"github.com/seternate/go-lanty-client/internal/diskspace"
	"github.com/seternate/go-lanty-client/internal/game"
	gameviewmodel "github.com/seternate/go-lanty-client/internal/ui/viewmodel/game"
)

type GameCatalogReader interface {
	GetAll(ctx context.Context) ([]game.CatalogItem, error)
	GetBySlug(ctx context.Context, slug string) (game.CatalogItem, error)
}

type GameCatalogIconFetcher interface {
	FetchIcon(ctx context.Context, slug string) (image.Image, error)
}

type GameInstallationReader interface {
	GetAll(ctx context.Context) ([]game.Installation, error)
	GetBySlug(ctx context.Context, slug string) (game.Installation, error)
}

type GameProgressFetcher interface {
	GetInstallationProgress(ctx context.Context, slug string) (game.InstallationProgress, error)
}

type GameScreenSyncer struct {
	mu         sync.RWMutex
	gameScreen *gameviewmodel.GameScreen

	catalogReader             GameCatalogReader
	catalogIconFetcher        GameCatalogIconFetcher
	installationReader        GameInstallationReader
	downloadProgressFetcher   GameProgressFetcher
	extractingProgressFetcher GameProgressFetcher
}

func NewGameScreenSyncer(gameScreen *gameviewmodel.GameScreen, catalogReader GameCatalogReader, catalogIconFetcher GameCatalogIconFetcher, installationReader GameInstallationReader, downloadProgressFetcher GameProgressFetcher, extractingProgressFetcher GameProgressFetcher) *GameScreenSyncer {
	return &GameScreenSyncer{
		gameScreen:                gameScreen,
		catalogReader:             catalogReader,
		catalogIconFetcher:        catalogIconFetcher,
		installationReader:        installationReader,
		downloadProgressFetcher:   downloadProgressFetcher,
		extractingProgressFetcher: extractingProgressFetcher,
	}
}

func (syncer *GameScreenSyncer) OnCatalogItemAdded(e game.CatalogEvent) {
	syncer.mu.Lock()
	defer syncer.mu.Unlock()

	catalogItem, err := syncer.catalogReader.GetBySlug(context.Background(), e.Slug)
	if err != nil {
		return
	}
	installation, err := syncer.installationReader.GetBySlug(context.Background(), e.Slug)
	if err != nil {
		return
	}
	icon, err := syncer.catalogIconFetcher.FetchIcon(context.Background(), e.Slug)
	if err != nil {
		return
	}

	syncer.gameScreen.AddGameTile(gameviewmodel.GameTileCreate{
		Icon:         icon,
		CatalogItem:  catalogItem,
		Installation: installation,
	})

	catalogItems, err := syncer.catalogReader.GetAll(context.Background())
	if err != nil {
		return
	}
	syncer.gameScreen.UpdateAvailableGames(len(catalogItems))
}

func (syncer *GameScreenSyncer) OnCatalogItemUpdated(e game.CatalogEvent) {
	syncer.mu.Lock()
	defer syncer.mu.Unlock()

	catalogItem, err := syncer.catalogReader.GetBySlug(context.Background(), e.Slug)
	if err != nil {
		return
	}
	installation, err := syncer.installationReader.GetBySlug(context.Background(), e.Slug)
	if err != nil {
		return
	}
	icon, err := syncer.catalogIconFetcher.FetchIcon(context.Background(), e.Slug)
	if err != nil {
		return
	}

	progress := game.InstallationProgress{}
	if installation.IsDownloading() {
		progress, _ = syncer.downloadProgressFetcher.GetInstallationProgress(context.Background(), e.Slug)
	} else if installation.IsExtracting() {
		progress, _ = syncer.extractingProgressFetcher.GetInstallationProgress(context.Background(), e.Slug)
	}

	syncer.gameScreen.UpdateGameTile(gameviewmodel.GameTileUpdate{
		Icon:         icon,
		CatalogItem:  catalogItem,
		Installation: installation,
		Progress:     progress,
	})

	catalogItems, err := syncer.catalogReader.GetAll(context.Background())
	if err != nil {
		return
	}
	syncer.gameScreen.UpdateAvailableGames(len(catalogItems))

	installations, err := syncer.installationReader.GetAll(context.Background())
	if err != nil {
		return
	}

	installedGames := 0
	for _, installation := range installations {
		if installation.IsInstalled() {
			installedGames++
		}
	}

	syncer.gameScreen.UpdateInstalledGames(installedGames)
}

func (syncer *GameScreenSyncer) OnCatalogItemRemoved(e game.CatalogEvent) {
	syncer.mu.Lock()
	defer syncer.mu.Unlock()

	syncer.gameScreen.RemoveGameTile(e.Slug)

	catalogItems, err := syncer.catalogReader.GetAll(context.Background())
	if err != nil {
		return
	}
	syncer.gameScreen.UpdateAvailableGames(len(catalogItems))

	installations, err := syncer.installationReader.GetAll(context.Background())
	if err != nil {
		return
	}

	installedGames := 0
	for _, installation := range installations {
		if installation.IsInstalled() {
			installedGames++
		}
	}

	syncer.gameScreen.UpdateInstalledGames(installedGames)
}

func (syncer *GameScreenSyncer) OnInstallationUpdated(e game.InstallationEvent) {
	syncer.mu.Lock()
	defer syncer.mu.Unlock()

	catalogItem, err := syncer.catalogReader.GetBySlug(context.Background(), e.Slug)
	if err != nil {
		return
	}
	installation, err := syncer.installationReader.GetBySlug(context.Background(), e.Slug)
	if err != nil {
		return
	}
	icon, err := syncer.catalogIconFetcher.FetchIcon(context.Background(), e.Slug)
	if err != nil {
		return
	}

	progress := game.InstallationProgress{}
	if installation.IsDownloading() {
		progress, _ = syncer.downloadProgressFetcher.GetInstallationProgress(context.Background(), e.Slug)
	} else if installation.IsExtracting() {
		progress, _ = syncer.extractingProgressFetcher.GetInstallationProgress(context.Background(), e.Slug)
	}

	syncer.gameScreen.UpdateGameTile(gameviewmodel.GameTileUpdate{
		Icon:         icon,
		CatalogItem:  catalogItem,
		Installation: installation,
		Progress:     progress,
	})

	catalogItems, err := syncer.catalogReader.GetAll(context.Background())
	if err != nil {
		return
	}
	syncer.gameScreen.UpdateAvailableGames(len(catalogItems))

	installations, err := syncer.installationReader.GetAll(context.Background())
	if err != nil {
		return
	}

	installedGames := 0
	for _, installation := range installations {
		if installation.IsInstalled() {
			installedGames++
		}
	}

	syncer.gameScreen.UpdateInstalledGames(installedGames)
}

func (syncer *GameScreenSyncer) OnDiskSpaceChanged(e diskspace.DiskSpaceEvent) {
	syncer.gameScreen.UpdateFreeDiskSpace(e.FreeDiskSpace)
}
