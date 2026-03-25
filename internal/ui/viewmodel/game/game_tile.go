package gameviewmodel

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"sync"

	"fyne.io/fyne/v2/data/binding"
	"github.com/dustin/go-humanize"
	"github.com/rs/zerolog"
	"github.com/seternate/go-lanty-client/internal/game"
	"github.com/seternate/go-lanty-client/internal/ui/theme"
)

type IconFetcher interface {
	FetchIcon(ctx context.Context, slug string) (image.Image, error)
}

type LaunchArgumentConfigurator interface {
	ShowLaunchArgumentScreen(title string, info string, arguments []game.LaunchParam, onSubmit func(values []game.LaunchArg))
}

type GameTileUpdate struct {
	Icon         image.Image
	CatalogItem  game.CatalogItem
	Installation game.Installation
	Progress     game.InstallationProgress
}

type GameTileCreate struct {
	Icon         image.Image
	CatalogItem  game.CatalogItem
	Installation game.Installation
}

type GameTile struct {
	slug string

	icon              binding.Untyped
	name              binding.String
	installedFileSize binding.String

	statusText  binding.String
	statusColor binding.Untyped

	showProgressIndicator binding.Bool
	Progress              binding.Float
	progressFrontText     binding.String
	progressEndText       binding.String

	enableSingleplayerButton    binding.Bool
	enableJoinMultiplayerButton binding.Bool
	enableHostMultiplayerButton binding.Bool
	enableOpenInExplorerButton  binding.Bool

	showInstallationStopIcon binding.Bool
	isInstalling             binding.Bool

	logger                      zerolog.Logger
	launchRunner                *game.LaunchRunner
	launchArgumentConfigurator  LaunchArgumentConfigurator
	installationRunner          *game.InstallationRunner
	installationDirectoryOpener *game.InstallationDirectoryOpener

	mu sync.RWMutex
}

func NewGameTileFromModel(create GameTileCreate, logger zerolog.Logger, launchRunner *game.LaunchRunner, installationRunner *game.InstallationRunner, installationDirectoryOpener *game.InstallationDirectoryOpener, launchArgumentConfigurator LaunchArgumentConfigurator) (*GameTile, error) {
	vm := &GameTile{
		slug:                        create.CatalogItem.Slug,
		logger:                      logger.With().Str("slug", create.CatalogItem.Slug).Logger(),
		icon:                        binding.NewUntyped(),
		name:                        binding.NewString(),
		installedFileSize:           binding.NewString(),
		statusText:                  binding.NewString(),
		statusColor:                 binding.NewUntyped(),
		showProgressIndicator:       binding.NewBool(),
		Progress:                    binding.NewFloat(),
		progressFrontText:           binding.NewString(),
		progressEndText:             binding.NewString(),
		enableSingleplayerButton:    binding.NewBool(),
		enableJoinMultiplayerButton: binding.NewBool(),
		enableHostMultiplayerButton: binding.NewBool(),
		enableOpenInExplorerButton:  binding.NewBool(),
		showInstallationStopIcon:    binding.NewBool(),
		isInstalling:                binding.NewBool(),
		launchRunner:                launchRunner,
		installationRunner:          installationRunner,
		installationDirectoryOpener: installationDirectoryOpener,
		launchArgumentConfigurator:  launchArgumentConfigurator,
	}

	err := vm.UpdateFromModel(GameTileUpdate{
		Icon:         create.Icon,
		CatalogItem:  create.CatalogItem,
		Installation: create.Installation,
		Progress:     game.InstallationProgress{},
	})
	if err != nil {
		return nil, fmt.Errorf("could not update game tile from model: %w", err)
	}

	return vm, nil
}

func (vm *GameTile) UpdateFromModel(update GameTileUpdate) error {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	if vm.slug != update.CatalogItem.Slug || vm.slug != update.Installation.Slug || (update.Progress.Slug != "" && vm.slug != update.Progress.Slug) {
		return fmt.Errorf("slug mismatch: want %s - got catalog item %s, installation %s, progress %s", vm.slug, update.CatalogItem.Slug, update.Installation.Slug, update.Progress.Slug)
	}

	vm.icon.Set(update.Icon)
	vm.name.Set(update.CatalogItem.Name)
	vm.installedFileSize.Set(humanize.Bytes(update.CatalogItem.InstalledFileSize))

	if update.Installation.IsDownloading() {
		vm.statusText.Set("Downloading")
		vm.statusColor.Set(theme.StatusColor(theme.StatusDownloading))
	} else if update.Installation.IsExtracting() {
		vm.statusText.Set("Extracting")
		vm.statusColor.Set(theme.StatusColor(theme.StatusExtracting))
	} else if update.Installation.HasFailed() {
		vm.statusText.Set("Installation failed")
		vm.statusColor.Set(theme.StatusColor(theme.StatusError))
	} else if update.Installation.WasCancelled() {
		vm.statusText.Set("Cancelled")
		vm.statusColor.Set(theme.StatusColor(theme.StatusCanceled))
	} else if update.Installation.IsInstalled() {
		vm.statusText.Set("Installed")
		vm.statusColor.Set(theme.StatusColor(theme.StatusReady))
	} else if !update.Installation.IsInstalled() {
		vm.statusText.Set("Not installed")
		vm.statusColor.Set(theme.StatusColor(theme.StatusNotReady))
	}

	vm.showProgressIndicator.Set(update.Installation.IsInstalling())
	if update.Installation.IsInstalling() {
		vm.Progress.Set(update.Progress.Progress())
		if update.Installation.IsDownloading() {
			vm.progressFrontText.Set(fmt.Sprintf("Downloading... %s / %s", humanize.SIWithDigits(float64(update.Progress.Completed), 2, "B"), humanize.SIWithDigits(float64(update.Progress.Total), 2, "B")))
			vm.progressEndText.Set(fmt.Sprintf("%s", humanize.SIWithDigits(float64(update.Progress.Speed()), 2, "B/s")))
		} else if update.Installation.IsExtracting() {
			vm.progressFrontText.Set(fmt.Sprintf("Extracting... %d / %d Files", update.Progress.Completed, update.Progress.Total))
			vm.progressEndText.Set(fmt.Sprintf("%d Files/s", update.Progress.Speed()))
		} else {
			vm.progressFrontText.Set("")
			vm.progressEndText.Set("")
		}
	}

	vm.enableSingleplayerButton.Set(update.Installation.IsIdle() && update.Installation.IsInstalled() && !update.Installation.HasFailed())
	vm.enableJoinMultiplayerButton.Set(update.Installation.IsIdle() && update.Installation.IsInstalled() && !update.Installation.HasFailed() && update.CatalogItem.Capabilities.JoiningMultiplayer)
	vm.enableHostMultiplayerButton.Set(update.Installation.IsIdle() && update.Installation.IsInstalled() && !update.Installation.HasFailed() && update.CatalogItem.Capabilities.HostingServer)
	vm.enableOpenInExplorerButton.Set(update.Installation.IsInstalled())

	vm.showInstallationStopIcon.Set(update.Installation.IsInstalling())
	vm.isInstalling.Set(update.Installation.IsInstalling())

	return nil
}

func (vm *GameTile) StartSingleplayer() {
	launchSpec, err := vm.launchRunner.GetLaunchSpec(context.Background(), vm.slug, game.LaunchSpecModePlay)
	if err != nil {
		vm.logger.Error().Err(err).Msg("get launch spec for singleplayer failed")
		return
	}

	vm.launchArgumentConfigurator.ShowLaunchArgumentScreen(fmt.Sprintf("Start · %s", vm.GetName()), launchSpec.ExecutablePathRelative, launchSpec.Params(), func(values []game.LaunchArg) {
		if err := vm.launchRunner.StartSingleplayer(context.Background(), vm.slug, values); err != nil {
			vm.logger.Error().Err(err).Msg("start singleplayer failed")
		}
	})
}

func (vm *GameTile) OpenUserSelectionToJoinMultiplayer() {
	launchSpec, err := vm.launchRunner.GetLaunchSpec(context.Background(), vm.slug, game.LaunchSpecModeJoin)
	if err != nil {
		vm.logger.Error().Err(err).Msg("get launch spec for join multiplayer failed")
		return
	}

	vm.launchArgumentConfigurator.ShowLaunchArgumentScreen(fmt.Sprintf("Join · %s", vm.GetName()), launchSpec.ExecutablePathRelative, launchSpec.Params(), func(values []game.LaunchArg) {
		if err := vm.launchRunner.JoinMultiplayer(context.Background(), vm.slug, values); err != nil {
			vm.logger.Error().Err(err).Msg("join multiplayer failed")
		}
	})
}

func (vm *GameTile) OpenArgumentConfigurationToHostMultiplayer() {
	launchSpec, err := vm.launchRunner.GetLaunchSpec(context.Background(), vm.slug, game.LaunchSpecModeHost)
	if err != nil {
		vm.logger.Error().Err(err).Msg("get launch spec for host multiplayer failed")
		return
	}

	vm.launchArgumentConfigurator.ShowLaunchArgumentScreen(fmt.Sprintf("Host · %s", vm.GetName()), launchSpec.ExecutablePathRelative, launchSpec.Params(), func(values []game.LaunchArg) {
		if err := vm.launchRunner.HostMultiplayer(context.Background(), vm.slug, values); err != nil {
			vm.logger.Error().Err(err).Msg("host multiplayer failed")
		}
	})
}

func (vm *GameTile) ToggleInstallation() {
	vm.mu.RLock()
	isInstalling, err := vm.isInstalling.Get()
	if err != nil {
		return
	}
	vm.mu.RUnlock()

	if isInstalling {
		if err := vm.installationRunner.CancelInstallation(context.Background(), vm.slug); err != nil {
			vm.logger.Error().Err(err).Msg("cancel installation failed")
		}
	} else {
		go func() {
			if err := vm.installationRunner.StartInstallation(context.Background(), vm.slug); err != nil {
				vm.logger.Error().Err(err).Msg("start installation failed")
			}
		}()
	}
}

func (vm *GameTile) OpenDirectoryInExplorer() {
	vm.installationDirectoryOpener.OpenDirectory(context.Background(), vm.slug)
}

func (vm *GameTile) AddChangeListener(fn func()) {
	l := binding.NewDataListener(fn)

	vm.icon.AddListener(l)
	vm.name.AddListener(l)
	vm.installedFileSize.AddListener(l)

	vm.statusText.AddListener(l)
	vm.statusColor.AddListener(l)

	vm.showProgressIndicator.AddListener(l)
	vm.Progress.AddListener(l)
	vm.progressFrontText.AddListener(l)
	vm.progressEndText.AddListener(l)

	vm.enableSingleplayerButton.AddListener(l)
	vm.enableJoinMultiplayerButton.AddListener(l)
	vm.enableHostMultiplayerButton.AddListener(l)
	vm.enableOpenInExplorerButton.AddListener(l)

	vm.showInstallationStopIcon.AddListener(l)
	vm.isInstalling.AddListener(l)
}

func (vm *GameTile) GetSlug() string {
	return vm.slug
}

func (vm *GameTile) GetIcon() image.Image {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	fallbackImage := image.NewRGBA(image.Rect(0, 0, 64, 64))
	draw.Draw(fallbackImage, fallbackImage.Bounds(), &image.Uniform{C: color.RGBA{R: 255, G: 0, B: 0, A: 255}}, image.Point{}, draw.Src)

	iv, err := vm.icon.Get()
	if err != nil || iv == nil {
		return fallbackImage
	}

	img, ok := iv.(image.Image)
	if !ok {
		return fallbackImage
	}

	return img
}

func (vm *GameTile) GetName() string {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	name, err := vm.name.Get()
	if err != nil || name == "" {
		return "N/A"
	}

	return name
}

func (vm *GameTile) GetInstalledFileSize() string {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	installedFileSize, err := vm.installedFileSize.Get()
	if err != nil || installedFileSize == "" {
		return "N/A GB"
	}

	return installedFileSize
}

func (vm *GameTile) GetStatusText() string {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	statusText, err := vm.statusText.Get()
	if err != nil || statusText == "" {
		return "N/A"
	}

	return statusText
}

func (vm *GameTile) GetStatusColor() color.Color {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	fallbackColor := color.RGBA{R: 255, G: 0, B: 0, A: 255}

	cv, err := vm.statusColor.Get()
	if err != nil || cv == nil {
		return fallbackColor
	}

	statusColor, ok := cv.(color.Color)
	if !ok {
		return fallbackColor
	}

	return statusColor
}

func (vm *GameTile) ShowProgressIndicator() bool {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	showProgressIndicator, err := vm.showProgressIndicator.Get()
	if err != nil {
		return false
	}

	return showProgressIndicator
}

func (vm *GameTile) GetProgressFrontText() string {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	progressFrontText, err := vm.progressFrontText.Get()
	if err != nil || progressFrontText == "" {
		return "N/A"
	}

	return progressFrontText
}

func (vm *GameTile) GetProgressEndText() string {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	progressEndText, err := vm.progressEndText.Get()
	if err != nil || progressEndText == "" {
		return "N/A"
	}

	return progressEndText
}

func (vm *GameTile) EnableSingleplayerButton() bool {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	enableSingleplayerButton, err := vm.enableSingleplayerButton.Get()
	if err != nil {
		return false
	}

	return enableSingleplayerButton
}

func (vm *GameTile) EnableJoinMultiplayerButton() bool {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	enableJoinMultiplayerButton, err := vm.enableJoinMultiplayerButton.Get()
	if err != nil {
		return false
	}

	return enableJoinMultiplayerButton
}

func (vm *GameTile) EnableHostMultiplayerButton() bool {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	enableHostMultiplayerButton, err := vm.enableHostMultiplayerButton.Get()
	if err != nil {
		return false
	}

	return enableHostMultiplayerButton
}

func (vm *GameTile) EnableOpenInExplorerButton() bool {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	enableOpenInExplorerButton, err := vm.enableOpenInExplorerButton.Get()
	if err != nil {
		return false
	}

	return enableOpenInExplorerButton
}

func (vm *GameTile) ShowInstallationStopIcon() bool {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	showInstallationStopIcon, err := vm.showInstallationStopIcon.Get()
	if err != nil {
		return false
	}

	return showInstallationStopIcon
}
