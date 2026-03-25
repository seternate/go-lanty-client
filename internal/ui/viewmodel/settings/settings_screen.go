package settingsviewmodel

import (
	"fmt"
	"slices"
	"strings"

	"fyne.io/fyne/v2/data/binding"
	"github.com/rs/zerolog"
	"github.com/seternate/go-lanty-client/internal/setting"
	"github.com/seternate/go-lanty-client/internal/ui/viewmodel"
)

type SettingsScreen struct {
	Header *viewmodel.HeaderScrollScreen

	serverurl     binding.Untyped
	gamedirectory binding.Untyped
	username      binding.Untyped

	settingStore setting.SettingsStore
	logger       zerolog.Logger
}

func NewSettingsScreen(logger zerolog.Logger, settingStore setting.SettingsStore, folderPicker FolderPicker) (*SettingsScreen, error) {
	vm := &SettingsScreen{
		Header:        viewmodel.NewHeaderScrollScreen("Settings"),
		serverurl:     binding.NewUntyped(),
		gamedirectory: binding.NewUntyped(),
		username:      binding.NewUntyped(),
		settingStore:  settingStore,
		logger:        logger,
	}

	settings, err := vm.settingStore.Get()
	if err != nil {
		return nil, fmt.Errorf("load settings for settings screen: %w", err)
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

	return vm, nil
}

func (vm *SettingsScreen) Save() {
	settings, err := vm.settingStore.Get()
	if err != nil {
		vm.logger.Error().Err(err).Msg("settings save: failed to load current settings")
		return
	}

	serverurl, err := vm.serverurl.Get()
	if err != nil {
		vm.logger.Error().Err(err).Msg("settings save: failed to read server URL binding")
		return
	}
	gamedirectory, err := vm.gamedirectory.Get()
	if err != nil {
		vm.logger.Error().Err(err).Msg("settings save: failed to read game directory binding")
		return
	}
	username, err := vm.username.Get()
	if err != nil {
		vm.logger.Error().Err(err).Msg("settings save: failed to read username binding")
		return
	}

	settings.ServerURL = serverurl.(*SettingTile).GetValue()
	settings.GameDirectory = gamedirectory.(*SettingTile).GetValue()
	settings.Username = username.(*SettingTile).GetValue()

	if err := vm.settingStore.Save(settings); err != nil {
		vm.logger.Error().Err(err).Msg("settings save: failed to persist settings")
	}
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
