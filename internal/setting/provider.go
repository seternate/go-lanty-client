package setting

type SettingsProvider interface {
	Get() (Settings, error)
}

type SettingsStore interface {
	Get() (Settings, error)
	Save(Settings) error
}
