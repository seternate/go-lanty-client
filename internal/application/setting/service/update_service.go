package service

import (
	eventbus "github.com/asaskevich/EventBus"
	"github.com/seternate/go-lanty-client/internal/application/setting"
	"github.com/seternate/go-lanty-client/internal/application/setting/event"
)

type SettingsUpdateService interface {
	ApplySettings(settings setting.Settings) error
}

type SettingsStore interface {
	Save(setting.Settings) error
}

type SettingsListener interface {
	ApplySettings(settings setting.Settings)
}

var _ SettingsUpdateService = (*UpdateService)(nil)

type UpdateService struct {
	eventbus          eventbus.Bus
	settingsStore     SettingsStore
	settingsListeners []SettingsListener
}

func NewUpdateService(bus eventbus.Bus, settingsStore SettingsStore, listeners ...SettingsListener) *UpdateService {
	return &UpdateService{
		eventbus:          bus,
		settingsStore:     settingsStore,
		settingsListeners: listeners,
	}
}

func (service *UpdateService) ApplySettings(settings setting.Settings) error {
	err := service.settingsStore.Save(settings)
	if err != nil {
		return err
	}

	for _, listener := range service.settingsListeners {
		listener.ApplySettings(settings)
	}

	service.eventbus.Publish(event.SettingsUpdatedEvent)

	return nil
}
