package gameview

import (
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/seternate/go-lanty-client/internal/ui/view"
	gameviewmodel "github.com/seternate/go-lanty-client/internal/ui/viewmodel/game"
)

type GameScreen struct {
	widget.BaseWidget

	vm *gameviewmodel.GameScreen

	view           *view.HeaderScrollScreen
	tilesContainer *fyne.Container
	gametiles      map[string]*GameTile
}

func NewGameScreen(vm *gameviewmodel.GameScreen) *GameScreen {
	gameContainer := container.NewVBox()

	view := &GameScreen{
		vm:             vm,
		view:           view.NewHeaderScrollScreen(vm.Header, gameContainer),
		tilesContainer: gameContainer,
		gametiles:      make(map[string]*GameTile),
	}

	vm.AddChangeListener(func() {
		view.Refresh()
	})

	view.ExtendBaseWidget(view)
	return view
}

func (view *GameScreen) Refresh() {
	view.refreshGameTiles()

	view.tilesContainer.Refresh()
	view.view.Refresh()

	view.BaseWidget.Refresh()
}

func (view *GameScreen) refreshGameTiles() {
	gametilesvm := view.vm.GetGameTiles()

	vmSlugs := make([]string, 0)
	for _, gametilevm := range gametilesvm {
		vmSlugs = append(vmSlugs, gametilevm.GetSlug())
	}

	for viewSlug := range view.gametiles {
		if !slices.Contains(vmSlugs, viewSlug) {
			delete(view.gametiles, viewSlug)
		}
	}

	for _, gametilevm := range gametilesvm {
		if _, found := view.gametiles[gametilevm.GetSlug()]; !found {
			view.gametiles[gametilevm.GetSlug()] = NewGameTile(gametilevm)
		}
	}

	view.tilesContainer.RemoveAll()
	for _, gametilevm := range gametilesvm {
		view.tilesContainer.Add(view.gametiles[gametilevm.GetSlug()])
	}
}

func (view *GameScreen) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(view.view)
}
