package widget

import (
	"errors"
	"regexp"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/seternate/go-lanty-client/pkg/controller"
)

type SettingsBrowser struct {
	widget.BaseWidget

	controller    *controller.Controller
	form          *Form
	serverurl     binding.String
	gamedirectory binding.String
	username      binding.String
}

func NewSettingsBrowser(controller *controller.Controller) (settingsbrowser *SettingsBrowser) {
	settingsbrowser = &SettingsBrowser{
		controller: controller,
		form:       NewForm(),
	}
	settingsbrowser.ExtendBaseWidget(settingsbrowser)

	serverurlbind := binding.NewString()
	serverurlbind.Set(controller.Settings.Settings().ServerURL)
	serverurlentry := widget.NewEntryWithData(serverurlbind)
	settingsbrowser.form.AppendItem(NewFormItem("Server URL", serverurlentry))

	gamedirectorybind := binding.NewString()
	gamedirectorybind.Set(controller.Settings.Settings().GameDirectory)
	gamedirectoryentry := widget.NewEntryWithData(gamedirectorybind)
	settingsbrowser.form.AppendItem(NewFormItem("Game Directory", gamedirectoryentry))

	usernamebind := binding.NewString()
	usernamebind.Set(controller.Settings.Settings().Username)
	usernameentry := widget.NewEntryWithData(usernamebind)
	usernameentry.Validator = func(s string) error {
		match, err := regexp.MatchString("^(?:[a-zA-Z]|[0-9]|-)+$", s)
		if !match || err != nil {
			return errors.New("only alphanumeric characters and \"-\" allowed")
		}
		return nil
	}
	settingsbrowser.form.AppendItem(NewFormItem("Username", usernameentry))
	settingsbrowser.form.SetSubmitText("Save")
	settingsbrowser.form.SetCancelText("Reset")
	settingsbrowser.form.OnSubmit = func() {
		serverurl, _ := serverurlbind.Get()
		controller.Settings.SetServerURL(serverurl)
		gamedirectory, _ := gamedirectorybind.Get()
		controller.Settings.SetGameDirectory(gamedirectory)
		username, _ := usernamebind.Get()
		controller.Settings.SetUsername(username)
		err := controller.Settings.Settings().Save()
		if err != nil {
			controller.Status.Error("Error saving settings: "+err.Error(), 3*time.Second)
		} else {
			controller.Status.Info("Settings successfully saved", 3*time.Second)
		}
	}
	settingsbrowser.form.OnCancel = func() {
		settingsbrowser.ResetData()
		controller.Status.Info("Settings resetted", 3*time.Second)
	}

	settingsbrowser.serverurl = serverurlbind
	settingsbrowser.gamedirectory = gamedirectorybind
	settingsbrowser.username = usernamebind

	return settingsbrowser
}

func (widget *SettingsBrowser) ResetData() {
	widget.serverurl.Set(widget.controller.Settings.Settings().ServerURL)
	widget.gamedirectory.Set(widget.controller.Settings.Settings().GameDirectory)
	widget.username.Set(widget.controller.Settings.Settings().Username)
	widget.Refresh()
}

func (w *SettingsBrowser) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(w.form)
}
