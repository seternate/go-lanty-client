package gamebrowser

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty-client/pkg/controller"
)

var Application fyne.App

type GameBrowser struct {
	//Container *container.Scroll
	Container *fyne.Container
	gametiles *fyne.Container

	controller        *controller.GameController
	controllerUpdated chan struct{}
}

func NewGameBrowser(controller *controller.GameController) (browser *GameBrowser) {
	browser = &GameBrowser{
		controller:        controller,
		controllerUpdated: make(chan struct{}),
	}

	browser.createContainer()
	controller.Subscribe(browser.controllerUpdated)
	browser.refresh()
	go browser.refreshLoop()

	return
}

func (browser *GameBrowser) createContainer() {
	browser.gametiles = container.NewGridWrap(fyne.NewSize(400, 130))
	overlay := container.NewMax(canvas.NewRectangle(color.RGBA{126, 126, 126, 255}),
		container.NewVBox(
			layout.NewSpacer(),
			widget.NewForm(
				widget.NewFormItem("Test", widget.NewLabel("TEST")),
				widget.NewFormItem("Test Entry:", widget.NewEntry()),
			),
			layout.NewSpacer(),
		),
	)
	//browser.Container = container.NewScroll(browser.gametiles)
	browser.Container = container.NewMax(
		container.NewScroll(browser.gametiles),
		overlay,
	)
}

func (browser *GameBrowser) refreshLoop() {
	for {
		<-browser.controllerUpdated
		log.Trace().Msg("Refreshing GameBrowser UI")
		browser.refresh()
	}
}

func (browser *GameBrowser) refresh() {
	browser.gametiles.RemoveAll()
	for _, slug := range browser.controller.GetGames().Slugs() {
		game, _ := browser.controller.GetGames().Get(slug)
		gametile := NewGametile(browser.controller, game)
		browser.gametiles.Add(gametile.Container)
	}
	browser.gametiles.Refresh()
}
