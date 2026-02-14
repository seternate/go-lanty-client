package model

import (
	"slices"
	"strings"

	"fyne.io/fyne/v2/data/binding"
	viewmodel "github.com/seternate/go-lanty-client/internal/ui/model"
)

var _ viewmodel.ChangeNotifier = (*GameListModel)(nil)

type GameListModel struct {
	games binding.UntypedList
}

func NewGameListModel() *GameListModel {
	return &GameListModel{
		games: binding.NewUntypedList(),
	}
}

func (m *GameListModel) GetGames() []*GameTileModel {
	gameviewmodels := make([]*GameTileModel, 0)

	games, err := m.games.Get()
	if err != nil {
		return gameviewmodels
	}
	for _, game := range games {
		gameviewmodels = append(gameviewmodels, game.(*GameTileModel))
	}

	return gameviewmodels
}

func (m *GameListModel) GetGameBySlug(slug string) *GameTileModel {
	games, err := m.games.Get()
	if err != nil {
		return nil
	}
	for _, g := range games {
		if tile, ok := g.(*GameTileModel); ok && tile.Slug == slug {
			return tile
		}
	}
	return nil
}

func (m *GameListModel) AddGame(game *GameTileModel) {
	m.games.Append(game)
}

func (m *GameListModel) RemoveGame(slug string) {
	games, err := m.games.Get()
	if err != nil {
		return
	}
	for _, g := range games {
		if tile, ok := g.(*GameTileModel); ok && tile.Slug == slug {
			m.games.Remove(g)
			return
		}
	}
}

func (m *GameListModel) GetGameSlugsOrdered() []string {
	slugs := make([]string, 0)
	for _, game := range m.GetGames() {
		slugs = append(slugs, game.Slug)
	}
	slices.SortStableFunc(slugs, func(a, b string) int {
		return strings.Compare(a, b)
	})
	return slugs
}

func (m *GameListModel) AddChangeListener(fn func()) {
	l := binding.NewDataListener(fn)
	m.games.AddListener(l)
}
