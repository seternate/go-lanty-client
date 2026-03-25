package app

import (
	"context"
	"image"
	"sync"

	"github.com/rs/zerolog"
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
	logger     zerolog.Logger
	mu         sync.RWMutex
	gameScreen *gameviewmodel.GameScreen

	catalogReader             GameCatalogReader
	catalogIconFetcher        GameCatalogIconFetcher
	installationReader        GameInstallationReader
	downloadProgressFetcher   GameProgressFetcher
	extractingProgressFetcher GameProgressFetcher
}

func NewGameScreenSyncer(logger zerolog.Logger, gameScreen *gameviewmodel.GameScreen, catalogReader GameCatalogReader, catalogIconFetcher GameCatalogIconFetcher, installationReader GameInstallationReader, downloadProgressFetcher GameProgressFetcher, extractingProgressFetcher GameProgressFetcher) *GameScreenSyncer {
	return &GameScreenSyncer{
		logger:                    logger,
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
	ctx := context.Background()
	const handler = "OnCatalogItemAdded"

	catalogItem, err := syncer.catalogReader.GetBySlug(ctx, e.Slug)
	if err != nil {
		syncer.logger.Warn().Err(err).Str("handler", handler).Str("slug", e.Slug).Msg("catalogReader.GetBySlug failed")
		return
	}
	installation, err := syncer.installationReader.GetBySlug(ctx, e.Slug)
	if err != nil {
		syncer.logger.Warn().Err(err).Str("handler", handler).Str("slug", e.Slug).Msg("installationReader.GetBySlug failed")
		return
	}
	icon, err := syncer.catalogIconFetcher.FetchIcon(ctx, e.Slug)
	if err != nil {
		syncer.logger.Warn().Err(err).Str("handler", handler).Str("slug", e.Slug).Msg("catalogIconFetcher.FetchIcon failed")
		return
	}

	if err := syncer.gameScreen.AddGameTile(gameviewmodel.GameTileCreate{
		Icon:         icon,
		CatalogItem:  catalogItem,
		Installation: installation,
	}); err != nil {
		syncer.logger.Error().Err(err).Str("handler", handler).Str("slug", e.Slug).Msg("AddGameTile failed")
	}

	catalogItems, err := syncer.catalogReader.GetAll(ctx)
	if err != nil {
		syncer.logger.Warn().Err(err).Str("handler", handler).Str("slug", e.Slug).Msg("catalogReader.GetAll failed")
		return
	}
	if err := syncer.gameScreen.UpdateAvailableGames(len(catalogItems)); err != nil {
		syncer.logger.Error().Err(err).Str("handler", handler).Str("slug", e.Slug).Msg("UpdateAvailableGames failed")
	}
}

func (syncer *GameScreenSyncer) OnCatalogItemUpdated(e game.CatalogEvent) {
	syncer.mu.Lock()
	defer syncer.mu.Unlock()
	ctx := context.Background()
	const handler = "OnCatalogItemUpdated"

	catalogItem, err := syncer.catalogReader.GetBySlug(ctx, e.Slug)
	if err != nil {
		syncer.logger.Warn().Err(err).Str("handler", handler).Str("slug", e.Slug).Msg("catalogReader.GetBySlug failed")
		return
	}
	installation, err := syncer.installationReader.GetBySlug(ctx, e.Slug)
	if err != nil {
		syncer.logger.Warn().Err(err).Str("handler", handler).Str("slug", e.Slug).Msg("installationReader.GetBySlug failed")
		return
	}
	icon, err := syncer.catalogIconFetcher.FetchIcon(ctx, e.Slug)
	if err != nil {
		syncer.logger.Warn().Err(err).Str("handler", handler).Str("slug", e.Slug).Msg("catalogIconFetcher.FetchIcon failed")
		return
	}

	progress := game.InstallationProgress{}
	if installation.IsDownloading() {
		var perr error
		progress, perr = syncer.downloadProgressFetcher.GetInstallationProgress(ctx, e.Slug)
		if perr != nil {
			syncer.logger.Warn().Err(perr).Str("handler", handler).Str("slug", e.Slug).Str("phase", "download").Msg("GetInstallationProgress failed")
		}
	} else if installation.IsExtracting() {
		var perr error
		progress, perr = syncer.extractingProgressFetcher.GetInstallationProgress(ctx, e.Slug)
		if perr != nil {
			syncer.logger.Warn().Err(perr).Str("handler", handler).Str("slug", e.Slug).Str("phase", "extract").Msg("GetInstallationProgress failed")
		}
	}

	if err := syncer.gameScreen.UpdateGameTile(gameviewmodel.GameTileUpdate{
		Icon:         icon,
		CatalogItem:  catalogItem,
		Installation: installation,
		Progress:     progress,
	}); err != nil {
		syncer.logger.Error().Err(err).Str("handler", handler).Str("slug", e.Slug).Msg("UpdateGameTile failed")
	}

	catalogItems, err := syncer.catalogReader.GetAll(ctx)
	if err != nil {
		syncer.logger.Warn().Err(err).Str("handler", handler).Str("slug", e.Slug).Msg("catalogReader.GetAll failed")
		return
	}
	if err := syncer.gameScreen.UpdateAvailableGames(len(catalogItems)); err != nil {
		syncer.logger.Error().Err(err).Str("handler", handler).Str("slug", e.Slug).Msg("UpdateAvailableGames failed")
	}

	installations, err := syncer.installationReader.GetAll(ctx)
	if err != nil {
		syncer.logger.Warn().Err(err).Str("handler", handler).Str("slug", e.Slug).Msg("installationReader.GetAll failed")
		return
	}

	installedGames := 0
	for _, installation := range installations {
		if installation.IsInstalled() {
			installedGames++
		}
	}

	if err := syncer.gameScreen.UpdateInstalledGames(installedGames); err != nil {
		syncer.logger.Error().Err(err).Str("handler", handler).Str("slug", e.Slug).Msg("UpdateInstalledGames failed")
	}
}

func (syncer *GameScreenSyncer) OnCatalogItemRemoved(e game.CatalogEvent) {
	syncer.mu.Lock()
	defer syncer.mu.Unlock()
	ctx := context.Background()
	const handler = "OnCatalogItemRemoved"

	if err := syncer.gameScreen.RemoveGameTile(e.Slug); err != nil {
		syncer.logger.Error().Err(err).Str("handler", handler).Str("slug", e.Slug).Msg("RemoveGameTile failed")
	}

	catalogItems, err := syncer.catalogReader.GetAll(ctx)
	if err != nil {
		syncer.logger.Warn().Err(err).Str("handler", handler).Str("slug", e.Slug).Msg("catalogReader.GetAll failed")
		return
	}
	if err := syncer.gameScreen.UpdateAvailableGames(len(catalogItems)); err != nil {
		syncer.logger.Error().Err(err).Str("handler", handler).Str("slug", e.Slug).Msg("UpdateAvailableGames failed")
	}

	installations, err := syncer.installationReader.GetAll(ctx)
	if err != nil {
		syncer.logger.Warn().Err(err).Str("handler", handler).Str("slug", e.Slug).Msg("installationReader.GetAll failed")
		return
	}

	installedGames := 0
	for _, installation := range installations {
		if installation.IsInstalled() {
			installedGames++
		}
	}

	if err := syncer.gameScreen.UpdateInstalledGames(installedGames); err != nil {
		syncer.logger.Error().Err(err).Str("handler", handler).Str("slug", e.Slug).Msg("UpdateInstalledGames failed")
	}
}

func (syncer *GameScreenSyncer) OnInstallationUpdated(e game.InstallationEvent) {
	syncer.mu.Lock()
	defer syncer.mu.Unlock()
	ctx := context.Background()
	const handler = "OnInstallationUpdated"

	catalogItem, err := syncer.catalogReader.GetBySlug(ctx, e.Slug)
	if err != nil {
		syncer.logger.Warn().Err(err).Str("handler", handler).Str("slug", e.Slug).Msg("catalogReader.GetBySlug failed")
		return
	}
	installation, err := syncer.installationReader.GetBySlug(ctx, e.Slug)
	if err != nil {
		syncer.logger.Warn().Err(err).Str("handler", handler).Str("slug", e.Slug).Msg("installationReader.GetBySlug failed")
		return
	}
	icon, err := syncer.catalogIconFetcher.FetchIcon(ctx, e.Slug)
	if err != nil {
		syncer.logger.Warn().Err(err).Str("handler", handler).Str("slug", e.Slug).Msg("catalogIconFetcher.FetchIcon failed")
		return
	}

	progress := game.InstallationProgress{}
	if installation.IsDownloading() {
		var perr error
		progress, perr = syncer.downloadProgressFetcher.GetInstallationProgress(ctx, e.Slug)
		if perr != nil {
			syncer.logger.Warn().Err(perr).Str("handler", handler).Str("slug", e.Slug).Str("phase", "download").Msg("GetInstallationProgress failed")
		}
	} else if installation.IsExtracting() {
		var perr error
		progress, perr = syncer.extractingProgressFetcher.GetInstallationProgress(ctx, e.Slug)
		if perr != nil {
			syncer.logger.Warn().Err(perr).Str("handler", handler).Str("slug", e.Slug).Str("phase", "extract").Msg("GetInstallationProgress failed")
		}
	}

	if err := syncer.gameScreen.UpdateGameTile(gameviewmodel.GameTileUpdate{
		Icon:         icon,
		CatalogItem:  catalogItem,
		Installation: installation,
		Progress:     progress,
	}); err != nil {
		syncer.logger.Error().Err(err).Str("handler", handler).Str("slug", e.Slug).Msg("UpdateGameTile failed")
	}

	catalogItems, err := syncer.catalogReader.GetAll(ctx)
	if err != nil {
		syncer.logger.Warn().Err(err).Str("handler", handler).Str("slug", e.Slug).Msg("catalogReader.GetAll failed")
		return
	}
	if err := syncer.gameScreen.UpdateAvailableGames(len(catalogItems)); err != nil {
		syncer.logger.Error().Err(err).Str("handler", handler).Str("slug", e.Slug).Msg("UpdateAvailableGames failed")
	}

	installations, err := syncer.installationReader.GetAll(ctx)
	if err != nil {
		syncer.logger.Warn().Err(err).Str("handler", handler).Str("slug", e.Slug).Msg("installationReader.GetAll failed")
		return
	}

	installedGames := 0
	for _, installation := range installations {
		if installation.IsInstalled() {
			installedGames++
		}
	}

	if err := syncer.gameScreen.UpdateInstalledGames(installedGames); err != nil {
		syncer.logger.Error().Err(err).Str("handler", handler).Str("slug", e.Slug).Msg("UpdateInstalledGames failed")
	}
}

func (syncer *GameScreenSyncer) OnDiskSpaceChanged(e diskspace.DiskSpaceEvent) {
	const handler = "OnDiskSpaceChanged"
	if err := syncer.gameScreen.UpdateFreeDiskSpace(e.FreeDiskSpace); err != nil {
		syncer.logger.Error().Err(err).Str("handler", handler).Uint64("free_disk_space", e.FreeDiskSpace).Msg("UpdateFreeDiskSpace failed")
	}
}
