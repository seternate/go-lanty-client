package widget

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty-client/pkg/controller"
	"github.com/seternate/go-lanty-client/pkg/layout"
	"github.com/seternate/go-lanty/pkg/game"
)

type GameBrowser struct {
	widget.BaseWidget

	controller *controller.Controller
	gametiles  []*GameTile
	//joinserver  *JoinServer
	//startserver *StartServer

	OnJoinServerTapped  func(game game.Game)
	OnStartServerTapped func(game game.Game)
	OnCancelTapped      func()

	gamesupdated          chan struct{}
	gametileContext       context.Context
	cancelGametileContext context.CancelFunc
}

func NewGameBrowser(controller *controller.Controller) (gamebrowser *GameBrowser) {
	gamebrowser = &GameBrowser{
		controller:   controller,
		gametiles:    make([]*GameTile, 0, 50),
		gamesupdated: make(chan struct{}, 50),
	}
	//gamebrowser.joinserver = NewJoinServer(controller, gamebrowser)
	//gamebrowser.startserver = NewStartServer(controller, gamebrowser)
	gamebrowser.gametileContext, gamebrowser.cancelGametileContext = context.WithCancel(controller.Context())
	gamebrowser.ExtendBaseWidget(gamebrowser)

	gamebrowser.updateGametiles()
	controller.Game.Subscribe(gamebrowser.gamesupdated)
	controller.WaitGroup().Add(1)
	go gamebrowser.gamesUpdater()

	return
}

func (widget *GameBrowser) setGametiles(gametiles ...*GameTile) {
	gametilesOld := widget.gametiles
	widget.gametiles = gametiles

	for index := range gametilesOld {
		gametilesOld[index] = nil
	}

	widget.Refresh()
}

func (widget *GameBrowser) CreateRenderer() fyne.WidgetRenderer {
	return newGameBrowserRenderer(widget)
}

func (widget *GameBrowser) gamesUpdater() {
	defer widget.controller.WaitGroup().Done()
	for {
		select {
		case <-widget.controller.Context().Done():
			log.Trace().Msg("exiting gamebrowser gamesUpdater()")
			return
		case <-widget.gamesupdated:
			widget.updateGametiles()
		}
	}
}

func (widget *GameBrowser) updateGametiles() {
	widget.cancelGametileContext()
	widget.gametileContext, widget.cancelGametileContext = context.WithCancel(widget.controller.Context())
	gametiles := make([]*GameTile, 0, len(widget.controller.Game.GetGames().Games()))
	for _, g := range widget.controller.Game.GetGames().Games() {
		gametile := NewGameTile(widget.gametileContext, widget.controller, g)
		gametiles = append(gametiles, gametile)
		gametile.OnJoinServerTapped = func(game game.Game) {
			if widget.OnJoinServerTapped != nil {
				widget.OnJoinServerTapped(game)
			}
		}
		gametile.OnStartServerTapped = func(game game.Game) {
			if widget.OnStartServerTapped != nil {
				widget.OnStartServerTapped(game)
			}
		}
		gametile.OnCancelTapped = func() {
			if widget.OnCancelTapped != nil {
				widget.OnCancelTapped()
			}
		}
	}
	widget.setGametiles(gametiles...)
}

type gameBrowserRenderer struct {
	widget    *GameBrowser
	gametiles *fyne.Container
}

func newGameBrowserRenderer(widget *GameBrowser) *gameBrowserRenderer {
	renderer := &gameBrowserRenderer{
		widget:    widget,
		gametiles: container.New(layout.NewGridScalingLayout(3)),
	}
	return renderer
}

func (renderer *gameBrowserRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{
		renderer.gametiles,
	}

	return objects
}

func (renderer *gameBrowserRenderer) Layout(size fyne.Size) {
	renderer.gametiles.Resize(size)
	renderer.gametiles.Move(fyne.NewPos(0, 0))
}

func (renderer *gameBrowserRenderer) MinSize() fyne.Size {
	return renderer.gametiles.MinSize()
}

func (renderer *gameBrowserRenderer) Refresh() {
	renderer.gametiles.RemoveAll()
	for _, gametile := range renderer.widget.gametiles {
		renderer.gametiles.Add(gametile)
	}
	renderer.gametiles.Refresh()
	renderer.Layout(renderer.widget.Size())
}

func (renderer *gameBrowserRenderer) Destroy() {}
