package lanty

import (
	"github.com/seternate/go-lanty/pkg/settings"
	lantyapi "github.com/seternate/lanty-api-golang/pkg/api"
	"github.com/seternate/lanty-api-golang/pkg/download"
	"github.com/seternate/lanty-api-golang/pkg/game"
)

type Lanty struct {
	Settings   *settings.Settings
	Client     *lantyapi.Client
	Downloader *download.Downloader
	Games      *game.Games
}

func NewLanty(settings *settings.Settings, client *lantyapi.Client, downloader *download.Downloader, games *game.Games) (*Lanty, error) {
	lanty := &Lanty{
		Settings:   settings,
		Client:     client,
		Downloader: downloader,
		Games:      games,
	}

	return lanty, nil
}
