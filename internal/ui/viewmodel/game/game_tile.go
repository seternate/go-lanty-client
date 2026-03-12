package gameviewmodel

import (
	"context"
	"fmt"
	"image"

	"fyne.io/fyne/v2/data/binding"
	"github.com/dustin/go-humanize"
	"github.com/seternate/go-lanty-client/internal/game"
)

type Status string

const (
	StatusNotInstalled Status = "Not Installed"
	StatusInstalled    Status = "Installed"
	StatusDownloading  Status = "Downloading"
	StatusExtracting   Status = "Extracting"
	StatusCanceled     Status = "Canceled"
	StatusError        Status = "Error"
)

type IconFetcher interface {
	FetchIcon(ctx context.Context, slug string) (image.Image, error)
}

// type JoinMultiplayerNavigator interface {
// 	OpenUserSelection(name string, onSelected func(ipAddress string))
// }

// type HostMultiplayerNavigator interface {
// 	ShowHostArgumentConfigForm(viewmodel *ArgumentConfigForm, onSubmit func(values []game.LaunchArg))
// }

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

	showDownloadStopIcon binding.Bool

	iconFetcher  IconFetcher
	launchRunner *game.LaunchRunner
	// joinGameNavigator JoinMultiplayerNavigator
	// hostGameNavigator HostMultiplayerNavigator
}

func NewGameTileFromModel(catalogItem game.CatalogItem, installation game.Installation, progress game.InstallationProgress, iconFetcher IconFetcher, launchRunner *game.LaunchRunner) *GameTile {
	vm := &GameTile{
		slug:         catalogItem.Slug,
		iconFetcher:  iconFetcher,
		launchRunner: launchRunner,
	}

	vm.UpdateFromModel(catalogItem, installation, progress)

	return vm
}

func (vm *GameTile) UpdateFromModel(catalogItem game.CatalogItem, installation game.Installation, progress game.InstallationProgress) error {
	if vm.slug != catalogItem.Slug || vm.slug != installation.Slug || vm.slug != progress.Slug {
		return fmt.Errorf("slug mismatch")
	}

	icon, err := vm.iconFetcher.FetchIcon(context.Background(), catalogItem.Slug)
	if err != nil {
		return fmt.Errorf("could not fetch icon: %w", err)
	}

	vm.icon.Set(icon)
	vm.name.Set(catalogItem.Name)
	vm.installedFileSize.Set(humanize.Bytes(catalogItem.InstalledFileSize))

	vm.statusText.Set()
	vm.statusColor.Set()

	vm.showProgressIndicator.Set(installation.IsInstalling())
	vm.Progress.Set(progress.Progress())
	if installation.IsDownloading() {
		vm.progressFrontText.Set(fmt.Sprintf("Downloading... %s / %s", humanize.SIWithDigits(float64(progress.Completed), 2, "B"), humanize.SIWithDigits(float64(progress.Total), 2, "B")))
		vm.progressEndText.Set(fmt.Sprintf("%s", humanize.SIWithDigits(float64(progress.Speed()), 2, "B/s")))
	} else if installation.IsExtracting() {
		vm.progressFrontText.Set(fmt.Sprintf("Extracting... %d / %d Files", progress.Completed, progress.Total))
		vm.progressEndText.Set(fmt.Sprintf("%d Files/s", progress.Speed()))
	} else {
		vm.progressFrontText.Set("")
		vm.progressEndText.Set("")
	}

	vm.enableSingleplayerButton.Set(installation.IsIdle() && installation.IsInstalled())
	vm.enableJoinMultiplayerButton.Set(installation.IsIdle() && installation.IsInstalled() && catalogItem.Capabilities.JoiningMultiplayer)
	vm.enableHostMultiplayerButton.Set(installation.IsIdle() && installation.IsInstalled() && catalogItem.Capabilities.HostingServer)
	vm.enableOpenInExplorerButton.Set(installation.IsIdle() && installation.IsInstalled())

	vm.showDownloadStopIcon.Set(installation.IsInstalling())

	return nil
}

func (vm *GameTile) StartSingleplayer() {
	// vm.launchRunner.StartSingleplayer(context.Background(), vm.Slug)
}

func (vm *GameTile) OpenUserSelectionToJoinMultiplayer() {
	// name, err := vm.Name.Get()
	// if err != nil {
	// 	return
	// }

	// vm.joinGameNavigator.OpenUserSelection(name, func(ipAddress string) {
	// 	vm.launchRunner.JoinMultiplayer(context.Background(), vm.Slug, ipAddress)
	// })
}

func (vm *GameTile) OpenArgumentConfigurationToHostMultiplayer() {
	// launchSpec, err := vm.launchRunner.GetLaunchSpec(context.Background(), vm.Slug, game.LaunchSpecModeHost)
	// if err != nil {
	// 	return
	// }

	// configForm, err := NewArgumentConfigFormFromLaunchParams("Host Game", launchSpec.Params())
	// if err != nil {
	// 	return
	// }
	// configForm.OnSubmit = func() {
	// 	for _, field := range configForm.Fields {
	// }

	// vm.hostGameNavigator.ShowHostArgumentConfigForm(configForm, func(values []game.LaunchArg) {
	// 	vm.launchRunner.HostMultiplayer(context.Background(), vm.Slug, values)
	// })
}

func (vm *GameTile) StartDownload() {
	// controller.mu.RLock()
	// tile := controller.GameListModel.GetGameBySlug(slug)
	// if tile == nil {
	// 	return
	// }
	// isProgressing, _ := tile.IsProgressing.Get()
	// controller.mu.RUnlock()

	// if isProgressing {
	// 	controller.installationService.StopInstallation(slug)
	// } else {
	// 	controller.installationService.StartInstallation(context.Background(), slug)
	// }
}

func (vm *GameTile) OpenDirectoryInExplorer() {
	// controller.directoryOpener.OpenDirectory(slug)
}

func (vm *GameTile) OpenExtensions() {
	// visible, err := controller.extensionService.ExtensionsVisible(slug)
	// if err != nil || !visible {
	// 	return
	// }
	// catalogItem, err := controller.catalogReader.GetBySlug(slug)
	// if err != nil {
	// 	return
	// }
	// installation, err := controller.installationReader.GetBySlug(slug)
	// if err != nil || !installation.IsInstalled() {
	// 	return
	// }
	// installPath, _ := filepath.Abs(filepath.Join(controller.installationRoot, installation.InstallationDirectoryRelative))

	// ctx := context.Background()
	// listModel := model.NewExtensionListModel()
	// extensions, err := controller.extensionService.ListExtensions(ctx, slug)
	// if err != nil {
	// 	return
	// }
	// for _, ext := range extensions {
	// 	listModel.AddExtension(model.NewExtensionModel(ext.Filename, ext.Name, ext.Size, ext.Downloaded))
	// }

	// onUpload := func(localPath string) {
	// 	_ = controller.extensionService.UploadExtension(context.Background(), slug, localPath)
	// }
	// onDownload := func(ext *model.ExtensionModel) {
	// 	if err := controller.extensionService.DownloadExtension(context.Background(), slug, ext.Filename); err == nil {
	// 		ext.SetDownloaded(true)
	// 	}
	// }
	// controller.extensionsAdapter.ShowExtensions(catalogItem.Name, slug, installPath, listModel, onUpload, onDownload)
}

func (vm *GameTile) AddChangeListener(fn func()) {
	l := binding.NewDataListener(fn)
	vm.Name.AddListener(l)
	vm.Icon.AddListener(l)
	vm.InstalledFileSize.AddListener(l)
	vm.SupportsJoiningMultiplayer.AddListener(l)
	vm.SupportsHostingServer.AddListener(l)
	vm.InstallationDetected.AddListener(l)
	vm.IsProgressing.AddListener(l)
	vm.Progress.AddListener(l)
	vm.ProgressFrontText.AddListener(l)
	vm.ProgressEndText.AddListener(l)
	vm.StatusText.AddListener(l)
	vm.StatusColor.AddListener(l)
	vm.HasNewExtensions.AddListener(l)
	vm.ExtensionsVisible.AddListener(l)
}

func (vm *GameTile) GetIcon() image.Image {
	iv, err := vm.Icon.Get()
	if err != nil {
		return nil
	}
	if img, ok := iv.(image.Image); ok {
		return img
	}
	return nil
}

func (vm *GameTile) GetName() string {
	name, err := vm.Name.Get()
	if err != nil {
		return "N/A"
	}
	return name
}

func (vm *GameTile) GetInstalledFileSize() string {
	installedFileSize, err := vm.InstalledFileSize.Get()
	if err != nil || installedFileSize == "" {
		return "N/A GB"
	}

	return installedFileSize
}

// func (vm *GameTile) SetStatusNotInstalled() {
// 	vm.StatusText.Set(string(StatusNotInstalled))
// 	vm.StatusColor.Set(theme.StatusColor(theme.StatusNotReady))
// 	vm.IsProgressing.Set(false)
// }

// func (model *GameTile) SetStatusDownloading() {
// 	model.StatusText.Set(string(StatusDownloading))
// 	model.StatusColor.Set(theme.StatusColor(theme.StatusDownloading))
// 	model.IsProgressing.Set(true)
// }

// func (model *GameTile) SetStatusExtracting() {
// 	model.StatusText.Set(string(StatusExtracting))
// 	model.StatusColor.Set(theme.StatusColor(theme.StatusExtracting))
// 	model.IsProgressing.Set(true)
// }

// func (model *GameTile) SetStatusInstalled() {
// 	model.StatusText.Set(string(StatusInstalled))
// 	model.StatusColor.Set(theme.StatusColor(theme.StatusReady))
// 	model.IsProgressing.Set(false)
// }

// func (model *GameTile) SetStatusError() {
// 	model.StatusText.Set(string(StatusError))
// 	model.StatusColor.Set(theme.StatusColor(theme.StatusError))
// 	model.IsProgressing.Set(false)
// }

// func (model *GameTile) SetStatusCanceled() {
// 	model.StatusText.Set(string(StatusCanceled))
// 	model.StatusColor.Set(theme.StatusColor(theme.StatusCanceled))
// 	model.IsProgressing.Set(false)
// }

// func (model *GameTile) MarkInstallationAsDetected() {
// 	model.InstallationDetected.Set(true)
// }

// func (model *GameTile) MarkInstallationAsRemoved() {
// 	model.InstallationDetected.Set(false)
// }

// func (model *GameTile) SetHasNewExtensions(hasNew bool) {
// 	model.HasNewExtensions.Set(hasNew)
// }

// func (model *GameTile) SetExtensionsVisible(visible bool) {
// 	model.ExtensionsVisible.Set(visible)
// }
