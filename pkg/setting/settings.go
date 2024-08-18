package setting

import (
	"path"

	"github.com/kardianos/osext"
	"github.com/seternate/go-lanty/pkg/filesystem"
)

const (
	APPLICATION_NAME = "Lanty"
	SETTINGS_PATH    = "settings.yaml"
	VERSION          = "v0.1.0-beta"
)

type Settings struct {
	ServerURL         string `yaml:"serverurl"`
	GameDirectory     string `yaml:"gamedirectory"`
	Username          string `yaml:"username"`
	DownloadDirectory string `yaml:"downloaddirectory"`
}

func LoadSettings() (s *Settings, err error) {
	root, err := osext.ExecutableFolder()
	if err != nil {
		return
	}
	err = filesystem.LoadFromYAMLFile(path.Join(root, SETTINGS_PATH), &s)
	if err != nil {
		err = filesystem.LoadFromYAMLFile(SETTINGS_PATH, &s)
	}
	return
}

func (settings Settings) Save() (err error) {
	root, err := osext.ExecutableFolder()
	if err != nil {
		return
	}
	err = filesystem.SaveToYAMLFile(path.Join(root, SETTINGS_PATH), settings)
	if err != nil {
		err = filesystem.LoadFromYAMLFile(SETTINGS_PATH, settings)
	}
	return
}
