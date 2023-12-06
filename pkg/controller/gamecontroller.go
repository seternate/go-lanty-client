package controller

import (
	"image"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/seternate/go-lanty-client/pkg/setting"
	"github.com/seternate/go-lanty/pkg/api"
	"github.com/seternate/go-lanty/pkg/filesystem"
	"github.com/seternate/go-lanty/pkg/game"
	"github.com/seternate/go-lanty/pkg/util"
)

type GameController struct {
	parent          *Controller
	games           game.Games
	gameIcons       map[string]image.Image
	subscriber      []chan struct{}
	refreshinterval time.Duration
	gamesmutex      sync.Mutex
	iconmutex       sync.Mutex

	Err error
}

func NewGameController(parent *Controller, refreshinterval time.Duration) (controller *GameController) {
	controller = &GameController{
		parent:          parent,
		refreshinterval: refreshinterval,
		gameIcons:       make(map[string]image.Image, 50),
		subscriber:      make([]chan struct{}, 0, 50),
	}

	controller.refresh()
	go controller.run()

	return
}

func (controller *GameController) GetGames() (games game.Games) {
	controller.gamesmutex.Lock()
	games = controller.games
	controller.gamesmutex.Unlock()

	return
}

func (controller *GameController) GetIcon(game game.Game) (image image.Image) {
	controller.iconmutex.Lock()
	image = controller.gameIcons[game.Slug]
	controller.iconmutex.Unlock()

	return
}

func (controller *GameController) StartGame(game game.Game) {
	paths, err := filesystem.SearchFileByName(game.ClientExecutable, controller.settings().GameDirectory)
	if err != nil {
		return
	}
	cmd := exec.Command("./" + filepath.Base(paths[0]))
	cmd.Dir = filepath.Dir(paths[0])
	cmd.Start()
}

func (controller *GameController) OpenGameInExplorer(game game.Game) {
	paths, err := filesystem.SearchFileByName(game.ClientExecutable, controller.settings().GameDirectory)
	if err != nil {
		return
	}
	exec.Command("explorer", "/select,", paths[0]).Start()
}

func (controller *GameController) DownloadGame(game game.Game) {
	controller.parent.Download.DownloadGame(game)
}

func (controller *GameController) SubscribeDownload(subscription Subscription) {
	controller.parent.Download.Subscribe(subscription)
}

func (controller *GameController) Subscribe(subscriber chan struct{}) {
	controller.subscriber = append(controller.subscriber, subscriber)
}

func (controller *GameController) client() *api.Client {
	return controller.parent.client
}

func (controller *GameController) settings() *setting.Settings {
	return controller.parent.settings
}

func (controller *GameController) notifySubcriber() {
	for _, subscriber := range controller.subscriber {
		util.ChannelWriteNonBlocking(subscriber, struct{}{})
	}
}

func (controller *GameController) run() {
	for {
		time.Sleep(controller.refreshinterval)
		controller.refresh()
	}
}

func (controller *GameController) refresh() {
	refreshed, err := controller.refreshGames()
	if err != nil || !refreshed {
		controller.Err = err
		return
	}
	_, err = controller.refreshGameIcons()
	if err != nil {
		controller.Err = err
	}
	controller.notifySubcriber()
}

func (controller *GameController) refreshGames() (refreshed bool, err error) {
	refreshed = false
	games, err := controller.getServerGames()
	if err != nil {
		return
	}

	if !controller.games.Equal(games) {
		refreshed = true
		controller.gamesmutex.Lock()
		controller.games = games
		controller.gamesmutex.Unlock()
	}

	return
}

func (controller *GameController) getServerGames() (games game.Games, err error) {
	slugs, err := controller.client().Game.GetGames()
	if err != nil {
		return
	}

	for _, slug := range slugs {
		game, err := controller.client().Game.GetGame(slug)
		if err != nil {
			return games, err
		}
		games.Add(game)
		if err != nil {
			return games, err
		}
	}

	return
}

func (controller *GameController) refreshGameIcons() (refreshed bool, err error) {
	refreshed = false
	controller.iconmutex.Lock()
	for _, slug := range controller.games.Slugs() {
		_, hasIcon := controller.gameIcons[slug]
		if hasIcon {
			continue
		}
		game, err := controller.games.Get(slug)
		if err != nil {
			return refreshed, err
		}
		image, err := controller.client().Game.GetIcon(game)
		if err != nil {
			return refreshed, err
		}
		refreshed = true
		controller.gameIcons[slug] = image
	}
	for slug := range controller.gameIcons {
		_, err = controller.games.Get(slug)
		if err != nil {
			err = nil
			delete(controller.gameIcons, slug)
		}
	}
	controller.iconmutex.Unlock()

	return
}
