package widget

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty-client/pkg/controller"
	"github.com/seternate/go-lanty-client/pkg/setting"
	"github.com/seternate/go-lanty-client/pkg/theme"
	"github.com/seternate/go-lanty/pkg/game"
	"github.com/seternate/go-lanty/pkg/user"
	"golang.design/x/clipboard"
)

type Lanty struct {
	widget.BaseWidget

	controller             *controller.Controller
	sidebar                *Sidebar
	gamebrowser            *ScrollWithState
	startserver            *ScrollWithState
	joinserver             *JoinServer
	downloadbrowser        *ScrollWithState
	userbrowser            *ScrollWithState
	settingsbrowser        *ScrollWithState
	statusbar              *StatusBar
	chatbrowser            *fyne.Container
	defaultusernamebrowser *ScrollWithState

	resetSettingsBrowser func()

	statusupdate chan struct{}
}

func NewLanty(controller *controller.Controller, window fyne.Window) *Lanty {
	gamebrowser := NewGameBrowser(controller)
	startserver := NewStartServer(controller)
	downloadbrowser := NewDownloadBrowser(controller)
	userbrowser := NewUserBrowser(controller)
	settingsbrowser := NewSettingsBrowser(controller, window)
	chatbrowser := NewChatBrowser(controller, window)
	defaultusernamebrowser := NewDefaultUsernameBrowser(controller, window)

	lanty := &Lanty{
		controller:             controller,
		sidebar:                NewSidebar(setting.APPLICATION_NAME),
		gamebrowser:            NewVScrollWithState(gamebrowser),
		startserver:            NewVScrollWithState(startserver),
		joinserver:             NewJoinServer(controller),
		downloadbrowser:        NewVScrollWithState(downloadbrowser),
		userbrowser:            NewVScrollWithState(userbrowser),
		settingsbrowser:        NewVScrollWithState(settingsbrowser),
		statusbar:              NewStatusBar(controller),
		chatbrowser:            container.NewStack(chatbrowser),
		defaultusernamebrowser: NewVScrollWithState(defaultusernamebrowser),
		resetSettingsBrowser: func() {
			settingsbrowser.ResetData()
		},
		statusupdate: make(chan struct{}, 50),
	}
	lanty.ExtendBaseWidget(lanty)

	lanty.sidebar.OnGamesTapped = func() {
		lanty.showGameBrowser()
	}
	lanty.sidebar.OnDownloadsTapped = func() {
		lanty.showDownloadBrowser()
	}
	lanty.sidebar.OnChatsTapped = func() {
		lanty.showChatBrowser()
	}
	lanty.sidebar.OnUsersTapped = func() {
		lanty.showUserBrowser()
	}
	lanty.sidebar.OnSettingsTapped = func() {
		lanty.showSettingsBrowser()
	}

	gamebrowser.OnJoinServerTapped = func(game game.Game) {
		lanty.joinserver.SetGame(game)
		lanty.showJoinServer()
	}
	gamebrowser.OnStartServerTapped = func(game game.Game) {
		startserver.SetGame(game)
		lanty.showStartServer()
	}

	startserver.OnSubmit = func(game game.Game) {
		lanty.showGameBrowser()
		controller.Game.StartServer(game)
	}
	startserver.OnCancel = func() {
		lanty.showGameBrowser()
	}

	lanty.joinserver.OnUserSelected = func(game game.Game, user user.User) {
		lanty.showGameBrowser()
		controller.Game.JoinServer(game, user)
	}
	lanty.joinserver.OnCancelTapped = func() {
		lanty.showGameBrowser()
	}

	userbrowser.SetOnUserTapped(func(user user.User) {
		clipboard.Write(clipboard.FmtText, []byte(user.IP))
		controller.Status.Info("IP copied", 3*time.Second)
	})
	userbrowser.SetOnUserDoubleTapped(func(user user.User) {
		clipboard.Write(clipboard.FmtText, []byte(user.IP))
		controller.Status.Info("IP copied", 3*time.Second)
	})

	defaultusernamebrowser.SetOnSubmit(func() {
		lanty.defaultusernamebrowser.Hide()
		lanty.showGameBrowser()
	})
	if controller.Settings.Settings().Username == setting.DEFAULT_USERNAME {
		lanty.defaultusernamebrowser.Show()
		lanty.hideAll()
	} else {
		lanty.defaultusernamebrowser.Hide()
		lanty.showGameBrowser()
	}

	controller.Status.Subscribe(lanty.statusupdate)
	controller.WaitGroup().Add(1)
	go lanty.statusbarUpdater()

	return lanty
}

func (widget *Lanty) hideAll() {
	widget.sidebar.Hide()
	widget.gamebrowser.Hide()
	widget.downloadbrowser.Hide()
	widget.chatbrowser.Hide()
	widget.userbrowser.Hide()
	widget.settingsbrowser.Hide()
	widget.startserver.Hide()
	widget.joinserver.Hide()
	widget.Refresh()
}

func (widget *Lanty) showGameBrowser() {
	widget.sidebar.Show()
	widget.gamebrowser.Show()
	widget.downloadbrowser.Hide()
	widget.chatbrowser.Hide()
	widget.userbrowser.Hide()
	widget.settingsbrowser.Hide()
	widget.startserver.Hide()
	widget.joinserver.Hide()
	widget.Refresh()
}

func (widget *Lanty) showDownloadBrowser() {
	widget.sidebar.Show()
	widget.gamebrowser.Hide()
	widget.downloadbrowser.Show()
	widget.chatbrowser.Hide()
	widget.userbrowser.Hide()
	widget.settingsbrowser.Hide()
	widget.startserver.Hide()
	widget.joinserver.Hide()
	widget.Refresh()
}

func (widget *Lanty) showChatBrowser() {
	widget.sidebar.Show()
	widget.gamebrowser.Hide()
	widget.downloadbrowser.Hide()
	widget.chatbrowser.Show()
	widget.userbrowser.Hide()
	widget.settingsbrowser.Hide()
	widget.startserver.Hide()
	widget.joinserver.Hide()
	widget.Refresh()
}

func (widget *Lanty) showUserBrowser() {
	widget.sidebar.Show()
	widget.gamebrowser.Hide()
	widget.downloadbrowser.Hide()
	widget.chatbrowser.Hide()
	widget.userbrowser.Show()
	widget.settingsbrowser.Hide()
	widget.startserver.Hide()
	widget.joinserver.Hide()
	widget.Refresh()
}

func (widget *Lanty) showSettingsBrowser() {
	widget.sidebar.Show()
	widget.gamebrowser.Hide()
	widget.downloadbrowser.Hide()
	widget.chatbrowser.Hide()
	widget.userbrowser.Hide()
	widget.settingsbrowser.Show()
	widget.resetSettingsBrowser()
	widget.startserver.Hide()
	widget.joinserver.Hide()
	widget.Refresh()
}

func (widget *Lanty) showStartServer() {
	widget.sidebar.Show()
	widget.gamebrowser.Hide()
	widget.downloadbrowser.Hide()
	widget.chatbrowser.Hide()
	widget.userbrowser.Hide()
	widget.settingsbrowser.Hide()
	widget.startserver.Show()
	widget.joinserver.Hide()
	widget.Refresh()
}

func (widget *Lanty) showJoinServer() {
	widget.sidebar.Show()
	widget.gamebrowser.Hide()
	widget.downloadbrowser.Hide()
	widget.chatbrowser.Hide()
	widget.userbrowser.Hide()
	widget.settingsbrowser.Hide()
	widget.startserver.Hide()
	widget.joinserver.Show()
	widget.Refresh()
}

func (widget *Lanty) statusbarUpdater() {
	defer widget.controller.WaitGroup().Done()
	for {
		select {
		case <-widget.controller.Context().Done():
			log.Trace().Msg("exiting lanty statusbarUpdater()")
			return
		case <-widget.statusupdate:
			status := widget.controller.Status.Next()
			widget.statusbar.ShowStatus(status)
			widget.Refresh()
		}
	}
}

func (widget *Lanty) CreateRenderer() fyne.WidgetRenderer {
	return newLantyRenderer(widget)
}

type lantyRenderer struct {
	widget *Lanty
	main   *fyne.Container
}

func newLantyRenderer(widget *Lanty) *lantyRenderer {
	return &lantyRenderer{
		widget: widget,
		main: container.NewStack(
			widget.gamebrowser,
			widget.downloadbrowser,
			widget.chatbrowser,
			widget.userbrowser,
			widget.settingsbrowser,
			widget.startserver,
			widget.joinserver,
		),
	}
}

func (renderer *lantyRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		renderer.widget.defaultusernamebrowser,
		renderer.widget.sidebar,
		renderer.main,
		renderer.widget.statusbar,
	}
}

func (renderer *lantyRenderer) Layout(size fyne.Size) {
	renderer.widget.defaultusernamebrowser.Move(fyne.NewPos(0, 0))
	renderer.widget.defaultusernamebrowser.Resize(size)
	renderer.widget.sidebar.Move(fyne.NewPos(0, 0))
	renderer.widget.sidebar.Resize(fyne.NewSize(fyne.Max(size.Width/7, renderer.widget.sidebar.MinSize().Width), size.Height))
	renderer.main.Resize(fyne.NewSize(size.Width-fyne.Max(size.Width/7, renderer.widget.sidebar.MinSize().Width), size.Height))
	renderer.main.Move(fyne.NewPos(fyne.Max(size.Width/7, renderer.widget.sidebar.MinSize().Width), 0))
	renderer.widget.statusbar.Resize(fyne.NewSize(renderer.main.Size().Width/2, 40))
	renderer.widget.statusbar.Move(fyne.NewPos(size.Width-theme.InnerPadding()-renderer.widget.statusbar.Size().Width, size.Height-theme.InnerPadding()-renderer.widget.statusbar.Size().Height))
}

func (renderer *lantyRenderer) MinSize() fyne.Size {
	size := fyne.NewSize(0, 0)
	size.Height = fyne.Max(renderer.widget.sidebar.MinSize().Height, renderer.main.MinSize().Height)
	size.Width = renderer.widget.gamebrowser.MinSize().Width + fyne.Max(renderer.widget.gamebrowser.MinSize().Width/6, renderer.widget.sidebar.MinSize().Width)
	return size
}

func (renderer *lantyRenderer) Refresh() {
	//In order to Hide() the widget at init a call of Refresh() of its parent is needed that it is redrawn to the canvas
	//(see https://github.com/fyne-io/fyne/issues/4494)
	//Needed for statusbar to be drawn correctly
	//renderer.Layout(renderer.widget.Size())
	renderer.main.Refresh()
	renderer.widget.sidebar.Refresh()
	renderer.widget.statusbar.Refresh()
	renderer.widget.defaultusernamebrowser.Refresh()
}

func (renderer *lantyRenderer) Destroy() {}
