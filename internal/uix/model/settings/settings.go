package model

import (
	"fyne.io/fyne/v2/data/binding"
	viewmodel "github.com/seternate/go-lanty-client/internal/ui/model"
)

var _ viewmodel.ChangeNotifier = (*SettingsModel)(nil)

type SettingsModel struct {
	Username      binding.String
	ServerURL     binding.String
	GameDirectory binding.String

	savedUsername      string
	savedServerURL     string
	savedGameDirectory string
}

func NewSettingsModel() *SettingsModel {
	return &SettingsModel{
		Username:           binding.NewString(),
		ServerURL:          binding.NewString(),
		GameDirectory:      binding.NewString(),
		savedUsername:      "",
		savedServerURL:     "",
		savedGameDirectory: "",
	}
}

func (model *SettingsModel) SetUsername(username string) {
	model.Username.Set(username)
	model.savedUsername = username
}

func (model *SettingsModel) SetServerURL(serverURL string) {
	model.ServerURL.Set(serverURL)
	model.savedServerURL = serverURL
}

func (model *SettingsModel) SetGameDirectory(gameDirectory string) {
	model.GameDirectory.Set(gameDirectory)
	model.savedGameDirectory = gameDirectory
}

func (model *SettingsModel) IsDirty() bool {
	username, err := model.Username.Get()
	if err != nil {
		return false
	}
	serverURL, err := model.ServerURL.Get()
	if err != nil {
		return false
	}
	gameDirectory, err := model.GameDirectory.Get()
	if err != nil {
		return false
	}
	return model.savedUsername != username ||
		model.savedServerURL != serverURL ||
		model.savedGameDirectory != gameDirectory
}

func (model *SettingsModel) AddChangeListener(fn func()) {
	l := binding.NewDataListener(fn)
	model.Username.AddListener(l)
	model.ServerURL.AddListener(l)
	model.GameDirectory.AddListener(l)
}
