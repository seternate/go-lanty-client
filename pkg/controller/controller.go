package controller

import (
	"context"
	"net/url"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty-client/pkg/setting"
	"github.com/seternate/go-lanty/pkg/api"
	"github.com/seternate/go-lanty/pkg/handler"
)

type Controller struct {
	Settings *SettingsController
	Status   *StatusController
	Download *DownloadController
	Game     *GameController
	User     *UserController

	settings  *setting.Settings
	client    *api.Client
	ctx       context.Context
	cancelCtx context.CancelFunc
	waitgrp   *sync.WaitGroup
}

func NewController(ctx context.Context) *Controller {
	settings, err := setting.LoadSettings()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load settings")
	}
	log.Debug().Interface("settings", settings).Msg("loaded settings successfully")

	timeout, err := time.ParseDuration("0s")
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	url, err := url.Parse(settings.ServerURL)
	if err != nil {
		log.Fatal().Err(err).Str("url", settings.ServerURL).Msg("failed to parse server URL")
	}
	client := api.NewClient(url, timeout)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create API client")
	}
	log.Debug().Msg("created API client")

	context, cancelContext := context.WithCancel(ctx)

	return &Controller{
		settings:  settings,
		client:    client,
		ctx:       context,
		cancelCtx: cancelContext,
		waitgrp:   &sync.WaitGroup{},
	}
}

func (controller *Controller) Context() context.Context {
	return controller.ctx
}

func (controller *Controller) Quit() {
	controller.cancelCtx()
}

func (controller *Controller) WaitGroup() *sync.WaitGroup {
	return controller.waitgrp
}

func (controller *Controller) WithSettingsController() *Controller {
	controller.Settings = NewSettingsController(controller, controller.settings)
	return controller
}

func (controller *Controller) WithStatusController() *Controller {
	controller.Status = NewStatusController()
	return controller
}

func (controller *Controller) WithGameController() *Controller {
	controller.Game = NewGameController(controller, 10*time.Second)
	return controller
}

func (controller *Controller) WithDownloadController() *Controller {
	controller.Download = NewDownloadController(controller)
	return controller
}

func (controller *Controller) WithUserController() *Controller {
	controller.User = NewUserController(controller, 3*handler.UserStaleDuration/4)
	return controller
}
