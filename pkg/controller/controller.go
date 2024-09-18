package controller

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty-client/pkg/setting"
	"github.com/seternate/go-lanty/pkg/api"
	"github.com/seternate/go-lanty/pkg/handler"
)

type Controller struct {
	Settings   *SettingsController
	Status     *StatusController
	Download   *DownloadController
	Game       *GameController
	User       *UserController
	Chat       *ChatController
	Connection *ConnectionController

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

	client, err := api.NewClient(settings.ServerURL, timeout)
	if err != nil {
		log.Error().Err(err).Msg("error creating API client")
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
	controller.Game = NewGameController(controller, 2*time.Second)
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

func (controller *Controller) WithChatController() *Controller {
	controller.Chat = NewChatController(controller)
	return controller
}

func (controller *Controller) WithConnectionController() *Controller {
	controller.Connection = NewConnectionController(controller, 1*time.Second)
	return controller
}
