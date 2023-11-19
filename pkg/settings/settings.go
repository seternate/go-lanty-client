package settings

import (
	"os"

	"gopkg.in/yaml.v3"
)

const (
	APPLICATION_NAME = "Lanty"
	SETTINGS_PATH    = "settings/lanty.yaml"
)

type Settings struct {
	ServerURL     string `json:"serverurl" yaml:"serverurl"`
	GameDirectory string `json:"gamedirectory" yaml:"gamedirectory"`

	filepath string
}

func LoadSettings(filepath string) (*Settings, error) {
	rawYAML, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	s := Settings{filepath: filepath}
	err = yaml.Unmarshal(rawYAML, &s)
	if err != nil {
		return nil, err
	}

	return &s, err
}
