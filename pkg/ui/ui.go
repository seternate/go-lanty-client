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
	"github.com/seternate/go-lanty-client/pkg/controller"
	"github.com/seternate/go-lanty-client/pkg/setting"
	"github.com/seternate/go-lanty-client/pkg/ui/gamebrowser"
	"github.com/seternate/go-lanty-client/pkg/ui/settingbrowser"
	"github.com/seternate/go-lanty/pkg/network"
)

type UI struct {
	application fyne.App
	mainWindow  fyne.Window
	sidebar     *fyne.Container

	gamebrowser    *gamebrowser.GameBrowser
	settingbrowser *settingbrowser.Settingbrowser

	controller *controller.Controller
}

func NewUI(controller *controller.Controller) *UI {
	lantyUI := &UI{
		controller: controller,
	}
	lantyUI.createApplication()

	return lantyUI
}

func (l *UI) createApplication() {
	l.application = app.New()
	l.createMainWindow()
}

func (l *UI) createMainWindow() {
	l.createSidebar()
	l.createBrowsers()

	ip, err := network.GetOutboundIP()
	var applicationName string
	if err != nil {
		applicationName = fmt.Sprintf("%s", setting.APPLICATION_NAME)
	} else {
		applicationName = fmt.Sprintf("%s - %s", setting.APPLICATION_NAME, ip.String())
	}

	l.mainWindow = l.CreateWindow(applicationName, container.NewBorder(nil, nil, l.sidebar, nil, l.gamebrowser.Container))
	l.mainWindow.Resize(fyne.Size{Width: 990, Height: 480})

}

func (l *UI) CreateWindow(title string, object fyne.CanvasObject) (window fyne.Window) {
	window = l.application.NewWindow(title)
	window.SetContent(object)
	return
}

func (l *UI) createSidebar() {
	appName := canvas.NewText("lanty App", color.White)
	appName.TextStyle.Bold = true
	appName.TextSize = 32

	gamesButton := widget.NewButtonWithIcon("Games", theme.MediaPlayIcon(), func() {
		l.mainWindow.SetContent(container.NewBorder(nil, nil, l.sidebar, nil, l.gamebrowser.Container))
	})
	gamesButton.Alignment = widget.ButtonAlignLeading

	usersButton := widget.NewButtonWithIcon("Users", theme.AccountIcon(), nil)
	usersButton.Alignment = widget.ButtonAlignLeading

	settingsButton := widget.NewButtonWithIcon("Settings", theme.SettingsIcon(), func() {
		l.mainWindow.SetContent(container.NewBorder(nil, nil, l.sidebar, nil, l.settingbrowser))
	})
	settingsButton.Alignment = widget.ButtonAlignLeading

	l.sidebar = container.NewMax(canvas.NewRectangle(color.RGBA{126, 126, 126, 255}), container.NewPadded(container.NewVBox(appName, gamesButton, usersButton, settingsButton)))
}

func (ui *UI) createBrowsers() {
	gamebrowser.Application = ui.application
	gb := gamebrowser.NewGameBrowser(ui.controller.Game)
	//sb, _ := settingbrowser.NewSettingBrowser(l.lanty.Settings)

	ui.gamebrowser = gb
	//l.settingbrowser = sb
}

func (l *UI) ShowAndRun() {
	l.mainWindow.ShowAndRun()
}
