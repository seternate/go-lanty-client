package controller

import (
	"image"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty/pkg/filesystem"
	"github.com/seternate/go-lanty/pkg/game"
	"github.com/seternate/go-lanty/pkg/user"
	"github.com/seternate/go-lanty/pkg/util"
)

type GameController struct {
	parent          *Controller
	games           game.Games
	gameIcons       map[string]image.Image
	subscriber      []chan struct{}
	refreshinterval time.Duration
	mutex           sync.RWMutex
	err             error
}

func NewGameController(parent *Controller, refreshinterval time.Duration) (controller *GameController) {
	controller = &GameController{
		parent:          parent,
		refreshinterval: refreshinterval,
		gameIcons:       make(map[string]image.Image, 50),
		subscriber:      make([]chan struct{}, 0, 50),
	}
	go controller.run()
	return
}

func (controller *GameController) GetGames() (games game.Games) {
	defer controller.mutex.RUnlock()
	controller.mutex.RLock()
	return controller.games
}

func (controller *GameController) GetIcon(game game.Game) (image image.Image) {
	defer controller.mutex.RUnlock()
	controller.mutex.RLock()
	return controller.gameIcons[game.Slug]
}

func (controller *GameController) Err() error {
	defer controller.mutex.RUnlock()
	controller.mutex.RLock()
	return controller.err
}

func (controller *GameController) StartGame(game game.Game) {
	paths, err := filesystem.SearchFileByName(game.ClientExecutable, controller.parent.settings.GameDirectory)
	if err != nil {
		return
	}
	cmd := exec.Command("./" + filepath.Base(paths[0]))
	cmd.Dir = filepath.Dir(paths[0])
	cmd.Start()
	log.Debug().Str("slug", game.Slug).Str("cmd", cmd.String()).Msg("started game")
}

func (controller *GameController) OpenGameInExplorer(game game.Game) {
	paths, err := filesystem.SearchFileByName(game.ClientExecutable, controller.parent.settings.GameDirectory)
	if err != nil {
		return
	}
	cmd := exec.Command("explorer", "/select,", paths[0])
	cmd.Start()
	log.Debug().Str("slug", game.Slug).Str("cmd", cmd.String()).Msg("open game in explorer")
}

func (controller *GameController) JoinServer(game game.Game, user user.User) {
	//TODO
	return
}

func (controller *GameController) StartServer(game game.Game) {
	paths, err := filesystem.SearchFileByName(game.ServerExecutable, controller.parent.settings.GameDirectory)
	if err != nil {
		return
	}
	cmd := exec.Command("./"+filepath.Base(paths[0]), strings.Split(game.ServerArgument, " ")...)
	cmd.Dir = filepath.Dir(paths[0])
	cmd.Start()
	log.Debug().Str("slug", game.Slug).Str("cmd", cmd.String()).Msg("started game server")
}

func (controller *GameController) Subscribe(subscriber chan struct{}) {
	defer controller.mutex.Unlock()
	controller.mutex.Lock()
	log.Trace().Msg("new subscriber to gamecontroller")
	controller.subscriber = append(controller.subscriber, subscriber)
}

func (controller *GameController) Unsubscribe(subscriber chan struct{}) {
	defer controller.mutex.Unlock()
	controller.mutex.Lock()
	index := slices.Index(controller.subscriber, subscriber)
	slices.Delete(controller.subscriber, index, index+1)
}

func (controller *GameController) notifySubcriber() {
	defer controller.mutex.RUnlock()
	controller.mutex.RLock()
	log.Trace().Msg("notify subscriber of usercontroller")
	for _, subscriber := range controller.subscriber {
		util.ChannelWriteNonBlocking(subscriber, struct{}{})
	}
}

func (controller *GameController) run() {
	for {
		log.Trace().Msg("gamecontroller update loop")
		controller.update()
		time.Sleep(controller.refreshinterval)
	}
}

func (controller *GameController) update() {
	refreshed, err := controller.updateGames()
	if err != nil || !refreshed {
		controller.mutex.Lock()
		controller.err = err
		controller.mutex.Unlock()
		return
	}
	_, err = controller.updateIcons()
	if err != nil {
		controller.mutex.Lock()
		controller.err = err
		controller.mutex.Unlock()
	}
	controller.notifySubcriber()
}

func (controller *GameController) updateGames() (refreshed bool, err error) {
	refreshed = false
	games, err := controller.getServerGames()
	if err != nil {
		return
	}
	if controller.games.Equal(games) {
		return
	}
	refreshed = true
	controller.mutex.Lock()
	controller.games = games
	controller.mutex.Unlock()
	log.Debug().Interface("games", games).Msg("updated games in gamescontroller")
	return
}

func (controller *GameController) updateIcons() (refreshed bool, err error) {
	defer controller.mutex.Unlock()
	controller.mutex.Lock()
	refreshed = false
	for _, game := range controller.games.Games() {
		_, hasIcon := controller.gameIcons[game.Slug]
		if hasIcon {
			log.Trace().Str("slug", game.Slug).Msg("game icon already cached")
			continue
		}
		image, err := controller.parent.client.Game.GetIcon(game)
		if err != nil {
			log.Error().Err(err).Str("slug", game.Slug).Msg("error retrieving game icon from server")
			return refreshed, err
		}
		refreshed = true
		controller.gameIcons[game.Slug] = image
		log.Debug().Str("slug", game.Slug).Msg("game icon updated")
	}
	for slug := range controller.gameIcons {
		_, err = controller.games.Get(slug)
		if err != nil {
			err = nil
			delete(controller.gameIcons, slug)
			log.Trace().Str("slug", slug).Msg("deleted game icon from cache")
		}
	}
	return
}

func (controller *GameController) getServerGames() (games game.Games, err error) {
	slugs, err := controller.parent.client.Game.GetGames()
	if err != nil {
		log.Error().Err(err).Msg("error retrieving gameslist from server")
		return
	}
	for _, slug := range slugs {
		game, err := controller.parent.client.Game.GetGame(slug)
		if err != nil {
			log.Error().Err(err).Str("slug", slug).Msg("error retrieving game from server")
			return games, err
		}
		err = games.Add(game)
		if err != nil {
			return games, err
		}
	}
	return
}
