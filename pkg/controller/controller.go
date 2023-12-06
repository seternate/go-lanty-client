package controller

import (
	"time"

	"github.com/seternate/go-lanty-client/pkg/setting"
	"github.com/seternate/go-lanty/pkg/api"
)

type Controller struct {
	Download *DownloadController
	Game     *GameController

	settings *setting.Settings
	client   *api.Client
}

func NewController(settings *setting.Settings, client *api.Client) *Controller {
	return &Controller{
		settings: settings,
		client:   client,
	}
}

func (controller *Controller) WithGameController() *Controller {
	controller.Game = NewGameController(controller, 10*time.Second)
	return controller
}

func (controller *Controller) WithDownloadController() *Controller {
	controller.Download = NewDownloadController(controller)
	return controller
}
