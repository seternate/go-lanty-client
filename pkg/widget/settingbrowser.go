package widget

import (
	"errors"
	"os"
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

	serverurl     *Entry
	gamedirectory *Entry
	username      *Entry

	settingschanged chan struct{}
}

func NewSettingsBrowser(controller *controller.Controller, window fyne.Window) (settingsbrowser *SettingsBrowser) {
	settingsbrowser = &SettingsBrowser{
		controller:      controller,
		window:          window,
		form:            NewForm(),
		serverurl:       NewEntry(),
		gamedirectory:   NewEntry(),
		username:        NewEntry(),
		settingschanged: make(chan struct{}, 50),
	}
	settingsbrowser.ExtendBaseWidget(settingsbrowser)

	settingsbrowser.serverurl.SetText(controller.Settings.Settings().ServerURL)
	settingsbrowser.serverurl.OnFocusChanged = func(b bool) {
		if !b && settingsbrowser.serverurl.Validate() == nil {
			controller.Settings.SetServerURL(settingsbrowser.serverurl.Text)
			controller.Settings.Save()
		}
	}
	settingsbrowser.serverurl.OnSubmitted = func(s string) {
		if settingsbrowser.serverurl.Validate() == nil {
			controller.Settings.SetServerURL(settingsbrowser.serverurl.Text)
			controller.Settings.Save()
		}
	}
	settingsbrowser.form.AppendItem(NewFormItem("Server URL", settingsbrowser.serverurl))

	settingsbrowser.gamedirectory.SetText(controller.Settings.Settings().GameDirectory)
	settingsbrowser.gamedirectory.Validator = func(path string) error {
		info, err := os.Stat(path)
		if errors.Is(err, os.ErrNotExist) {
			return errors.New("folder does not exist")
		}
		if !info.IsDir() {
			return errors.New("path is not a folder")
		}
		return nil
	}
	settingsbrowser.gamedirectory.OnFocusChanged = func(b bool) {
		if !b && settingsbrowser.gamedirectory.Validate() == nil {
			controller.Settings.SetGameDirectory(settingsbrowser.gamedirectory.Text)
			controller.Settings.Save()
		}
	}
	settingsbrowser.gamedirectory.OnSubmitted = func(s string) {
		if settingsbrowser.gamedirectory.Validate() == nil {
			controller.Settings.SetGameDirectory(settingsbrowser.gamedirectory.Text)
			controller.Settings.Save()
		}
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
			controller.Settings.Save()
		}
	}
	settingsbrowser.username.OnSubmitted = func(s string) {
		if settingsbrowser.username.Validate() == nil {
			controller.Settings.SetUsername(settingsbrowser.username.Text)
			controller.Settings.Save()
		}
	}
	settingsbrowser.form.AppendItem(NewFormItem("Username", settingsbrowser.username))

	settingsbrowser.form.SetSubmitText("Save")
	settingsbrowser.form.SetCancelText("Reset")
	settingsbrowser.form.OnSubmit = func() {
		controller.Settings.SetServerURL(settingsbrowser.serverurl.Text)
		controller.Settings.SetGameDirectory(settingsbrowser.gamedirectory.Text)
		controller.Settings.SetUsername(settingsbrowser.username.Text)
		if controller.Settings.Save() == nil {
			controller.Status.Info("Settings successfully saved", 3*time.Second)
		}
	}
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

func (widget *SettingsBrowser) ResetData() {
	widget.serverurl.SetText(widget.controller.Settings.Settings().ServerURL)
	widget.gamedirectory.SetText(widget.controller.Settings.Settings().GameDirectory)
	widget.username.SetText(widget.controller.Settings.Settings().Username)
	widget.Refresh()
}

func (w *SettingsBrowser) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(w.form)
}
