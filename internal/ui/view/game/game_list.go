package gameview

// import (
// 	"fyne.io/fyne/v2"
// 	"fyne.io/fyne/v2/container"
// 	"fyne.io/fyne/v2/widget"
// 	model "github.com/seternate/go-lanty-client/internal/ui/model/game"
// )

// type GameList struct {
// 	widget.BaseWidget

// 	games *model.GameListModel

// 	tilesContainer *fyne.Container
// 	tiles          map[string]*GameTile

// 	OnPlayButtonPressed       func(slug string)
// 	OnJoinButtonPressed       func(slug string)
// 	OnServerButtonPressed     func(slug string)
// 	OnDownloadButtonPressed   func(slug string)
// 	OnOpenButtonPressed       func(slug string)
// 	OnExtensionBadgePressed   func(slug string)
// }

// func NewGameList(games *model.GameListModel) *GameList {
// 	view := &GameList{
// 		games:          games,
// 		tilesContainer: container.NewVBox(),
// 		tiles:          make(map[string]*GameTile),
// 	}

// 	view.RefreshGameTiles()

// 	games.AddChangeListener(func() {
// 		view.Refresh()
// 	})

// 	view.ExtendBaseWidget(view)
// 	return view
// }

// func (view *GameList) RefreshGameTiles() {
// 	slugsOrdered := view.games.GetGameSlugsOrdered()
// 	games := view.games.GetGames()

// 	currentSlugs := make(map[string]bool)
// 	for _, game := range games {
// 		currentSlugs[game.Slug] = true
// 	}

// 	for slug := range view.tiles {
// 		if !currentSlugs[slug] {
// 			delete(view.tiles, slug)
// 		}
// 	}

// 	for _, game := range games {
// 		if _, found := view.tiles[game.Slug]; !found {
// 			view.tiles[game.Slug] = NewGameTile(game)
// 		}
// 		tile := view.tiles[game.Slug]
// 		tile.OnPlayButtonPressed = func(slug string) {
// 			if view.OnPlayButtonPressed != nil {
// 				view.OnPlayButtonPressed(slug)
// 			}
// 		}
// 		tile.OnJoinButtonPressed = func(slug string) {
// 			if view.OnJoinButtonPressed != nil {
// 				view.OnJoinButtonPressed(slug)
// 			}
// 		}
// 		tile.OnServerButtonPressed = func(slug string) {
// 			if view.OnServerButtonPressed != nil {
// 				view.OnServerButtonPressed(slug)
// 			}
// 		}
// 		tile.OnDownloadButtonPressed = func(slug string) {
// 			if view.OnDownloadButtonPressed != nil {
// 				view.OnDownloadButtonPressed(slug)
// 			}
// 		}
// 		tile.OnOpenButtonPressed = func(slug string) {
// 			if view.OnOpenButtonPressed != nil {
// 				view.OnOpenButtonPressed(slug)
// 			}
// 		}
// 		tile.OnExtensionBadgePressed = func(slug string) {
// 			if view.OnExtensionBadgePressed != nil {
// 				view.OnExtensionBadgePressed(slug)
// 			}
// 		}
// 	}

// 	view.tilesContainer.RemoveAll()
// 	for _, slug := range slugsOrdered {
// 		view.tilesContainer.Add(view.tiles[slug])
// 	}
// }

// func (view *GameList) Refresh() {
// 	view.RefreshGameTiles()
// 	view.BaseWidget.Refresh()
// }

// func (view *GameList) CreateRenderer() fyne.WidgetRenderer {
// 	c := container.NewScroll(view.tilesContainer)
// 	c.SetMinSize(fyne.NewSize(0, 0))
// 	return widget.NewSimpleRenderer(c)
// }
