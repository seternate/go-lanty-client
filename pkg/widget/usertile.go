package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/seternate/go-lanty-client/pkg/theme"
	"github.com/seternate/go-lanty/pkg/user"
	"golang.design/x/clipboard"
)

type UserTile struct {
	widget.BaseWidget

	user user.User
}

func NewUserTile(user user.User) (usertile *UserTile) {
	usertile = &UserTile{
		user: user,
	}
	usertile.ExtendBaseWidget(usertile)
	return usertile
}

func (widget *UserTile) Tapped(event *fyne.PointEvent) {
	clipboard.Write(clipboard.FmtText, []byte(widget.user.IP))
}

func (widget *UserTile) TappedSecondary(event *fyne.PointEvent) {
	clipboard.Write(clipboard.FmtText, []byte(widget.user.IP))
}

func (widget *UserTile) CreateRenderer() fyne.WidgetRenderer {
	return newUsertileRenderer(widget)
}

type usertileRenderer struct {
	widget     *UserTile
	background *canvas.Rectangle
	name       *canvas.Text
	ip         *canvas.Text
}

func newUsertileRenderer(widget *UserTile) *usertileRenderer {
	renderer := &usertileRenderer{
		widget:     widget,
		background: canvas.NewRectangle(theme.BackgroundColor()),
		name:       canvas.NewText(widget.user.Name, theme.ForegroundColor()),
		ip:         canvas.NewText(widget.user.IP, theme.ForegroundColor()),
	}
	renderer.name.TextSize = theme.TextSize()
	renderer.ip.TextSize = theme.TextSize()

	return renderer
}

func (renderer *usertileRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		renderer.background,
		renderer.name,
		renderer.ip,
	}
}

func (renderer *usertileRenderer) Layout(size fyne.Size) {
	renderer.background.Resize(size)

	nametextsize := fyne.MeasureText(renderer.name.Text, renderer.name.TextSize, renderer.name.TextStyle)
	renderer.name.Move(fyne.NewPos(theme.InnerPadding(), (size.Height-nametextsize.Height)/2))
	iptextsize := fyne.MeasureText(renderer.ip.Text, renderer.ip.TextSize, renderer.ip.TextStyle)
	renderer.ip.Move(fyne.NewPos(size.Width-iptextsize.Width-theme.InnerPadding(), (size.Height-iptextsize.Height)/2))
}

func (renderer *usertileRenderer) MinSize() fyne.Size {
	nametextsize := fyne.MeasureText(renderer.name.Text, renderer.name.TextSize, renderer.name.TextStyle)
	iptextsize := fyne.MeasureText(renderer.ip.Text, renderer.ip.TextSize, renderer.ip.TextStyle)

	return fyne.NewSize(nametextsize.Width+iptextsize.Width+3*theme.InnerPadding(), iptextsize.Height+2*theme.InnerPadding())
}

func (renderer *usertileRenderer) Refresh() {
	renderer.background.Refresh()
	renderer.name.Refresh()
	renderer.ip.Refresh()
}

func (renderer *usertileRenderer) Destroy() {}
