package main

import (
	"fmt"
	"net/url"
	"time"

	lanty "github.com/seternate/go-lanty/pkg"
	"github.com/seternate/go-lanty/pkg/settings"
	lantyUI "github.com/seternate/go-lanty/pkg/ui"
	lantyapi "github.com/seternate/lanty-api-golang/pkg/api"
	"github.com/seternate/lanty-api-golang/pkg/download"
	"github.com/seternate/lanty-api-golang/pkg/game"
)

func main() {
	settings, err := settings.LoadSettings("settings/lanty.yaml")
	if err != nil {
		fmt.Printf("Error loading the settings: %s", err)
	}

	timeout, _ := time.ParseDuration("0s")
	baseURL, _ := url.Parse(settings.ServerURL)
	client, _ := lantyapi.NewClient(baseURL, "", "", timeout)

	games, _ := client.Game.GetList()

	downloader := &download.Downloader{Download: map[game.Game]*download.Download{}}

	lanty, _ := lanty.NewLanty(settings, client, downloader, &games)

	lantyui := lantyUI.NewLantyUI(lanty)
	lantyui.ShowAndRun()
}
