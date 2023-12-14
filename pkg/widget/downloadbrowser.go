package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/seternate/go-lanty-client/pkg/controller"
	"github.com/seternate/go-lanty-client/pkg/theme"
)

type DownloadBrowser struct {
	widget.BaseWidget

	controller    *controller.Controller
	downloadtiles []*DownloadTile

	newdownload           chan struct{}
	downloadstatusupdated chan struct{}
}

func NewDownloadBrowser(controller *controller.Controller) (downloadbrowser *DownloadBrowser) {
	downloadbrowser = &DownloadBrowser{
		controller:            controller,
		downloadtiles:         make([]*DownloadTile, 0),
		newdownload:           make(chan struct{}),
		downloadstatusupdated: make(chan struct{}),
	}
	downloadbrowser.ExtendBaseWidget(downloadbrowser)

	controller.Download.Subscribe(downloadbrowser.newdownload)
	downloadbrowser.run()

	return downloadbrowser
}

func (widget *DownloadBrowser) run() {
	go widget.downloadUpdater()
	go widget.downloadStatusUpdater()
}

func (widget *DownloadBrowser) downloadUpdater() {
	for {
		<-widget.newdownload
		download := widget.controller.Download.GetLastQueued()
		if len(widget.downloadtiles) > 0 && widget.downloadtiles[len(widget.downloadtiles)-1].download == download {
			continue
		}
		widget.downloadtiles = append(widget.downloadtiles, NewDownloadTile(download))
		download.Subscribe(widget.downloadstatusupdated)
		widget.Refresh()
	}
}

func (widget *DownloadBrowser) downloadStatusUpdater() {
	for {
		<-widget.downloadstatusupdated
		widget.Refresh()
	}
}

func (browser *DownloadBrowser) CreateRenderer() fyne.WidgetRenderer {
	return newDownloadBrowserRenderer(browser)
}

type downloadBrowserRenderer struct {
	widget          *DownloadBrowser
	queuedText      *canvas.Text
	downloadingText *canvas.Text
	unzippingText   *canvas.Text
	finishedText    *canvas.Text
	queued          []*DownloadTile
	downloading     []*DownloadTile
	unzipping       []*DownloadTile
	finished        []*DownloadTile
}

func newDownloadBrowserRenderer(widget *DownloadBrowser) fyne.WidgetRenderer {
	renderer := &downloadBrowserRenderer{
		widget:          widget,
		queuedText:      canvas.NewText("Queued", theme.ForegroundColor()),
		downloadingText: canvas.NewText("Downloading", theme.ForegroundColor()),
		unzippingText:   canvas.NewText("Unzipping", theme.ForegroundColor()),
		finishedText:    canvas.NewText("Finished", theme.ForegroundColor()),
	}
	for _, text := range []*canvas.Text{renderer.queuedText, renderer.downloadingText, renderer.unzippingText, renderer.finishedText} {
		text.TextSize = 18
		text.TextStyle.Bold = true
	}
	return renderer
}

func (renderer *downloadBrowserRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{
		renderer.queuedText,
		renderer.downloadingText,
		renderer.unzippingText,
		renderer.finishedText,
	}
	for _, downloadtile := range renderer.widget.downloadtiles {
		objects = append(objects, downloadtile)
	}
	return objects
}

func (renderer *downloadBrowserRenderer) Layout(size fyne.Size) {
	renderer.queuedText.Move(fyne.NewPos(theme.InnerPadding(), theme.InnerPadding()))
	queuedtextsize := fyne.MeasureText(renderer.queuedText.Text, renderer.queuedText.TextSize, renderer.queuedText.TextStyle)
	for index, queued := range renderer.queued {
		queued.Resize(fyne.NewSize(size.Width-2*theme.InnerPadding(), queued.MinSize().Height))
		queued.Move(fyne.NewPos(theme.InnerPadding(), queuedtextsize.Height+theme.InnerPadding()+float32(index)*(queued.MinSize().Height+theme.InnerPadding())))
	}
	queuedbottom := renderer.queuedText.Position().Y + queuedtextsize.Height
	if len(renderer.queued) > 0 {
		queuedbottom = renderer.queued[len(renderer.queued)-1].Position().Y + renderer.queued[len(renderer.queued)-1].Size().Height
	}
	queuedbottom += theme.InnerPadding()

	renderer.downloadingText.Move(fyne.NewPos(theme.InnerPadding(), queuedbottom))
	downloadingtextsize := fyne.MeasureText(renderer.downloadingText.Text, renderer.downloadingText.TextSize, renderer.downloadingText.TextStyle)
	for index, downloading := range renderer.downloading {
		downloading.Resize(fyne.NewSize(size.Width-2*theme.InnerPadding(), downloading.MinSize().Height))
		downloading.Move(fyne.NewPos(theme.InnerPadding(), queuedbottom+downloadingtextsize.Height+theme.InnerPadding()+float32(index)*(downloading.MinSize().Height+theme.InnerPadding())))
	}
	downloadingbottom := queuedbottom + downloadingtextsize.Height
	if len(renderer.downloading) > 0 {
		downloadingbottom = renderer.downloading[len(renderer.downloading)-1].Position().Y + renderer.downloading[len(renderer.downloading)-1].Size().Height
	}
	downloadingbottom += theme.InnerPadding()

	renderer.unzippingText.Move(fyne.NewPos(theme.InnerPadding(), downloadingbottom))
	unzippingtextsize := fyne.MeasureText(renderer.unzippingText.Text, renderer.unzippingText.TextSize, renderer.unzippingText.TextStyle)
	for index, unzipping := range renderer.unzipping {
		unzipping.Resize(fyne.NewSize(size.Width-2*theme.InnerPadding(), unzipping.MinSize().Height))
		unzipping.Move(fyne.NewPos(theme.InnerPadding(), downloadingbottom+unzippingtextsize.Height+theme.InnerPadding()+float32(index)*(unzipping.MinSize().Height+theme.InnerPadding())))
	}
	unzippingbottom := downloadingbottom + unzippingtextsize.Height
	if len(renderer.unzipping) > 0 {
		unzippingbottom = renderer.unzipping[len(renderer.unzipping)-1].Position().Y + renderer.unzipping[len(renderer.unzipping)-1].Size().Height
	}
	unzippingbottom += theme.InnerPadding()

	renderer.finishedText.Move(fyne.NewPos(theme.InnerPadding(), unzippingbottom))
	finishedtextsize := fyne.MeasureText(renderer.finishedText.Text, renderer.finishedText.TextSize, renderer.finishedText.TextStyle)
	for index, finished := range renderer.finished {
		finished.Resize(fyne.NewSize(size.Width-2*theme.InnerPadding(), finished.MinSize().Height))
		finished.Move(fyne.NewPos(theme.InnerPadding(), unzippingbottom+finishedtextsize.Height+theme.InnerPadding()+float32(index)*(finished.MinSize().Height+theme.InnerPadding())))
	}
}

func (renderer *downloadBrowserRenderer) MinSize() fyne.Size {
	minwidth := float32(0)
	minheight := theme.InnerPadding()
	for _, text := range []*canvas.Text{renderer.queuedText, renderer.downloadingText, renderer.unzippingText, renderer.finishedText} {
		textsize := fyne.MeasureText(text.Text, text.TextSize, text.TextStyle)
		minwidth = fyne.Max(minwidth, textsize.Width+2*theme.InnerPadding())
		minheight += textsize.Height + theme.InnerPadding()
	}
	for _, downloadtile := range renderer.widget.downloadtiles {
		minwidth = fyne.Max(minwidth, downloadtile.MinSize().Width+2*theme.InnerPadding())
		minheight += downloadtile.MinSize().Height + theme.InnerPadding()
	}
	return fyne.NewSize(minwidth, minheight)
}

func (renderer *downloadBrowserRenderer) Refresh() {
	renderer.queued = make([]*DownloadTile, 0)
	renderer.downloading = make([]*DownloadTile, 0)
	renderer.unzipping = make([]*DownloadTile, 0)
	renderer.finished = make([]*DownloadTile, 0)

	for _, downloadtile := range renderer.widget.downloadtiles {
		if !downloadtile.download.IsStarted() {
			renderer.queued = append(renderer.queued, downloadtile)
		} else if downloadtile.download.IsDownloading() {
			renderer.downloading = append(renderer.downloading, downloadtile)
		} else if downloadtile.download.IsUnzipping() {
			renderer.unzipping = append(renderer.unzipping, downloadtile)
		} else if downloadtile.download.IsComplete() {
			renderer.finished = append(renderer.finished, downloadtile)
		}
		downloadtile.Refresh()
	}
}

func (renderer *downloadBrowserRenderer) Destroy() {}
