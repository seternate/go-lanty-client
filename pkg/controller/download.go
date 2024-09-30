package controller

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
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
	subscriberprogress []chan struct{}
	started            bool
	stopped            bool
	running            bool
	downloading        bool
	err                error
	retries            uint64
	mutex              sync.RWMutex
	context            context.Context
	cancelContext      context.CancelFunc
}

func NewDownload(controller *Controller, game game.Game) (download *Download) {
	download = &Download{
		controller:  controller,
		game:        game,
		subscriber:  make([]chan struct{}, 0, 50),
		started:     false,
		stopped:     false,
		running:     false,
		downloading: false,
	}
	return
}

func (controller *Download) Start(ctx context.Context, waitgrp *sync.WaitGroup) (err error) {
	if controller.IsStarted() {
		log.Debug().Str("slug", controller.Game().Slug).Msg("download already started")
		return errors.New("download already started")
	}
	controller.context, controller.cancelContext = context.WithCancel(ctx)
	download, err := controller.controller.client.Game.Download(controller.context, controller.game, controller.controller.settings.GameDirectory)
	if err != nil {
		controller.mutex.Lock()
		if strings.Contains(err.Error(), "connectex: No connection") {
			controller.err = errors.New("error connecting to server")
		} else {
			controller.err = err
		}
		controller.retries += 1
		controller.mutex.Unlock()
		log.Error().Err(err).Str("slug", controller.Game().Slug).Msg("error starting game download from server")
		return
	}
	controller.mutex.Lock()
	controller.download = download
	controller.started = true
	controller.running = true
	controller.downloading = true
	controller.mutex.Unlock()
	waitgrp.Add(1)
	go controller.watch(controller.context, waitgrp)
	log.Debug().Str("slug", controller.Game().Slug).Msg("game download started")
	return
}

func (controller *Download) Controller() *Controller {
	return controller.controller
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

func (controller *Download) IsStopped() bool {
	defer controller.mutex.RUnlock()
	controller.mutex.RLock()
	return controller.stopped
}

func (controller *Download) Filesize() int64 {
	defer controller.mutex.RUnlock()
	controller.mutex.RLock()
	if controller.unzip != nil {
		return int64(controller.unzip.Filesize())
	}
	if controller.download != nil {
		return int64(controller.download.Filesize())
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

func (controller *Download) Progress() float64 {
	defer controller.mutex.RUnlock()
	controller.mutex.RLock()
	if controller.unzip != nil {
		return controller.unzip.Progress()
	}
	if controller.download != nil {
		return controller.download.Progress()
	}
	return 0
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

func (controller *Download) Retries() uint64 {
	defer controller.mutex.RUnlock()
	controller.mutex.RLock()
	return controller.retries
}

func (controller *Download) Stop() {
	if !controller.IsStopped() {
		controller.cancelContext()
		controller.mutex.Lock()
		controller.stopped = true
		controller.mutex.Unlock()
		//time.Sleep(5 * time.Second)
		controller.notifySubcriber()
	}
}

func (controller *Download) watch(ctx context.Context, waitgrp *sync.WaitGroup) {
	defer waitgrp.Done()
	controller.notifySubcriber()
	controller.subscribeSubscriber(controller.download)
	<-controller.download.Done
	if controller.download.Err != nil {
		controller.mutex.Lock()
		controller.running = false
		controller.downloading = false
		if controller.download.Err == context.Canceled {
			controller.err = controller.download.Err
			log.Debug().Err(controller.download.Err).Str("slug", controller.game.Slug).Msg("download canceled")
		} else {
			controller.err = errors.New("error downloading")
			log.Error().Err(controller.download.Err).Str("slug", controller.game.Slug).Msg("error downloading game")
			controller.controller.Status.Error(fmt.Sprintf("Error downloading game: %s", controller.game.Name), 8*time.Second)
		}
		controller.mutex.Unlock()
		controller.notifySubcriber()
		controller.removeGameData(controller.gameDataFilepath())
		log.Trace().Str("slug", controller.game.Slug).Msg("exiting download watch()")
		return
	} else {
		log.Debug().Str("slug", controller.game.Slug).Msg("game downloads download part finished")
	}
	controller.notifySubcriber()
	controller.unsubscribeSubscriber(controller.download)
	controller.mutex.Lock()
	controller.unzip = filesystem.NewUnzip(
		controller.gameDataFilepath(),
		controller.gameDataDestination(),
	)
	controller.unzip.StartUnzip(ctx)
	controller.downloading = false
	controller.mutex.Unlock()
	log.Debug().Str("slug", controller.game.Slug).Msg("game downloads unzip part started")
	controller.notifySubcriber()
	controller.subscribeSubscriber(controller.unzip)
	<-controller.unzip.Done
	controller.mutex.Lock()
	if controller.unzip.Err != nil {
		controller.err = controller.unzip.Err
		if controller.unzip.Err == context.Canceled {
			log.Debug().Err(controller.download.Err).Str("slug", controller.game.Slug).Msg("unzip canceled")
		} else {
			log.Error().Err(controller.unzip.Err).Str("slug", controller.game.Slug).Msg("error unzipping game")
		}
	} else {
		log.Debug().Str("slug", controller.game.Slug).Msg("game downloads unzip part finished")
	}
	controller.running = false
	controller.mutex.Unlock()
	controller.notifySubcriber()
	controller.unsubscribeSubscriber(controller.unzip)
	controller.removeGameData(controller.gameDataFilepath())
	log.Trace().Str("slug", controller.game.Slug).Msg("exiting download watch()")
}

func (controller *Download) removeGameData(filepath string) {
	err := os.Remove(filepath)
	if err != nil {
		log.Error().Err(err).Str("slug", controller.game.Slug).Msg("error removing game data")
	}
}

func (controller *Download) gameDataFilepath() string {
	return path.Join(controller.controller.settings.GameDirectory, controller.download.Filename())
}

func (controller *Download) gameDataDestination() string {
	gamedir := path.Join(controller.controller.settings.GameDirectory, controller.game.Slug)
	paths, err := filesystem.SearchFilesBreadthFirst(controller.controller.Settings.Settings().GameDirectory, controller.game.Client.Executable, 3, 1)
	if len(paths) > 0 && err == nil {
		gamedir = filepath.Dir(paths[0])
	}
	fmt.Println(gamedir)
	return gamedir
}

func (controller *Download) Subscribe(subscriber chan struct{}) {
	defer controller.mutex.Unlock()
	controller.mutex.Lock()
	controller.subscriber = append(controller.subscriber, subscriber)
}

func (controller *Download) Unsubscribe(subscriber chan struct{}) {
	defer controller.mutex.Unlock()
	controller.mutex.Lock()
	index := slices.Index(controller.subscriber, subscriber)
	_ = slices.Delete(controller.subscriber, index, index+1)
}

func (controller *Download) notifySubcriber() {
	defer controller.mutex.RUnlock()
	controller.mutex.RLock()
	for _, subscriber := range controller.subscriber {
		util.ChannelWriteNonBlocking(subscriber, struct{}{})
	}
}

func (controller *Download) SubscribeProgress(subscriber chan struct{}) {
	defer controller.mutex.Unlock()
	controller.mutex.Lock()
	controller.subscriberprogress = append(controller.subscriberprogress, subscriber)
}

func (controller *Download) UnsubscribeProgress(subscriber chan struct{}) {
	defer controller.mutex.Unlock()
	controller.mutex.Lock()
	controller.subscriberprogress = append(controller.subscriberprogress, subscriber)
}

func (controller *Download) subscribeSubscriber(publisher util.Publisher[struct{}]) {
	defer controller.mutex.RUnlock()
	controller.mutex.RLock()
	for _, subscriber := range controller.subscriberprogress {
		publisher.Subscribe(subscriber)
	}
}

func (controller *Download) unsubscribeSubscriber(publisher util.Publisher[struct{}]) {
	defer controller.mutex.RUnlock()
	controller.mutex.RLock()
	for _, subscriber := range controller.subscriberprogress {
		publisher.Unsubscribe(subscriber)
	}
}
