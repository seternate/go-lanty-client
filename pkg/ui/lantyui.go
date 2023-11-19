package ui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	lanty "github.com/seternate/go-lanty/pkg"
	"github.com/seternate/go-lanty/pkg/settings"
	"github.com/seternate/go-lanty/pkg/ui/gamebrowser"
	"github.com/seternate/go-lanty/pkg/ui/settingbrowser"
	"github.com/seternate/lanty-api-golang/pkg/util"
)

type LantyUI struct {
	application fyne.App
	mainWindow  fyne.Window
	sidebar     *fyne.Container

	gamebrowser    *gamebrowser.Gamebrowser
	settingbrowser *settingbrowser.Settingbrowser

	lanty *lanty.Lanty
}

func NewLantyUI(lanty *lanty.Lanty) *LantyUI {
	lantyUI := &LantyUI{
		lanty: lanty,
	}
	lantyUI.createApplication()

	return lantyUI
}

func (l *LantyUI) createApplication() {
	l.application = app.New()
	l.createMainWindow()
}

func (l *LantyUI) createMainWindow() {
	l.mainWindow = l.createBlankWindow()
	l.mainWindow.Resize(fyne.Size{Width: 990, Height: 480})
	l.createSidebar()
	l.createBrowsers()

	l.mainWindow.SetContent(container.NewBorder(nil, nil, l.sidebar, nil, l.gamebrowser))
}

func (l *LantyUI) createBlankWindow() fyne.Window {
	applicationName := fmt.Sprintf("%s - %s", settings.APPLICATION_NAME, util.GetOutboundIP().String())
	return l.application.NewWindow(applicationName)
}

func (l *LantyUI) createSidebar() {
	appName := canvas.NewText("lanty App", color.White)
	appName.TextStyle.Bold = true
	appName.TextSize = 32

	gamesButton := widget.NewButtonWithIcon("Games", theme.MediaPlayIcon(), func() {
		l.mainWindow.SetContent(container.NewBorder(nil, nil, l.sidebar, nil, l.gamebrowser))
	})
	gamesButton.Alignment = widget.ButtonAlignLeading

	usersButton := widget.NewButtonWithIcon("Users", theme.AccountIcon(), nil)
	usersButton.Alignment = widget.ButtonAlignLeading

	serverbrowserButton := widget.NewButtonWithIcon("Servers", theme.ListIcon(), nil)
	serverbrowserButton.Alignment = widget.ButtonAlignLeading

	settingsButton := widget.NewButtonWithIcon("Settings", theme.SettingsIcon(), func() {
		l.mainWindow.SetContent(container.NewBorder(nil, nil, l.sidebar, nil, l.settingbrowser))
	})
	settingsButton.Alignment = widget.ButtonAlignLeading

	l.sidebar = container.NewMax(canvas.NewRectangle(color.RGBA{126, 126, 126, 255}), container.NewPadded(container.NewVBox(appName, gamesButton, serverbrowserButton, usersButton, settingsButton)))
}

func (l *LantyUI) createBrowsers() {
	gb, _ := gamebrowser.NewGameBrowser(l.lanty.Client, l.lanty.Settings, *l.lanty.Games, l.lanty.Downloader)
	sb, _ := settingbrowser.NewSettingBrowser(l.lanty.Settings)

	l.gamebrowser = gb
	l.settingbrowser = sb
}

func (l *LantyUI) ShowAndRun() {
	l.mainWindow.ShowAndRun()
}
