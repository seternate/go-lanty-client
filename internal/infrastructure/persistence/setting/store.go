package settings

import (
	"fmt"
	"os"
	"sync"

	"github.com/seternate/go-lanty-client/internal/application/setting"
	"github.com/seternate/go-lanty-client/internal/application/setting/service"
	"gopkg.in/yaml.v3"
)

type fileSettings struct {
	ServerURL     string `yaml:"serverurl"`
	GameDirectory string `yaml:"gamedirectory"`
	Username      string `yaml:"username"`
}

var _ service.SettingsStore = (*FileSettingsStore)(nil)

type FileSettingsStore struct {
	settingsFilePath string
	settings         fileSettings
	mu               sync.RWMutex
}

func NewFileSettingsStore(settingsPath string) (*FileSettingsStore, error) {
	store := &FileSettingsStore{
		settingsFilePath: settingsPath,
		settings:         fileSettings{},
	}

	if err := store.loadFromFilesystem(); err != nil {
		return nil, fmt.Errorf("failed to load settings from filesystem: %w", err)
	}

	return store, nil
}

func (store *FileSettingsStore) loadFromFilesystem() error {
	data, err := os.ReadFile(store.settingsFilePath)
	if err != nil {
		return fmt.Errorf("failed to read settings file: %w", err)
	}

	var parsed fileSettings
	if err := yaml.Unmarshal(data, &parsed); err != nil {
		return fmt.Errorf("failed to decode settings file: %w", err)
	}

	store.settings = parsed

	return nil
}

func (store *FileSettingsStore) Save(settings setting.Settings) error {
	data := fileSettings{
		ServerURL:     settings.ServerURL,
		GameDirectory: settings.GameDirectory,
		Username:      settings.Username,
	}

	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	if err := os.WriteFile(store.settingsFilePath, yamlData, 0644); err != nil {
		return fmt.Errorf("failed to write settings file: %w", err)
	}

	store.mu.Lock()
	store.settings = data
	store.mu.Unlock()

	return nil
}

func (store *FileSettingsStore) GetSettings() (setting.Settings, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	return setting.Settings{
		ServerURL:     store.settings.ServerURL,
		GameDirectory: store.settings.GameDirectory,
		Username:      store.settings.Username,
	}, nil
}
