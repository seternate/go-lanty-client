package setting

import (
	"github.com/seternate/go-lanty-client/internal/application/setting"
	settingsappsrv "github.com/seternate/go-lanty-client/internal/application/setting/service"
	model "github.com/seternate/go-lanty-client/internal/ui/model/settings"
)

type SettingsFetcher interface {
	GetSettings() (setting.Settings, error)
}

type SettingsController struct {
	SettingsModel *model.SettingsModel

	settingsUpdater settingsappsrv.SettingsUpdateService
	settingsFetcher SettingsFetcher
}

func NewSettingsController(settingsUpdater settingsappsrv.SettingsUpdateService, settingsFetcher SettingsFetcher) (*SettingsController, error) {
	controller := &SettingsController{
		SettingsModel:   model.NewSettingsModel(),
		settingsUpdater: settingsUpdater,
		settingsFetcher: settingsFetcher,
	}

	settings, err := settingsFetcher.GetSettings()
	if err != nil {
		return nil, err
	}

	controller.SettingsModel.SetServerURL(settings.ServerURL)
	controller.SettingsModel.SetGameDirectory(settings.GameDirectory)
	controller.SettingsModel.SetUsername(settings.Username)

	return controller, nil
}

func (controller *SettingsController) OnSaveClicked() {
	serverURL, err := controller.SettingsModel.ServerURL.Get()
	if err != nil {
		return
	}
	gameDirectory, err := controller.SettingsModel.GameDirectory.Get()
	if err != nil {
		return
	}
	username, err := controller.SettingsModel.Username.Get()
	if err != nil {
		return
	}

	controller.SettingsModel.SetServerURL(serverURL)
	controller.SettingsModel.SetGameDirectory(gameDirectory)
	controller.SettingsModel.SetUsername(username)

	controller.settingsUpdater.ApplySettings(setting.Settings{
		ServerURL:     serverURL,
		GameDirectory: gameDirectory,
		Username:      username,
	})
}

func (controller *SettingsController) OnResetClicked() {
	settings, err := controller.settingsFetcher.GetSettings()
	if err != nil {
		return
	}

	controller.SettingsModel.SetServerURL(settings.ServerURL)
	controller.SettingsModel.SetGameDirectory(settings.GameDirectory)
	controller.SettingsModel.SetUsername(settings.Username)
}
