package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/seternate/go-lanty-client/pkg/theme"
	"github.com/seternate/go-lanty/pkg/user"
)

type UserTile struct {
	widget.BaseWidget

	parent  *UserBrowser
	user    user.User
	hovered bool

	OnTapped       func(user.User)
	OnDoubleTapped func(user.User)
}

func NewUserTile(parent *UserBrowser, user user.User) (usertile *UserTile) {
	usertile = &UserTile{
		parent:  parent,
		user:    user,
		hovered: false,
	}
	usertile.ExtendBaseWidget(usertile)
	return usertile
}

func (widget *UserTile) Tapped(event *fyne.PointEvent) {
	if widget.OnTapped != nil {
		widget.OnTapped(widget.user)
	}
}

func (widget *UserTile) DoubleTapped(event *fyne.PointEvent) {
	if widget.OnDoubleTapped != nil {
		widget.OnDoubleTapped(widget.user)
	}
}

func (widget *UserTile) MouseIn(event *desktop.MouseEvent) {
	widget.hovered = true
	widget.Refresh()
}

func (widget *UserTile) MouseMoved(event *desktop.MouseEvent) {

}

func (widget *UserTile) MouseOut() {
	widget.hovered = false
	widget.Refresh()
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
		background: canvas.NewRectangle(fynetheme.InputBorderColor()),
		name:       canvas.NewText(widget.user.Name, theme.ForegroundColor()),
		ip:         canvas.NewText(widget.user.IP, theme.ForegroundColor()),
	}
	renderer.name.TextSize = theme.TextSize()
	renderer.ip.TextSize = theme.TextSize()
	renderer.background.CornerRadius = fynetheme.SelectionRadiusSize()

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

	return fyne.NewSize(nametextsize.Width+iptextsize.Width+3*theme.InnerPadding(), 1.5*iptextsize.Height+2*theme.InnerPadding())
}

func (renderer *usertileRenderer) Refresh() {
	if renderer.widget.hovered {
		renderer.background.FillColor = fynetheme.InputBorderColor()
	} else {
		renderer.background.FillColor = fynetheme.InputBackgroundColor()
	}
	renderer.background.Refresh()
	renderer.name.Refresh()
	renderer.ip.Refresh()
}

func (renderer *usertileRenderer) Destroy() {}
