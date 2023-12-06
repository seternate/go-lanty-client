package widget

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty-client/pkg/controller"
	"github.com/seternate/go-lanty-client/pkg/theme"
)

type DownloadTile struct {
	widget.BaseWidget

	controller  *controller.Controller
	download    *controller.Download
	progressbar *widget.ProgressBar

	downloadstatusupdated chan struct{}
	progress              chan struct{}
}

func NewDownloadTile(download *controller.Download) (downloadtile *DownloadTile) {
	downloadtile = &DownloadTile{
		controller:            download.Controller(),
		download:              download,
		progressbar:           widget.NewProgressBar(),
		downloadstatusupdated: make(chan struct{}, 50),
		progress:              make(chan struct{}, 50),
	}
	downloadtile.ExtendBaseWidget(downloadtile)

	downloadtile.progressbar.TextFormatter = func() string {
		if !download.IsStarted() && !download.IsStopped() {
			return "Queued"
		} else if download.Err() != nil {
			return download.Err().Error()
		} else if download.IsComplete() {
			return fmt.Sprintf("%.0f%%", downloadtile.progressbar.Value*100)
		} else if download.BytesPerSecond() < 1024*1024 {
			return fmt.Sprintf("%.0f%% (%.0f KB/s)", downloadtile.progressbar.Value*100, download.BytesPerSecond()/1024)
		}
		return fmt.Sprintf("%.0f%% (%.0f MB/s)", downloadtile.progressbar.Value*100, download.BytesPerSecond()/(1024*1024))
	}
	download.Subscribe(downloadtile.downloadstatusupdated)
	download.SubscribeProgress(downloadtile.progress)

	downloadtile.run()

	return downloadtile
}

func (widget *DownloadTile) run() {
	widget.controller.WaitGroup().Add(2)
	go widget.downloadStatusUpdater()
	go widget.progressUpdater()
}

func (widget *DownloadTile) downloadStatusUpdater() {
	defer widget.controller.WaitGroup().Done()
	for !widget.download.IsComplete() {
		select {
		case <-widget.controller.Context().Done():
			log.Trace().Msg("exiting downloadtile downloadStatusUpdater()")
			return
		case <-widget.downloadstatusupdated:
			widget.Refresh()
		}
	}
	widget.Refresh()
}

func (widget *DownloadTile) progressUpdater() {
	defer widget.controller.WaitGroup().Done()
	for !widget.download.IsComplete() {
		select {
		case <-widget.controller.Context().Done():
			log.Trace().Msg("exiting downloadtile downloadStatusUpdater()")
			return
		case <-widget.progress:
			widget.progressbar.SetValue(widget.download.Progress())
			widget.Refresh()
		}
	}
}

func (widget *DownloadTile) CreateRenderer() fyne.WidgetRenderer {
	return newDownloadTileRenderer(widget)
}

type downloadTileRenderer struct {
	widget     *DownloadTile
	background *canvas.Rectangle
	name       *canvas.Text
	filesize   *canvas.Text
	starttime  *canvas.Text
	duration   *canvas.Text
}

func newDownloadTileRenderer(widget *DownloadTile) (renderer *downloadTileRenderer) {
	renderer = &downloadTileRenderer{
		widget:     widget,
		background: canvas.NewRectangle(fynetheme.InputBackgroundColor()),
		name:       canvas.NewText(widget.download.Game().Name, theme.ForegroundColor()),
		filesize:   canvas.NewText("---", theme.ForegroundColor()),
		starttime:  canvas.NewText("---", theme.ForegroundColor()),
		duration:   canvas.NewText("---", theme.ForegroundColor()),
	}
	renderer.background.CornerRadius = fynetheme.SelectionRadiusSize()
	return
}

func (renderer *downloadTileRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		renderer.background,
		renderer.name,
		renderer.filesize,
		renderer.starttime,
		renderer.duration,
		renderer.widget.progressbar,
	}
}

func (renderer *downloadTileRenderer) Layout(size fyne.Size) {
	renderer.background.Resize(size)
	renderer.background.Move(fyne.NewPos(0, 0))

	nametextsize := fyne.MeasureText(renderer.name.Text, renderer.name.TextSize, renderer.name.TextStyle)
	renderer.widget.progressbar.Resize(fyne.NewSize((size.Width-2*theme.InnerPadding())/3, 1.5*nametextsize.Height))
	renderer.widget.progressbar.Move(fyne.NewPos(size.Width-theme.InnerPadding()-renderer.widget.progressbar.Size().Width, theme.InnerPadding()))

	renderer.name.Move(fyne.NewPos(theme.InnerPadding(), (size.Height-nametextsize.Height)/2))

	textfieldwidth := (size.Width - 3*theme.InnerPadding() - renderer.widget.progressbar.Size().Width - fyne.Max(nametextsize.Width, 250)) / 3
	for index, text := range []*canvas.Text{renderer.duration, renderer.starttime, renderer.filesize} {
		textsize := fyne.MeasureText(text.Text, text.TextSize, text.TextStyle)
		text.Move(fyne.NewPos(fyne.Max(nametextsize.Width, 250)+2*theme.InnerPadding()+float32(index)*(textfieldwidth+theme.InnerPadding()), (size.Height-textsize.Height)/2))
	}

}

func (renderer *downloadTileRenderer) MinSize() fyne.Size {
	minsize := fyne.NewSize(
		2*theme.InnerPadding()+renderer.widget.progressbar.MinSize().Width,
		1.5*fyne.MeasureText(renderer.name.Text, renderer.name.TextSize, renderer.name.TextStyle).Height+2*theme.InnerPadding(),
	)
	for _, text := range []*canvas.Text{renderer.name, renderer.filesize, renderer.starttime, renderer.duration} {
		minsize = minsize.AddWidthHeight(fyne.MeasureText(text.Text, text.TextSize, text.TextStyle).Width+theme.InnerPadding(), 0)
	}
	return minsize.AddWidthHeight(250, 0)
}

func (renderer *downloadTileRenderer) Refresh() {
	if renderer.widget.download.Filesize() < 1024*1024*1024 {
		renderer.filesize.Text = fmt.Sprintf("%.0f MB", float32(renderer.widget.download.Filesize())/float32(1024*1024))
	} else {
		renderer.filesize.Text = fmt.Sprintf("%.2f GB", float32(renderer.widget.download.Filesize())/float32(1024*1024*1024))
	}
	renderer.starttime.Text = renderer.widget.download.StartTime().Format("15:04:05")
	renderer.duration.Text = renderer.widget.download.Duration().Truncate(time.Second).String()
	renderer.filesize.Refresh()
	renderer.starttime.Refresh()
	renderer.duration.Refresh()
	renderer.widget.progressbar.Refresh()
	renderer.background.Refresh()
}

func (renderer *downloadTileRenderer) Destroy() {}
