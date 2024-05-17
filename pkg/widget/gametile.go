package widget

import (
	"context"
	"fmt"
	"image"
	"runtime"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty-client/pkg/controller"
	"github.com/seternate/go-lanty-client/pkg/theme"
	"github.com/seternate/go-lanty/pkg/filesystem"
	"github.com/seternate/go-lanty/pkg/game"
)

type GameTile struct {
	widget.BaseWidget

	controller  *controller.Controller
	game        game.Game
	icon        image.Image
	progressbar *widget.ProgressBar
	buttons     map[string]*widget.Button

	download              *controller.Download
	newdownload           chan struct{}
	downloadstatusupdated chan struct{}
	progress              chan struct{}
	context               context.Context

	OnJoinServerTapped  func(game game.Game)
	OnStartServerTapped func(game game.Game)
	OnCancelTapped      func()
}

func NewGameTile(context context.Context, controller *controller.Controller, game game.Game) (gametile *GameTile) {
	gametile = &GameTile{
		controller:  controller,
		game:        game,
		icon:        controller.Game.GetIcon(game),
		progressbar: widget.NewProgressBar(),
		buttons: map[string]*widget.Button{
			"play":        widget.NewButtonWithIcon("Play", fynetheme.MediaPlayIcon(), nil),
			"open":        widget.NewButtonWithIcon("Open", fynetheme.FolderIcon(), nil),
			"download":    widget.NewButtonWithIcon("Download", fynetheme.DownloadIcon(), nil),
			"configure":   widget.NewButtonWithIcon("Configure", fynetheme.SettingsIcon(), nil),
			"start":       widget.NewButtonWithIcon("Start", fynetheme.MediaPlayIcon(), nil),
			"joinserver":  widget.NewButtonWithIcon("Join Server", fynetheme.LoginIcon(), nil),
			"startserver": widget.NewButtonWithIcon("Start Server", fynetheme.MailForwardIcon(), nil),
			"cancel":      widget.NewButtonWithIcon("Cancel", fynetheme.CancelIcon(), nil),
		},
		newdownload:           make(chan struct{}, 50),
		downloadstatusupdated: make(chan struct{}, 50),
		progress:              make(chan struct{}, 50),
		context:               context,
	}
	gametile.ExtendBaseWidget(gametile)

	gametile.buttons["play"].OnTapped = func() { gametile.showPlayControls() }
	gametile.buttons["open"].OnTapped = func() { controller.Game.OpenGameInExplorer(game) }
	gametile.updatePlayAndOpenButtonStatus()
	if runtime.GOOS != "windows" {
		gametile.buttons["play"].Disable()
		gametile.buttons["open"].Disable()
	}
	gametile.buttons["download"].OnTapped = func() { controller.Download.Download(game) }
	gametile.buttons["configure"].OnTapped = func() {
		//TODO
	}
	//DELETE
	gametile.buttons["configure"].Disable()
	gametile.buttons["start"].OnTapped = func() {
		controller.Game.StartGame(game)
		gametile.showDefaultControls()
	}
	gametile.buttons["joinserver"].OnTapped = func() {
		if gametile.OnJoinServerTapped != nil {
			gametile.OnJoinServerTapped(gametile.game)
		}
		gametile.showDefaultControls()
	}
	if !game.Client.CanConnect() {
		gametile.buttons["joinserver"].Disable()
	}
	gametile.buttons["startserver"].OnTapped = func() {
		if gametile.OnStartServerTapped != nil {
			gametile.OnStartServerTapped(gametile.game)
		}
		gametile.showDefaultControls()
	}
	if !game.CanStartServer() {
		gametile.buttons["startserver"].Disable()
	}
	gametile.buttons["cancel"].OnTapped = func() { gametile.showDefaultControls() }

	gametile.progressbar.TextFormatter = func() string {
		if gametile.download != nil {
			if gametile.download.IsDownloading() {
				return fmt.Sprintf("%.0f%% (%.0f MB/s)", gametile.progressbar.Value*100, gametile.download.BytesPerSecond()/(1024*1024))
			} else if gametile.download.IsUnzipping() {
				return fmt.Sprintf("Extracting (%.0f%% - %.0f MB/s)", gametile.progressbar.Value*100, gametile.download.BytesPerSecond()/(1024*1024))
			}
		}
		return fmt.Sprintf("%.0f%%", gametile.progressbar.Value*100)
	}

	//Needed because hiding the playControls at the beginning the MinSize() of the buttons of the are reported wrongly (too little)
	//https://github.com/fyne-io/fyne/issues/4453
	gametile.Refresh()

	gametile.showDefaultControls()
	controller.Download.Subscribe(gametile.newdownload)
	gametile.run()

	return gametile
}

func (widget *GameTile) showDefaultControls() {
	widget.buttons["play"].Show()
	widget.buttons["open"].Show()
	widget.buttons["download"].Show()
	widget.buttons["configure"].Show()
	widget.buttons["start"].Hide()
	widget.buttons["joinserver"].Hide()
	widget.buttons["startserver"].Hide()
	widget.buttons["cancel"].Hide()
	widget.Refresh()
}

func (widget *GameTile) showPlayControls() {
	widget.buttons["play"].Hide()
	widget.buttons["open"].Hide()
	widget.buttons["download"].Hide()
	widget.buttons["configure"].Hide()
	widget.buttons["start"].Show()
	widget.buttons["joinserver"].Show()
	widget.buttons["startserver"].Show()
	widget.buttons["cancel"].Show()
	widget.Refresh()
}

func (widget *GameTile) CreateRenderer() fyne.WidgetRenderer {
	return newGametileRenderer(widget)
}

func (widget *GameTile) run() {
	widget.controller.WaitGroup().Add(4)
	go widget.downloadUpdater()
	go widget.downloadStatusUpdater()
	go widget.progressUpdater()
	go widget.gameAvailabilityUpdater()
}

// TODO event based model should be used
func (widget *GameTile) gameAvailabilityUpdater() {
	defer widget.controller.WaitGroup().Done()
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-widget.context.Done():
			log.Trace().Str("slug", widget.game.Slug).Msg("exiting gametile gameAvailabilityUpdater()")
			return
		case <-ticker.C:
			widget.updatePlayAndOpenButtonStatus()
		}
	}
}

func (widget *GameTile) updatePlayAndOpenButtonStatus() {
	if runtime.GOOS != "windows" {
		return
	}
	paths, err := filesystem.SearchFilesBreadthFirst(widget.controller.Settings.Settings().GameDirectory, widget.game.Client.Executable, 3, 1)
	if (err != nil || len(paths) == 0) && !widget.buttons["play"].Disabled() {
		widget.buttons["play"].Disable()
		widget.buttons["open"].Disable()
		widget.Refresh()
	} else if widget.buttons["play"].Disabled() && err == nil && len(paths) > 0 {
		widget.buttons["play"].Enable()
		widget.buttons["open"].Enable()
		widget.Refresh()
	}
}

func (widget *GameTile) downloadUpdater() {
	defer widget.controller.WaitGroup().Done()
	for {
		select {
		case <-widget.context.Done():
			log.Trace().Str("slug", widget.game.Slug).Msg("exiting gametile downloadUpdater()")
			return
		case <-widget.newdownload:
			download, err := widget.controller.Download.GetLatest(widget.game)
			if err != nil {
				continue
			}
			if widget.download != nil {
				widget.download.Unsubscribe(widget.downloadstatusupdated)
				widget.download.UnsubscribeProgress(widget.progress)
			}
			download.Subscribe(widget.downloadstatusupdated)
			download.SubscribeProgress(widget.progress)
			widget.download = download
			widget.Refresh()
		}
	}
}

func (widget *GameTile) downloadStatusUpdater() {
	defer widget.controller.WaitGroup().Done()
	for {
		select {
		case <-widget.context.Done():
			log.Trace().Str("slug", widget.game.Slug).Msg("exiting gametile downloadStatusUpdater()")
			return
		case <-widget.downloadstatusupdated:
			widget.Refresh()
		}
	}
}

func (widget *GameTile) progressUpdater() {
	defer widget.controller.WaitGroup().Done()
	for {
		select {
		case <-widget.context.Done():
			log.Trace().Str("slug", widget.game.Slug).Msg("exiting gametile progressUpdater()")
			return
		case <-widget.progress:
			widget.progressbar.SetValue(widget.download.Progress())
		}
	}
}

type gametileRenderer struct {
	widget     *GameTile
	background *canvas.Rectangle
	icon       *canvas.Image
	name       *canvas.Text
	objects    []fyne.CanvasObject
}

func newGametileRenderer(widget *GameTile) *gametileRenderer {

	renderer := &gametileRenderer{
		widget:     widget,
		background: canvas.NewRectangle(fynetheme.InputBorderColor()),
		icon:       canvas.NewImageFromImage(widget.icon),
		name:       canvas.NewText(widget.game.Name, theme.ForegroundColor()),
	}
	renderer.background.CornerRadius = fynetheme.SelectionRadiusSize()
	renderer.objects = []fyne.CanvasObject{
		renderer.background,
		renderer.icon,
		renderer.name,
		renderer.widget.progressbar,
	}
	for _, button := range renderer.widget.buttons {
		renderer.objects = append(renderer.objects, button)
	}
	renderer.name.TextSize = 14

	return renderer
}

func (renderer *gametileRenderer) Objects() []fyne.CanvasObject {
	return renderer.objects
}

func (renderer *gametileRenderer) Layout(size fyne.Size) {
	renderer.background.Resize(size)

	renderer.icon.Resize(fyne.NewSize(size.Height-theme.InnerPadding(), size.Height-2*theme.InnerPadding()))
	renderer.icon.Move(fyne.NewPos(theme.InnerPadding(), theme.InnerPadding()))
	iconright := renderer.icon.Size().Width + renderer.icon.Position().X + theme.InnerPadding()

	renderer.widget.progressbar.Resize(fyne.NewSize(size.Width-iconright-theme.InnerPadding(), size.Height/3-2*theme.InnerPadding()))
	renderer.widget.progressbar.Move(fyne.NewPos(iconright, theme.InnerPadding()))
	progressbarbottom := renderer.widget.progressbar.Size().Height + renderer.widget.progressbar.Position().Y + theme.InnerPadding()

	textsize := fyne.MeasureText(renderer.name.Text, renderer.name.TextSize, renderer.name.TextStyle)
	renderer.name.Move(fyne.NewPos(iconright, (progressbarbottom-textsize.Height)/2))

	buttonsize := fyne.NewSize(
		(size.Width-iconright-2*theme.InnerPadding())/2,
		(size.Height-progressbarbottom-2*theme.InnerPadding())/2,
	)
	for _, button := range renderer.widget.buttons {
		button.Resize(buttonsize)
	}

	renderer.widget.buttons["play"].Move(fyne.NewPos(iconright, progressbarbottom))
	renderer.widget.buttons["open"].Move(fyne.NewPos(size.Width-buttonsize.Width-theme.InnerPadding(), progressbarbottom))
	renderer.widget.buttons["download"].Move(fyne.NewPos(iconright, size.Height-buttonsize.Height-theme.InnerPadding()))
	renderer.widget.buttons["configure"].Move(fyne.NewPos(size.Width-buttonsize.Width-theme.InnerPadding(), size.Height-buttonsize.Height-theme.InnerPadding()))

	renderer.widget.buttons["start"].Move(renderer.widget.buttons["play"].Position())
	renderer.widget.buttons["joinserver"].Move(renderer.widget.buttons["open"].Position())
	renderer.widget.buttons["startserver"].Move(renderer.widget.buttons["download"].Position())
	renderer.widget.buttons["cancel"].Move(renderer.widget.buttons["configure"].Position())
}

func (renderer *gametileRenderer) MinSize() fyne.Size {
	nameheight := fyne.MeasureText(renderer.name.Text, renderer.name.TextSize, renderer.name.TextStyle).Height
	minbuttonwidth := float32(0)
	for _, button := range renderer.widget.buttons {
		minbuttonwidth = fyne.Max(minbuttonwidth, button.MinSize().Width)
	}
	minheight := fyne.Max(renderer.widget.progressbar.MinSize().Height, nameheight) + 4*theme.InnerPadding() + 2*renderer.widget.buttons["play"].MinSize().Height
	minwidth := minheight + 4*theme.InnerPadding() + 2*minbuttonwidth

	return fyne.NewSize(minwidth, minheight)
}

func (renderer *gametileRenderer) Refresh() {
	if renderer.widget.download != nil && renderer.widget.download.IsRunning() {
		renderer.name.Hide()
		renderer.widget.progressbar.Show()
	} else {
		renderer.widget.progressbar.Hide()
		renderer.name.Show()
	}

	if renderer.widget.download != nil && (!renderer.widget.download.IsStarted() && !renderer.widget.download.IsStopped()) {
		renderer.name.Text = "Queued"
	} else {
		renderer.name.Text = renderer.widget.game.Name
	}

	renderer.background.Refresh()
	renderer.icon.Refresh()
	renderer.name.Refresh()
	renderer.widget.progressbar.Refresh()
	for _, button := range renderer.widget.buttons {
		button.Refresh()
	}
}

func (renderer *gametileRenderer) Destroy() {}
