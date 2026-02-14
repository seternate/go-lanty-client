package service

import (
	"sync"

	eventbus "github.com/asaskevich/EventBus"
	"github.com/seternate/go-lanty-client/internal/application/game"
	"github.com/seternate/go-lanty-client/internal/application/game/event"
)

type GameInstallationProgressFetcher interface {
	GetAllProgress() (map[string]game.InstallationProgress, error)
	GetProgress(slug string) (game.InstallationProgress, error)
}

type GameInstallationProgressEmitter interface {
	EmitProgress()
}

var _ GameInstallationProgressEmitter = (*InstallationProgressEmitter)(nil)

type InstallationProgressEmitter struct {
	bus                       eventbus.Bus
	downloadProgressFetcher   GameInstallationProgressFetcher
	extractingProgressFetcher GameInstallationProgressFetcher

	snapshot map[string]game.InstallationProgress
	mu       sync.RWMutex
}

func NewGameInstallationProgressEmitter(bus eventbus.Bus, downloadProgressFetcher GameInstallationProgressFetcher, extractingProgressFetcher GameInstallationProgressFetcher) *InstallationProgressEmitter {
	return &InstallationProgressEmitter{
		bus:                       bus,
		downloadProgressFetcher:   downloadProgressFetcher,
		extractingProgressFetcher: extractingProgressFetcher,
		snapshot:                  make(map[string]game.InstallationProgress),
	}
}

func (emitter *InstallationProgressEmitter) EmitProgress() {
	downloadProgress, err := emitter.downloadProgressFetcher.GetAllProgress()
	if err != nil {
		return
	}
	extractingProgress, err := emitter.extractingProgressFetcher.GetAllProgress()
	if err != nil {
		return
	}

	slugs := make(map[string]bool)
	for slug := range downloadProgress {
		slugs[slug] = true
	}
	for slug := range extractingProgress {
		slugs[slug] = true
	}

	for slug := range slugs {
		latestProgress := game.InstallationProgress{}

		downloadProgress, downloadFound := downloadProgress[slug]
		extractingProgress, extractingFound := extractingProgress[slug]
		if (downloadFound && extractingFound) && (downloadProgress.StartedAt.Equal(extractingProgress.StartedAt) || downloadProgress.StartedAt.Before(extractingProgress.StartedAt)) {
			latestProgress = extractingProgress
		} else if (downloadFound && extractingFound) && downloadProgress.StartedAt.After(extractingProgress.StartedAt) {
			latestProgress = downloadProgress
		} else if downloadFound {
			latestProgress = downloadProgress
		} else if extractingFound {
			latestProgress = extractingProgress
		}

		if !latestProgress.EndedAt.IsZero() {
			delete(emitter.snapshot, slug)
			continue
		}

		emitter.mu.RLock()
		lastProgress, found := emitter.snapshot[slug]
		emitter.mu.RUnlock()
		if !found || lastProgress.StartedAt.Before(latestProgress.StartedAt) {
			emitter.mu.Lock()
			emitter.snapshot[slug] = latestProgress
			emitter.mu.Unlock()
			emitter.bus.Publish(event.InstallationProgressEvent, event.InstallationEvent{Slug: slug})
			return
		}

		if lastProgress.Completed < latestProgress.Completed {
			emitter.mu.Lock()
			emitter.snapshot[slug] = latestProgress
			emitter.mu.Unlock()
			emitter.bus.Publish(event.InstallationProgressEvent, event.InstallationEvent{Slug: slug})
		}
	}
}
