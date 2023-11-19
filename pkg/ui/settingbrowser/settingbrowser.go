package settingbrowser

import (
	"errors"
	"net/url"
	"os"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/seternate/go-lanty/pkg/settings"
)

type Settingbrowser struct {
	*widget.Form

	settings *settings.Settings
}

func NewSettingBrowser(settings *settings.Settings) (*Settingbrowser, error) {
	settingbrowser := &Settingbrowser{settings: settings}

	serverurlbinding := binding.BindString(&settings.ServerURL)
	gamedirectorybinding := binding.BindString(&settings.GameDirectory)

	serverurl := widget.NewEntryWithData(serverurlbinding)
	gamedirectory := widget.NewEntryWithData(gamedirectorybinding)

	serverurl.Validator = func(s string) error {
		u, err := url.Parse(s)
		if err != nil || u.Scheme == "" || u.Host == "" {
			return errors.New("No valid URL")
		}
		return nil
	}
	gamedirectory.Validator = func(s string) error {
		if _, err := os.Stat(s); os.IsNotExist(err) {
			return errors.New("Path does not exist")
		}
		return nil
	}

	settingbrowser.Form = &widget.Form{}

	settingbrowser.Form.Append("Server URL", serverurl)
	settingbrowser.Form.Append("Game Directory", gamedirectory)

	return settingbrowser, nil
}
