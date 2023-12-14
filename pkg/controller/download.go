package controller

import (
	"errors"
	"os"
	"path"
	"slices"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty/pkg/filesystem"
	"github.com/seternate/go-lanty/pkg/game"
	"github.com/seternate/go-lanty/pkg/network"
	"github.com/seternate/go-lanty/pkg/util"
)

type Download struct {
	controller         *Controller
	game               game.Game
	download           *network.Download
	unzip              *filesystem.Unzip
	subscriber         []chan struct{}
	subscriberprogress []chan float64
	started            bool
	running            bool
	downloading        bool
	err                error
	mutex              sync.RWMutex
}

func NewDownload(controller *Controller, game game.Game) (download *Download) {
	download = &Download{
		controller:  controller,
		game:        game,
		subscriber:  make([]chan struct{}, 0),
		started:     false,
		running:     false,
		downloading: false,
	}
	return
}

func (controller *Download) Start() (err error) {
	if controller.IsStarted() {
		log.Debug().Str("slug", controller.Game().Slug).Msg("download already started")
		return errors.New("download already started")
	}
	download, err := controller.controller.client.Game.Download(controller.game, controller.controller.settings.GameDirectory)
	if err != nil {
		log.Error().Err(err).Str("slug", controller.Game().Slug).Msg("error starting game download from server")
		return
	}
	controller.mutex.Lock()
	controller.download = download
	controller.started = true
	controller.running = true
	controller.downloading = true
	controller.mutex.Unlock()
	go controller.watch()
	log.Debug().Str("slug", controller.Game().Slug).Msg("game download started")
	return
}

func (controller *Download) Game() game.Game {
	defer controller.mutex.RUnlock()
	controller.mutex.RLock()
	return controller.game
}

func (controller *Download) IsDownloading() bool {
	defer controller.mutex.RUnlock()
	controller.mutex.RLock()
	return controller.downloading
}

func (controller *Download) IsUnzipping() bool {
	return !controller.IsDownloading() && controller.IsRunning()
}

func (controller *Download) IsComplete() bool {
	return controller.IsStarted() && !controller.IsRunning()
}

func (controller *Download) IsRunning() bool {
	defer controller.mutex.RUnlock()
	controller.mutex.RLock()
	return controller.running
}

func (controller *Download) IsStarted() bool {
	defer controller.mutex.RUnlock()
	controller.mutex.RLock()
	return controller.started
}

func (controller *Download) Filesize() int64 {
	defer controller.mutex.RUnlock()
	controller.mutex.RLock()
	if controller.unzip != nil {
		return controller.unzip.Filesize()
	}
	if controller.download != nil {
		return controller.download.Filesize()
	}
	return 0
}

func (controller *Download) StartTime() time.Time {
	defer controller.mutex.RUnlock()
	controller.mutex.RLock()
	if controller.download != nil {
		return controller.download.StartTime()
	}
	return time.Time{}
}

func (controller *Download) EndTime() time.Time {
	defer controller.mutex.RUnlock()
	controller.mutex.RLock()
	if controller.unzip != nil && controller.unzip.IsComplete() {
		return controller.unzip.EndTime()
	}
	if controller.download != nil && controller.download.IsComplete() {
		return controller.download.EndTime()
	}
	return time.Time{}
}

func (controller *Download) Duration() (duration time.Duration) {
	defer controller.mutex.RUnlock()
	controller.mutex.RLock()
	if controller.unzip != nil {
		if controller.unzip.EndTime().IsZero() {
			return time.Since(controller.download.StartTime())
		}
		return controller.unzip.EndTime().Sub(controller.download.StartTime())
	}
	if controller.download != nil {
		return controller.download.Duration()
	}
	return
}

func (controller *Download) BytesPerSecond() float64 {
	defer controller.mutex.RUnlock()
	controller.mutex.RLock()
	if controller.unzip != nil {
		return controller.unzip.BytesPerSecond()
	}
	if controller.download != nil {
		return controller.download.BytesPerSecond()
	}
	return 0
}

func (controller *Download) Err() error {
	defer controller.mutex.RUnlock()
	controller.mutex.RLock()
	return controller.err
}

func (controller *Download) watch() {
	controller.notifySubcriber()
	controller.subscribeSubscriber(controller.download)
	<-controller.download.Done
	log.Debug().Str("slug", controller.Game().Slug).Msg("game downloads download part finished")
	if controller.download.Err != nil {
		controller.mutex.Lock()
		controller.running = false
		controller.downloading = false
		controller.err = controller.download.Err
		controller.mutex.Unlock()
		log.Error().Err(controller.download.Err).Str("slug", controller.Game().Slug).Msg("error downloading game")
		return
	}
	controller.notifySubcriber()
	controller.unsubscribeSubscriber(controller.download)

	controller.mutex.Lock()
	controller.unzip = filesystem.NewUnzip(
		path.Join(controller.controller.settings.GameDirectory, controller.download.Filename()),
		path.Join(controller.controller.settings.GameDirectory, controller.game.Slug),
	)
	controller.unzip.StartUnzip()
	controller.downloading = false
	controller.mutex.Unlock()
	log.Debug().Str("slug", controller.Game().Slug).Msg("game downloads unzip part started")
	controller.notifySubcriber()
	controller.subscribeSubscriber(controller.unzip)
	<-controller.unzip.Done
	log.Debug().Str("slug", controller.Game().Slug).Msg("game downloads unzip part finished")
	controller.mutex.Lock()
	if controller.unzip.Err != nil {
		controller.err = controller.unzip.Err
		log.Error().Err(controller.unzip.Err).Str("slug", controller.game.Slug).Msg("error unzipping game")
	}
	controller.running = false
	controller.mutex.Unlock()
	controller.notifySubcriber()
	controller.unsubscribeSubscriber(controller.unzip)
	err := os.Remove(controller.unzip.Filename())
	if err != nil {
		log.Error().Err(err).Str("slug", controller.Game().Slug).Msg("error removing game download file")
	}
}

func (controller *Download) Subscribe(subscriber chan struct{}) {
	defer controller.mutex.Unlock()
	controller.mutex.Lock()
	log.Trace().Msg("new subscriber to download")
	controller.subscriber = append(controller.subscriber, subscriber)
}

func (controller *Download) Unsubscribe(subscriber chan struct{}) {
	defer controller.mutex.Unlock()
	controller.mutex.Lock()
	index := slices.Index(controller.subscriber, subscriber)
	slices.Delete(controller.subscriber, index, index+1)
}

func (controller *Download) notifySubcriber() {
	defer controller.mutex.RUnlock()
	controller.mutex.RLock()
	log.Trace().Msg("notify subscriber of download")
	for _, subscriber := range controller.subscriber {
		util.ChannelWriteNonBlocking(subscriber, struct{}{})
	}
}

func (controller *Download) SubscribeProgress(subscriber chan float64) {
	defer controller.mutex.Unlock()
	controller.mutex.Lock()
	log.Trace().Msg("new subscriber for progress of download")
	controller.subscriberprogress = append(controller.subscriberprogress, subscriber)
}

func (controller *Download) UnsubscribeProgress(subscriber chan float64) {
	defer controller.mutex.Unlock()
	controller.mutex.Lock()
	controller.subscriberprogress = append(controller.subscriberprogress, subscriber)
}

func (controller *Download) subscribeSubscriber(publisher util.Publisher[float64]) {
	defer controller.mutex.RUnlock()
	controller.mutex.RLock()
	for _, subscriber := range controller.subscriberprogress {
		publisher.Subscribe(subscriber)
	}
}

func (controller *Download) unsubscribeSubscriber(publisher util.Publisher[float64]) {
	defer controller.mutex.RUnlock()
	controller.mutex.RLock()
	for _, subscriber := range controller.subscriberprogress {
		publisher.Unsubscribe(subscriber)
	}
}
