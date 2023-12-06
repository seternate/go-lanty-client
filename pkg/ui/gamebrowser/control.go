package gamebrowser

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/seternate/go-lanty-client/pkg/controller"
	"github.com/seternate/go-lanty/pkg/game"
)

type Control struct {
	Container        *fyne.Container
	defaultcontainer *fyne.Container
	playcontainer    *fyne.Container

	controller *controller.GameController
	game       game.Game
}

func NewControl(controller *controller.GameController, game game.Game) (control *Control) {
	control = &Control{
		controller: controller,
		game:       game,
	}

	control.createDefaultView()
	control.createPlayView()

	control.Container = container.NewMax(
		control.playcontainer,
		control.defaultcontainer,
	)

	return
}

func (control *Control) createDefaultView() {
	play := widget.NewButton("Play", func() { control.showPlayContainer() })
	play.SetIcon(theme.MediaPlayIcon())

	open := widget.NewButton("Open", func() { control.controller.OpenGameInExplorer(control.game) })
	open.SetIcon(theme.FolderIcon())

	download := widget.NewButton("Download", func() { control.controller.DownloadGame(control.game) })
	download.SetIcon(theme.DownloadIcon())

	configure := widget.NewButton("Configure", func() {
		//TODO
	})
	configure.SetIcon(theme.SettingsIcon())

	control.defaultcontainer = container.NewGridWithRows(
		2,
		play,
		download,
		open,
		configure,
	)
}

func (control *Control) createPlayView() {
	start := widget.NewButton("Start", func() {
		control.controller.StartGame(control.game)
		control.showDefaultContainer()
	})
	start.SetIcon(theme.MediaPlayIcon())

	joinserver := widget.NewButton("Join server", func() {
		//TODO
	})
	joinserver.SetIcon(theme.LoginIcon())

	startserver := widget.NewButton("Start server", func() {
		//TODO
		//ui.Application.CreateWindow("Start server", NewStartServerView())
		// window := Application.NewWindow("Start Server")
		// window.SetContent(NewStartServerView(control.game))
		// window.
		// window.Show()
		control.showDefaultContainer()
	})
	startserver.SetIcon(theme.MailForwardIcon())

	cancel := widget.NewButton("Cancel", func() { control.showDefaultContainer() })
	cancel.SetIcon(theme.CancelIcon())

	control.playcontainer = container.NewGridWithRows(
		2,
		start,
		joinserver,
		startserver,
		cancel,
	)
}

func (control *Control) showPlayContainer() {
	control.defaultcontainer.Hide()
	control.playcontainer.Show()
}

func (control *Control) showDefaultContainer() {
	control.defaultcontainer.Show()
	control.playcontainer.Hide()
}
