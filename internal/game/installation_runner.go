package game

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"sync"

	eventbus "github.com/asaskevich/EventBus"
	"github.com/seternate/go-lanty-client/internal/setting"
)

type ArchiveDownloader interface {
	Download(ctx context.Context, slug string, baseDir string) (filePath string, err error)
}

type ArchiveExtractor interface {
	Extract(ctx context.Context, slug string, archivePath string, extractionPath string, removeArchive bool) (err error)
}

type InstallationRunner struct {
	bus        eventbus.Bus
	repo       InstallationRepository
	downloader ArchiveDownloader
	extractor  ArchiveExtractor

	mu               sync.Mutex
	installations    map[string]context.CancelFunc
	settingsProvider setting.SettingsProvider
}

func NewInstallationRunner(bus eventbus.Bus, repo InstallationRepository, downloader ArchiveDownloader, extractor ArchiveExtractor, settingsProvider setting.SettingsProvider) *InstallationRunner {
	return &InstallationRunner{
		bus:              bus,
		repo:             repo,
		downloader:       downloader,
		extractor:        extractor,
		settingsProvider: settingsProvider,
		installations:    make(map[string]context.CancelFunc),
	}
}

func (runner *InstallationRunner) StartInstallation(ctx context.Context, slug string) error {
	installation, err := runner.repo.GetBySlug(ctx, slug)
	if err != nil {
		return fmt.Errorf("failed to get installation: %w", err)
	}

	if installation.IsInstalling() {
		return fmt.Errorf("installation already in progress")
	}

	ctx, cancel := context.WithCancel(ctx)
	runner.mu.Lock()
	runner.installations[installation.Slug] = cancel
	runner.mu.Unlock()

	settings, err := runner.settingsProvider.Get()
	if err != nil {
		return fmt.Errorf("failed to get settings: %w", err)
	}

	err = runner.runInstallation(ctx, settings.GameDirectory, installation)
	runner.mu.Lock()
	delete(runner.installations, installation.Slug)
	runner.mu.Unlock()
	if err != nil {
		return fmt.Errorf("failed to run installation: %w", err)
	}

	return nil
}

func (runner *InstallationRunner) runInstallation(ctx context.Context, baseDir string, installation Installation) error {
	err := installation.MarkAsDownloading()
	if err != nil {
		return fmt.Errorf("failed to mark as downloading: %w", err)
	}
	err = runner.repo.Store(ctx, installation)
	if err != nil {
		return fmt.Errorf("failed to store installation: %w", err)
	}
	runner.bus.Publish(InstallationStartedEvent, InstallationEvent{Slug: installation.Slug})

	downloadedFilePath, err := runner.downloader.Download(ctx, installation.Slug, baseDir)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return fmt.Errorf("download canceled: %w", err)
		}

		err = installation.MarkAsFailed()
		if err != nil {
			return fmt.Errorf("failed to mark as failed: %w", err)
		}
		err = runner.repo.Store(ctx, installation)
		if err != nil {
			return fmt.Errorf("failed to store installation: %w", err)
		}
		runner.bus.Publish(InstallationFailedEvent, InstallationEvent{Slug: installation.Slug})

		return fmt.Errorf("failed to download blob: %w", err)
	}

	err = installation.MarkAsExtracting()
	if err != nil {
		return fmt.Errorf("failed to mark as extracting: %w", err)
	}
	err = runner.repo.Store(ctx, installation)
	if err != nil {
		return fmt.Errorf("failed to store installation: %w", err)
	}
	runner.bus.Publish(InstallationProgressedEvent, InstallationEvent{Slug: installation.Slug})

	extractionPath := filepath.Join(baseDir, installation.Slug)
	err = runner.extractor.Extract(ctx, installation.Slug, downloadedFilePath, extractionPath, true)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			err = installation.MarkAsCancelled()
			if err != nil {
				return fmt.Errorf("failed to mark as cancelled: %w", err)
			}
			err = runner.repo.Store(ctx, installation)
			if err != nil {
				return fmt.Errorf("failed to store installation: %w", err)
			}
			runner.bus.Publish(InstallationCancelledEvent, InstallationEvent{Slug: installation.Slug})

			return fmt.Errorf("extraction canceled: %w", err)
		}

		err = installation.MarkAsFailed()
		if err != nil {
			return fmt.Errorf("failed to mark as failed: %w", err)
		}
		err = runner.repo.Store(ctx, installation)
		if err != nil {
			return fmt.Errorf("failed to store installation: %w", err)
		}
		runner.bus.Publish(InstallationFailedEvent, InstallationEvent{Slug: installation.Slug})

		return fmt.Errorf("failed to extract blob: %w", err)
	}

	err = installation.MarkAsCompleted(extractionPath)
	if err != nil {
		return fmt.Errorf("failed to mark as completed: %w", err)
	}
	err = runner.repo.Store(ctx, installation)
	if err != nil {
		return fmt.Errorf("failed to store installation: %w", err)
	}
	runner.bus.Publish(InstallationFinishedEvent, InstallationEvent{Slug: installation.Slug})

	return nil
}

func (runner *InstallationRunner) CancelInstallation(ctx context.Context, slug string) error {
	installation, err := runner.repo.GetBySlug(ctx, slug)
	if err != nil {
		return fmt.Errorf("failed to get installation: %w", err)
	}

	if installation.IsIdle() {
		return fmt.Errorf("installation is not in progress")
	}

	runner.mu.Lock()
	cancelInstallation, found := runner.installations[slug]
	if found {
		delete(runner.installations, slug)
		cancelInstallation()
	}
	runner.mu.Unlock()

	err = installation.MarkAsCancelled()
	if err != nil {
		return fmt.Errorf("failed to mark as cancelled: %w", err)
	}
	err = runner.repo.Store(ctx, installation)
	if err != nil {
		return fmt.Errorf("failed to store installation: %w", err)
	}
	runner.bus.Publish(InstallationCancelledEvent, InstallationEvent{Slug: slug})

	return nil
}
