package app

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"github.com/rs/zerolog"
	"github.com/seternate/go-lanty-client/internal/game"
	"github.com/seternate/go-lanty-client/internal/network"
	gameview "github.com/seternate/go-lanty-client/internal/ui/view/game"
	settingsview "github.com/seternate/go-lanty-client/internal/ui/view/setting"
	"github.com/seternate/go-lanty-client/internal/ui/view/sidebar"
	userview "github.com/seternate/go-lanty-client/internal/ui/view/user"
	gameviewmodel "github.com/seternate/go-lanty-client/internal/ui/viewmodel/game"
	settingsviewmodel "github.com/seternate/go-lanty-client/internal/ui/viewmodel/settings"
	userviewmodel "github.com/seternate/go-lanty-client/internal/ui/viewmodel/user"
)

type AppShell struct {
	logger zerolog.Logger
	window fyne.Window

	icon    fyne.Resource
	version string

	contentScreen  *fyne.Container
	gameScreen     fyne.CanvasObject
	argumentScreen fyne.CanvasObject
	userScreen     fyne.CanvasObject
	settingsScreen fyne.CanvasObject
}

func NewAppShell(logger zerolog.Logger, name string, icon fyne.Resource, version string) *AppShell {
	fyneApp := app.NewWithID("com.seternate.lanty")

	window := fyneApp.NewWindow(appTitle(name, logger))
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
		logger:        logger,
		window:        window,
		icon:          icon,
		version:       version,
		contentScreen: contentScreen,
	}

	sidebarView.OnGamesButtonPressed = func() {
		appShell.ShowGameScreen()
	}
	sidebarView.OnUsersButtonPressed = func() {
		appShell.ShowUserScreen()
	}
	sidebarView.OnSettingsButtonPressed = func() {
		appShell.ShowSettingsScreen()
	}

	return appShell
}

func (appShell *AppShell) Bootstrap(gameScreen *gameviewmodel.GameScreen, userScreen *userviewmodel.UserScreen, settingsscreen *settingsviewmodel.SettingsScreen) {
	appShell.gameScreen = gameview.NewGameScreen(gameScreen)
	appShell.userScreen = userview.NewUserScreen(userScreen)
	appShell.settingsScreen = settingsview.NewSettingsScreen(settingsscreen)

	appShell.contentScreen.Objects = []fyne.CanvasObject{appShell.gameScreen, appShell.userScreen, appShell.settingsScreen}

	appShell.ShowGameScreen()
}

func (appShell *AppShell) ShowGameScreen() {
	appShell.userScreen.Hide()
	appShell.settingsScreen.Hide()
	appShell.gameScreen.Show()
	if appShell.argumentScreen != nil {
		appShell.argumentScreen.Hide()
	}
}

func (appShell *AppShell) ShowUserScreen() {
	appShell.gameScreen.Hide()
	appShell.settingsScreen.Hide()
	appShell.userScreen.Show()
	if appShell.argumentScreen != nil {
		appShell.argumentScreen.Hide()
	}
}

func (appShell *AppShell) ShowSettingsScreen() {
	appShell.gameScreen.Hide()
	appShell.userScreen.Hide()
	appShell.settingsScreen.Show()
	if appShell.argumentScreen != nil {
		appShell.argumentScreen.Hide()
	}
}

func (appShell *AppShell) ShowFolderPickerDialog(location string, callback func(path string)) {
	folderDialog := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
		if err != nil {
			appShell.logger.Warn().Err(err).Msg("folder picker callback returned an error")
			return
		}
		if uri != nil {
			callback(uri.Path())
		}
	}, appShell.window)

	listabelLocation, err := storage.ListerForURI(storage.NewFileURI(location))
	if err != nil {
		appShell.logger.Error().Err(err).Str("location", location).Msg("failed to list URI for folder picker")
		return
	}

	folderDialog.SetLocation(listabelLocation)
	folderDialog.Resize(fyne.NewSize(800, 600))
	folderDialog.Show()
}

func (appShell *AppShell) ShowLaunchArgumentScreen(title string, info string, arguments []game.LaunchParam, onSubmit func(values []game.LaunchArg)) {
	cleanup := func() {
		if appShell.argumentScreen != nil {
			appShell.contentScreen.Remove(appShell.argumentScreen)
			appShell.argumentScreen = nil
		}
		appShell.ShowGameScreen()
	}

	wrappedSubmit := func(values []game.LaunchArg) {
		cleanup()

		if onSubmit != nil {
			onSubmit(values)
		}
	}

	onCancel := func() {
		cleanup()
	}

	launchArgumentScreen, err := gameviewmodel.NewLaunchArgumentScreen(title, info, arguments, wrappedSubmit, onCancel)
	if err != nil {
		appShell.logger.Error().Err(err).Str("title", title).Msg("failed to create launch argument screen")
		return
	}

	appShell.argumentScreen = gameview.NewLaunchArgumentScreen(launchArgumentScreen)

	appShell.gameScreen.Hide()
	appShell.userScreen.Hide()
	appShell.settingsScreen.Hide()

	appShell.contentScreen.Add(appShell.argumentScreen)
}

func (appShell *AppShell) ShowAndRun() {
	appShell.window.ShowAndRun()
}

func (appShell *AppShell) Quit() {
	fyne.CurrentApp().Quit()
}

func appTitle(name string, logger zerolog.Logger) string {
	ip, err := network.GetOutboundIP()
	if err != nil {
		logger.Error().Err(err).Msg("failed to get outbound IP for app title")
		return name
	}
	return fmt.Sprintf("%s - %s", name, ip.String())
}
