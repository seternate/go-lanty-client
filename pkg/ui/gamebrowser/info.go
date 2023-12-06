package gamebrowser

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/seternate/go-lanty-client/pkg/controller"
	"github.com/seternate/go-lanty/pkg/filesystem"
	"github.com/seternate/go-lanty/pkg/game"
	"github.com/seternate/go-lanty/pkg/network"
)

type Info struct {
	Container   *fyne.Container
	name        *widget.Label
	progressbar *widget.ProgressBar

	controller   *controller.GameController
	game         game.Game
	subscription controller.Subscription
	download     *network.Download
	unzip        *filesystem.Unzip
}

func NewInfo(controller *controller.GameController, game game.Game) (info *Info) {
	info = &Info{
		game:        game,
		controller:  controller,
		name:        widget.NewLabel(game.Name),
		progressbar: widget.NewProgressBar(),
	}
	info.subscription.Game = game
	info.subscription.Download = make(chan *network.Download)
	info.subscription.Unzip = make(chan *filesystem.Unzip)
	info.subscription.DownloadProgress = make(chan float64)
	info.subscription.UnzipProgress = make(chan float64)

	info.name.Alignment = fyne.TextAlignLeading

	info.progressbar.Hide()
	info.progressbar.TextFormatter = func() string {
		if info.download == nil {
			return string("nil")
		} else if info.download.IsComplete() {
			return fmt.Sprintf("Extract - %.0f%%", info.progressbar.Value*100)
		}
		return fmt.Sprintf("%.0f%% - %.0f MB/s", info.progressbar.Value*100, info.download.BytesPerSecond()/(1024*1024))
	}

	info.Container = container.NewMax(
		info.name,
		container.NewPadded(info.progressbar),
	)
	info.controller.SubscribeDownload(info.subscription)

	go info.run()

	return
}

func (info *Info) run() {
	for {
		info.download = <-info.subscription.Download
		info.showProgressbar()
		info.updateDownloadProgress()
		info.unzip = <-info.subscription.Unzip
		info.updateUnzipProgress()
		info.showName()
	}
}

func (info *Info) updateDownloadProgress() {
	for !info.download.IsComplete() {
		info.progressbar.SetValue(<-info.subscription.DownloadProgress)
	}
}

func (info *Info) updateUnzipProgress() {
	for !info.unzip.IsComplete() {
		info.progressbar.SetValue(<-info.subscription.UnzipProgress)
	}
}

func (info *Info) showProgressbar() {
	info.name.Hide()
	info.progressbar.Show()
}

func (info *Info) showName() {
	info.progressbar.Hide()
	info.name.Show()
}
