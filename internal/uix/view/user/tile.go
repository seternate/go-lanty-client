package userview

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	model "github.com/seternate/go-lanty-client/internal/ui/model/user"
	"github.com/seternate/go-lanty-client/internal/ui/theme"
)

type UserTile struct {
	widget.BaseWidget

	user       *model.UserTileModel
	copyButton *widget.Button
}

func NewUserTile(user *model.UserTileModel) *UserTile {
	tile := &UserTile{
		user: user,
	}

	tile.copyButton = widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		if tile.user.IPAddress != "" {
			fyne.CurrentApp().Clipboard().SetContent(tile.user.IPAddress)
		}
	})

	tile.copyButton.Importance = widget.LowImportance

	user.AddChangeListener(func() {
		tile.Refresh()
	})

	tile.ExtendBaseWidget(tile)

	return tile
}

func NewUserTileWithoutModel(ipAddress string, name string) *UserTile {
	tile := &UserTile{
		user: model.NewUserTileModel(ipAddress),
	}

	tile.user.Name.Set(name)

	tile.copyButton = widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		if tile.user.IPAddress != "" {
			fyne.CurrentApp().Clipboard().SetContent(tile.user.IPAddress)
		}
	})

	tile.user.AddChangeListener(func() {
		tile.Refresh()
	})

	tile.ExtendBaseWidget(tile)
	return tile
}

func (tile *UserTile) SetIPAddress(ipAddress string) {
	tile.user.IPAddress = ipAddress
	tile.Refresh()
}

func (tile *UserTile) SetName(name string) {
	tile.user.SetName(name)
	tile.Refresh()
}

func (tile *UserTile) HideCopyButton() {
	tile.copyButton.Hide()
	tile.Refresh()
}

func (tile *UserTile) CreateRenderer() fyne.WidgetRenderer {
	return newTileRenderer(tile)
}

type tileRenderer struct {
	tile       *UserTile
	background *canvas.Rectangle
	name       *canvas.Text
	ipAddress  *canvas.Text
}

func newTileRenderer(tile *UserTile) *tileRenderer {
	renderer := &tileRenderer{
		tile: tile,
	}

	renderer.background = canvas.NewRectangle(theme.BackgroundColor2())
	renderer.background.CornerRadius = theme.CornerRadius()

	renderer.name = canvas.NewText("", color.White)
	renderer.name.TextSize = theme.TextSize()

	renderer.ipAddress = canvas.NewText("", color.RGBA{180, 180, 180, 255})
	renderer.ipAddress.TextSize = theme.TextSize() * 0.85

	if name, err := tile.user.Name.Get(); err == nil {
		renderer.name.Text = name
	}
	renderer.ipAddress.Text = tile.user.IPAddress

	return renderer
}

func (renderer *tileRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{
		renderer.background,
		renderer.name,
		renderer.ipAddress,
	}

	if renderer.tile.copyButton.Visible() {
		objects = append(objects, renderer.tile.copyButton)
	}

	return objects
}

func (renderer *tileRenderer) Layout(size fyne.Size) {
	renderer.background.Resize(size)

	padding := theme.InnerPadding()

	nameHeight := fyne.MeasureText(renderer.name.Text, renderer.name.TextSize, renderer.name.TextStyle).Height
	ipHeight := fyne.MeasureText(renderer.ipAddress.Text, renderer.ipAddress.TextSize, renderer.ipAddress.TextStyle).Height

	textSpacing := padding * 0.5
	totalTextHeight := nameHeight + textSpacing + ipHeight

	centerY := size.Height / 2
	nameY := centerY - totalTextHeight/2
	ipY := nameY + nameHeight + textSpacing

	renderer.name.Move(fyne.NewPos(padding, nameY))

	ipX := padding
	renderer.ipAddress.Move(fyne.NewPos(ipX, ipY))

	buttonSize := float32(32) // Small button size
	ipWidth := fyne.MeasureText(renderer.ipAddress.Text, renderer.ipAddress.TextSize, renderer.ipAddress.TextStyle).Width
	buttonX := ipX + ipWidth + padding*0.5
	buttonY := ipY - (buttonSize-ipHeight)/2 // Center button vertically with IP text
	renderer.tile.copyButton.Resize(fyne.NewSize(buttonSize, buttonSize))
	renderer.tile.copyButton.Move(fyne.NewPos(buttonX, buttonY))
}

func (renderer *tileRenderer) MinSize() fyne.Size {
	padding := theme.InnerPadding()

	nameSize := fyne.MeasureText(renderer.name.Text, renderer.name.TextSize, renderer.name.TextStyle)
	ipSize := fyne.MeasureText(renderer.ipAddress.Text, renderer.ipAddress.TextSize, renderer.ipAddress.TextStyle)

	textSpacing := padding * 0.5
	totalTextHeight := nameSize.Height + textSpacing + ipSize.Height

	buttonSize := float32(32)
	ipRowWidth := ipSize.Width + padding*0.5 + buttonSize

	minWidth := nameSize.Width
	if ipRowWidth > minWidth {
		minWidth = ipRowWidth
	}
	minWidth += 2 * padding

	minHeight := totalTextHeight + 2*padding

	return fyne.NewSize(minWidth, minHeight)
}

func (renderer *tileRenderer) Refresh() {
	renderer.ipAddress.Text = renderer.tile.user.IPAddress
	name, _ := renderer.tile.user.Name.Get()
	renderer.name.Text = name

	renderer.background.Refresh()
	renderer.name.Refresh()
	renderer.ipAddress.Refresh()
	renderer.tile.copyButton.Refresh()
}

func (renderer *tileRenderer) Destroy() {}
