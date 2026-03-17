package userview

import (
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/seternate/go-lanty-client/internal/ui/view"
	userviewmodel "github.com/seternate/go-lanty-client/internal/ui/viewmodel/user"
)

type UserScreen struct {
	widget.BaseWidget

	vm *userviewmodel.UserScreen

	view           *view.HeaderScrollScreen
	tilesContainer *fyne.Container
	usertiles      map[string]*UserTile
}

func NewUserScreen(vm *userviewmodel.UserScreen) *UserScreen {
	userContainer := container.NewVBox()

	view := &UserScreen{
		vm:             vm,
		view:           view.NewHeaderScrollScreen(vm.Header, userContainer),
		tilesContainer: userContainer,
		usertiles:      make(map[string]*UserTile),
	}

	vm.AddChangeListener(func() {
		view.Refresh()
	})

	view.ExtendBaseWidget(view)

	return view
}

func (view *UserScreen) Refresh() {
	view.refreshUserTiles()

	view.tilesContainer.Refresh()
	view.view.Refresh()

	view.BaseWidget.Refresh()
}

func (view *UserScreen) refreshUserTiles() {
	usertilesvm := view.vm.GetUserTiles()

	vmIPAddresses := make([]string, 0)
	for _, usertilevm := range usertilesvm {
		vmIPAddresses = append(vmIPAddresses, usertilevm.GetIPAddress())
	}

	for viewIPAddress := range view.usertiles {
		if !slices.Contains(vmIPAddresses, viewIPAddress) {
			delete(view.usertiles, viewIPAddress)
		}
	}

	for _, usertilevm := range usertilesvm {
		if _, found := view.usertiles[usertilevm.GetIPAddress()]; !found {
			view.usertiles[usertilevm.GetIPAddress()] = NewUserTile(usertilevm)
		}
	}

	view.tilesContainer.RemoveAll()
	for _, usertilevm := range usertilesvm {
		view.tilesContainer.Add(view.usertiles[usertilevm.GetIPAddress()])
	}
}

func (view *UserScreen) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(view.view)
}
