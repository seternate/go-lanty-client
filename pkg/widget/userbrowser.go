package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/seternate/go-lanty-client/pkg/controller"
	"github.com/seternate/go-lanty-client/pkg/theme"
)

type UserBrowser struct {
	widget.BaseWidget

	controller *controller.Controller
	usertiles  []*UserTile

	usersupdated chan struct{}
}

func NewUserBrowser(controller *controller.Controller) (userbrowser *UserBrowser) {
	userbrowser = &UserBrowser{
		controller:   controller,
		usertiles:    make([]*UserTile, 0, 50),
		usersupdated: make(chan struct{}),
	}
	userbrowser.ExtendBaseWidget(userbrowser)
	userbrowser.updateUsertiles()
	controller.User.Subscribe(userbrowser.usersupdated)
	go userbrowser.usersUpdater()

	return
}

func (widget *UserBrowser) setUsertiles(usertiles ...*UserTile) {
	usertilesOld := widget.usertiles
	widget.usertiles = usertiles

	for index := range usertilesOld {
		usertilesOld[index] = nil
	}

	widget.Refresh()
}

func (widget *UserBrowser) CreateRenderer() fyne.WidgetRenderer {
	return newUserBrowserRenderer(widget)
}

func (widget *UserBrowser) usersUpdater() {
	for {
		<-widget.usersupdated
		widget.updateUsertiles()
	}
}

func (widget *UserBrowser) updateUsertiles() {
	usertiles := make([]*UserTile, 0, len(widget.controller.User.GetUsers()))
	for _, user := range widget.controller.User.GetUsers() {
		usertile := NewUserTile(user)
		usertiles = append(usertiles, usertile)
	}
	widget.setUsertiles(usertiles...)
}

type userBrowserRenderer struct {
	widget    *UserBrowser
	usertiles *fyne.Container
}

func newUserBrowserRenderer(widget *UserBrowser) *userBrowserRenderer {
	usertiles := container.NewVBox()
	renderer := &userBrowserRenderer{
		widget:    widget,
		usertiles: usertiles,
	}

	return renderer
}

func (renderer *userBrowserRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		renderer.usertiles,
	}
}

func (renderer *userBrowserRenderer) Layout(size fyne.Size) {
	renderer.usertiles.Resize(size.SubtractWidthHeight(2*theme.InnerPadding(), 0))
	renderer.usertiles.Move(fyne.NewPos(theme.InnerPadding(), 0))
}

func (renderer *userBrowserRenderer) MinSize() fyne.Size {
	return renderer.usertiles.MinSize().AddWidthHeight(2*theme.InnerPadding(), 0)
}

func (renderer *userBrowserRenderer) Refresh() {
	renderer.usertiles.RemoveAll()
	for _, usertile := range renderer.widget.usertiles {
		renderer.usertiles.Add(usertile)
		usertile.Refresh()
	}
	renderer.usertiles.Refresh()
}

func (renderer *userBrowserRenderer) Destroy() {}
