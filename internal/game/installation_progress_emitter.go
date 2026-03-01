package game

import (
	"context"
	"fmt"
	"sync"

	eventbus "github.com/asaskevich/EventBus"
)

type InstallationProgressProvider interface {
	GetAllInstallationProgresses(ctx context.Context) (map[string]InstallationProgress, error)
}
type InstallationProgressEmitter struct {
	bus         eventbus.Bus
	downloads   InstallationProgressProvider
	extractions InstallationProgressProvider

	snapshot map[string]InstallationProgress
	mu       sync.RWMutex
}

func NewGameInstallationProgressEmitter(bus eventbus.Bus, downloads InstallationProgressProvider, extractions InstallationProgressProvider) *InstallationProgressEmitter {
	return &InstallationProgressEmitter{
		bus:         bus,
		downloads:   downloads,
		extractions: extractions,
		snapshot:    make(map[string]InstallationProgress),
	}
}

func (emitter *InstallationProgressEmitter) Emit(ctx context.Context) error {
	downloads, err := emitter.downloads.GetAllInstallationProgresses(ctx)
	if err != nil {
		return fmt.Errorf("failed to get all download progresses: %w", err)
	}
	extractions, err := emitter.extractions.GetAllInstallationProgresses(ctx)
	if err != nil {
		return fmt.Errorf("failed to get all extraction progresses: %w", err)
	}

	slugs := make(map[string]bool)
	for slug := range downloads {
		slugs[slug] = true
	}
	for slug := range extractions {
		slugs[slug] = true
	}

	for slug := range slugs {
		latestProgress := InstallationProgress{}

		download, downloadFound := downloads[slug]
		extraction, extractionFound := extractions[slug]
		if (downloadFound && extractionFound) && (extraction.StartedAt.Equal(download.StartedAt) || extraction.StartedAt.After(download.StartedAt)) {
			latestProgress = extraction
		} else if (downloadFound && extractionFound) && download.StartedAt.After(extraction.StartedAt) {
			latestProgress = download
		} else if downloadFound {
			latestProgress = download
		} else if extractionFound {
			latestProgress = extraction
		}

		if latestProgress.HasFinished() {
			delete(emitter.snapshot, slug)
			continue
		}

		emitter.mu.RLock()
		lastProgress, found := emitter.snapshot[slug]
		emitter.mu.RUnlock()
		if !found || latestProgress.StartedAt.After(lastProgress.StartedAt) || (lastProgress.Completed < latestProgress.Completed) {
			emitter.mu.Lock()
			emitter.snapshot[slug] = latestProgress
			emitter.mu.Unlock()
			emitter.bus.Publish(InstallationProgressedEvent, InstallationEvent{Slug: slug})
		}
	}

	return nil
}
