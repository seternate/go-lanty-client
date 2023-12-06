package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty-client/pkg/controller"
	"github.com/seternate/go-lanty-client/pkg/theme"
	"github.com/seternate/go-lanty/pkg/user"
)

type UserBrowser struct {
	widget.BaseWidget

	controller *controller.Controller
	usertiles  []*UserTile

	usersupdated chan struct{}

	onUserTapped       func(user.User)
	onUserDoubleTapped func(user.User)
}

func NewUserBrowser(controller *controller.Controller) (userbrowser *UserBrowser) {
	userbrowser = &UserBrowser{
		controller:   controller,
		usertiles:    make([]*UserTile, 0, 50),
		usersupdated: make(chan struct{}, 50),
	}
	userbrowser.ExtendBaseWidget(userbrowser)
	userbrowser.updateUsertiles()
	controller.User.Subscribe(userbrowser.usersupdated)
	controller.WaitGroup().Add(1)
	go userbrowser.usersUpdater()

	return
}

func (widget *UserBrowser) SetOnUserTapped(cb func(user.User)) {
	for _, usertile := range widget.usertiles {
		usertile.OnTapped = cb
	}
	widget.onUserTapped = cb
}

func (widget *UserBrowser) SetOnUserDoubleTapped(cb func(user.User)) {
	for _, usertile := range widget.usertiles {
		usertile.OnDoubleTapped = cb
	}
	widget.onUserDoubleTapped = cb
}

func (widget *UserBrowser) CreateRenderer() fyne.WidgetRenderer {
	return newUserBrowserRenderer(widget)
}

func (widget *UserBrowser) usersUpdater() {
	defer widget.controller.WaitGroup().Done()
	for {
		select {
		case <-widget.controller.Context().Done():
			log.Trace().Msg("exiting userbrowser usersUpdater()")
			return
		case <-widget.usersupdated:
			widget.updateUsertiles()
		}
	}
}

func (widget *UserBrowser) updateUsertiles() {
	usertiles := make([]*UserTile, 0, len(widget.controller.User.GetUsers()))
	for _, user := range widget.controller.User.GetUsers() {
		usertile := NewUserTile(widget, user)
		usertile.OnTapped = widget.onUserTapped
		usertile.OnDoubleTapped = widget.onUserDoubleTapped
		usertiles = append(usertiles, usertile)
	}
	widget.setUsertiles(usertiles...)
}

func (widget *UserBrowser) setUsertiles(usertiles ...*UserTile) {
	usertilesOld := widget.usertiles
	widget.usertiles = usertiles

	for index := range usertilesOld {
		usertilesOld[index] = nil
	}

	widget.Refresh()
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
	renderer.usertiles.Resize(size.SubtractWidthHeight(2*theme.InnerPadding(), 2*theme.InnerPadding()))
	renderer.usertiles.Move(fyne.NewPos(theme.InnerPadding(), theme.InnerPadding()))
}

func (renderer *userBrowserRenderer) MinSize() fyne.Size {
	return renderer.usertiles.MinSize().AddWidthHeight(2*theme.InnerPadding(), 2*theme.InnerPadding())
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
