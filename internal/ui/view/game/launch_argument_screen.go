package gameview

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/seternate/go-lanty-client/internal/ui/view"
	gameviewmodel "github.com/seternate/go-lanty-client/internal/ui/viewmodel/game"
)

type LaunchArgumentScreen struct {
	widget.BaseWidget

	vm *gameviewmodel.LaunchArgumentScreen

	view           *view.HeaderScrollScreen
	tilesContainer *fyne.Container
}

func NewLaunchArgumentScreen(vm *gameviewmodel.LaunchArgumentScreen) *LaunchArgumentScreen {
	argumentContainer := container.NewVBox()

	view := &LaunchArgumentScreen{
		vm:   vm,
		view: view.NewHeaderScrollScreen(vm.Header, argumentContainer),
	}

	for _, argumentTile := range vm.GetArgumentTiles() {
		argumentContainer.Add(NewLaunchArgumentTile(argumentTile))
	}

	view.ExtendBaseWidget(view)

	return view
}

func (view *LaunchArgumentScreen) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(view.view)
}
