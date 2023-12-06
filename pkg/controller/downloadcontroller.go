package controller

import (
	"errors"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty/pkg/game"
	"github.com/seternate/go-lanty/pkg/util"
)

type DownloadController struct {
	parent     *Controller
	downloads  []*Download
	subscriber []chan struct{}
	mutex      sync.RWMutex
}

func NewDownloadController(parent *Controller) (controller *DownloadController) {
	controller = &DownloadController{
		parent:     parent,
		downloads:  make([]*Download, 0),
		subscriber: make([]chan struct{}, 0, 50),
	}
	parent.WaitGroup().Add(1)
	go controller.run()
	return
}

func (controller *DownloadController) Download(game game.Game) {
	if controller.isDownloading(game) {
		log.Debug().Str("slug", game.Slug).Msg("game already downloading")
		return
	}
	controller.mutex.Lock()
	controller.downloads = append(controller.downloads, NewDownload(controller.parent, game))
	controller.mutex.Unlock()
	controller.notifySubcriber()
	log.Debug().Str("slug", game.Slug).Msg("added game to download queue")
}

func (controller *DownloadController) isDownloading(game game.Game) bool {
	download, err := controller.GetLatest(game)
	return err == nil && (!download.IsComplete() && !download.IsStopped())
}

func (controller *DownloadController) GetLatest(game game.Game) (*Download, error) {
	controller.mutex.Lock()
	downloads := controller.downloads
	controller.mutex.Unlock()
	for i := len(downloads) - 1; i >= 0; i-- {
		download := downloads[i]
		if download.Game().Equal(game) {
			return download, nil
		}
	}
	return nil, errors.New("game is not being downloaded")
}

func (controller *DownloadController) GetLastQueued() *Download {
	defer controller.mutex.RUnlock()
	controller.mutex.RLock()
	return controller.downloads[len(controller.downloads)-1]
}

func (controller *DownloadController) Subscribe(subscriber chan struct{}) {
	defer controller.mutex.Unlock()
	controller.mutex.Lock()
	controller.subscriber = append(controller.subscriber, subscriber)
}

func (controller *DownloadController) Unsubscribe(subscriber chan struct{}) {
	defer controller.mutex.Unlock()
	controller.mutex.Lock()
	index := slices.Index(controller.subscriber, subscriber)
	slices.Delete(controller.subscriber, index, index+1)
}

func (controller *DownloadController) notifySubcriber() {
	defer controller.mutex.RUnlock()
	controller.mutex.RLock()
	for _, subscriber := range controller.subscriber {
		util.ChannelWriteNonBlocking(subscriber, struct{}{})
	}
}

func (controller *DownloadController) run() {
	defer controller.parent.WaitGroup().Done()
	ticker := time.NewTicker(150 * time.Millisecond)
	for {
		select {
		case <-controller.parent.Context().Done():
			log.Trace().Err(controller.parent.Context().Err()).Msg("exiting downloadcontroller run()")
			return
		case <-ticker.C:
			controller.startQueuedDownloads()
		}
	}
}

func (controller *DownloadController) startQueuedDownloads() {
	controller.mutex.Lock()
	downloads := controller.downloads
	controller.mutex.Unlock()
	for _, download := range downloads {
		if !download.IsStarted() && !download.IsStopped() {
			if download.Retries() > 10 {
				download.Stop()
				controller.parent.Status.Error(fmt.Sprintf("Error starting download of game: %s", download.Game().Name), 8*time.Second)
				continue
			}
			err := download.Start(controller.parent.Context(), controller.parent.WaitGroup())
			if err != nil {
				log.Error().Err(err).Str("slug", download.Game().Slug).Uint64("retries", download.Retries()).Msg("failed to start download of game")
				continue
			}
			log.Debug().Str("slug", download.Game().Slug).Msg("started game download")
		}
	}
}
