package widget

import (
	"fyne.io/fyne/v2"
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/seternate/go-lanty-client/pkg/controller"
	"github.com/seternate/go-lanty-client/pkg/theme"
	"github.com/seternate/go-lanty/pkg/game"
	"github.com/seternate/go-lanty/pkg/user"
)

type JoinServer struct {
	widget.BaseWidget

	userbrowser *UserBrowser
	cancel      *widget.Button
	game        game.Game

	OnUserSelected func(game game.Game, user user.User)
	OnCancelTapped func()
}

func NewJoinServer(controller *controller.Controller) *JoinServer {
	joinserver := &JoinServer{
		userbrowser: NewUserBrowser(controller),
		cancel:      widget.NewButtonWithIcon("Cancel", fynetheme.CancelIcon(), nil),
	}
	joinserver.ExtendBaseWidget(joinserver)

	joinserver.userbrowser.SetOnUserDoubleTapped(func(user user.User) {
		if joinserver.OnUserSelected != nil {
			joinserver.OnUserSelected(joinserver.game, user)
		}
	})
	joinserver.cancel.OnTapped = func() {
		if joinserver.OnCancelTapped != nil {
			joinserver.OnCancelTapped()
		}
	}

	return joinserver
}

func (widget *JoinServer) SetGame(game game.Game) {
	widget.game = game
}

func (widget *JoinServer) CreateRenderer() fyne.WidgetRenderer {
	return newJoinServerRenderer(widget)
}

type joinServerRenderer struct {
	widget *JoinServer
}

func newJoinServerRenderer(widget *JoinServer) *joinServerRenderer {
	return &joinServerRenderer{
		widget: widget,
	}
}

func (renderer *joinServerRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		renderer.widget.userbrowser,
		renderer.widget.cancel,
	}
}

func (renderer *joinServerRenderer) Layout(size fyne.Size) {
	renderer.widget.userbrowser.Resize(fyne.NewSize(size.Width, size.Height-renderer.widget.cancel.Size().Height-2*theme.InnerPadding()))
	renderer.widget.userbrowser.Move(fyne.NewPos(0, 0))
	renderer.widget.cancel.Resize(fyne.NewSize(renderer.widget.cancel.MinSize().Width, renderer.widget.cancel.MinSize().Height))
	renderer.widget.cancel.Move(fyne.NewPos(renderer.widget.userbrowser.Size().Width-renderer.widget.cancel.Size().Width-theme.InnerPadding(), size.Height-renderer.widget.cancel.Size().Height-theme.InnerPadding()))

}

func (renderer *joinServerRenderer) MinSize() fyne.Size {
	return renderer.widget.userbrowser.MinSize().AddWidthHeight(0, renderer.widget.cancel.MinSize().Height)
}

func (renderer *joinServerRenderer) Refresh() {
	renderer.widget.userbrowser.Refresh()
	renderer.widget.cancel.Refresh()
}

func (renderer *joinServerRenderer) Destroy() {}
