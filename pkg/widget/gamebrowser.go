package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/seternate/go-lanty-client/pkg/controller"
	"github.com/seternate/go-lanty-client/pkg/layout"
	"github.com/seternate/go-lanty/pkg/user"
)

type GameBrowser struct {
	widget.BaseWidget

	controller  *controller.Controller
	gametiles   []*GameTile
	joinserver  *JoinServer
	startserver *StartServer

	gamesupdated chan struct{}
}

func NewGameBrowser(controller *controller.Controller) (gamebrowser *GameBrowser) {
	gamebrowser = &GameBrowser{
		controller: controller,
		gametiles:  make([]*GameTile, 0, 50),
		joinserver: NewJoinServer(),

		gamesupdated: make(chan struct{}),
	}
	gamebrowser.startserver = NewStartServer(controller, gamebrowser)
	gamebrowser.ExtendBaseWidget(gamebrowser)

	//TODO
	gamebrowser.joinserver.OnSubmit = func() {
		gamebrowser.ShowGametiles()
		controller.Game.JoinServer(gamebrowser.joinserver.game, user.User{})
	}
	gamebrowser.joinserver.OnCancel = func() {
		gamebrowser.ShowGametiles()
	}
	//TODO

	gamebrowser.ShowGametiles()
	gamebrowser.updateGametiles()
	controller.Game.Subscribe(gamebrowser.gamesupdated)
	go gamebrowser.gamesUpdater()

	return
}

func (widget *GameBrowser) showJoinServer() {
	widget.joinserver.Show()
	widget.Refresh()
}

func (widget *GameBrowser) showStartServer() {
	widget.startserver.Show()
	widget.Refresh()
}

func (widget *GameBrowser) ShowGametiles() {
	widget.joinserver.Hide()
	widget.startserver.Hide()
	widget.Refresh()
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
	for {
		<-widget.gamesupdated
		widget.updateGametiles()
	}
}

func (widget *GameBrowser) updateGametiles() {
	gametiles := make([]*GameTile, 0, len(widget.controller.Game.GetGames().Games()))
	for _, game := range widget.controller.Game.GetGames().Games() {
		gametile := NewGameTile(widget.controller, game)
		gametiles = append(gametiles, gametile)
		gametile.OnJoinServerTapped = func() {
			widget.joinserver.Update(gametile.game)
			widget.showJoinServer()
		}
		gametile.OnStartServerTapped = func() {
			widget.startserver.Update(gametile.game)
			widget.showStartServer()
		}
	}
	widget.setGametiles(gametiles...)
}

type gameBrowserRenderer struct {
	widget          *GameBrowser
	gametiles       *fyne.Container
	gametilesScroll *container.Scroll
}

func newGameBrowserRenderer(widget *GameBrowser) *gameBrowserRenderer {
	gametiles := container.New(layout.NewGridScalingLayout(3))
	renderer := &gameBrowserRenderer{
		widget:          widget,
		gametiles:       gametiles,
		gametilesScroll: container.NewVScroll(gametiles),
	}
	return renderer
}

func (renderer *gameBrowserRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{
		renderer.gametilesScroll,
		renderer.widget.joinserver,
		renderer.widget.startserver,
	}

	return objects
}

func (renderer *gameBrowserRenderer) Layout(size fyne.Size) {
	renderer.gametilesScroll.Resize(size)
	renderer.gametilesScroll.Move(fyne.NewPos(0, 0))
	renderer.widget.joinserver.Resize(size)
	renderer.widget.joinserver.Move(fyne.NewPos(0, 0))
	renderer.widget.startserver.Resize(size)
	renderer.widget.startserver.Move(fyne.NewPos(0, 0))
}

func (renderer *gameBrowserRenderer) MinSize() fyne.Size {
	if !renderer.widget.joinserver.Hidden {
		return renderer.widget.joinserver.MinSize()
	} else if !renderer.widget.startserver.Hidden {
		return renderer.widget.startserver.MinSize()
	} else {
		return renderer.gametilesScroll.MinSize()
	}
}

func (renderer *gameBrowserRenderer) Refresh() {
	renderer.gametiles.RemoveAll()
	for _, gametile := range renderer.widget.gametiles {
		renderer.gametiles.Add(gametile)
	}
	if !renderer.widget.joinserver.Hidden || !renderer.widget.startserver.Hidden {
		renderer.gametilesScroll.Hide()
	} else {
		renderer.gametilesScroll.Show()
	}
	renderer.gametiles.Refresh()
	renderer.gametilesScroll.Refresh()
	renderer.widget.joinserver.Refresh()
	renderer.widget.startserver.Refresh()
}

func (renderer *gameBrowserRenderer) Destroy() {}
