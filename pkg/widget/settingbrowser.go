package widget

import (
	"errors"
	"regexp"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty-client/pkg/controller"
)

type SettingsBrowser struct {
	widget.BaseWidget

	controller *controller.Controller
	window     fyne.Window
	form       *Form

	serverurl         *Entry
	gamedirectory     *Entry
	username          *Entry
	downloaddirectory *Entry

	settingschanged chan struct{}
}

func NewSettingsBrowser(controller *controller.Controller, window fyne.Window) (settingsbrowser *SettingsBrowser) {
	settingsbrowser = &SettingsBrowser{
		controller:        controller,
		window:            window,
		form:              NewForm(),
		serverurl:         NewEntry(),
		gamedirectory:     NewEntry(),
		username:          NewEntry(),
		downloaddirectory: NewEntry(),
		settingschanged:   make(chan struct{}, 50),
	}
	settingsbrowser.ExtendBaseWidget(settingsbrowser)

	settingsbrowser.serverurl.SetText(controller.Settings.Settings().ServerURL)
	settingsbrowser.serverurl.OnFocusChanged = func(b bool) {
		if !b {
			controller.Settings.SetServerURL(settingsbrowser.serverurl.Text)
		}
	}
	settingsbrowser.serverurl.OnSubmitted = func(s string) {
		controller.Settings.SetServerURL(settingsbrowser.serverurl.Text)
	}
	settingsbrowser.form.AppendItem(NewFormItem("Server URL", settingsbrowser.serverurl))

	settingsbrowser.gamedirectory.SetText(controller.Settings.Settings().GameDirectory)
	settingsbrowser.gamedirectory.OnFocusChanged = func(b bool) {
		if !b {
			controller.Settings.SetGameDirectory(settingsbrowser.gamedirectory.Text)
		}
	}
	settingsbrowser.gamedirectory.OnSubmitted = func(s string) {
		controller.Settings.SetGameDirectory(settingsbrowser.gamedirectory.Text)
	}
	gamedirectoryexplorer := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), settingsbrowser.gamedirectoryExplorerCallback)
	gamedirectory := container.NewBorder(nil, nil, nil, gamedirectoryexplorer, settingsbrowser.gamedirectory)
	settingsbrowser.form.AppendItem(NewFormItem("Game Directory", gamedirectory))

	settingsbrowser.username.SetText(controller.Settings.Settings().Username)
	settingsbrowser.username.Validator = func(username string) error {
		match, err := regexp.MatchString("^(?:[a-zA-Z]|[0-9]|-)+$", username)
		if !match || err != nil {
			return errors.New("only alphanumeric characters and \"-\" allowed")
		}
		return nil
	}
	settingsbrowser.username.OnFocusChanged = func(b bool) {
		if !b && settingsbrowser.username.Validate() == nil {
			controller.Settings.SetUsername(settingsbrowser.username.Text)
		}
	}
	settingsbrowser.username.OnSubmitted = func(s string) {
		if settingsbrowser.username.Validate() == nil {
			controller.Settings.SetUsername(settingsbrowser.username.Text)
		}
	}
	settingsbrowser.form.AppendItem(NewFormItem("Username", settingsbrowser.username))

	settingsbrowser.downloaddirectory.SetText(controller.Settings.Settings().DownloadDirectory)
	settingsbrowser.downloaddirectory.OnFocusChanged = func(b bool) {
		if !b {
			controller.Settings.SetDownloadDirectory(settingsbrowser.downloaddirectory.Text)
		}
	}
	settingsbrowser.downloaddirectory.OnSubmitted = func(s string) {
		controller.Settings.SetDownloadDirectory(settingsbrowser.downloaddirectory.Text)
	}
	downloaddirectoryexplorer := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), settingsbrowser.downloaddirectoryExplorerCallback)
	downloaddirectory := container.NewBorder(nil, nil, nil, downloaddirectoryexplorer, settingsbrowser.downloaddirectory)
	settingsbrowser.form.AppendItem(NewFormItem("Download Directory", downloaddirectory))

	settingsbrowser.form.HideSubmit()
	settingsbrowser.form.SetCancelText("Reset")
	settingsbrowser.form.OnCancel = func() {
		settingsbrowser.ResetData()
		controller.Status.Info("Settings resetted", 3*time.Second)
	}

	controller.Settings.Subscribe(settingsbrowser.settingschanged)
	settingsbrowser.run()

	return settingsbrowser
}

func (widget *SettingsBrowser) run() {
	widget.controller.WaitGroup().Add(1)
	go widget.settingsUpdater()
}

func (widget *SettingsBrowser) settingsUpdater() {
	defer widget.controller.WaitGroup().Done()
	for {
		select {
		case <-widget.controller.Context().Done():
			log.Trace().Msg("exiting SettingsBrowser settingsUpdater()")
			return
		case <-widget.settingschanged:
			widget.serverurl.SetText(widget.controller.Settings.Settings().ServerURL)
			widget.gamedirectory.SetText(widget.controller.Settings.Settings().GameDirectory)
			widget.username.SetText(widget.controller.Settings.Settings().Username)
			widget.downloaddirectory.SetText(widget.controller.Settings.Settings().DownloadDirectory)
			widget.Refresh()
		}
	}
}

func (widget *SettingsBrowser) gamedirectoryExplorerCallback() {
	folderdialog := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
		if uri == nil || err != nil {
			return
		}
		widget.gamedirectory.SetText(uri.Path())
		if widget.gamedirectory.Validate() == nil {
			widget.controller.Settings.SetGameDirectory(widget.gamedirectory.Text)
			widget.controller.Settings.Save()
		}
	}, widget.window)

	dialogStartURI, err := storage.ListerForURI(storage.NewFileURI(widget.controller.Settings.Settings().GameDirectory))
	if err == nil {
		folderdialog.SetLocation(dialogStartURI)
	}

	//This will make the folderopen dialog to be "fullscreen" inside the app
	folderdialog.Resize(fyne.NewSize(10000, 10000))
	folderdialog.Show()
}

func (widget *SettingsBrowser) downloaddirectoryExplorerCallback() {
	folderdialog := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
		if uri == nil || err != nil {
			return
		}
		widget.downloaddirectory.SetText(uri.Path())
		if widget.downloaddirectory.Validate() == nil {
			widget.controller.Settings.SetDownloadDirectory(widget.downloaddirectory.Text)
			widget.controller.Settings.Save()
		}
	}, widget.window)

	dialogStartURI, err := storage.ListerForURI(storage.NewFileURI(widget.controller.Settings.Settings().DownloadDirectory))
	if err == nil {
		folderdialog.SetLocation(dialogStartURI)
	}

	//This will make the folderopen dialog to be "fullscreen" inside the app
	folderdialog.Resize(fyne.NewSize(10000, 10000))
	folderdialog.Show()
}

func (widget *SettingsBrowser) ResetData() {
	widget.serverurl.SetText(widget.controller.Settings.Settings().ServerURL)
	widget.gamedirectory.SetText(widget.controller.Settings.Settings().GameDirectory)
	widget.username.SetText(widget.controller.Settings.Settings().Username)
	widget.downloaddirectory.SetText(widget.controller.Settings.Settings().DownloadDirectory)
	widget.Refresh()
}

func (w *SettingsBrowser) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(w.form)
}
