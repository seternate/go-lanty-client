package event

type SettingsUpdateEventHandler func()

const SettingsUpdatedEvent string = "settings:updated"
