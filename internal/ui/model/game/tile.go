package model

import (
	"fmt"
	"image"

	"fyne.io/fyne/v2/data/binding"
	"github.com/dustin/go-humanize"
	viewmodel "github.com/seternate/go-lanty-client/internal/ui/model"
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

var _ viewmodel.ChangeNotifier = (*GameTileModel)(nil)

type GameTileModel struct {
	Slug                       string
	Name                       binding.String
	Icon                       binding.Untyped
	BlobSizeText               binding.String
	InstallationDetected       binding.Bool
	SupportsJoiningMultiplayer binding.Bool
	SupportsHostingServer      binding.Bool

	IsProgressing     binding.Bool
	Progress          binding.Float
	ProgressFrontText binding.String
	ProgressEndText   binding.String
	StatusText        binding.String
	StatusColor       binding.Untyped
}

func NewGameTileModel(slug string) *GameTileModel {
	tile := &GameTileModel{
		Slug:                       slug,
		Name:                       binding.NewString(),
		Icon:                       binding.NewUntyped(),
		BlobSizeText:               binding.NewString(),
		SupportsJoiningMultiplayer: binding.NewBool(),
		SupportsHostingServer:      binding.NewBool(),
		InstallationDetected:       binding.NewBool(),
		IsProgressing:              binding.NewBool(),
		Progress:                   binding.NewFloat(),
		ProgressFrontText:          binding.NewString(),
		ProgressEndText:            binding.NewString(),
		StatusText:                 binding.NewString(),
		StatusColor:                binding.NewUntyped(),
	}

	tile.Name.Set("N/A")
	tile.BlobSizeText.Set("N/A GB")
	tile.SetStatusNotInstalled()

	return tile
}

func (model *GameTileModel) SetName(name string) {
	model.Name.Set(name)
}

func (model *GameTileModel) SetIcon(icon image.Image) {
	model.Icon.Set(icon)
}

func (model *GameTileModel) SetBlobSize(blobSize uint64) {
	if blobSize == 0 {
		return
	}
	model.BlobSizeText.Set(humanize.Bytes(blobSize))
}

func (model *GameTileModel) SetSupportsJoiningMultiplayer(supportsJoiningMultiplayer bool) {
	model.SupportsJoiningMultiplayer.Set(supportsJoiningMultiplayer)
}

func (model *GameTileModel) SetSupportsHostingServer(supportsHostingServer bool) {
	model.SupportsHostingServer.Set(supportsHostingServer)
}

func (model *GameTileModel) UpdateDownloadProgress(totalBytes uint64, downloadedBytes uint64, speed uint64) {
	model.Progress.Set(float64(downloadedBytes) / float64(totalBytes))
	model.ProgressFrontText.Set(fmt.Sprintf("Downloading... %s / %s", humanize.SIWithDigits(float64(downloadedBytes), 2, "B"), humanize.SIWithDigits(float64(totalBytes), 2, "B")))
	model.ProgressEndText.Set(fmt.Sprintf("%s", humanize.SIWithDigits(float64(speed), 2, "B/s")))
}

func (model *GameTileModel) UpdateExtractingProgress(totalFiles uint64, extractedFiles uint64, speed uint64) {
	model.Progress.Set(float64(extractedFiles) / float64(totalFiles))
	model.ProgressFrontText.Set(fmt.Sprintf("Extracting... %d / %d Files", extractedFiles, totalFiles))
	model.ProgressEndText.Set(fmt.Sprintf("%d Files/s", speed))
}

func (model *GameTileModel) SetStatusNotInstalled() {
	model.StatusText.Set(string(StatusNotInstalled))
	model.StatusColor.Set(theme.StatusColor(theme.StatusNotReady))
	model.IsProgressing.Set(false)
}

func (model *GameTileModel) SetStatusDownloading() {
	model.StatusText.Set(string(StatusDownloading))
	model.StatusColor.Set(theme.StatusColor(theme.StatusDownloading))
	model.IsProgressing.Set(true)
}

func (model *GameTileModel) SetStatusExtracting() {
	model.StatusText.Set(string(StatusExtracting))
	model.StatusColor.Set(theme.StatusColor(theme.StatusExtracting))
	model.IsProgressing.Set(true)
}

func (model *GameTileModel) SetStatusInstalled() {
	model.StatusText.Set(string(StatusInstalled))
	model.StatusColor.Set(theme.StatusColor(theme.StatusReady))
	model.IsProgressing.Set(false)
}

func (model *GameTileModel) SetStatusError() {
	model.StatusText.Set(string(StatusError))
	model.StatusColor.Set(theme.StatusColor(theme.StatusError))
	model.IsProgressing.Set(false)
}

func (model *GameTileModel) SetStatusCanceled() {
	model.StatusText.Set(string(StatusCanceled))
	model.StatusColor.Set(theme.StatusColor(theme.StatusCanceled))
	model.IsProgressing.Set(false)
}

func (model *GameTileModel) MarkInstallationAsDetected() {
	model.InstallationDetected.Set(true)
}

func (model *GameTileModel) MarkInstallationAsRemoved() {
	model.InstallationDetected.Set(false)
}

func (model *GameTileModel) AddChangeListener(fn func()) {
	l := binding.NewDataListener(fn)
	model.Name.AddListener(l)
	model.Icon.AddListener(l)
	model.InstallationDetected.AddListener(l)
	model.IsProgressing.AddListener(l)
	model.Progress.AddListener(l)
	model.ProgressFrontText.AddListener(l)
	model.ProgressEndText.AddListener(l)
	model.StatusText.AddListener(l)
	model.StatusColor.AddListener(l)
	model.BlobSizeText.AddListener(l)
	model.SupportsJoiningMultiplayer.AddListener(l)
	model.SupportsHostingServer.AddListener(l)
}
