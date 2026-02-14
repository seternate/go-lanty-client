package api

import (
	"context"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"io"
	"maps"
	"os"
	"path/filepath"
	"sync"
	"time"

	gameapp "github.com/seternate/go-lanty-client/internal/application/game"
	gameappsrv "github.com/seternate/go-lanty-client/internal/application/game/service"
	"github.com/seternate/go-lanty-client/internal/domain/game"
	gamecontroller "github.com/seternate/go-lanty-client/internal/ui/controller/game"
)

const alternatingSlug = "z-alternating"

var _ gameappsrv.GameCatalogSource = (*gameMockAPIClient)(nil)
var _ gameappsrv.BlobDownloader = (*gameMockAPIClient)(nil)
var _ gameappsrv.GameInstallationProgressFetcher = (*gameMockAPIClient)(nil)
var _ gamecontroller.CatalogIconFetcher = (*gameMockAPIClient)(nil)
var _ gameappsrv.LaunchSpecSource = (*gameMockAPIClient)(nil)
var _ gamecontroller.GameLauncherSpecReader = (*gameMockAPIClient)(nil)

type gameMockAPIClient struct {
	mu                 sync.RWMutex
	catalog            map[string]game.GameCatalogItem
	icons              map[string]image.Image
	order              []string
	alternatingPresent bool
	DownloadPath       string
	progress           map[string]gameapp.InstallationProgress
}

func NewGameMockAPIClient(path string) *gameMockAPIClient {
	client := &gameMockAPIClient{
		catalog:      make(map[string]game.GameCatalogItem),
		icons:        make(map[string]image.Image),
		progress:     make(map[string]gameapp.InstallationProgress),
		DownloadPath: path,
	}
	client.seedCatalog()
	return client
}

func (client *gameMockAPIClient) FetchCatalog() ([]game.GameCatalogItem, error) {
	client.mu.Lock()
	defer client.mu.Unlock()
	client.toggleAlternatingGameLocked()

	out := make([]game.GameCatalogItem, 0, len(client.order))
	for _, slug := range client.order {
		if item, ok := client.catalog[slug]; ok {
			out = append(out, item)
		}
	}
	return out, nil
}

func (client *gameMockAPIClient) GetIcon(slug string) (image.Image, error) {
	client.mu.RLock()
	icon, ok := client.icons[slug]
	client.mu.RUnlock()
	if !ok {
		return nil, errors.New("icon not found")
	}
	return icon, nil
}

func (client *gameMockAPIClient) Download(ctx context.Context, slug string) (filePath string, err error) {
	client.mu.RLock()
	downloadPath := client.DownloadPath
	client.mu.RUnlock()
	if downloadPath == "" {
		downloadPath = "."
	}

	err = os.MkdirAll(downloadPath, 0755)
	if err != nil {
		return "", err
	}

	sourcePath := filepath.Join(".", slug+"-mock.zip")
	destPath := filepath.Join(downloadPath, slug+".zip")

	src, err := os.Open(sourcePath)
	if err != nil {
		return "", err
	}
	defer src.Close()

	info, err := src.Stat()
	if err != nil {
		return "", err
	}

	dst, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	client.mu.Lock()
	client.progress[slug] = gameapp.InstallationProgress{
		Slug:      slug,
		Total:     uint64(info.Size()),
		Completed: 0,
		StartedAt: time.Now(),
	}
	client.mu.Unlock()

	const maxDownloadBytesPerSecond = 4 * 120 * 1024 * 1024

	buffer := make([]byte, 32*1024)
	var copyErr error
	rateStart := time.Now()
	var rateBytes uint64
	for {
		select {
		case <-ctx.Done():
			copyErr = ctx.Err()
		default:
		}

		if copyErr != nil {
			break
		}

		n, readErr := src.Read(buffer)
		if n > 0 {
			_, writeErr := dst.Write(buffer[:n])
			if writeErr != nil {
				copyErr = writeErr
				break
			}
			client.mu.Lock()
			progress := client.progress[slug]
			progress.Completed += uint64(n)
			client.progress[slug] = progress
			client.mu.Unlock()

			rateBytes += uint64(n)
			elapsed := time.Since(rateStart)
			if elapsed < time.Second {
				expected := time.Duration(rateBytes * uint64(time.Second) / uint64(maxDownloadBytesPerSecond))
				if elapsed < expected {
					time.Sleep(expected - elapsed)
					rateStart = time.Now()
					rateBytes = 0
				}
			} else {
				rateStart = time.Now()
				rateBytes = 0
			}
		}

		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				break
			}
			copyErr = readErr
			break
		}
	}

	client.mu.Lock()
	progress := client.progress[slug]
	progress.EndedAt = time.Now()
	if copyErr != nil {
		progress.Failed = true
	} else {
		progress.Succeeded = true
		if progress.Total > 0 {
			progress.Completed = progress.Total
		}
	}
	client.progress[slug] = progress
	client.mu.Unlock()

	if copyErr != nil {
		return "", copyErr
	}
	return destPath, nil
}

func (client *gameMockAPIClient) GetAllProgress() (map[string]gameapp.InstallationProgress, error) {
	client.mu.RLock()
	defer client.mu.RUnlock()

	return maps.Clone(client.progress), nil
}

func (client *gameMockAPIClient) GetProgress(slug string) (gameapp.InstallationProgress, error) {
	client.mu.RLock()
	defer client.mu.RUnlock()

	progress, ok := client.progress[slug]
	if !ok {
		return gameapp.InstallationProgress{}, errors.New("progress not found")
	}

	return progress, nil
}

func (client *gameMockAPIClient) GetLaunchSpec(slug string, mode gameappsrv.LaunchSpecMode) (game.GameLaunchSpec, error) {
	switch mode {
	case gameappsrv.LaunchSpecModePlay:
		return client.getSingleplayerLaunchSpec()
	case gameappsrv.LaunchSpecModeJoin:
		return client.getJoinLaunchSpec()
	case gameappsrv.LaunchSpecModeHost:
		return client.getHostLaunchSpec()

	}
	return game.GameLaunchSpec{}, errors.New("invalid launch spec mode")
}

func (client *gameMockAPIClient) getSingleplayerLaunchSpec() (game.GameLaunchSpec, error) {
	nointroParam, err := game.NewGameLaunchParam(
		game.GameLaunchParamInput{
			Type:              game.GameLaunchParamFlag,
			Name:              "No Intro",
			Description:       "Whether to skip the intro",
			Argument:          "-novid",
			Enabled:           true,
			ArgumentSeparator: " ",
		})
	if err != nil {
		return game.GameLaunchSpec{}, err
	}

	consoleParam, err := game.NewGameLaunchParam(
		game.GameLaunchParamInput{
			Type:              game.GameLaunchParamFlag,
			Name:              "Console",
			Description:       "Whether to open the console",
			Argument:          "-console",
			Enabled:           true,
			ArgumentSeparator: " ",
		})
	if err != nil {
		return game.GameLaunchSpec{}, err
	}

	netgraphParam, err := game.NewGameLaunchParam(
		game.GameLaunchParamInput{
			Type:              game.GameLaunchParamEnum,
			Name:              "Netgraph",
			Description:       "Whether to show the netgraph",
			Argument:          "+net_graph",
			Value:             "1",
			EnumValues:        []game.EnumValueOption{{Value: "1", Label: "Enable"}, {Value: "0", Label: "Disable"}},
			Enabled:           true,
			ArgumentSeparator: " ",
			ValueSeparator:    " ",
		})
	if err != nil {
		return game.GameLaunchSpec{}, err
	}

	spec, err := game.NewGameLaunchSpec("left4dead.exe", []game.GameLaunchParam{
		*nointroParam,
		*consoleParam,
		*netgraphParam,
	})
	if err != nil {
		return game.GameLaunchSpec{}, err
	}

	return *spec, nil
}

func (client *gameMockAPIClient) getJoinLaunchSpec() (game.GameLaunchSpec, error) {
	connectParam, err := game.NewGameLaunchParam(
		game.GameLaunchParamInput{
			Type:              game.GameLaunchParamString,
			Name:              "Connect",
			Description:       "The address to connect to",
			Argument:          "+connect",
			Value:             "",
			Required:          true,
			Enabled:           true,
			ArgumentSeparator: " ",
			ValueSeparator:    " ",
		})
	if err != nil {
		return game.GameLaunchSpec{}, err
	}

	nointroParam, err := game.NewGameLaunchParam(
		game.GameLaunchParamInput{
			Type:              game.GameLaunchParamFlag,
			Name:              "No Intro",
			Description:       "Whether to skip the intro",
			Argument:          "-novid",
			Enabled:           true,
			ArgumentSeparator: " ",
		})
	if err != nil {
		return game.GameLaunchSpec{}, err
	}

	consoleParam, err := game.NewGameLaunchParam(
		game.GameLaunchParamInput{
			Type:              game.GameLaunchParamFlag,
			Name:              "Console",
			Description:       "Whether to open the console",
			Argument:          "-console",
			Enabled:           true,
			ArgumentSeparator: " ",
		})
	if err != nil {
		return game.GameLaunchSpec{}, err
	}

	netgraphParam, err := game.NewGameLaunchParam(
		game.GameLaunchParamInput{
			Type:              game.GameLaunchParamEnum,
			Name:              "Netgraph",
			Description:       "Whether to show the netgraph",
			Argument:          "+net_graph",
			Value:             "1",
			EnumValues:        []game.EnumValueOption{{Value: "1", Label: "Enable"}, {Value: "0", Label: "Disable"}},
			Enabled:           true,
			ArgumentSeparator: " ",
			ValueSeparator:    " ",
		})
	if err != nil {
		return game.GameLaunchSpec{}, err
	}

	spec, err := game.NewGameLaunchSpec("left4dead.exe", []game.GameLaunchParam{
		*connectParam,
		*nointroParam,
		*consoleParam,
		*netgraphParam,
	})
	if err != nil {
		return game.GameLaunchSpec{}, err
	}

	return *spec, nil
}

func (client *gameMockAPIClient) getHostLaunchSpec() (game.GameLaunchSpec, error) {
	minint := int64(2)
	maxint := int64(64)

	minfloat := float64(0.1)
	maxfloat := float64(10)

	floatprecision := int64(10)

	serverbaseParam, err := game.NewGameLaunchParam(
		game.GameLaunchParamInput{
			Type:              game.GameLaunchParamFlag,
			Name:              "Server Base",
			Description:       "The server base to use",
			Argument:          "-serverbase",
			Enabled:           true,
			Required:          true,
			ArgumentSeparator: " ",
		})
	if err != nil {
		return game.GameLaunchSpec{}, err
	}

	nointroParam, err := game.NewGameLaunchParam(
		game.GameLaunchParamInput{
			Type:              game.GameLaunchParamFlag,
			Name:              "No Intro",
			Description:       "Whether to skip the intro",
			Argument:          "-novid",
			Enabled:           true,
			ArgumentSeparator: " ",
		})
	if err != nil {
		return game.GameLaunchSpec{}, err
	}

	mapParam, err := game.NewGameLaunchParam(
		game.GameLaunchParamInput{
			Type:              game.GameLaunchParamEnum,
			Name:              "Map",
			Description:       "The map to play",
			Argument:          "",
			Value:             "dust2",
			EnumValues:        []game.EnumValueOption{{Value: "dust2", Label: "Dust 2"}, {Value: "inferno", Label: "Inferno"}},
			Required:          true,
			Enabled:           true,
			ArgumentSeparator: " ;",
			ValueSeparator:    " ",
		})
	if err != nil {
		return game.GameLaunchSpec{}, err
	}

	maxPlayersParam, err := game.NewGameLaunchParam(
		game.GameLaunchParamInput{
			Type:              game.GameLaunchParamInt,
			Name:              "Max Players",
			Description:       "The maximum number of players",
			Argument:          "+maxplayers",
			Value:             "16",
			Required:          true,
			Enabled:           true,
			ArgumentSeparator: "?",
			ValueSeparator:    "=",
			MinInt:            &minint,
			MaxInt:            &maxint,
		})
	if err != nil {
		return game.GameLaunchSpec{}, err
	}

	gravityParam, err := game.NewGameLaunchParam(
		game.GameLaunchParamInput{
			Type:              game.GameLaunchParamFloat,
			Name:              "Gravity Multiplier",
			Description:       "The gravity multiplier",
			Argument:          "+gravity",
			Value:             "1",
			Enabled:           true,
			ArgumentSeparator: "; ",
			ValueSeparator:    "=",
			MinFloat:          &minfloat,
			MaxFloat:          &maxfloat,
			FloatPrecision:    &floatprecision,
		})
	if err != nil {
		return game.GameLaunchSpec{}, err
	}

	friendlyFireParam, err := game.NewGameLaunchParam(
		game.GameLaunchParamInput{
			Type:        game.GameLaunchParamEnum,
			Name:        "Friendly Fire",
			Description: "Whether friendly fire is enabled",
			Argument:    "+friendlyfire",
			Value:       "false",
			EnumValues:  []game.EnumValueOption{{Value: "false", Label: "Disable"}, {Value: "true", Label: "Enable"}},
		})
	if err != nil {
		return game.GameLaunchSpec{}, err
	}

	spec, err := game.NewGameLaunchSpec("unrealtournament.exe", []game.GameLaunchParam{
		*serverbaseParam,
		*nointroParam,
		*mapParam,
		*maxPlayersParam,
		*gravityParam,
		*friendlyFireParam,
	})
	if err != nil {
		return game.GameLaunchSpec{}, err
	}

	return *spec, nil
}

func (c *gameMockAPIClient) seedCatalog() {
	c.addCatalogItem("example-game", "Example Game", makeIcon(100, 150, 200), 2500000000, true, false)
	c.addCatalogItem("downloading-game", "Downloading Game", makeIcon(200, 100, 150), 5000000000, true, true)
	c.addCatalogItem("extracting-game", "Extracting Game", makeIcon(156, 39, 176), 3200000000, false, true)
	c.addCatalogItem("not-installed-game", "Not Installed Game", makeIcon(0, 188, 212), 4800000000, true, true)
	c.addCatalogItem("error-game", "Error Game", makeIcon(244, 67, 54), 1200000000, true, true)
	c.addCatalogItem("not-installed-game-2", "Not Installed Game 2", makeIcon(255, 152, 0), 1800000000, true, true)
	c.addCatalogItem("installed-game", "Installed Game", makeIcon(97, 97, 97), 850000000, true, true)
}

func (c *gameMockAPIClient) addCatalogItem(slug, name string, icon image.Image, blobSize uint64, supportsJoiningMultiplayer bool, supportsHostingServer bool) {
	if icon == nil {
		icon = makeIcon(120, 120, 120)
	}
	if _, exists := c.catalog[slug]; !exists {
		c.order = append(c.order, slug)
	}
	if _, exists := c.icons[slug]; !exists {
		c.icons[slug] = icon
	}
	c.catalog[slug] = game.NewGameCatalogItem(
		slug,
		name,
		"mock-icon-"+slug,
		blobSize,
		slug+".exe",
		supportsJoiningMultiplayer,
		supportsHostingServer,
	)
}

func (c *gameMockAPIClient) toggleAlternatingGameLocked() {
	if c.alternatingPresent {
		delete(c.catalog, alternatingSlug)
		c.removeFromOrder(alternatingSlug)
		c.alternatingPresent = false
		return
	}

	c.addCatalogItem(alternatingSlug, "Z Alternating Game", makeIcon(60, 179, 113), 2100000000, true, true)
	c.alternatingPresent = true
}

func (c *gameMockAPIClient) removeFromOrder(slug string) {
	for i, candidate := range c.order {
		if candidate == slug {
			c.order = append(c.order[:i], c.order[i+1:]...)
			return
		}
	}
}

func makeIcon(r, g, b uint8) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 64, 64))
	c := color.RGBA{R: r, G: g, B: b, A: 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{C: c}, image.Pt(0, 0), draw.Src)
	return img
}
