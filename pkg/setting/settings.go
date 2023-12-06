package setting

import "github.com/seternate/go-lanty/pkg/filesystem"

const (
	APPLICATION_NAME = "Lanty"
	SETTINGS_PATH    = "settings.yaml"
)

type Settings struct {
	ServerURL     string `yaml:"serverurl"`
	GameDirectory string `yaml:"gamedirectory"`
}

func LoadSettings() (s *Settings, err error) {
	err = filesystem.LoadFromYAMLFile(SETTINGS_PATH, &s)
	return
}
