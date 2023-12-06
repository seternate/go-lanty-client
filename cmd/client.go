package main

import (
	"flag"
	"net/url"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty-client/pkg/controller"
	"github.com/seternate/go-lanty-client/pkg/setting"
	"github.com/seternate/go-lanty-client/pkg/ui"
	"github.com/seternate/go-lanty/pkg/api"
)

func main() {
	parseFlags()

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

	controller := controller.NewController(settings, client).
		WithGameController().
		WithDownloadController()

	ui := ui.NewUI(controller)
	ui.ShowAndRun()
}

func parseFlags() {
	logLevel := flag.String("loglevel", "info", "Sets the log level of the application")
	flag.Parse()

	switch *logLevel {
	case "disable":
		zerolog.SetGlobalLevel(zerolog.Disabled)
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warning":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	}
}
