package gamebrowser

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/seternate/go-lanty/pkg/settings"
	lanty "github.com/seternate/lanty-api-golang/pkg/api"
	"github.com/seternate/lanty-api-golang/pkg/download"
	"github.com/seternate/lanty-api-golang/pkg/game"
)

type Gamebrowser struct {
	*container.Scroll

	client     *lanty.Client
	settings   *settings.Settings
	downloader *download.Downloader
}

func NewGameBrowser(client *lanty.Client, settings *settings.Settings, games game.Games, downloader *download.Downloader) (*Gamebrowser, error) {
	gamebrowser := &Gamebrowser{
		client:     client,
		settings:   settings,
		downloader: downloader,
	}

	gamebrowser.RefreshUI(games)

	return gamebrowser, nil
}

func (gb *Gamebrowser) RefreshUI(games game.Games) {
	gridWrap := container.NewGridWrap(fyne.NewSize(400, 130))
	for _, game := range games {
		gametile, _ := NewGametile(gb.client, gb.settings, game, gb.downloader)
		gridWrap.Add(gametile.MainContainer)
	}
	gb.Scroll = container.NewVScroll(gridWrap)
}
