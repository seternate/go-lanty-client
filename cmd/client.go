package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty-client/pkg/controller"
	"github.com/seternate/go-lanty-client/pkg/setting"
	"github.com/seternate/go-lanty-client/pkg/widget"
	"github.com/seternate/go-lanty/pkg/logging"
	"github.com/seternate/go-lanty/pkg/network"
	"golang.design/x/clipboard"
)

//go:generate go run github.com/tc-hib/go-winres@v0.3.1 make --in ./winres.json --arch amd64
//go:generate go run fyne.io/fyne/v2/cmd/fyne@v2.4.3 bundle -o ./bundled.go ./icon.png

func main() {
	signalCtx, cancelSignalCtx := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancelSignalCtx()

	logconfig := logging.Config{
		ConsoleLoggingEnabled: true,
		FileLoggingEnabled:    true,
		Filename:              "lanty.log",
		Directory:             "log",
	}
	parseFlags(&logconfig)
	log.Logger = logging.Configure(logconfig)

	err := clipboard.Init()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init clipboard package")
	}

	controller := controller.NewController(signalCtx).
		WithSettingsController().
		WithStatusController().
		WithGameController().
		WithDownloadController().
		WithUserController().
		WithChatController()

	app := app.New()
	window := app.NewWindow(getApplicationTitle())
	lanty := widget.NewLanty(controller, window)
	window.SetContent(lanty)
	window.SetPadded(false)
	window.Resize(fyne.NewSize(1024, 600))
	window.SetIcon(resourceIconPng)

	window.ShowAndRun()

	log.Debug().Msg("quit application")
	controller.Quit()
	controller.WaitGroup().Wait()
	log.Debug().Msg("application stopped")
}

func getApplicationTitle() string {
	ip, err := network.GetOutboundIP()
	if err != nil {
		return setting.APPLICATION_NAME
	}
	return fmt.Sprintf("%s - %s", setting.APPLICATION_NAME, ip.String())
}

func parseFlags(config *logging.Config) {
	flag.StringVar(&config.LogLevel, "loglevel", "info", "Sets the log level")
	flag.IntVar(&config.MaxBackups, "logbackups", 0, "Sets the number of old logs to remain")
	flag.IntVar(&config.MaxSize, "logfilesize", 10, "Sets the size of the logs before rotating to new file")
	flag.IntVar(&config.MaxAge, "logage", 0, "Sets the maximum number of days to retain old logs")
	flag.Parse()
}
