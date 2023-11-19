package gamebrowser

import (
	"fmt"
	"image/color"
	"os"
	"path"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/seternate/go-lanty/pkg/settings"
	lanty "github.com/seternate/lanty-api-golang/pkg/api"
	"github.com/seternate/lanty-api-golang/pkg/download"
	"github.com/seternate/lanty-api-golang/pkg/game"
	"github.com/seternate/lanty-api-golang/pkg/util"
)

type Gametile struct {
	MainContainer    *fyne.Container
	infoContainer    *fyne.Container
	controlContainer *fyne.Container
	game             game.Game
	name             *widget.Label
	version          *widget.Label
	progressbar      *widget.ProgressBar
	icon             *canvas.Image
	client           *lanty.Client
	settings         *settings.Settings
	downloader       *download.Downloader
}

func NewGametile(client *lanty.Client, settings *settings.Settings, game game.Game, downloader *download.Downloader) (*Gametile, error) {
	gametile := &Gametile{
		game:       game,
		client:     client,
		settings:   settings,
		downloader: downloader,
	}

	gametile.createInfoContainer()
	gametile.createControlContainer()
	gametile.getIcon()

	gametile.MainContainer = container.NewMax(canvas.NewRectangle(color.RGBA{126, 126, 126, 255}),
		container.NewBorder(nil,
			nil,
			gametile.icon,
			nil,
			container.NewVBox(
				gametile.infoContainer,
				layout.NewSpacer(),
				container.NewPadded(
					gametile.controlContainer,
				),
			),
		),
	)

	return gametile, nil
}

func (gt *Gametile) getIcon() {
	image, _ := gt.client.Game.GetIcon(gt.game)
	icon := canvas.NewImageFromImage(image)
	icon.SetMinSize(fyne.NewSize(130, 130))

	gt.icon = icon
}

func (gt *Gametile) createInfoContainer() {
	gt.name = widget.NewLabel(gt.game.Name)
	gt.name.Alignment = fyne.TextAlignLeading

	gt.version = widget.NewLabel(gt.game.Version)
	gt.version.Alignment = fyne.TextAlignTrailing

	gt.progressbar = widget.NewProgressBar()
	gt.progressbar.Hide()

	gt.infoContainer = container.NewMax(
		container.NewHBox(gt.name,
			layout.NewSpacer(),
			gt.version,
		),
		container.NewPadded(gt.progressbar),
	)
}

func (gt *Gametile) createControlContainer() {
	buttonDownload := widget.NewButton("Download", func() { go gt.downloadGame() })
	buttonDownload.SetIcon(theme.DownloadIcon())

	buttonStart := widget.NewButton("Play", func() { gt.game.Start(gt.settings.GameDirectory) })
	buttonStart.SetIcon(theme.MediaPlayIcon())

	buttonOpenFile := widget.NewButton("Open", func() { gt.game.OpenFile(gt.settings.GameDirectory) })
	buttonOpenFile.SetIcon(theme.FolderIcon())

	buttonConfigure := widget.NewButton("Configure", func() {})
	buttonConfigure.SetIcon(theme.SettingsIcon())

	gt.controlContainer = container.NewGridWithRows(
		2,
		buttonStart,
		buttonDownload,
		buttonOpenFile,
		buttonConfigure,
	)
}

func (gt *Gametile) downloadGame() {
	if gt.downloader.IsDownloading(gt.game) {
		return
	}

	download, _ := gt.client.Game.GetFile(gt.game, gt.settings.GameDirectory)

	gt.downloader.Download[gt.game] = download

	gt.progressbar.SetValue(download.Progress())
	gt.progressbar.TextFormatter = func() string {
		return fmt.Sprintf("%.0f%% - %.0f MB/s", gt.progressbar.Value*100, download.BytesPerSecond()/(1024*1024))
	}
	gt.name.Hide()
	gt.version.Hide()
	gt.progressbar.Show()
	for !download.IsComplete() {
		time.Sleep(100 * time.Millisecond)
		gt.progressbar.SetValue(download.Progress())
	}
	gt.progressbar.TextFormatter = func() string { return "Extracting ..." }
	gt.progressbar.SetValue(1)
	util.ExtractZipFile(path.Join(gt.settings.GameDirectory, download.Filename), path.Join(gt.settings.GameDirectory, gt.game.Slug))
	e := os.Remove(path.Join(gt.settings.GameDirectory, download.Filename))
	if e != nil {
		fmt.Println(e)
	}
	gt.progressbar.Hide()
	gt.name.Show()
	gt.version.Show()
}
