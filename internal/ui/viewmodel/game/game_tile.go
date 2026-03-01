package gameviewmodel

import (
	"fyne.io/fyne/v2/data/binding"
	"github.com/seternate/go-lanty-client/internal/game"
	"github.com/seternate/go-lanty-client/internal/ui/theme"
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

// type JoinMultiplayerNavigator interface {
// 	OpenUserSelection(name string, onSelected func(ipAddress string))
// }

// type HostMultiplayerNavigator interface {
// 	ShowHostArgumentConfigForm(viewmodel *ArgumentConfigForm, onSubmit func(values []game.LaunchArg))
// }

type GameTile struct {
	Slug                       string
	Name                       binding.String
	Icon                       binding.Untyped
	InstalledFileSize          binding.String
	InstallationDetected       binding.Bool
	SupportsJoiningMultiplayer binding.Bool
	SupportsHostingServer      binding.Bool

	IsProgressing     binding.Bool
	Progress          binding.Float
	ProgressFrontText binding.String
	ProgressEndText   binding.String
	StatusText        binding.String
	StatusColor       binding.Untyped

	HasNewExtensions  binding.Bool
	ExtensionsVisible binding.Bool

	launchRunner *game.LaunchRunner
	// joinGameNavigator JoinMultiplayerNavigator
	// hostGameNavigator HostMultiplayerNavigator
}

func NewGameTile(slug string, launchRunner *game.LaunchRunner) *GameTile {
	vm := &GameTile{
		Slug:                       slug,
		Name:                       binding.NewString(),
		Icon:                       binding.NewUntyped(),
		InstalledFileSize:          binding.NewString(),
		SupportsJoiningMultiplayer: binding.NewBool(),
		SupportsHostingServer:      binding.NewBool(),
		InstallationDetected:       binding.NewBool(),
		IsProgressing:              binding.NewBool(),
		Progress:                   binding.NewFloat(),
		ProgressFrontText:          binding.NewString(),
		ProgressEndText:            binding.NewString(),
		StatusText:                 binding.NewString(),
		StatusColor:                binding.NewUntyped(),
		HasNewExtensions:           binding.NewBool(),
		ExtensionsVisible:          binding.NewBool(),
		launchRunner:               launchRunner,
	}

	vm.Name.Set("N/A")
	vm.InstalledFileSize.Set("N/A GB")
	vm.SetStatusNotInstalled()

	return vm
}

func (vm *GameTile) StartSingleplayer() {
	// vm.launchRunner.StartSingleplayer(context.Background(), vm.Slug)
}

func (vm *GameTile) JoinMultiplayer() {
	// name, err := vm.Name.Get()
	// if err != nil {
	// 	return
	// }

	// vm.joinGameNavigator.OpenUserSelection(name, func(ipAddress string) {
	// 	vm.launchRunner.JoinMultiplayer(context.Background(), vm.Slug, ipAddress)
	// })
}

func (vm *GameTile) HostMultiplayer() {
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

func (vm *GameTile) DownloadGame(slug string) {
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

func (vm *GameTile) OpenGame(slug string) {
	// controller.directoryOpener.OpenDirectory(slug)
}

func (vm *GameTile) OpenExtensions(slug string) {
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

// func (model *GameTile) SetName(name string) {
// 	model.Name.Set(name)
// }

// func (model *GameTile) SetIcon(icon image.Image) {
// 	model.Icon.Set(icon)
// }

// func (model *GameTile) SetInstalledFileSize(installedFileSize uint64) {
// 	if installedFileSize == 0 {
// 		return
// 	}
// 	model.InstalledFileSize.Set(humanize.Bytes(installedFileSize))
// }

// func (model *GameTile) SetSupportsJoiningMultiplayer(supportsJoiningMultiplayer bool) {
// 	model.SupportsJoiningMultiplayer.Set(supportsJoiningMultiplayer)
// }

// func (model *GameTile) SetSupportsHostingServer(supportsHostingServer bool) {
// 	model.SupportsHostingServer.Set(supportsHostingServer)
// }

// func (model *GameTile) UpdateDownloadProgress(totalBytes uint64, downloadedBytes uint64, speed uint64) {
// 	model.Progress.Set(float64(downloadedBytes) / float64(totalBytes))
// 	model.ProgressFrontText.Set(fmt.Sprintf("Downloading... %s / %s", humanize.SIWithDigits(float64(downloadedBytes), 2, "B"), humanize.SIWithDigits(float64(totalBytes), 2, "B")))
// 	model.ProgressEndText.Set(fmt.Sprintf("%s", humanize.SIWithDigits(float64(speed), 2, "B/s")))
// }

// func (model *GameTile) UpdateExtractingProgress(totalFiles uint64, extractedFiles uint64, speed uint64) {
// 	model.Progress.Set(float64(extractedFiles) / float64(totalFiles))
// 	model.ProgressFrontText.Set(fmt.Sprintf("Extracting... %d / %d Files", extractedFiles, totalFiles))
// 	model.ProgressEndText.Set(fmt.Sprintf("%d Files/s", speed))
// }

func (vm *GameTile) SetStatusNotInstalled() {
	vm.StatusText.Set(string(StatusNotInstalled))
	vm.StatusColor.Set(theme.StatusColor(theme.StatusNotReady))
	vm.IsProgressing.Set(false)
}

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
