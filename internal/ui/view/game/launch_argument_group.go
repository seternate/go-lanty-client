package gameview

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/seternate/go-lanty-client/internal/ui/theme"
	gameviewmodel "github.com/seternate/go-lanty-client/internal/ui/viewmodel/game"
)

type LaunchArgumentGroup struct {
	widget.BaseWidget

	vm *gameviewmodel.LaunchArgumentGroup

	view *fyne.Container
}

func NewLaunchArgumentGroup(vm *gameviewmodel.LaunchArgumentGroup) *LaunchArgumentGroup {
	headerText := canvas.NewText(vm.Header, theme.ForegroundColor())
	headerText.TextSize = theme.TextSize()
	headerText.TextStyle = fyne.TextStyle{Bold: true}

	background := canvas.NewRectangle(color.RGBA{104, 104, 104, 255})
	background.CornerRadius = theme.CornerRadius()

	padLeft := canvas.NewRectangle(color.Transparent)
	padLeft.SetMinSize(fyne.NewSize(theme.InnerPadding()/2, 0))
	padRight := canvas.NewRectangle(color.Transparent)
	padRight.SetMinSize(fyne.NewSize(theme.InnerPadding()/2, 0))
	padTop := canvas.NewRectangle(color.Transparent)
	padTop.SetMinSize(fyne.NewSize(0, theme.InnerPadding()/2))
	padBottom := canvas.NewRectangle(color.Transparent)
	padBottom.SetMinSize(fyne.NewSize(0, theme.InnerPadding()/2))

	textContainer := container.NewBorder(padTop, padBottom, padLeft, padRight, headerText)
	headerContainer := container.NewStack(background, textContainer)

	tileContainer := container.NewVBox()
	for _, tile := range vm.Tiles {
		tileContainer.Add(NewLaunchArgumentTile(tile))
	}

	viewContainer := container.NewVBox()
	if vm.ShowHeader {
		viewContainer.Add(headerContainer)
		tilePadLeft := canvas.NewRectangle(color.Transparent)
		tilePadLeft.SetMinSize(fyne.NewSize(theme.InnerPadding(), 0))
		tileContainer = container.NewBorder(nil, nil, tilePadLeft, nil, tileContainer)
	}
	viewContainer.Add(tileContainer)

	view := &LaunchArgumentGroup{
		vm:   vm,
		view: viewContainer,
	}

	view.ExtendBaseWidget(view)

	return view
}

func (view *LaunchArgumentGroup) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(view.view)
}
