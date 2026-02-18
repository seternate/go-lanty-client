package game

import (
	"context"
	"image"
	"math"
	"path/filepath"
	"sync"
	"time"

	"github.com/seternate/go-lanty-client/internal/application/game"
	gameapp "github.com/seternate/go-lanty-client/internal/application/game"
	gameevent "github.com/seternate/go-lanty-client/internal/application/game/event"
	gameappsrv "github.com/seternate/go-lanty-client/internal/application/game/service"
	systemevent "github.com/seternate/go-lanty-client/internal/application/system/event"
	gamedomain "github.com/seternate/go-lanty-client/internal/domain/game"
	model "github.com/seternate/go-lanty-client/internal/ui/model/game"
)

type CatalogReader interface {
	GetAll() ([]gamedomain.GameCatalogItem, error)
	GetBySlug(slug string) (gamedomain.GameCatalogItem, error)
}

type CatalogIconFetcher interface {
	GetIcon(slug string) (image.Image, error)
}

type InstallationReader interface {
	GetAll() ([]gamedomain.GameInstallation, error)
	GetBySlug(slug string) (gamedomain.GameInstallation, error)
}

type GameLauncherSpecReader interface {
	GetLaunchSpec(slug string, mode gameappsrv.LaunchSpecMode) (gamedomain.GameLaunchSpec, error)
}

type ProgressFetcher interface {
	GetProgress(slug string) (game.InstallationProgress, error)
}

type DiskSpaceProvider interface {
	GetFreeDiskSpace() (uint64, error)
}

type JoinUserAdapter interface {
	OpenUserSelection(gamename string, onUserSelected func(ipAddress string))
}

type HostConfigAdapter interface {
	ShowHostConfig(gamename string, configFields []model.HostConfigField, onSubmit func(values []gameapp.ArgInput))
}

type ExtensionsAdapter interface {
	ShowExtensions(gameName, gameSlug, installPath string, extensions *model.ExtensionListModel, onUpload func(localPath string), onDownload func(ext *model.ExtensionModel))
}

type GameController struct {
	GameListModel  *model.GameListModel
	GameStatsModel *model.GameStatsModel

	catalogReader             CatalogReader
	catalogIconFetcher        CatalogIconFetcher
	installationReader        InstallationReader
	spaceProvider             DiskSpaceProvider
	downloadProgressFetcher   ProgressFetcher
	extractingProgressFetcher ProgressFetcher
	installationService       gameappsrv.GameInstallationService
	directoryOpener           gameappsrv.GameDirectoryOpener
	joinUserAdapter           JoinUserAdapter
	hostConfigAdapter         HostConfigAdapter
	launcherService           gameappsrv.GameLauncherService
	launchSpecReader          GameLauncherSpecReader
	extensionService          gameappsrv.ExtensionService
	extensionsAdapter         ExtensionsAdapter
	installationRoot          string
	mu                        sync.RWMutex
}

func NewGameController(
	catalogReader CatalogReader,
	catalogIconFetcher CatalogIconFetcher,
	installationReader InstallationReader,
	spaceProvider DiskSpaceProvider,
	downloadProgressFetcher ProgressFetcher,
	extractingProgressFetcher ProgressFetcher,
	installationService gameappsrv.GameInstallationService,
	directoryOpener gameappsrv.GameDirectoryOpener,
	joinUserAdapter JoinUserAdapter,
	launcherService gameappsrv.GameLauncherService,
	hostConfigAdapter HostConfigAdapter,
	launchSpecReader GameLauncherSpecReader,
	extensionService gameappsrv.ExtensionService,
	extensionsAdapter ExtensionsAdapter,
	installationRoot string,
) *GameController {
	return &GameController{
		GameListModel:             model.NewGameListModel(),
		GameStatsModel:            model.NewGameStatsModel(),
		catalogReader:             catalogReader,
		catalogIconFetcher:        catalogIconFetcher,
		installationReader:        installationReader,
		spaceProvider:             spaceProvider,
		downloadProgressFetcher:   downloadProgressFetcher,
		extractingProgressFetcher: extractingProgressFetcher,
		installationService:       installationService,
		directoryOpener:           directoryOpener,
		joinUserAdapter:           joinUserAdapter,
		launcherService:           launcherService,
		hostConfigAdapter:         hostConfigAdapter,
		launchSpecReader:          launchSpecReader,
		extensionService:         extensionService,
		extensionsAdapter:         extensionsAdapter,
		installationRoot:         installationRoot,
	}
}

func (controller *GameController) OnCatalogItemAdded(e gameevent.CatalogEvent) {
	controller.mu.Lock()
	defer controller.mu.Unlock()

	catalogItem, err := controller.catalogReader.GetBySlug(e.Slug)
	if err != nil {
		return
	}
	installation, err := controller.installationReader.GetBySlug(catalogItem.Slug)
	if err != nil {
		return
	}
	icon, err := controller.catalogIconFetcher.GetIcon(catalogItem.Slug)
	if err != nil {
		return
	}

	tile := model.NewGameTileModel(e.Slug)
	controller.updateGameModel(tile, catalogItem, installation, icon, game.InstallationProgress{})
	controller.GameListModel.AddGame(tile)

	freeDiskSpace, err := controller.spaceProvider.GetFreeDiskSpace()
	if err != nil {
		return
	}
	catalog, err := controller.catalogReader.GetAll()
	if err != nil {
		return
	}
	installations, err := controller.installationReader.GetAll()
	if err != nil {
		return
	}

	installedGames := 0
	for _, installation := range installations {
		if installation.IsInstalled() {
			installedGames++
		}
	}

	controller.updateStats(len(catalog), installedGames, freeDiskSpace)
}

func (controller *GameController) OnCatalogItemUpdated(e gameevent.CatalogEvent) {
	controller.mu.Lock()
	defer controller.mu.Unlock()

	catalogItem, err := controller.catalogReader.GetBySlug(e.Slug)
	if err != nil {
		return
	}
	installation, err := controller.installationReader.GetBySlug(catalogItem.Slug)
	if err != nil {
		return
	}
	icon, err := controller.catalogIconFetcher.GetIcon(catalogItem.Slug)
	if err != nil {
		return
	}

	tile := controller.GameListModel.GetGameBySlug(e.Slug)
	controller.updateGameModel(tile, catalogItem, installation, icon, game.InstallationProgress{})

	freeDiskSpace, err := controller.spaceProvider.GetFreeDiskSpace()
	if err != nil {
		return
	}
	catalog, err := controller.catalogReader.GetAll()
	if err != nil {
		return
	}
	installations, err := controller.installationReader.GetAll()
	if err != nil {
		return
	}

	installedGames := 0
	for _, installation := range installations {
		if installation.IsInstalled() {
			installedGames++
		}
	}

	controller.updateStats(len(catalog), installedGames, freeDiskSpace)
}

func (controller *GameController) OnCatalogItemRemoved(e gameevent.CatalogEvent) {
	controller.mu.Lock()
	defer controller.mu.Unlock()

	controller.GameListModel.RemoveGame(e.Slug)

	freeDiskSpace, err := controller.spaceProvider.GetFreeDiskSpace()
	if err != nil {
		return
	}
	catalog, err := controller.catalogReader.GetAll()
	if err != nil {
		return
	}
	installations, err := controller.installationReader.GetAll()
	if err != nil {
		return
	}

	installedGames := 0
	for _, installation := range installations {
		if installation.IsInstalled() {
			installedGames++
		}
	}

	controller.updateStats(len(catalog), installedGames, freeDiskSpace)
}

func (controller *GameController) OnInstallationStarted(e gameevent.InstallationEvent) {
	controller.mu.Lock()
	defer controller.mu.Unlock()

	catalogItem, err := controller.catalogReader.GetBySlug(e.Slug)
	if err != nil {
		return
	}
	installation, err := controller.installationReader.GetBySlug(catalogItem.Slug)
	if err != nil {
		return
	}
	icon, err := controller.catalogIconFetcher.GetIcon(catalogItem.Slug)
	if err != nil {
		return
	}

	tile := controller.GameListModel.GetGameBySlug(e.Slug)
	controller.updateGameModel(tile, catalogItem, installation, icon, game.InstallationProgress{})

	freeDiskSpace, err := controller.spaceProvider.GetFreeDiskSpace()
	if err != nil {
		return
	}
	catalog, err := controller.catalogReader.GetAll()
	if err != nil {
		return
	}
	installations, err := controller.installationReader.GetAll()
	if err != nil {
		return
	}

	installedGames := 0
	for _, installation := range installations {
		if installation.IsInstalled() {
			installedGames++
		}
	}

	controller.updateStats(len(catalog), installedGames, freeDiskSpace)
}

func (controller *GameController) OnInstallationProgress(e gameevent.InstallationEvent) {
	controller.mu.Lock()
	defer controller.mu.Unlock()

	catalogItem, err := controller.catalogReader.GetBySlug(e.Slug)
	if err != nil {
		return
	}
	installation, err := controller.installationReader.GetBySlug(catalogItem.Slug)
	if err != nil {
		return
	}
	icon, err := controller.catalogIconFetcher.GetIcon(catalogItem.Slug)
	if err != nil {
		return
	}

	progress := game.InstallationProgress{}
	if installation.IsDownloading() {
		progress, _ = controller.downloadProgressFetcher.GetProgress(e.Slug)
	} else if installation.IsExtracting() {
		progress, _ = controller.extractingProgressFetcher.GetProgress(e.Slug)
	}

	tile := controller.GameListModel.GetGameBySlug(e.Slug)
	controller.updateGameModel(tile, catalogItem, installation, icon, progress)

	freeDiskSpace, err := controller.spaceProvider.GetFreeDiskSpace()
	if err != nil {
		return
	}
	catalog, err := controller.catalogReader.GetAll()
	if err != nil {
		return
	}
	installations, err := controller.installationReader.GetAll()
	if err != nil {
		return
	}

	installedGames := 0
	for _, installation := range installations {
		if installation.IsInstalled() {
			installedGames++
		}
	}

	controller.updateStats(len(catalog), installedGames, freeDiskSpace)
}

func (controller *GameController) OnInstallationFinished(e gameevent.InstallationEvent) {
	controller.mu.Lock()
	defer controller.mu.Unlock()

	catalogItem, err := controller.catalogReader.GetBySlug(e.Slug)
	if err != nil {
		return
	}
	installation, err := controller.installationReader.GetBySlug(catalogItem.Slug)
	if err != nil {
		return
	}
	icon, err := controller.catalogIconFetcher.GetIcon(catalogItem.Slug)
	if err != nil {
		return
	}

	tile := controller.GameListModel.GetGameBySlug(e.Slug)
	controller.updateGameModel(tile, catalogItem, installation, icon, game.InstallationProgress{})

	freeDiskSpace, err := controller.spaceProvider.GetFreeDiskSpace()
	if err != nil {
		return
	}
	catalog, err := controller.catalogReader.GetAll()
	if err != nil {
		return
	}
	installations, err := controller.installationReader.GetAll()
	if err != nil {
		return
	}

	installedGames := 0
	for _, installation := range installations {
		if installation.IsInstalled() {
			installedGames++
		}
	}

	controller.updateStats(len(catalog), installedGames, freeDiskSpace)
}

func (controller *GameController) OnInstallationCanceled(e gameevent.InstallationEvent) {
	controller.mu.Lock()
	defer controller.mu.Unlock()

	catalogItem, err := controller.catalogReader.GetBySlug(e.Slug)
	if err != nil {
		return
	}
	installation, err := controller.installationReader.GetBySlug(catalogItem.Slug)
	if err != nil {
		return
	}
	icon, err := controller.catalogIconFetcher.GetIcon(catalogItem.Slug)
	if err != nil {
		return
	}

	tile := controller.GameListModel.GetGameBySlug(e.Slug)
	controller.updateGameModel(tile, catalogItem, installation, icon, game.InstallationProgress{})

	freeDiskSpace, err := controller.spaceProvider.GetFreeDiskSpace()
	if err != nil {
		return
	}
	catalog, err := controller.catalogReader.GetAll()
	if err != nil {
		return
	}
	installations, err := controller.installationReader.GetAll()
	if err != nil {
		return
	}

	installedGames := 0
	for _, installation := range installations {
		if installation.IsInstalled() {
			installedGames++
		}
	}

	controller.updateStats(len(catalog), installedGames, freeDiskSpace)
}

func (controller *GameController) OnInstallationFailed(e gameevent.InstallationEvent) {
	controller.mu.Lock()
	defer controller.mu.Unlock()

	catalogItem, err := controller.catalogReader.GetBySlug(e.Slug)
	if err != nil {
		return
	}
	installation, err := controller.installationReader.GetBySlug(catalogItem.Slug)
	if err != nil {
		return
	}
	icon, err := controller.catalogIconFetcher.GetIcon(catalogItem.Slug)
	if err != nil {
		return
	}

	tile := controller.GameListModel.GetGameBySlug(e.Slug)
	controller.updateGameModel(tile, catalogItem, installation, icon, game.InstallationProgress{})

	freeDiskSpace, err := controller.spaceProvider.GetFreeDiskSpace()
	if err != nil {
		return
	}
	catalog, err := controller.catalogReader.GetAll()
	if err != nil {
		return
	}
	installations, err := controller.installationReader.GetAll()
	if err != nil {
		return
	}

	installedGames := 0
	for _, installation := range installations {
		if installation.IsInstalled() {
			installedGames++
		}
	}

	controller.updateStats(len(catalog), installedGames, freeDiskSpace)
}

func (controller *GameController) OnInstallationDetected(e gameevent.InstallationEvent) {
	controller.mu.Lock()
	defer controller.mu.Unlock()

	catalogItem, err := controller.catalogReader.GetBySlug(e.Slug)
	if err != nil {
		return
	}
	installation, err := controller.installationReader.GetBySlug(catalogItem.Slug)
	if err != nil {
		return
	}
	icon, err := controller.catalogIconFetcher.GetIcon(catalogItem.Slug)
	if err != nil {
		return
	}

	tile := controller.GameListModel.GetGameBySlug(e.Slug)
	controller.updateGameModel(tile, catalogItem, installation, icon, game.InstallationProgress{})

	freeDiskSpace, err := controller.spaceProvider.GetFreeDiskSpace()
	if err != nil {
		return
	}
	catalog, err := controller.catalogReader.GetAll()
	if err != nil {
		return
	}
	installations, err := controller.installationReader.GetAll()
	if err != nil {
		return
	}

	installedGames := 0
	for _, installation := range installations {
		if installation.IsInstalled() {
			installedGames++
		}
	}

	controller.updateStats(len(catalog), installedGames, freeDiskSpace)
}

func (controller *GameController) OnInstallationRemoved(e gameevent.InstallationEvent) {
	controller.mu.Lock()
	defer controller.mu.Unlock()

	catalogItem, err := controller.catalogReader.GetBySlug(e.Slug)
	if err != nil {
		return
	}
	installation, err := controller.installationReader.GetBySlug(catalogItem.Slug)
	if err != nil {
		return
	}
	icon, err := controller.catalogIconFetcher.GetIcon(catalogItem.Slug)
	if err != nil {
		return
	}

	tile := controller.GameListModel.GetGameBySlug(e.Slug)
	controller.updateGameModel(tile, catalogItem, installation, icon, game.InstallationProgress{})

	freeDiskSpace, err := controller.spaceProvider.GetFreeDiskSpace()
	if err != nil {
		return
	}
	catalog, err := controller.catalogReader.GetAll()
	if err != nil {
		return
	}
	installations, err := controller.installationReader.GetAll()
	if err != nil {
		return
	}

	installedGames := 0
	for _, installation := range installations {
		if installation.IsInstalled() {
			installedGames++
		}
	}

	controller.updateStats(len(catalog), installedGames, freeDiskSpace)
}

func (controller *GameController) OnExtensionListRefreshed(e gameevent.ExtensionEvent) {
	controller.mu.Lock()
	defer controller.mu.Unlock()

	tile := controller.GameListModel.GetGameBySlug(e.Slug)
	if tile == nil {
		return
	}
	hasNew, _ := controller.extensionService.HasNewExtensions(context.Background(), e.Slug)
	tile.SetHasNewExtensions(hasNew)
}

func (controller *GameController) OnExtensionDownloaded(e gameevent.ExtensionEvent) {
	controller.mu.Lock()
	defer controller.mu.Unlock()

	tile := controller.GameListModel.GetGameBySlug(e.Slug)
	if tile == nil {
		return
	}
	hasNew, _ := controller.extensionService.HasNewExtensions(context.Background(), e.Slug)
	tile.SetHasNewExtensions(hasNew)
}

func (controller *GameController) OnDiskSpaceChanged(e systemevent.DiskSpaceEvent) {
	controller.mu.Lock()
	defer controller.mu.Unlock()

	freeDiskSpace, err := controller.spaceProvider.GetFreeDiskSpace()
	if err != nil {
		return
	}
	catalog, err := controller.catalogReader.GetAll()
	if err != nil {
		return
	}
	installations, err := controller.installationReader.GetAll()
	if err != nil {
		return
	}

	installedGames := 0
	for _, installation := range installations {
		if installation.IsInstalled() {
			installedGames++
		}
	}

	controller.updateStats(len(catalog), installedGames, freeDiskSpace)
}

func (controller *GameController) updateGameModel(gamemodel *model.GameTileModel, catalogItem gamedomain.GameCatalogItem, installation gamedomain.GameInstallation, icon image.Image, progress game.InstallationProgress) {
	gamemodel.SetName(catalogItem.Name)
	gamemodel.SetIcon(icon)
	gamemodel.SetBlobSize(catalogItem.BlobSize)
	gamemodel.SetSupportsJoiningMultiplayer(catalogItem.SupportsJoiningMultiplayer)
	gamemodel.SetSupportsHostingServer(catalogItem.SupportsHostingServer)

	if installation.IsNotInstalled() {
		gamemodel.SetStatusNotInstalled()
	} else if installation.IsDownloading() {
		gamemodel.SetStatusDownloading()
		speed := uint64(0)
		progressRuntime := time.Since(progress.StartedAt).Seconds()
		if progressRuntime > 0 {
			speed = uint64(math.Round(float64(progress.Completed) / progressRuntime))
		}
		gamemodel.UpdateDownloadProgress(progress.Total, progress.Completed, speed)
	} else if installation.IsExtracting() {
		gamemodel.SetStatusExtracting()
		speed := uint64(0)
		progressRuntime := time.Since(progress.StartedAt).Seconds()
		if progressRuntime > 0 {
			speed = uint64(math.Round(float64(progress.Completed) / progressRuntime))
		}
		gamemodel.UpdateExtractingProgress(progress.Total, progress.Completed, speed)
	} else if installation.IsInstallationError() {
		gamemodel.SetStatusError()
	} else if installation.IsInstalled() {
		gamemodel.SetStatusInstalled()
	} else if installation.IsInstallationCanceled() {
		gamemodel.SetStatusCanceled()
	}

	if installation.InstallationAvailable() {
		gamemodel.MarkInstallationAsDetected()
	} else {
		gamemodel.MarkInstallationAsRemoved()
	}

	visible, _ := controller.extensionService.ExtensionsVisible(catalogItem.Slug)
	gamemodel.SetExtensionsVisible(visible)
	hasNew, _ := controller.extensionService.HasNewExtensions(context.Background(), catalogItem.Slug)
	gamemodel.SetHasNewExtensions(hasNew)
}

func (controller *GameController) updateStats(availableGames int, installedGames int, freeDiskSpace uint64) {
	controller.GameStatsModel.SetAvailableGames(availableGames)
	controller.GameStatsModel.SetInstalledGames(installedGames)
	controller.GameStatsModel.SetFreeDiskSpace(freeDiskSpace)
}

func (controller *GameController) PlayGame(slug string) {
	controller.launcherService.StartSingleplayer(slug)
}

func (controller *GameController) JoinGame(slug string) {
	catalogItem, err := controller.catalogReader.GetBySlug(slug)
	if err != nil {
		return
	}

	controller.joinUserAdapter.OpenUserSelection(catalogItem.Name, func(ipAddress string) {
		controller.launcherService.JoinMultiplayer(slug, ipAddress)
	})
}

func (controller *GameController) StartGameServer(slug string) {
	catalogItem, err := controller.catalogReader.GetBySlug(slug)
	if err != nil {
		return
	}

	launchSpec, err := controller.launchSpecReader.GetLaunchSpec(slug, gameappsrv.LaunchSpecModeHost)
	if err != nil {
		return
	}

	configFields := []model.HostConfigField{}
	for _, param := range launchSpec.Params {
		if param.Required && param.Type == gamedomain.GameLaunchParamFlag {
			continue
		}

		var fieldType model.FieldType

		switch param.Type {
		case gamedomain.GameLaunchParamFlag:
			fieldType = model.FieldFlag
		case gamedomain.GameLaunchParamString:
			fieldType = model.FieldString
		case gamedomain.GameLaunchParamInt:
			fieldType = model.FieldInt
		case gamedomain.GameLaunchParamFloat:
			fieldType = model.FieldFloat
		case gamedomain.GameLaunchParamEnum:
			fieldType = model.FieldEnum
		}

		enumValues := map[string]string{}
		for _, option := range param.EnumValues {
			enumValues[option.Label] = option.Value
		}

		configFields = append(configFields, model.HostConfigField{
			Name:        param.Name,
			Description: param.Description,
			Argument:    param.Argument,
			Required:    param.Required,
			Enabled:     param.Enabled,
			Type:        fieldType,
			Value:       param.Value,
			EnumValues:  enumValues,
			MinInt:      param.MinInt,
			MaxInt:      param.MaxInt,
			MinFloat:    param.MinFloat,
			MaxFloat:    param.MaxFloat,
		})
	}

	controller.hostConfigAdapter.ShowHostConfig(catalogItem.Name, configFields, func(values []gameapp.ArgInput) {
		controller.launcherService.HostMultiplayer(slug, values)
	})
}

func (controller *GameController) DownloadGame(slug string) {
	controller.mu.RLock()
	tile := controller.GameListModel.GetGameBySlug(slug)
	if tile == nil {
		return
	}
	isProgressing, _ := tile.IsProgressing.Get()
	controller.mu.RUnlock()

	if isProgressing {
		controller.installationService.StopInstallation(slug)
	} else {
		controller.installationService.StartInstallation(context.Background(), slug)
	}
}

func (controller *GameController) OpenGame(slug string) {
	controller.directoryOpener.OpenDirectory(slug)
}

func (controller *GameController) OpenExtensions(slug string) {
	visible, err := controller.extensionService.ExtensionsVisible(slug)
	if err != nil || !visible {
		return
	}
	catalogItem, err := controller.catalogReader.GetBySlug(slug)
	if err != nil {
		return
	}
	installation, err := controller.installationReader.GetBySlug(slug)
	if err != nil || !installation.IsInstalled() {
		return
	}
	installPath, _ := filepath.Abs(filepath.Join(controller.installationRoot, installation.InstallationDirectoryRelative))

	ctx := context.Background()
	listModel := model.NewExtensionListModel()
	extensions, err := controller.extensionService.ListExtensions(ctx, slug)
	if err != nil {
		return
	}
	for _, ext := range extensions {
		listModel.AddExtension(model.NewExtensionModel(ext.Filename, ext.Name, ext.Size, ext.Downloaded))
	}

	onUpload := func(localPath string) {
		_ = controller.extensionService.UploadExtension(context.Background(), slug, localPath)
	}
	onDownload := func(ext *model.ExtensionModel) {
		if err := controller.extensionService.DownloadExtension(context.Background(), slug, ext.Filename); err == nil {
			ext.SetDownloaded(true)
		}
	}
	controller.extensionsAdapter.ShowExtensions(catalogItem.Name, slug, installPath, listModel, onUpload, onDownload)
}
