package app

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty-client/internal/network"
	gameview "github.com/seternate/go-lanty-client/internal/ui/view/game"
	"github.com/seternate/go-lanty-client/internal/ui/view/sidebar"
)

type AppShell struct {
	window fyne.Window

	icon    fyne.Resource
	version string

	contentScreen *fyne.Container
	gameScreen    fyne.CanvasObject
}

func NewAppShell(name string, icon fyne.Resource, version string) *AppShell {
	fyneApp := app.NewWithID("com.seternate.lanty")

	window := fyneApp.NewWindow(appTitle(name))
	window.SetPadded(false)
	window.Resize(fyne.NewSize(1024, 600))

	if icon != nil {
		window.SetIcon(icon)
	}

	contentScreen := container.NewStack()
	sidebarView := sidebar.NewSidebar(icon, version)
	mainScreen := container.NewBorder(nil, nil, sidebarView, nil, contentScreen)
	window.SetContent(mainScreen)

	appShell := &AppShell{
		window:        window,
		icon:          icon,
		version:       version,
		contentScreen: contentScreen,
	}

	return appShell
}

func (appShell *AppShell) SetGameTile(gameTile *gameview.GameTile) {
	appShell.gameScreen = gameTile
	appShell.contentScreen.Objects = []fyne.CanvasObject{container.NewVBox(appShell.gameScreen)}
	appShell.contentScreen.Refresh()
}

func (appShell *AppShell) ShowAndRun() {
	appShell.window.ShowAndRun()
}

func (appShell *AppShell) Quit() {
	fyne.CurrentApp().Quit()
}

func appTitle(name string) string {
	ip, err := network.GetOutboundIP()
	if err != nil {
		log.Error().Err(err).Msg("failed to get outbound IP for app title")
		return name
	}
	return fmt.Sprintf("%s - %s", name, ip.String())
}

// import (
// 	"fmt"

// 	"fyne.io/fyne/v2"
// 	"fyne.io/fyne/v2/app"
// 	"fyne.io/fyne/v2/container"
// 	"github.com/rs/zerolog/log"
// 	"github.com/seternate/go-lanty-client/internal/infrastructure/network"
// 	gamecontroller "github.com/seternate/go-lanty-client/internal/ui/controller/game"
// 	settingcontroller "github.com/seternate/go-lanty-client/internal/ui/controller/setting"
// 	usercontroller "github.com/seternate/go-lanty-client/internal/ui/controller/user"
// 	gameview "github.com/seternate/go-lanty-client/internal/ui/view/game"
// 	settingsview "github.com/seternate/go-lanty-client/internal/ui/view/settings"
// 	"github.com/seternate/go-lanty-client/internal/ui/view/sidebar"
// 	userview "github.com/seternate/go-lanty-client/internal/ui/view/user"
// )

// type App struct {
// 	Window fyne.Window

// 	icon    fyne.Resource
// 	version string

// 	contentArea     *fyne.Container
// 	gameContent     fyne.CanvasObject
// 	userContent     fyne.CanvasObject
// 	settingsContent fyne.CanvasObject
// }

// func NewApp(
// 	name string,
// 	icon fyne.Resource,
// 	version string,
// ) *App {
// 	fyneApp := app.NewWithID("com.seternate.lanty")
// 	window := fyneApp.NewWindow(appTitle(name))

// 	window.SetPadded(false)
// 	window.Resize(fyne.NewSize(1024, 600))
// 	if icon != nil {
// 		window.SetIcon(icon)
// 	}

// 	appInstance := &App{
// 		Window:      window,
// 		icon:        icon,
// 		version:     version,
// 		contentArea: container.NewStack(),
// 	}

// 	return appInstance
// }

// func (app *App) Bootstrap(gameController *gamecontroller.GameController, userController *usercontroller.UserController, settingsController *settingcontroller.SettingsController) {
// 	gameList := gameview.NewGameList(gameController.GameListModel)
// 	gameList.OnPlayButtonPressed = func(slug string) {
// 		gameController.PlayGame(slug)
// 	}
// 	gameList.OnJoinButtonPressed = func(slug string) {
// 		gameController.JoinGame(slug)
// 	}
// 	gameList.OnServerButtonPressed = func(slug string) {
// 		gameController.StartGameServer(slug)
// 	}
// 	gameList.OnDownloadButtonPressed = func(slug string) {
// 		gameController.DownloadGame(slug)
// 	}
// 	gameList.OnOpenButtonPressed = func(slug string) {
// 		gameController.OpenGame(slug)
// 	}
// 	gameList.OnExtensionBadgePressed = func(slug string) {
// 		gameController.OpenExtensions(slug)
// 	}
// 	gameBrowser := gameview.NewGameBrowser(gameList, gameController.GameStatsModel)
// 	app.gameContent = container.NewStack(gameBrowser)

// 	userView := userview.NewUserView(userController.UserListModel)
// 	app.userContent = container.NewStack(userView)

// 	settingsView := settingsview.NewSettingsView(app.Window, settingsController.SettingsModel)
// 	settingsView.OnSavePressed = func() {
// 		settingsController.OnSaveClicked()
// 	}
// 	settingsView.OnResetPressed = func() {
// 		settingsController.OnResetClicked()
// 	}
// 	app.settingsContent = container.NewStack(settingsView)

// 	sidebarWidget := sidebar.NewSidebar(app.icon, app.version)
// 	sidebarWidget.OnGamesButtonPressed = func() {
// 		app.showGamesView()
// 	}
// 	sidebarWidget.OnUsersButtonPressed = func() {
// 		app.showUsersView()
// 	}
// 	sidebarWidget.OnSettingsButtonPressed = func() {
// 		app.showSettingsView()
// 	}

// 	contentArea := container.NewBorder(nil, nil, sidebarWidget, nil, app.contentArea)
// 	app.Window.SetContent(contentArea)
// 	app.showGamesView()
// }

// func (app *App) showGamesView() {
// 	app.contentArea.Objects = []fyne.CanvasObject{app.gameContent}
// 	app.contentArea.Refresh()
// }

// func (app *App) showUsersView() {
// 	app.contentArea.Objects = []fyne.CanvasObject{app.userContent}
// 	app.contentArea.Refresh()
// }

// func (app *App) showSettingsView() {
// 	app.contentArea.Objects = []fyne.CanvasObject{app.settingsContent}
// 	app.contentArea.Refresh()
// }

// func (app *App) ShowAndRun() {
// 	app.Window.ShowAndRun()
// }

// func (app *App) Quit() {
// 	fyne.CurrentApp().Quit()
// }

// func appTitle(name string) string {
// 	ip, err := network.GetOutboundIP()
// 	if err != nil {
// 		log.Error().Err(err).Msg("failed to get outbound IP for app title")
// 		return name
// 	}
// 	return fmt.Sprintf("%s - %s", name, ip.String())
// }
