package settingsview

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/seternate/go-lanty-client/internal/ui/view"
	settingsviewmodel "github.com/seternate/go-lanty-client/internal/ui/viewmodel/settings"
)

type SettingsScreen struct {
	widget.BaseWidget

	vm *settingsviewmodel.SettingsScreen

	view           *view.HeaderScrollScreen
	tilesContainer *fyne.Container
	settingstiles  map[string]*SettingTile
}

func NewSettingsScreen(vm *settingsviewmodel.SettingsScreen) *SettingsScreen {
	settingsContainer := container.NewVBox()

	view := &SettingsScreen{
		vm:             vm,
		view:           view.NewHeaderScrollScreen(vm.Header, settingsContainer),
		tilesContainer: settingsContainer,
		settingstiles:  make(map[string]*SettingTile),
	}

	settingstilesvm := view.vm.GetSettingsTiles()

	for _, settingstilevm := range settingstilesvm {
		view.tilesContainer.Add(NewSettingTile(settingstilevm))
	}

	view.ExtendBaseWidget(view)
	return view
}

func (view *SettingsScreen) Hide() {
	view.vm.Save()

	view.BaseWidget.Hide()
}

func (view *SettingsScreen) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(view.view)
}
