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
	gameevent "github.com/seternate/go-lanty-client/internal/application/game/event"
	gameappsrv "github.com/seternate/go-lanty-client/internal/application/game/service"
	settingsappsrv "github.com/seternate/go-lanty-client/internal/application/setting/service"
	systemevent "github.com/seternate/go-lanty-client/internal/application/system/event"
	systemappsrv "github.com/seternate/go-lanty-client/internal/application/system/service"
	userevent "github.com/seternate/go-lanty-client/internal/application/user/event"
	userappsrv "github.com/seternate/go-lanty-client/internal/application/user/service"
	apiadapter "github.com/seternate/go-lanty-client/internal/infrastructure/api"
	"github.com/seternate/go-lanty-client/internal/infrastructure/archive"
	"github.com/seternate/go-lanty-client/internal/infrastructure/detector"
	"github.com/seternate/go-lanty-client/internal/infrastructure/filesystem"
	gamerepository "github.com/seternate/go-lanty-client/internal/infrastructure/persistence/game"
	settingsrepo "github.com/seternate/go-lanty-client/internal/infrastructure/persistence/setting"
	userrepository "github.com/seternate/go-lanty-client/internal/infrastructure/persistence/user"
	"github.com/seternate/go-lanty-client/internal/infrastructure/process"
	"github.com/seternate/go-lanty-client/internal/ui/adapter"
	"github.com/seternate/go-lanty-client/internal/ui/app"
	gamecontroller "github.com/seternate/go-lanty-client/internal/ui/controller/game"
	settingcontroller "github.com/seternate/go-lanty-client/internal/ui/controller/setting"
	usercontroller "github.com/seternate/go-lanty-client/internal/ui/controller/user"
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

	fileSettingsStore, _ := settingsrepo.NewFileSettingsStore(SettingsPath)
	settings, _ := fileSettingsStore.GetSettings()
	gamedirectory := settings.GameDirectory

	bus := eventbus.New()
	gameMockAPIClient := apiadapter.NewGameMockAPIClient(gamedirectory)
	userMockAPIClient := apiadapter.NewUserMockAPIClient()
	gameCatalogRepo := gamerepository.NewInMemoryCatalogRepository()
	userCatalogRepo := userrepository.NewInMemoryCatalogRepository()
	gameInstallationRepo := gamerepository.NewInMemoryInstallationRepository()
	gameInstallationDetector := detector.NewGameInstallationDetector(gamedirectory)
	diskSpaceProvider := detector.NewDiskSpaceProvider(gamedirectory)
	archiveExtractor := archive.NewZipExtractor(gamedirectory)
	explorerOpener := filesystem.NewFileExplorerOpener()
	commandLauncher := process.NewCommandLauncher()

	gameCatalogRefreshService := gameappsrv.NewGameCatalogRefreshService(bus, gameMockAPIClient, gameCatalogRepo)
	gameInstallationProgressEmitter := gameappsrv.NewGameInstallationProgressEmitter(bus, gameMockAPIClient, archiveExtractor)
	gameInstallationService := gameappsrv.NewInstallationService(bus, gameInstallationRepo, gameCatalogRepo, gameMockAPIClient, archiveExtractor, gameInstallationDetector)
	gameDirectoryOpener := gameappsrv.NewGameDirectoryOpener(gamedirectory, gameInstallationRepo, explorerOpener)
	gameLauncherService := gameappsrv.NewLauncherService(gamedirectory, gameInstallationRepo, commandLauncher, gameMockAPIClient)
	diskSpaceChangeDetector := systemappsrv.NewDiskSpaceChangeDetector(bus, DiskSpaceChangeDetectionThresholdBytes, diskSpaceProvider)
	userCatalogRefreshService := userappsrv.NewUserCatalogRefreshService(bus, userMockAPIClient, userCatalogRepo)
	settingsUpdateService := settingsappsrv.NewUpdateService(
		bus,
		fileSettingsStore,
		gameInstallationDetector,
		diskSpaceProvider,
		archiveExtractor,
		gameDirectoryOpener,
		gameLauncherService,
	)

	scheduler, _ := gocron.NewScheduler(gocron.WithGlobalJobOptions(gocron.WithStartAt(gocron.WithStartImmediately())))
	scheduler.NewJob(gocron.DurationJob(5*time.Second), gocron.NewTask(gameCatalogRefreshService.Refresh))
	scheduler.NewJob(gocron.DurationJob(250*time.Millisecond), gocron.NewTask(gameInstallationProgressEmitter.EmitProgress))
	scheduler.NewJob(gocron.DurationJob(1*time.Second), gocron.NewTask(gameInstallationService.DetectInstalledGames))
	scheduler.NewJob(gocron.DurationJob(1*time.Second), gocron.NewTask(diskSpaceChangeDetector.DetectChange))
	scheduler.NewJob(gocron.DurationJob(5*time.Second), gocron.NewTask(userCatalogRefreshService.Refresh))

	application := app.NewApp(AppName, resourceIconPng, Version)

	joinUserAdapter := adapter.NewJoinUserAdapter(application.Window, userCatalogRepo)
	hostConfigAdapter := adapter.NewHostConfigAdapter(application.Window)

	gameController := gamecontroller.NewGameController(
		gameCatalogRepo,
		gameMockAPIClient,
		gameInstallationRepo,
		diskSpaceProvider,
		gameMockAPIClient,
		archiveExtractor,
		gameInstallationService,
		gameDirectoryOpener,
		joinUserAdapter,
		gameLauncherService,
		hostConfigAdapter,
		gameMockAPIClient,
	)
	userController := usercontroller.NewUserController(userCatalogRepo)
	settingsController, _ := settingcontroller.NewSettingsController(settingsUpdateService, fileSettingsStore)

	bus.Subscribe(gameevent.CatalogAddedEvent, gameController.OnCatalogItemAdded)
	bus.Subscribe(gameevent.CatalogUpdatedEvent, gameController.OnCatalogItemUpdated)
	bus.Subscribe(gameevent.CatalogRemovedEvent, gameController.OnCatalogItemRemoved)
	bus.Subscribe(gameevent.InstallationStartedEvent, gameController.OnInstallationStarted)
	bus.Subscribe(gameevent.InstallationProgressEvent, gameController.OnInstallationProgress)
	bus.Subscribe(gameevent.InstallationFinishedEvent, gameController.OnInstallationFinished)
	bus.Subscribe(gameevent.InstallationCanceledEvent, gameController.OnInstallationCanceled)
	bus.Subscribe(gameevent.InstallationFailedEvent, gameController.OnInstallationFailed)
	bus.Subscribe(gameevent.InstallationDetectedEvent, gameController.OnInstallationDetected)
	bus.Subscribe(gameevent.InstallationRemovedEvent, gameController.OnInstallationRemoved)
	bus.Subscribe(systemevent.DiskSpaceChangedEvent, gameController.OnDiskSpaceChanged)

	bus.Subscribe(userevent.CatalogAddedEvent, userController.OnUserAdded)
	bus.Subscribe(userevent.CatalogUpdatedEvent, userController.OnUserUpdated)
	bus.Subscribe(userevent.CatalogRemovedEvent, userController.OnUserRemoved)

	application.Bootstrap(gameController, userController, settingsController)

	scheduler.Start()

	errgrp.Go(func() error {
		<-errCtx.Done()
		application.Quit()
		return nil
	})

	application.ShowAndRun()

	cancelSignalCtx()

	errgrp.Wait()
	scheduler.Shutdown()
	log.Debug().Msg("application shutdown")
}
