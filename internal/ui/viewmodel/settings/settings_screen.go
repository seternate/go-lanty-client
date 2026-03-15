package settingsviewmodel

import (
	"slices"
	"strings"

	"fyne.io/fyne/v2/data/binding"
	"github.com/seternate/go-lanty-client/internal/setting"
	"github.com/seternate/go-lanty-client/internal/ui/viewmodel"
)

type SettingsScreen struct {
	Header *viewmodel.HeaderScrollScreen

	serverurl     binding.Untyped
	gamedirectory binding.Untyped
	username      binding.Untyped

	settingStore setting.SettingsStore
}

func NewSettingsScreen(settingStore setting.SettingsStore, folderPicker FolderPicker) *SettingsScreen {
	vm := &SettingsScreen{
		Header:        viewmodel.NewHeaderScrollScreen("Settings"),
		serverurl:     binding.NewUntyped(),
		gamedirectory: binding.NewUntyped(),
		username:      binding.NewUntyped(),
		settingStore:  settingStore,
	}

	settings, err := vm.settingStore.Get()
	if err != nil {
		return nil
	}

	serverurl := binding.NewString()
	gamedirectory := binding.NewString()
	username := binding.NewString()
	serverurl.Set(settings.ServerURL)
	gamedirectory.Set(settings.GameDirectory)
	username.Set(settings.Username)

	serverurlTile := NewSettingTile("Server URL", "Base URL to the API server", serverurl, false, folderPicker, vm.Save)
	gamedirectoryTile := NewSettingTile("Game Directory", "Directory for game installations and downloads", gamedirectory, true, folderPicker, vm.Save)
	usernameTile := NewSettingTile("Username", "Your username seen by others", username, false, folderPicker, vm.Save)

	vm.serverurl.Set(serverurlTile)
	vm.gamedirectory.Set(gamedirectoryTile)
	vm.username.Set(usernameTile)

	return vm
}

func (vm *SettingsScreen) Save() {
	settings, err := vm.settingStore.Get()
	if err != nil {
		return
	}

	serverurl, err := vm.serverurl.Get()
	if err != nil {
		return
	}
	gamedirectory, err := vm.gamedirectory.Get()
	if err != nil {
		return
	}
	username, err := vm.username.Get()
	if err != nil {
		return
	}

	settings.ServerURL = serverurl.(*SettingTile).GetValue()
	settings.GameDirectory = gamedirectory.(*SettingTile).GetValue()
	settings.Username = username.(*SettingTile).GetValue()

	vm.settingStore.Save(settings)
}

func (vm *SettingsScreen) GetSettingsTiles() []*SettingTile {
	settingsTiles := make([]*SettingTile, 0)

	serverurl, err := vm.serverurl.Get()
	if err != nil {
		return nil
	}
	gamedirectory, err := vm.gamedirectory.Get()
	if err != nil {
		return nil
	}
	username, err := vm.username.Get()
	if err != nil {
		return nil
	}

	if tile, ok := serverurl.(*SettingTile); ok {
		settingsTiles = append(settingsTiles, tile)
	}
	if tile, ok := gamedirectory.(*SettingTile); ok {
		settingsTiles = append(settingsTiles, tile)
	}
	if tile, ok := username.(*SettingTile); ok {
		settingsTiles = append(settingsTiles, tile)
	}

	slices.SortStableFunc(settingsTiles, func(a, b *SettingTile) int {
		return strings.Compare(a.GetLabel(), b.GetLabel())
	})

	return settingsTiles
}
