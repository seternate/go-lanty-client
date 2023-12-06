package controller

import (
	"os"
	"path"
	"slices"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty-client/pkg/setting"
	"github.com/seternate/go-lanty/pkg/api"
	"github.com/seternate/go-lanty/pkg/filesystem"
	"github.com/seternate/go-lanty/pkg/game"
	"github.com/seternate/go-lanty/pkg/network"
	"github.com/seternate/go-lanty/pkg/util"
)

type DownloadController struct {
	parent     *Controller
	downloads  map[game.Game]*network.Download
	unzips     map[game.Game]*filesystem.Unzip
	queue      []game.Game
	subscriber map[game.Game][]Subscription
}

type Subscription struct {
	Game             game.Game
	Download         chan *network.Download
	Unzip            chan *filesystem.Unzip
	DownloadProgress chan float64
	UnzipProgress    chan float64
}

func NewDownloadController(parent *Controller) (controller *DownloadController) {
	controller = &DownloadController{
		parent:     parent,
		downloads:  make(map[game.Game]*network.Download),
		unzips:     make(map[game.Game]*filesystem.Unzip),
		queue:      make([]game.Game, 0),
		subscriber: make(map[game.Game][]Subscription),
	}

	go controller.run()

	return
}

func (controller *DownloadController) DownloadGame(game game.Game) {
	if slices.Contains(controller.queue, game) {
		log.Trace().Str("slug", game.Slug).Msg("game already in download queue")
		return
	}
	controller.queue = append(controller.queue, game)
	log.Trace().Str("slug", game.Slug).Msg("added game to download queue")
}

func (controller DownloadController) Subscribe(subscriber Subscription) {
	controller.subscriber[subscriber.Game] = append(controller.subscriber[subscriber.Game], subscriber)
}

func (controller *DownloadController) run() {
	for {
		controller.startDownloadsFromQueue()
		controller.extractFinishedDownloads()
		controller.deleteFinishedDownloads()
		time.Sleep(150 * time.Millisecond)
	}
}

func (controller *DownloadController) startDownloadsFromQueue() {
	for index, game := range controller.queue {
		if controller.isDownloading(game) {
			controller.queue = slices.Delete(controller.queue, 0, 1)
			log.Trace().Str("slug", game.Slug).Msg("game already downloading - removed from queue")
			continue
		}

		download, err := controller.client().Game.Download(game, controller.settings().GameDirectory)
		if err != nil {
			return
		}

		controller.queue = slices.Delete(controller.queue, index, index+1)
		controller.downloads[game] = download
		log.Trace().Str("slug", game.Slug).Msg("download of game started")

		for _, subscriber := range controller.subscriber[game] {
			util.ChannelWriteNonBlocking(subscriber.Download, download)
			download.Subscribe(subscriber.DownloadProgress)
		}
	}
}

func (controller *DownloadController) extractFinishedDownloads() {
	for game, download := range controller.downloads {
		if download.IsComplete() {
			log.Trace().Err(download.Err).Str("slug", game.Slug).Msg("game download error status")
			unzip := filesystem.NewUnzip(
				path.Join(controller.settings().GameDirectory, download.Filename),
				path.Join(controller.settings().GameDirectory, game.Slug),
			)
			unzip.StartUnzip()
			delete(controller.downloads, game)
			controller.unzips[game] = unzip
			log.Trace().Str("slug", game.Slug).Msg("start unzip of game")
			for _, subscriber := range controller.subscriber[game] {
				util.ChannelWriteNonBlocking(subscriber.Unzip, unzip)
				unzip.Subscribe(subscriber.UnzipProgress)
			}
		}
	}
}

func (controller *DownloadController) deleteFinishedDownloads() {
	for game, unzip := range controller.unzips {
		if unzip.IsComplete() {
			log.Trace().Err(unzip.Err).Str("slug", game.Slug).Msg("game unzip error status")
			os.Remove(unzip.Filename)
			delete(controller.unzips, game)
			log.Trace().Str("slug", game.Slug).Msg("removed game zip file")
		}
	}
}

func (controller *DownloadController) isDownloading(game game.Game) bool {
	_, foundDownload := controller.downloads[game]
	_, foundUnzip := controller.unzips[game]
	return foundDownload || foundUnzip
}

func (controller *DownloadController) client() *api.Client {
	return controller.parent.client
}

func (controller *DownloadController) settings() *setting.Settings {
	return controller.parent.settings
}
