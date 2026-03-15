package settingsview

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/seternate/go-lanty-client/internal/ui/theme"
	settingsviewmodel "github.com/seternate/go-lanty-client/internal/ui/viewmodel/settings"
)

type SettingTile struct {
	widget.BaseWidget

	vm *settingsviewmodel.SettingTile

	entry              *widget.Entry
	folderPickerButton *widget.Button
}

func NewSettingTile(vm *settingsviewmodel.SettingTile) *SettingTile {
	view := &SettingTile{
		vm:                 vm,
		entry:              widget.NewEntryWithData(vm.Value),
		folderPickerButton: widget.NewButtonWithIcon("", fynetheme.FolderOpenIcon(), vm.OpenFolderPicker),
	}

	view.entry.Scroll = fyne.ScrollNone
	view.entry.Wrapping = fyne.TextWrapOff
	view.entry.OnSubmitted = func(s string) {
		vm.Save()
	}

	view.folderPickerButton.Importance = widget.LowImportance

	view.ExtendBaseWidget(view)

	return view
}

func (view *SettingTile) Refresh() {
	view.entry.Refresh()
	view.folderPickerButton.Refresh()

	view.BaseWidget.Refresh()
}

func (view *SettingTile) CreateRenderer() fyne.WidgetRenderer {
	return newSettingTileRenderer(view)
}

type settingTileRenderer struct {
	view *SettingTile

	background *canvas.Rectangle

	label       *canvas.Text
	description *canvas.Text
}

func newSettingTileRenderer(view *SettingTile) *settingTileRenderer {
	renderer := &settingTileRenderer{
		view: view,
	}

	renderer.background = canvas.NewRectangle(theme.BackgroundColor2())
	renderer.background.CornerRadius = theme.CornerRadius()

	renderer.label = canvas.NewText(view.vm.GetLabel(), color.White)
	renderer.label.TextSize = theme.TextSize()
	renderer.label.TextStyle = fyne.TextStyle{Bold: true}

	renderer.description = canvas.NewText(view.vm.GetDescription(), color.RGBA{180, 180, 180, 255})
	renderer.description.TextSize = theme.TextSize() * 0.85

	return renderer
}

func (renderer *settingTileRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		renderer.background,
		renderer.label,
		renderer.description,
		renderer.view.entry,
		renderer.view.folderPickerButton,
	}
}

func (renderer *settingTileRenderer) Layout(size fyne.Size) {
	renderer.background.Resize(size)
	renderer.background.Move(fyne.NewPos(0, 0))

	labelSize := fyne.MeasureText(renderer.label.Text, renderer.label.TextSize, renderer.label.TextStyle)
	renderer.label.Move(fyne.NewPos(theme.InnerPadding(), theme.InnerPadding()))

	descriptionSize := fyne.MeasureText(renderer.description.Text, renderer.description.TextSize, renderer.description.TextStyle)
	renderer.description.Move(fyne.NewPos(renderer.label.Position().X, renderer.label.Position().Y+labelSize.Height+theme.InnerPadding()*0.5))

	entryHeight := float32(40)
	entryWidth := size.Width - 2*theme.InnerPadding()
	if renderer.view.vm.GetShowFolderPicker() {
		entryWidth -= entryHeight + theme.InnerPadding()*0.5
	}
	renderer.view.entry.Resize(fyne.NewSize(entryWidth, entryHeight))
	renderer.view.entry.Move(fyne.NewPos(renderer.label.Position().X, renderer.description.Position().Y+descriptionSize.Height+theme.InnerPadding()))

	renderer.view.folderPickerButton.Resize(fyne.NewSize(entryHeight, entryHeight))
	renderer.view.folderPickerButton.Move(fyne.NewPos(renderer.view.entry.Position().X+entryWidth+theme.InnerPadding()*0.5, renderer.view.entry.Position().Y))

}

func (renderer *settingTileRenderer) MinSize() fyne.Size {
	labelSize := fyne.MeasureText(renderer.label.Text, renderer.label.TextSize, renderer.label.TextStyle)
	descriptionSize := fyne.MeasureText(renderer.description.Text, renderer.description.TextSize, renderer.description.TextStyle)

	minWidth := fyne.Max(labelSize.Width, descriptionSize.Width) + 2*theme.InnerPadding()
	minHeight := theme.InnerPadding() + labelSize.Height + theme.InnerPadding()*0.5 + descriptionSize.Height + theme.InnerPadding() + renderer.view.entry.Size().Height + theme.InnerPadding()

	return fyne.NewSize(minWidth, minHeight)
}

func (renderer *settingTileRenderer) Refresh() {
	renderer.background.Refresh()
	renderer.label.Refresh()
	renderer.description.Refresh()
}

func (renderer *settingTileRenderer) Destroy() {}
