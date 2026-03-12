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
	"github.com/seternate/go-lanty-client/internal/setting"
	"github.com/seternate/go-lanty-client/internal/ui/app"
	gameview "github.com/seternate/go-lanty-client/internal/ui/view/game"
	gameviewmodel "github.com/seternate/go-lanty-client/internal/ui/viewmodel/game"
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
	//explorerOpener := filesystem.NewFileExplorerOpener()
	archiveExtractor := archive.NewZipExtractor()
	//processLauncher := process.NewProcessLauncher()
	diskSpaceProvider := filesystem.NewDiskSpaceProvider()
	// OLD gameExtensionRepo := gamerepository.NewInMemoryExtensionRepository()
	// OLD gameExtensionMock := apiadapter.NewGameExtensionMock()

	diskmonitor := diskspace.NewMonitor(bus, DiskSpaceChangeDetectionThresholdBytes, diskSpaceProvider)
	gameCatalogSyncer := game.NewCatalogSyncer(bus, gameMockAPIClient, gamecatalogrepo)
	//gameInstallationDirectoryOpener := game.NewInstallationDirectoryOpener(gameinstallationrepo, explorerOpener)
	gameInstallationProgressEmitter := game.NewInstallationProgressEmitter(bus, gameMockAPIClient, archiveExtractor)
	gameInstallationReconciler := game.NewInstallationReconciler(bus, gameinstallationrepo, gamecatalogrepo, gameInstallationPathFinder)
	//gameInstallationRunner := game.NewInstallationRunner(bus, gameinstallationrepo, gameMockAPIClient, archiveExtractor)
	//gameLaunchRunner := game.NewLaunchRunner(gameinstallationrepo, processLauncher, gameMockAPIClient)
	userCatalogSyncer := user.NewCatalogSyncer(bus, userMockAPIClient, usercatalogrepo)
	// OLD extensionService := gameappsrv.NewExtensionService(bus, gameExtensionRepo, gameInstallationRepo, gamedirectory, gameExtensionMock, gameExtensionMock, gameExtensionMock)

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

	// TEST CODE
	games, _ := gameMockAPIClient.FetchCatalog(context.Background())
	game := games[0]
	gametileviewmodel := gameviewmodel.NewGameTile(game.Slug, nil)
	icon, _ := gameMockAPIClient.GetIcon(game.Slug)
	gametileviewmodel.Icon.Set(icon)
	gametileviewmodel.Name.Set(game.Name)
	gametileviewmodel.InstalledFileSize.Set(fmt.Sprintf("%d GB", game.InstalledFileSize))
	gametileviewmodel.IsProgressing.Set(true)
	gametileviewmodel.Progress.Set(0.5)
	gametileviewmodel.ProgressFrontText.Set("Downloading... 100 MB / 200 MB")
	gametileviewmodel.ProgressEndText.Set("100 MB/s")
	gametileviewmodel.ExtensionsVisible.Set(true)
	gametileviewmodel.HasNewExtensions.Set(true)
	gametileview := gameview.NewGameTile(gametileviewmodel)
	//

	appShell.SetGameTile(gametileview)

	// OLD

	// joinUserAdapter := adapter.NewJoinUserAdapter(application.Window, userCatalogRepo)
	// hostConfigAdapter := adapter.NewHostConfigAdapter(application.Window)
	// extensionsAdapter := adapter.NewExtensionsAdapter(application.Window)

	// gameController := gamecontroller.NewGameController(
	// 	gameCatalogRepo,
	// 	gameMockAPIClient,
	// 	gameInstallationRepo,
	// 	diskSpaceProvider,
	// 	gameMockAPIClient,
	// 	archiveExtractor,
	// 	gameInstallationService,
	// 	gameDirectoryOpener,
	// 	joinUserAdapter,
	// 	gameLauncherService,
	// 	hostConfigAdapter,
	// 	gameMockAPIClient,
	// 	extensionService,
	// 	extensionsAdapter,
	// 	gamedirectory,
	// )
	// userController := usercontroller.NewUserController(userCatalogRepo)
	// settingsController, _ := settingcontroller.NewSettingsController(settingsUpdateService, fileSettingsStore)

	// bus.Subscribe(gameevent.CatalogAddedEvent, gameController.OnCatalogItemAdded)
	// bus.Subscribe(gameevent.CatalogUpdatedEvent, gameController.OnCatalogItemUpdated)
	// bus.Subscribe(gameevent.CatalogRemovedEvent, gameController.OnCatalogItemRemoved)
	// bus.Subscribe(gameevent.InstallationStartedEvent, gameController.OnInstallationStarted)
	// bus.Subscribe(gameevent.InstallationProgressEvent, gameController.OnInstallationProgress)
	// bus.Subscribe(gameevent.InstallationFinishedEvent, gameController.OnInstallationFinished)
	// bus.Subscribe(gameevent.InstallationCanceledEvent, gameController.OnInstallationCanceled)
	// bus.Subscribe(gameevent.InstallationFailedEvent, gameController.OnInstallationFailed)
	// bus.Subscribe(gameevent.InstallationDetectedEvent, gameController.OnInstallationDetected)
	// bus.Subscribe(gameevent.InstallationRemovedEvent, gameController.OnInstallationRemoved)
	// bus.Subscribe(gameevent.ExtensionListRefreshedEvent, gameController.OnExtensionListRefreshed)
	// bus.Subscribe(gameevent.ExtensionDownloadedEvent, gameController.OnExtensionDownloaded)
	// bus.Subscribe(systemevent.DiskSpaceChangedEvent, gameController.OnDiskSpaceChanged)

	// bus.Subscribe(userevent.CatalogAddedEvent, userController.OnUserAdded)
	// bus.Subscribe(userevent.CatalogUpdatedEvent, userController.OnUserUpdated)
	// bus.Subscribe(userevent.CatalogRemovedEvent, userController.OnUserRemoved)

	// application.Bootstrap(gameController, userController, settingsController)

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
