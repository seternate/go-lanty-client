package service

import (
	"context"
	"errors"
	"sync"

	eventbus "github.com/asaskevich/EventBus"
	"github.com/seternate/go-lanty-client/internal/application/game/event"
	"github.com/seternate/go-lanty-client/internal/domain/game"
)

type GameInstallationService interface {
	StartInstallation(ctx context.Context, slug string) error
	StopInstallation(slug string) error
	DetectInstalledGames() error
}

type BlobDownloader interface {
	Download(ctx context.Context, slug string) (filePath string, err error)
}

type BlobExtractor interface {
	Extract(ctx context.Context, slug string, filePath string, removeArchive bool) (extractionPath string, err error)
}

type InstallationDetector interface {
	DetectInstalledGame(slug string, executableRelativePath string) (detected bool, installationPath string, err error)
}

var _ GameInstallationService = (*InstallationService)(nil)

type InstallationService struct {
	bus                  eventbus.Bus
	installationRepo     game.GameInstallationRepository
	catalogRepo          game.GameCatalogRepository
	blobDownloader       BlobDownloader
	blobExtractor        BlobExtractor
	installationDetector InstallationDetector

	mu            sync.Mutex
	installations map[string]context.CancelFunc
}

func NewInstallationService(bus eventbus.Bus, installationRepo game.GameInstallationRepository, catalogRepo game.GameCatalogRepository, blobDownloader BlobDownloader, blobExtractor BlobExtractor, installationDetector InstallationDetector) *InstallationService {
	return &InstallationService{
		bus:                  bus,
		installationRepo:     installationRepo,
		catalogRepo:          catalogRepo,
		blobDownloader:       blobDownloader,
		blobExtractor:        blobExtractor,
		installationDetector: installationDetector,
		installations:        make(map[string]context.CancelFunc),
	}
}

func (service *InstallationService) StartInstallation(ctx context.Context, slug string) error {
	installation, err := service.installationRepo.GetBySlug(slug)
	if err != nil {
		return err
	}

	if installation.IsInstalling() {
		return nil
	}

	ctx, cancel := context.WithCancel(ctx)
	service.mu.Lock()
	service.installations[installation.Slug] = cancel
	service.mu.Unlock()

	go service.runInstallation(ctx, installation)

	return nil
}

func (service *InstallationService) runInstallation(ctx context.Context, installation game.GameInstallation) {
	defer func() {
		service.mu.Lock()
		delete(service.installations, installation.Slug)
		service.mu.Unlock()
	}()

	installation.MarkInstallationStarted()
	service.installationRepo.Store(installation)
	service.bus.Publish(event.InstallationStartedEvent, event.InstallationEvent{Slug: installation.Slug})

	filePath, err := service.blobDownloader.Download(ctx, installation.Slug)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		service.installationFailed(installation)
		return
	}

	installation.MarkExtractionStarted()
	service.installationRepo.Store(installation)
	service.bus.Publish(event.InstallationProgressEvent, event.InstallationEvent{Slug: installation.Slug})

	extractionPath, err := service.blobExtractor.Extract(ctx, installation.Slug, filePath, true)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		service.installationFailed(installation)
		return
	}

	installation.MarkInstallationFinished(extractionPath)
	service.installationRepo.Store(installation)
	service.bus.Publish(event.InstallationFinishedEvent, event.InstallationEvent{Slug: installation.Slug})
}

func (service *InstallationService) installationFailed(installation game.GameInstallation) {
	installation.MarkInstallationFailed()
	service.installationRepo.Store(installation)
	service.bus.Publish(event.InstallationFailedEvent, event.InstallationEvent{Slug: installation.Slug})
}

func (service *InstallationService) StopInstallation(slug string) error {
	installation, err := service.installationRepo.GetBySlug(slug)
	if err != nil {
		return err
	}

	if !installation.IsInstalling() {
		return nil
	}

	service.mu.Lock()
	cancelInstallation, found := service.installations[slug]
	if found {
		delete(service.installations, slug)
	}
	service.mu.Unlock()

	if found {
		cancelInstallation()
	}

	installation.MarkInstallationCanceled()
	service.installationRepo.Store(installation)
	service.bus.Publish(event.InstallationCanceledEvent, event.InstallationEvent{Slug: slug})

	return nil
}

func (service *InstallationService) DetectInstalledGames() error {
	catalog, err := service.catalogRepo.GetAll()
	if err != nil {
		return err
	}

	installations, err := service.installationRepo.GetAll()
	if err != nil {
		return err
	}

	instaBySlug := make(map[string]game.GameInstallation)
	for _, installation := range installations {
		instaBySlug[installation.Slug] = installation
	}

	catalogMap := make(map[string]game.GameCatalogItem)
	for _, catalogItem := range catalog {
		catalogMap[catalogItem.Slug] = catalogItem
	}

	for _, catalogItem := range catalog {
		if _, found := instaBySlug[catalogItem.Slug]; !found {
			newInstallation := game.NewGameInstallation(catalogItem.Slug)
			service.installationRepo.Store(*newInstallation)
			instaBySlug[catalogItem.Slug] = *newInstallation
		}
	}

	for _, installation := range installations {
		if _, found := catalogMap[installation.Slug]; !found {
			service.installationRepo.Remove(installation.Slug)
			delete(instaBySlug, installation.Slug)
		}
	}

	installations, err = service.installationRepo.GetAll()
	if err != nil {
		return err
	}

	for slug, installation := range instaBySlug {
		detected, installationPath, _ := service.installationDetector.DetectInstalledGame(slug, catalogMap[slug].ExecutableRelativePath)
		oldInstallation := installation

		if detected && !installation.InstallationAvailable() {
			installation.MarkInstallationDetected(installationPath)
		} else if !detected && installation.InstallationAvailable() {
			installation.MarkInstallationNotFound()
		}

		service.installationRepo.Store(installation)

		if detected && !oldInstallation.InstallationAvailable() {
			service.bus.Publish(event.InstallationDetectedEvent, event.InstallationEvent{Slug: slug})
		} else if !detected && oldInstallation.InstallationAvailable() {
			service.bus.Publish(event.InstallationRemovedEvent, event.InstallationEvent{Slug: slug})
		}

	}

	return nil
}
