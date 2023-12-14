package controller

import (
	"net/url"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty-client/pkg/setting"
	"github.com/seternate/go-lanty/pkg/api"
)

type Controller struct {
	Download *DownloadController
	Game     *GameController
	User     *UserController

	settings *setting.Settings
	client   *api.Client
}

func NewController() *Controller {
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

	return &Controller{
		settings: settings,
		client:   client,
	}
}

func (controller *Controller) WithGameController() *Controller {
	controller.Game = NewGameController(controller, 1*time.Minute)
	return controller
}

func (controller *Controller) WithDownloadController() *Controller {
	controller.Download = NewDownloadController(controller)
	return controller
}

func (controller *Controller) WithUserController() *Controller {
	controller.User = NewUserController(controller, 10*time.Second)
	return controller
}
