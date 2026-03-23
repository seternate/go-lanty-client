package gameview

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/seternate/go-lanty-client/internal/ui/theme"
	gameviewmodel "github.com/seternate/go-lanty-client/internal/ui/viewmodel/game"
)

type LaunchArgumentTile struct {
	widget.BaseWidget

	vm *gameviewmodel.LaunchArgumentTile

	enabled *widget.Check
	entry   *widget.Entry
	list    *widget.Select
}

func NewLaunchArgumentTile(vm *gameviewmodel.LaunchArgumentTile) *LaunchArgumentTile {
	options := make([]string, 0, len(vm.GetEnumValues()))
	for key := range vm.GetEnumValues() {
		options = append(options, key)
	}

	view := &LaunchArgumentTile{
		vm:      vm,
		enabled: widget.NewCheckWithData("", vm.Enabled),
		entry:   widget.NewEntryWithData(vm.Value),
		list: widget.NewSelect(options, func(selected string) {
			vm.SetValueFromList(selected)
		}),
	}

	view.list.SetSelected(vm.GetValueForList())

	if !vm.ShowEntry() {
		view.entry.Disable()
		view.entry.Hide()
	}
	if !vm.ShowList() {
		view.list.Disable()
		view.list.Hide()
	}

	if !vm.TileEnabled() {
		view.enabled.Disable()
	}

	vm.AddChangeListener(func() {
		view.Refresh()
	})

	view.ExtendBaseWidget(view)

	return view
}

func (view *LaunchArgumentTile) Refresh() {
	if view.vm.GetEnabled() {
		view.entry.Enable()
		view.list.Enable()
	} else {
		view.entry.Disable()
		view.list.Disable()
	}

	view.enabled.Refresh()
	view.entry.Refresh()
	view.list.Refresh()

	view.BaseWidget.Refresh()
}

func (view *LaunchArgumentTile) CreateRenderer() fyne.WidgetRenderer {
	return newLaunchArgumentTileRenderer(view)
}

type launchArgumentTileRenderer struct {
	view *LaunchArgumentTile

	background *canvas.Rectangle

	name              *canvas.Text
	description       *canvas.Text
	argument          *canvas.Text
	validationMessage *canvas.Text
}

func newLaunchArgumentTileRenderer(view *LaunchArgumentTile) *launchArgumentTileRenderer {
	renderer := &launchArgumentTileRenderer{
		view: view,
	}

	renderer.background = canvas.NewRectangle(theme.BackgroundColor2())
	renderer.background.CornerRadius = theme.CornerRadius()

	renderer.name = canvas.NewText(view.vm.GetName(), color.White)
	renderer.name.TextSize = theme.TextSize()
	renderer.name.TextStyle = fyne.TextStyle{Bold: true}

	renderer.description = canvas.NewText(view.vm.GetDescription(), color.RGBA{180, 180, 180, 255})
	renderer.description.TextSize = theme.TextSize() * 0.85

	renderer.argument = canvas.NewText(view.vm.GetArgument(), color.RGBA{180, 180, 180, 255})
	renderer.argument.TextSize = theme.TextSize() * 0.85

	renderer.validationMessage = canvas.NewText(view.vm.GetValidationError(), theme.StatusColor(theme.StatusError))
	renderer.validationMessage.TextSize = theme.TextSize() * 0.85

	return renderer
}

func (renderer *launchArgumentTileRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		renderer.background,
		renderer.name,
		renderer.description,
		renderer.argument,
		renderer.view.enabled,
		renderer.view.entry,
		renderer.view.list,
		renderer.validationMessage,
	}
}

func (renderer *launchArgumentTileRenderer) Layout(size fyne.Size) {
	renderer.background.Resize(size)
	renderer.background.Move(fyne.NewPos(0, 0))

	nameSize := fyne.MeasureText(renderer.name.Text, renderer.name.TextSize, renderer.name.TextStyle)
	renderer.name.Move(fyne.NewPos(theme.InnerPadding(), theme.InnerPadding()))

	argumentSize := fyne.MeasureText(renderer.argument.Text, renderer.argument.TextSize, renderer.argument.TextStyle)
	renderer.argument.Move(fyne.NewPos(renderer.name.Position().X+nameSize.Width+theme.InnerPadding(), renderer.name.Position().Y+nameSize.Height/2-argumentSize.Height/2))

	renderer.view.enabled.Resize(renderer.view.enabled.MinSize())
	renderer.view.enabled.Move(fyne.NewPos(size.Width-renderer.view.enabled.MinSize().Width/2-2*theme.InnerPadding(), 0))

	descriptionSize := fyne.MeasureText(renderer.description.Text, renderer.description.TextSize, renderer.description.TextStyle)
	renderer.description.Move(fyne.NewPos(renderer.name.Position().X, renderer.name.Position().Y+nameSize.Height))

	controlPosition := fyne.NewPos(renderer.name.Position().X, renderer.description.Position().Y+descriptionSize.Height+theme.InnerPadding())
	controlSize := fyne.NewSize(size.Width-2*theme.InnerPadding(), renderer.view.list.MinSize().Height)

	renderer.view.list.Move(controlPosition)
	renderer.view.list.Resize(controlSize)

	renderer.view.entry.Move(controlPosition)
	renderer.view.entry.Resize(controlSize)

	renderer.validationMessage.Move(fyne.NewPos(renderer.name.Position().X, controlPosition.Y+controlSize.Height))
}

func (renderer *launchArgumentTileRenderer) MinSize() fyne.Size {
	nameSize := fyne.MeasureText(renderer.name.Text, renderer.name.TextSize, renderer.name.TextStyle)
	argumentSize := fyne.MeasureText(renderer.argument.Text, renderer.argument.TextSize, renderer.argument.TextStyle)
	descriptionSize := fyne.MeasureText(renderer.description.Text, renderer.description.TextSize, renderer.description.TextStyle)

	minWidth := fyne.Max(
		renderer.name.Position().X+nameSize.Width+theme.InnerPadding()+argumentSize.Width+theme.InnerPadding(),
		renderer.description.Position().X+descriptionSize.Width+theme.InnerPadding(),
	)
	minHeight := renderer.description.Position().Y + descriptionSize.Height + theme.InnerPadding()

	if renderer.view.vm.ShowList() || renderer.view.vm.ShowEntry() {
		minHeight += renderer.view.list.MinSize().Height + theme.InnerPadding()
	}

	if renderer.view.vm.GetShowValidationError() {
		minHeight += fyne.MeasureText(renderer.view.vm.GetValidationError(), renderer.validationMessage.TextSize, renderer.validationMessage.TextStyle).Height
	}

	return fyne.Size{
		Width:  minWidth,
		Height: minHeight,
	}
}

func (renderer *launchArgumentTileRenderer) Refresh() {
	if renderer.view.vm.GetEnabled() {
		renderer.background.FillColor = theme.BackgroundColor2()
		renderer.name.Color = color.White
	} else {
		renderer.background.FillColor = color.RGBA{32, 32, 32, 255}
		renderer.name.Color = color.RGBA{128, 128, 128, 255}
	}

	renderer.validationMessage.Text = renderer.view.vm.GetValidationError()
	if renderer.view.vm.GetShowValidationError() {
		renderer.validationMessage.Show()
	} else {
		renderer.validationMessage.Hide()
	}

	renderer.background.Refresh()
	renderer.name.Refresh()
	renderer.description.Refresh()
	renderer.argument.Refresh()
	renderer.validationMessage.Refresh()
}

func (renderer *launchArgumentTileRenderer) Destroy() {}
