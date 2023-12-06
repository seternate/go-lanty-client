package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/rs/zerolog/log"
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
		newdownload:           make(chan struct{}, 50),
		downloadstatusupdated: make(chan struct{}, 50),
	}
	downloadbrowser.ExtendBaseWidget(downloadbrowser)

	controller.Download.Subscribe(downloadbrowser.newdownload)
	downloadbrowser.run()

	return downloadbrowser
}

func (widget *DownloadBrowser) run() {
	widget.controller.WaitGroup().Add(2)
	go widget.downloadUpdater()
	go widget.downloadStatusUpdater()
}

func (widget *DownloadBrowser) downloadUpdater() {
	defer widget.controller.WaitGroup().Done()
	for {
		select {
		case <-widget.controller.Context().Done():
			log.Trace().Msg("exiting downloadbrowser downloadUpdater()")
			return
		case <-widget.newdownload:
			download := widget.controller.Download.GetLastQueued()
			if len(widget.downloadtiles) > 0 && widget.downloadtiles[len(widget.downloadtiles)-1].download == download {
				continue
			}
			widget.downloadtiles = append(widget.downloadtiles, NewDownloadTile(download))
			download.Subscribe(widget.downloadstatusupdated)
			widget.Refresh()
		}
	}
}

func (widget *DownloadBrowser) downloadStatusUpdater() {
	defer widget.controller.WaitGroup().Done()
	for {
		select {
		case <-widget.controller.Context().Done():
			log.Trace().Msg("exiting downloadbrowser downloadStatusUpdater()")
			return
		case <-widget.downloadstatusupdated:
			widget.Refresh()
		}
	}
}

func (browser *DownloadBrowser) CreateRenderer() fyne.WidgetRenderer {
	return newDownloadBrowserRenderer(browser)
}

type downloadBrowserRenderer struct {
	widget                *DownloadBrowser
	downloadtileoffset    int
	queuedText            *canvas.Text
	downloadingText       *canvas.Text
	unzippingText         *canvas.Text
	finishedText          *canvas.Text
	queuedBackground      *canvas.Rectangle
	downloadingBackground *canvas.Rectangle
	unzippingBackground   *canvas.Rectangle
	finishedBackground    *canvas.Rectangle
	queued                []*DownloadTile
	downloading           []*DownloadTile
	unzipping             []*DownloadTile
	finished              []*DownloadTile
}

func newDownloadBrowserRenderer(widget *DownloadBrowser) fyne.WidgetRenderer {
	renderer := &downloadBrowserRenderer{
		widget:                widget,
		downloadtileoffset:    25,
		queuedText:            canvas.NewText("Queued", theme.ForegroundColor()),
		downloadingText:       canvas.NewText("Downloading", theme.ForegroundColor()),
		unzippingText:         canvas.NewText("Unzipping", theme.ForegroundColor()),
		finishedText:          canvas.NewText("Finished", theme.ForegroundColor()),
		queuedBackground:      canvas.NewRectangle(fynetheme.InputBorderColor()),
		downloadingBackground: canvas.NewRectangle(fynetheme.InputBorderColor()),
		unzippingBackground:   canvas.NewRectangle(fynetheme.InputBorderColor()),
		finishedBackground:    canvas.NewRectangle(fynetheme.InputBorderColor()),
	}
	renderer.queuedBackground.CornerRadius = fynetheme.SelectionRadiusSize()
	renderer.downloadingBackground.CornerRadius = fynetheme.SelectionRadiusSize()
	renderer.unzippingBackground.CornerRadius = fynetheme.SelectionRadiusSize()
	renderer.finishedBackground.CornerRadius = fynetheme.SelectionRadiusSize()
	for _, text := range []*canvas.Text{renderer.queuedText, renderer.downloadingText, renderer.unzippingText, renderer.finishedText} {
		text.TextSize = 18
		text.TextStyle.Bold = true
	}
	return renderer
}

func (renderer *downloadBrowserRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{
		renderer.queuedBackground,
		renderer.downloadingBackground,
		renderer.unzippingBackground,
		renderer.finishedBackground,
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
	queuedtextsize := fyne.MeasureText(renderer.queuedText.Text, renderer.queuedText.TextSize, renderer.queuedText.TextStyle)
	renderer.queuedBackground.Move(fyne.NewPos(theme.InnerPadding(), theme.InnerPadding()))
	renderer.queuedBackground.Resize(fyne.NewSize(size.Width-2*theme.InnerPadding(), queuedtextsize.Height+2*theme.InnerPadding()))
	renderer.queuedText.Move(fyne.NewPos(renderer.queuedBackground.Position().X+theme.InnerPadding(), renderer.queuedBackground.Position().Y+((renderer.queuedBackground.Size().Height-queuedtextsize.Height)/2)))
	for index, queued := range renderer.queued {
		queued.Resize(fyne.NewSize(size.Width-2*theme.InnerPadding()-float32(renderer.downloadtileoffset), queued.MinSize().Height))
		queued.Move(fyne.NewPos(renderer.queuedBackground.Position().X+float32(renderer.downloadtileoffset), renderer.queuedBackground.Size().Height+2*theme.InnerPadding()+float32(index)*(queued.MinSize().Height+theme.InnerPadding())))
	}
	queuedbottom := renderer.queuedBackground.Position().Y + renderer.queuedBackground.Size().Height
	if len(renderer.queued) > 0 {
		queuedbottom = renderer.queued[len(renderer.queued)-1].Position().Y + renderer.queued[len(renderer.queued)-1].Size().Height
	}
	queuedbottom += theme.InnerPadding()

	downloadingtextsize := fyne.MeasureText(renderer.downloadingText.Text, renderer.downloadingText.TextSize, renderer.downloadingText.TextStyle)
	renderer.downloadingBackground.Move(fyne.NewPos(theme.InnerPadding(), queuedbottom))
	renderer.downloadingBackground.Resize(fyne.NewSize(size.Width-2*theme.InnerPadding(), downloadingtextsize.Height+2*theme.InnerPadding()))
	renderer.downloadingText.Move(fyne.NewPos(renderer.downloadingBackground.Position().X+theme.InnerPadding(), renderer.downloadingBackground.Position().Y+((renderer.downloadingBackground.Size().Height-downloadingtextsize.Height)/2)))
	for index, downloading := range renderer.downloading {
		downloading.Resize(fyne.NewSize(size.Width-2*theme.InnerPadding()-float32(renderer.downloadtileoffset), downloading.MinSize().Height))
		downloading.Move(fyne.NewPos(renderer.downloadingBackground.Position().X+float32(renderer.downloadtileoffset), queuedbottom+renderer.downloadingBackground.Size().Height+theme.InnerPadding()+float32(index)*(downloading.MinSize().Height+theme.InnerPadding())))
	}
	downloadingbottom := queuedbottom + renderer.downloadingBackground.Size().Height
	if len(renderer.downloading) > 0 {
		downloadingbottom = renderer.downloading[len(renderer.downloading)-1].Position().Y + renderer.downloading[len(renderer.downloading)-1].Size().Height
	}
	downloadingbottom += theme.InnerPadding()

	unzippingtextsize := fyne.MeasureText(renderer.unzippingText.Text, renderer.unzippingText.TextSize, renderer.unzippingText.TextStyle)
	renderer.unzippingBackground.Move(fyne.NewPos(theme.InnerPadding(), downloadingbottom))
	renderer.unzippingBackground.Resize(fyne.NewSize(size.Width-2*theme.InnerPadding(), unzippingtextsize.Height+2*theme.InnerPadding()))
	renderer.unzippingText.Move(fyne.NewPos(renderer.unzippingBackground.Position().X+theme.InnerPadding(), renderer.unzippingBackground.Position().Y+((renderer.unzippingBackground.Size().Height-unzippingtextsize.Height)/2)))
	for index, unzipping := range renderer.unzipping {
		unzipping.Resize(fyne.NewSize(size.Width-2*theme.InnerPadding()-float32(renderer.downloadtileoffset), unzipping.MinSize().Height))
		unzipping.Move(fyne.NewPos(renderer.unzippingBackground.Position().X+float32(renderer.downloadtileoffset), downloadingbottom+renderer.unzippingBackground.Size().Height+theme.InnerPadding()+float32(index)*(unzipping.MinSize().Height+theme.InnerPadding())))
	}
	unzippingbottom := downloadingbottom + renderer.unzippingBackground.Size().Height
	if len(renderer.unzipping) > 0 {
		unzippingbottom = renderer.unzipping[len(renderer.unzipping)-1].Position().Y + renderer.unzipping[len(renderer.unzipping)-1].Size().Height
	}
	unzippingbottom += theme.InnerPadding()

	finishedtextsize := fyne.MeasureText(renderer.finishedText.Text, renderer.finishedText.TextSize, renderer.finishedText.TextStyle)
	renderer.finishedBackground.Move(fyne.NewPos(theme.InnerPadding(), unzippingbottom))
	renderer.finishedBackground.Resize(fyne.NewSize(size.Width-2*theme.InnerPadding(), finishedtextsize.Height+2*theme.InnerPadding()))
	renderer.finishedText.Move(fyne.NewPos(renderer.finishedBackground.Position().X+theme.InnerPadding(), renderer.finishedBackground.Position().Y+((renderer.finishedBackground.Size().Height-finishedtextsize.Height)/2)))
	for index, finished := range renderer.finished {
		finished.Resize(fyne.NewSize(size.Width-2*theme.InnerPadding()-float32(renderer.downloadtileoffset), finished.MinSize().Height))
		finished.Move(fyne.NewPos(renderer.finishedBackground.Position().X+float32(renderer.downloadtileoffset), unzippingbottom+renderer.finishedBackground.Size().Height+theme.InnerPadding()+float32(index)*(finished.MinSize().Height+theme.InnerPadding())))
	}
}

func (renderer *downloadBrowserRenderer) MinSize() fyne.Size {
	minwidth := float32(0)
	minheight := theme.InnerPadding()
	for _, text := range []*canvas.Text{renderer.queuedText, renderer.downloadingText, renderer.unzippingText, renderer.finishedText} {
		textsize := fyne.MeasureText(text.Text, text.TextSize, text.TextStyle)
		minwidth = fyne.Max(minwidth, textsize.Width+2*theme.InnerPadding())
		minheight += textsize.Height + 3*theme.InnerPadding()
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
		if !downloadtile.download.IsStarted() && !downloadtile.download.IsStopped() {
			renderer.queued = append(renderer.queued, downloadtile)
		} else if downloadtile.download.IsDownloading() {
			renderer.downloading = append(renderer.downloading, downloadtile)
		} else if downloadtile.download.IsUnzipping() {
			renderer.unzipping = append(renderer.unzipping, downloadtile)
		} else if downloadtile.download.IsComplete() || downloadtile.download.IsStopped() {
			renderer.finished = append(renderer.finished, downloadtile)
		}
		downloadtile.Refresh()
	}
}

func (renderer *downloadBrowserRenderer) Destroy() {}
