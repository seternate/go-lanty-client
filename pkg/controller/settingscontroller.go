package controller

import (
	"errors"
	"os"
	"slices"
	"sync"
	"time"

	"github.com/seternate/go-lanty-client/pkg/setting"
	"github.com/seternate/go-lanty/pkg/util"
)

type SettingsController struct {
	parent     *Controller
	settings   *setting.Settings
	subscriber []chan struct{}
	mutex      sync.RWMutex
}

func NewSettingsController(parent *Controller, settings *setting.Settings) (controller *SettingsController) {
	controller = &SettingsController{
		parent:     parent,
		settings:   settings,
		subscriber: make([]chan struct{}, 0, 50),
	}

	return
}

func (controller *SettingsController) SetServerURL(serverurl string) {
	err := controller.parent.client.SetBaseURL(serverurl)
	if err != nil {
		controller.parent.Status.Error("Invalid server URL: "+err.Error(), 3*time.Second)
	}
	controller.mutex.Lock()
	controller.settings.ServerURL = serverurl
	controller.mutex.Unlock()
	controller.notifySubcriber()
	controller.Save()
}

func (controller *SettingsController) SetGameDirectory(gamedirectory string) {
	info, err := os.Stat(gamedirectory)
	if errors.Is(err, os.ErrNotExist) {
		controller.parent.Status.Error("Invalid gamedirectory path: does not exist", 3*time.Second)
		return
	}
	if !info.IsDir() {
		controller.parent.Status.Error("Invalid gamedirectory path: not a folder", 3*time.Second)
		return
	}
	controller.mutex.Lock()
	controller.settings.GameDirectory = gamedirectory
	controller.mutex.Unlock()
	controller.notifySubcriber()
	controller.Save()
}

func (controller *SettingsController) SetUsername(username string) {
	controller.mutex.Lock()
	controller.settings.Username = username
	controller.mutex.Unlock()
	controller.notifySubcriber()
	controller.Save()
}

func (controller *SettingsController) Settings() setting.Settings {
	defer controller.mutex.RUnlock()
	controller.mutex.RLock()
	return *controller.settings
}

func (controller *SettingsController) Save() (err error) {
	err = controller.Settings().Save()
	if err != nil {
		controller.parent.Status.Error("Error saving settings: "+err.Error(), 3*time.Second)
	}
	return
}

func (controller *SettingsController) Subscribe(subscriber chan struct{}) {
	defer controller.mutex.Unlock()
	controller.mutex.Lock()
	controller.subscriber = append(controller.subscriber, subscriber)
}

func (controller *SettingsController) Unsubscribe(subscriber chan struct{}) {
	defer controller.mutex.Unlock()
	controller.mutex.Lock()
	index := slices.Index(controller.subscriber, subscriber)
	slices.Delete(controller.subscriber, index, index+1)
}

func (controller *SettingsController) notifySubcriber() {
	defer controller.mutex.RUnlock()
	controller.mutex.RLock()
	for _, subscriber := range controller.subscriber {
		util.ChannelWriteNonBlocking(subscriber, struct{}{})
	}
}
