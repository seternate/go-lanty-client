package main

import (
	"context"
	"fmt"
	_ "image/png"
	"os"
	"os/signal"
	"syscall"
	"time"

	eventbus "github.com/asaskevich/EventBus"
	"github.com/go-co-op/gocron/v2"
	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty-client/internal/backend"
	"github.com/seternate/go-lanty-client/internal/diskspace"
	"github.com/seternate/go-lanty-client/internal/game"
	"github.com/seternate/go-lanty-client/internal/platform/archive"
	"github.com/seternate/go-lanty-client/internal/platform/filesystem"
	"github.com/seternate/go-lanty-client/internal/platform/process"
	"github.com/seternate/go-lanty-client/internal/setting"
	"github.com/seternate/go-lanty-client/internal/ui/app"
	gameviewmodel "github.com/seternate/go-lanty-client/internal/ui/viewmodel/game"
	settingsviewmodel "github.com/seternate/go-lanty-client/internal/ui/viewmodel/settings"
	userviewmodel "github.com/seternate/go-lanty-client/internal/ui/viewmodel/user"
	"github.com/seternate/go-lanty-client/internal/user"
	"github.com/seternate/go-lanty/pkg/logging"
	"golang.org/x/sync/errgroup"
)

var AppName = "Lanty"
var SettingsPath = "settings.yaml"
var Version = "dev-build"
var DiskSpaceChangeDetectionThresholdBytes = uint64(100 * 1024 * 1024)

//go:generate go run github.com/tc-hib/go-winres@v0.3.1 make --in ./winres.json --arch amd64
//go:generate go run fyne.io/fyne/v2/cmd/fyne@v2.4.3 bundle -o ./bundled.go ./icon.png
//go:generate go run fyne.io/fyne/v2/cmd/fyne@v2.4.3 bundle -o ../pkg/widget/bundled.go ./game-icon.png

func main() {
	signalCtx, cancelSignalCtx := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancelSignalCtx()
	errgrp, errCtx := errgroup.WithContext(signalCtx)

	config, flagset, err := ParseConfig(os.Args...)
	if err != nil {
		if flagset != nil {
			fmt.Printf("%s\n\n", err)
			flagset.Usage()
			os.Exit(1)
		}
		log.Fatal().Err(err).Msg("error parsing flagset for configuration")
	}

	log.Logger = logging.Configure(logging.Config{
		LogLevel:           config.LogLevel,
		FileLoggingEnabled: true,
		Directory:          ".",
		Filename:           "lanty-client.log",
	})

	// apiClient, _ := api.New("http://localhost:8080")
	// client := apiadapter.NewAPIClient(apiClient, decoder.NewImageDecoder())

	settingStore, err := setting.NewFileStore(SettingsPath)
	if err != nil {
		log.Fatal().Err(err).Msg("error creating settings store")
	}

	bus := eventbus.New()
	gameMockAPIClient := backend.NewGameMockAPIClient()
	userMockAPIClient := backend.NewUserMockAPIClient()
	gamecatalogrepo := game.NewInMemoryCatalogRepository()
	usercatalogrepo := user.NewInMemoryCatalogRepository()
	gameinstallationrepo := game.NewInMemoryInstallationRepository()
	gameInstallationPathFinder := game.NewInstallationPathFinder()
	explorerOpener := filesystem.NewFileExplorerOpener()
	archiveExtractor := archive.NewZipExtractor()
	processLauncher := process.NewProcessLauncher()
	diskSpaceProvider := filesystem.NewDiskSpaceProvider()

	diskmonitor := diskspace.NewMonitor(bus, DiskSpaceChangeDetectionThresholdBytes, diskSpaceProvider)
	gameCatalogSyncer := game.NewCatalogSyncer(bus, gameMockAPIClient, gamecatalogrepo)
	gameInstallationDirectoryOpener := game.NewInstallationDirectoryOpener(gameinstallationrepo, explorerOpener)
	gameInstallationProgressEmitter := game.NewInstallationProgressEmitter(bus, gameMockAPIClient, archiveExtractor)
	gameInstallationReconciler := game.NewInstallationReconciler(bus, gameinstallationrepo, gamecatalogrepo, gameInstallationPathFinder)
	gameInstallationRunner := game.NewInstallationRunner(bus, gameinstallationrepo, gameMockAPIClient, archiveExtractor, settingStore)
	gameLaunchRunner := game.NewLaunchRunner(gameinstallationrepo, usercatalogrepo, processLauncher, gameMockAPIClient)
	userCatalogSyncer := user.NewCatalogSyncer(bus, userMockAPIClient, usercatalogrepo)

	scheduler, err := gocron.NewScheduler(gocron.WithGlobalJobOptions(gocron.WithStartAt(gocron.WithStartImmediately())))
	if err != nil {
		log.Fatal().Err(err).Msg("error creating scheduler")
	}
	scheduler.NewJob(gocron.DurationJob(5*time.Second), gocron.NewTask(gameCatalogSyncer.Sync, context.Background()))
	scheduler.NewJob(gocron.DurationJob(250*time.Millisecond), gocron.NewTask(gameInstallationProgressEmitter.Emit, context.Background()))
	scheduler.NewJob(gocron.DurationJob(1*time.Second), gocron.NewTask(
		func(ctx context.Context) error {
			settings, err := settingStore.Get()
			if err != nil {
				return err
			}
			return gameInstallationReconciler.Reconcile(ctx, settings.GameDirectory)
		},
		context.Background(),
	),
	)
	scheduler.NewJob(gocron.DurationJob(1*time.Second), gocron.NewTask(
		func(ctx context.Context) error {
			settings, err := settingStore.Get()
			if err != nil {
				return err
			}
			return diskmonitor.Detect(ctx, settings.GameDirectory)
		},
		context.Background(),
	),
	)
	scheduler.NewJob(gocron.DurationJob(5*time.Second), gocron.NewTask(userCatalogSyncer.Sync, context.Background()))

	appShell := app.NewAppShell(AppName, resourceIconPng, Version)

	gamescreenviewmodel := gameviewmodel.NewGameScreen(gameLaunchRunner, gameInstallationRunner, gameInstallationDirectoryOpener, appShell)
	userscreenviewmodel := userviewmodel.NewUserScreen()
	settingsscreenviewmodel := settingsviewmodel.NewSettingsScreen(settingStore, appShell)

	gameScreenSyncer := app.NewGameScreenSyncer(gamescreenviewmodel, gamecatalogrepo, gameMockAPIClient, gameinstallationrepo, gameMockAPIClient, archiveExtractor)
	userScreenSyncer := app.NewUserScreenSyncer(userscreenviewmodel, usercatalogrepo)

	bus.Subscribe(game.CatalogAddedEvent, gameScreenSyncer.OnCatalogItemAdded)
	bus.Subscribe(game.CatalogUpdatedEvent, gameScreenSyncer.OnCatalogItemUpdated)
	bus.Subscribe(game.CatalogRemovedEvent, gameScreenSyncer.OnCatalogItemRemoved)
	bus.Subscribe(game.InstallationStartedEvent, gameScreenSyncer.OnInstallationUpdated)
	bus.Subscribe(game.InstallationProgressedEvent, gameScreenSyncer.OnInstallationUpdated)
	bus.Subscribe(game.InstallationFinishedEvent, gameScreenSyncer.OnInstallationUpdated)
	bus.Subscribe(game.InstallationCancelledEvent, gameScreenSyncer.OnInstallationUpdated)
	bus.Subscribe(game.InstallationFailedEvent, gameScreenSyncer.OnInstallationUpdated)
	bus.Subscribe(game.InstallationDetectedEvent, gameScreenSyncer.OnInstallationUpdated)
	bus.Subscribe(game.InstallationRemovedEvent, gameScreenSyncer.OnInstallationUpdated)
	bus.Subscribe(diskspace.DiskSpaceChangedEvent, gameScreenSyncer.OnDiskSpaceChanged)

	bus.Subscribe(user.CatalogAddedEvent, userScreenSyncer.OnCatalogItemAdded)
	bus.Subscribe(user.CatalogUpdatedEvent, userScreenSyncer.OnCatalogItemUpdated)
	bus.Subscribe(user.CatalogRemovedEvent, userScreenSyncer.OnCatalogItemRemoved)

	appShell.Bootstrap(gamescreenviewmodel, userscreenviewmodel, settingsscreenviewmodel)

	// OLD

	// joinUserAdapter := adapter.NewJoinUserAdapter(application.Window, userCatalogRepo)
	// hostConfigAdapter := adapter.NewHostConfigAdapter(application.Window)
	// extensionsAdapter := adapter.NewExtensionsAdapter(application.Window)

	scheduler.Start()

	errgrp.Go(func() error {
		<-errCtx.Done()
		appShell.Quit()
		return nil
	})

	appShell.ShowAndRun()

	cancelSignalCtx()
	errgrp.Wait()
	scheduler.Shutdown()
	log.Debug().Msg("application shutdown")
}
