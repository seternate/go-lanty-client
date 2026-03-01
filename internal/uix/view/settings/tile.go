package settingsview

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/seternate/go-lanty-client/internal/ui/theme"
)

type SettingTile struct {
	widget.BaseWidget

	header        string
	description   string
	entryBinding  binding.String
	window        fyne.Window
	hasFilePicker bool

	entry         *widget.Entry
	filePickerBtn *widget.Button
}

func NewSettingTile(header, description string, entryBinding binding.String, window fyne.Window, hasFilePicker bool) *SettingTile {
	tile := &SettingTile{
		header:        header,
		description:   description,
		entryBinding:  entryBinding,
		window:        window,
		hasFilePicker: hasFilePicker,
	}

	tile.entry = widget.NewEntryWithData(entryBinding)
	tile.entry.Scroll = fyne.ScrollNone
	tile.entry.Wrapping = fyne.TextWrapOff

	if hasFilePicker {
		tile.filePickerBtn = widget.NewButtonWithIcon("", fynetheme.FolderOpenIcon(), func() {
			fileDialog := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
				if err != nil {
					return
				}
				if uri != nil {
					tile.entryBinding.Set(uri.Path())
				}
			}, window)
			fileDialog.Resize(fyne.NewSize(800, 600))
			fileDialog.Show()
		})
		tile.filePickerBtn.Importance = widget.LowImportance
	}

	tile.ExtendBaseWidget(tile)
	return tile
}

func (tile *SettingTile) CreateRenderer() fyne.WidgetRenderer {
	return newSettingTileRenderer(tile)
}

type settingTileRenderer struct {
	tile       *SettingTile
	background *canvas.Rectangle
	headerText *canvas.Text
	descText   *canvas.Text
	objects    []fyne.CanvasObject
}

func newSettingTileRenderer(tile *SettingTile) *settingTileRenderer {
	renderer := &settingTileRenderer{
		tile: tile,
	}

	renderer.background = canvas.NewRectangle(theme.BackgroundColor2())
	renderer.background.CornerRadius = theme.CornerRadius()

	renderer.headerText = canvas.NewText(tile.header, color.White)
	renderer.headerText.TextSize = theme.TextSize()
	renderer.headerText.TextStyle = fyne.TextStyle{Bold: true}

	renderer.descText = canvas.NewText(tile.description, color.RGBA{180, 180, 180, 255})
	renderer.descText.TextSize = theme.TextSize() * 0.85

	renderer.objects = []fyne.CanvasObject{
		renderer.background,
		renderer.headerText,
		renderer.descText,
		tile.entry,
	}

	if tile.hasFilePicker && tile.filePickerBtn != nil {
		renderer.objects = append(renderer.objects, tile.filePickerBtn)
	}

	return renderer
}

func (r *settingTileRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *settingTileRenderer) Layout(size fyne.Size) {
	padding := theme.InnerPadding() * 1.5
	spacing := theme.InnerPadding() * 0.5

	r.background.Resize(size)
	r.background.Move(fyne.NewPos(0, 0))

	headerY := padding
	r.headerText.Move(fyne.NewPos(padding, headerY))
	headerHeight := fyne.MeasureText(r.headerText.Text, r.headerText.TextSize, r.headerText.TextStyle).Height

	descY := headerY + headerHeight + spacing
	r.descText.Move(fyne.NewPos(padding, descY))
	descHeight := fyne.MeasureText(r.descText.Text, r.descText.TextSize, r.descText.TextStyle).Height

	entryY := descY + descHeight + spacing
	entryHeight := float32(40)

	entryWidth := size.Width - 2*padding
	buttonWidth := float32(0)
	if r.tile.hasFilePicker && r.tile.filePickerBtn != nil {
		buttonWidth = entryHeight
		entryWidth -= buttonWidth + spacing
	}

	r.tile.entry.Resize(fyne.NewSize(entryWidth, entryHeight))
	r.tile.entry.Move(fyne.NewPos(padding, entryY))

	if r.tile.hasFilePicker && r.tile.filePickerBtn != nil {
		r.tile.filePickerBtn.Resize(fyne.NewSize(buttonWidth, entryHeight))
		r.tile.filePickerBtn.Move(fyne.NewPos(padding+entryWidth+spacing, entryY))
	}
}

func (r *settingTileRenderer) MinSize() fyne.Size {
	padding := theme.InnerPadding() * 1.5
	spacing := theme.InnerPadding() * 0.5

	headerHeight := fyne.MeasureText(r.headerText.Text, r.headerText.TextSize, r.headerText.TextStyle).Height
	descHeight := fyne.MeasureText(r.descText.Text, r.descText.TextSize, r.descText.TextStyle).Height
	entryHeight := float32(40)

	minWidth := float32(500)
	buttonWidth := float32(0)
	if r.tile.hasFilePicker && r.tile.filePickerBtn != nil {
		buttonWidth = entryHeight + spacing
	}
	minWidth = fyne.Max(minWidth, entryHeight*8+buttonWidth+2*padding)

	minHeight := padding*2 + headerHeight + spacing + descHeight + spacing + entryHeight

	return fyne.NewSize(minWidth, minHeight)
}

func (r *settingTileRenderer) Refresh() {
	r.background.Refresh()
	r.headerText.Refresh()
	r.descText.Refresh()
	r.tile.entry.Refresh()
	if r.tile.hasFilePicker && r.tile.filePickerBtn != nil {
		r.tile.filePickerBtn.Refresh()
	}
}

func (r *settingTileRenderer) Destroy() {}
