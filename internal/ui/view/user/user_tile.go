package userview

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/seternate/go-lanty-client/internal/ui/theme"
	userviewmodel "github.com/seternate/go-lanty-client/internal/ui/viewmodel/user"
)

type UserTile struct {
	widget.BaseWidget

	vm *userviewmodel.UserTile

	copyButton *widget.Button
}

func NewUserTile(vm *userviewmodel.UserTile) *UserTile {
	view := &UserTile{
		vm:         vm,
		copyButton: widget.NewButtonWithIcon("", theme.ContentCopyIcon(), vm.CopyIPAddressToClipboard),
	}

	view.copyButton.Importance = widget.LowImportance

	vm.AddChangeListener(func() {
		view.Refresh()
	})

	view.ExtendBaseWidget(view)

	return view
}

func (view *UserTile) Refresh() {
	view.copyButton.Refresh()
	view.BaseWidget.Refresh()
}

func (view *UserTile) CreateRenderer() fyne.WidgetRenderer {
	return newUserTileRenderer(view)
}

type userTileRenderer struct {
	view *UserTile

	background *canvas.Rectangle

	name      *canvas.Text
	ipAddress *canvas.Text
}

func newUserTileRenderer(view *UserTile) *userTileRenderer {
	renderer := &userTileRenderer{
		view: view,
	}

	renderer.background = canvas.NewRectangle(theme.BackgroundColor2())
	renderer.background.CornerRadius = theme.CornerRadius()

	renderer.name = canvas.NewText(view.vm.GetName(), color.White)
	renderer.name.TextSize = theme.TextSize()

	renderer.ipAddress = canvas.NewText(view.vm.GetIPAddress(), color.RGBA{180, 180, 180, 255})
	renderer.ipAddress.TextSize = theme.TextSize() * 0.85

	return renderer
}

func (renderer *userTileRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		renderer.background,
		renderer.name,
		renderer.ipAddress,
		renderer.view.copyButton,
	}
}

func (renderer *userTileRenderer) Layout(size fyne.Size) {
	renderer.background.Resize(size)

	nameHeight := fyne.MeasureText(renderer.name.Text, renderer.name.TextSize, renderer.name.TextStyle).Height
	ipTextSize := fyne.MeasureText(renderer.ipAddress.Text, renderer.ipAddress.TextSize, renderer.ipAddress.TextStyle)

	textSpacing := theme.InnerPadding() * 0.5
	totalTextHeight := nameHeight + textSpacing + ipTextSize.Height
	centerY := size.Height / 2

	nameY := centerY - totalTextHeight/2
	renderer.name.Move(fyne.NewPos(theme.InnerPadding(), nameY))

	ipY := nameY + nameHeight + textSpacing
	renderer.ipAddress.Move(fyne.NewPos(theme.InnerPadding(), ipY))

	buttonSize := fyne.NewSize(32, 32)
	buttonX := renderer.ipAddress.Position().X + ipTextSize.Width + theme.InnerPadding()*0.5
	buttonY := renderer.ipAddress.Position().Y + ipTextSize.Height/2 - buttonSize.Height/2
	renderer.view.copyButton.Resize(buttonSize)
	renderer.view.copyButton.Move(fyne.NewPos(buttonX, buttonY))
}

func (renderer *userTileRenderer) MinSize() fyne.Size {
	nameSize := fyne.MeasureText(renderer.name.Text, renderer.name.TextSize, renderer.name.TextStyle)
	ipSize := fyne.MeasureText(renderer.ipAddress.Text, renderer.ipAddress.TextSize, renderer.ipAddress.TextStyle)

	textSpacing := theme.InnerPadding() * 0.5
	totalTextHeight := nameSize.Height + textSpacing + ipSize.Height

	buttonSize := float32(32)
	ipRowWidth := ipSize.Width + theme.InnerPadding()*0.5 + buttonSize

	minWidth := nameSize.Width
	if ipRowWidth > minWidth {
		minWidth = ipRowWidth
	}

	minWidth += 2 * theme.InnerPadding()
	minHeight := totalTextHeight + 2*theme.InnerPadding()

	return fyne.NewSize(minWidth, minHeight)
}

func (renderer *userTileRenderer) Refresh() {
	renderer.name.Text = renderer.view.vm.GetName()
	renderer.ipAddress.Text = renderer.view.vm.GetIPAddress()

	renderer.background.Refresh()
	renderer.name.Refresh()
	renderer.ipAddress.Refresh()
}

func (renderer *userTileRenderer) Destroy() {}
