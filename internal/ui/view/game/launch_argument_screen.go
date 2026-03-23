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

	view *view.HeaderScrollScreen
}

func NewLaunchArgumentScreen(vm *gameviewmodel.LaunchArgumentScreen) *LaunchArgumentScreen {
	groupsContainer := container.NewVBox()

	view := &LaunchArgumentScreen{
		vm:   vm,
		view: view.NewHeaderScrollScreen(vm.Header, groupsContainer),
	}

	for _, group := range vm.GetArgumentGroups() {
		groupsContainer.Add(NewLaunchArgumentGroup(group))
	}

	view.ExtendBaseWidget(view)

	return view
}

func (view *LaunchArgumentScreen) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(view.view)
}
