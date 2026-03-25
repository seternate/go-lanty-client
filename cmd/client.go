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
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty-client/internal/backend"
	"github.com/seternate/go-lanty-client/internal/diskspace"
	"github.com/seternate/go-lanty-client/internal/game"
	"github.com/seternate/go-lanty-client/internal/platform/archive"
	"github.com/seternate/go-lanty-client/internal/platform/filesystem"
	"github.com/seternate/go-lanty-client/internal/paths"
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

	exeDir, exeDirErr := paths.ExecutableDir()
	if exeDirErr != nil {
		log.Warn().Err(exeDirErr).Msg("could not resolve executable directory - using current working directory for log file")
	}
	appLogger := logging.Configure(logging.Config{
		LogLevel:           config.LogLevel,
		FileLoggingEnabled: true,
		Directory:          exeDir,
		Filename:           "lanty-client.log",
	})
	log.Logger = appLogger

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

	schedLogger := appLogger.With().Str("component", "gocron").Logger()
	scheduler, err := gocron.NewScheduler(gocron.WithGlobalJobOptions(
		gocron.WithStartAt(gocron.WithStartImmediately()),
		gocron.WithEventListeners(
			gocron.AfterJobRunsWithError(func(jobID uuid.UUID, jobName string, err error) {
				schedLogger.Error().Err(err).Str("job_id", jobID.String()).Str("job_name", jobName).Msg("scheduled job failed")
			}),
			gocron.AfterJobRunsWithPanic(func(jobID uuid.UUID, jobName string, recoverData any) {
				schedLogger.Error().Str("job_id", jobID.String()).Str("job_name", jobName).Interface("recover", recoverData).Msg("scheduled job panicked")
			}),
		),
	))
	if err != nil {
		log.Fatal().Err(err).Msg("error creating scheduler")
	}
	scheduler.NewJob(gocron.DurationJob(5*time.Second), gocron.NewTask(gameCatalogSyncer.Sync, context.Background()), gocron.WithName("game_catalog_sync"))
	scheduler.NewJob(gocron.DurationJob(250*time.Millisecond), gocron.NewTask(gameInstallationProgressEmitter.Emit, context.Background()), gocron.WithName("game_installation_progress_emit"))
	scheduler.NewJob(gocron.DurationJob(3*time.Second), gocron.NewTask(
		func(ctx context.Context) error {
			settings, err := settingStore.Get()
			if err != nil {
				return err
			}
			return gameInstallationReconciler.Reconcile(ctx, settings.GameDirectory)
		},
		context.Background(),
	),
		gocron.WithName("game_installation_reconcile"),
	)
	scheduler.NewJob(gocron.DurationJob(5*time.Second), gocron.NewTask(
		func(ctx context.Context) error {
			settings, err := settingStore.Get()
			if err != nil {
				return err
			}
			return diskmonitor.Detect(ctx, settings.GameDirectory)
		},
		context.Background(),
	),
		gocron.WithName("disk_space_detect"),
	)
	scheduler.NewJob(gocron.DurationJob(5*time.Second), gocron.NewTask(userCatalogSyncer.Sync, context.Background()), gocron.WithName("user_catalog_sync"))

	appShellLogger := appLogger.With().Str("component", "AppShell").Logger()
	appShell := app.NewAppShell(appShellLogger, AppName, resourceIconPng, Version)

	gamescreenviewmodel := gameviewmodel.NewGameScreen(gameLaunchRunner, gameInstallationRunner, gameInstallationDirectoryOpener, appShell)
	userscreenviewmodel := userviewmodel.NewUserScreen()
	settingsScreenLogger := appLogger.With().Str("component", "SettingsScreen").Logger()
	settingsscreenviewmodel, err := settingsviewmodel.NewSettingsScreen(settingsScreenLogger, settingStore, appShell)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize settings screen")
	}

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
