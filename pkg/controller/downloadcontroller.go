package controller

import (
	"errors"
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
		subscriber: make([]chan struct{}, 0),
	}
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
	return err == nil && !download.IsComplete()
}

func (controller *DownloadController) GetLatest(game game.Game) (*Download, error) {
	controller.mutex.Lock()
	downloads := controller.downloads
	controller.mutex.Unlock()
	for i := len(downloads) - 1; i >= 0; i-- {
		download := downloads[i]
		if download.Game() == game {
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
	log.Trace().Msg("new subscriber to downloadcontroller")
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
	log.Trace().Msg("notify subscriber of downloadcontroller")
	for _, subscriber := range controller.subscriber {
		util.ChannelWriteNonBlocking(subscriber, struct{}{})
	}
}

func (controller *DownloadController) run() {
	for {
		log.Trace().Msg("downloadcontroller update loop")
		controller.startQueuedDownloads()
		time.Sleep(150 * time.Millisecond)
	}
}

func (controller *DownloadController) startQueuedDownloads() {
	controller.mutex.Lock()
	downloads := controller.downloads
	controller.mutex.Unlock()
	for _, download := range downloads {
		if !download.IsStarted() {
			err := download.Start()
			if err != nil {
				log.Error().Err(err).Str("slug", download.Game().Slug).Msg("failed to start download of game")
				return
			}
			controller.notifySubcriber()
			log.Debug().Str("slug", download.Game().Slug).Msg("started game download")
		}
	}
}
