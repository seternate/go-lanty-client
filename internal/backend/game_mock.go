package backend

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
	"slices"
	"sync"
	"time"

	"github.com/seternate/go-lanty-client/internal/game"
)

const alternatingSlug = "z-alternating"
const extensionMockGameSlug = "installed-game" // game that has mock extensions for UI testing

type gameMockAPIClient struct {
	mu                 sync.RWMutex
	catalog            map[string]game.CatalogItem
	icons              map[string]image.Image
	order              []string
	alternatingPresent bool
	progress           map[string]game.InstallationProgress
}

func NewGameMockAPIClient() *gameMockAPIClient {
	client := &gameMockAPIClient{
		catalog:  make(map[string]game.CatalogItem),
		icons:    make(map[string]image.Image),
		progress: make(map[string]game.InstallationProgress),
	}
	client.seedCatalog()
	return client
}

func (client *gameMockAPIClient) FetchCatalog(ctx context.Context) ([]game.CatalogItem, error) {
	client.mu.Lock()
	defer client.mu.Unlock()
	client.toggleAlternatingGameLocked()

	out := make([]game.CatalogItem, 0, len(client.order))
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

func (client *gameMockAPIClient) Download(ctx context.Context, slug string, baseDir string) (filePath string, err error) {
	if baseDir == "" {
		baseDir = "."
	}

	err = os.MkdirAll(baseDir, 0755)
	if err != nil {
		return "", err
	}

	sourcePath := filepath.Join(".", slug+"-mock.zip")
	destPath := filepath.Join(baseDir, slug+".zip")

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
	client.progress[slug] = game.InstallationProgress{
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

func (client *gameMockAPIClient) GetAllInstallationProgresses(ctx context.Context) (map[string]game.InstallationProgress, error) {
	client.mu.RLock()
	defer client.mu.RUnlock()

	return maps.Clone(client.progress), nil
}

func (client *gameMockAPIClient) GetProgress(slug string) (game.InstallationProgress, error) {
	client.mu.RLock()
	defer client.mu.RUnlock()

	progress, ok := client.progress[slug]
	if !ok {
		return game.InstallationProgress{}, errors.New("progress not found")
	}

	return progress, nil
}

func (client *gameMockAPIClient) GetLaunchSpec(ctx context.Context, slug string, mode game.LaunchSpecMode) (game.LaunchSpec, error) {
	switch mode {
	case game.LaunchSpecModePlay:
		return client.getSingleplayerLaunchSpec(ctx)
	case game.LaunchSpecModeJoin:
		return client.getJoinLaunchSpec(ctx)
	case game.LaunchSpecModeHost:
		return client.getHostLaunchSpec(ctx)
	}
	return game.LaunchSpec{}, errors.New("invalid launch spec mode")
}

func (client *gameMockAPIClient) getSingleplayerLaunchSpec(ctx context.Context) (game.LaunchSpec, error) {
	nointroParam, err := game.NewLaunchParam(
		game.LaunchParamInput{
			Type:              game.LaunchParamFlag,
			Name:              "No Intro",
			Description:       "Whether to skip the intro",
			Argument:          "-novid",
			Enabled:           true,
			ArgumentSeparator: " ",
		})
	if err != nil {
		return game.LaunchSpec{}, err
	}

	consoleParam, err := game.NewLaunchParam(
		game.LaunchParamInput{
			Type:              game.LaunchParamFlag,
			Name:              "Console",
			Description:       "Whether to open the console",
			Argument:          "-console",
			Enabled:           true,
			ArgumentSeparator: " ",
		})
	if err != nil {
		return game.LaunchSpec{}, err
	}

	netgraphParam, err := game.NewLaunchParam(
		game.LaunchParamInput{
			Type:              game.LaunchParamEnum,
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
		return game.LaunchSpec{}, err
	}

	spec, err := game.NewLaunchSpec("left4dead.exe", []game.LaunchParam{
		*nointroParam,
		*consoleParam,
		*netgraphParam,
	})
	if err != nil {
		return game.LaunchSpec{}, err
	}

	return *spec, nil
}

func (client *gameMockAPIClient) getJoinLaunchSpec(ctx context.Context) (game.LaunchSpec, error) {
	connectParam, err := game.NewLaunchParam(
		game.LaunchParamInput{
			Type:              game.LaunchParamString,
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
		return game.LaunchSpec{}, err
	}

	nointroParam, err := game.NewLaunchParam(
		game.LaunchParamInput{
			Type:              game.LaunchParamFlag,
			Name:              "No Intro",
			Description:       "Whether to skip the intro",
			Argument:          "-novid",
			Enabled:           true,
			ArgumentSeparator: " ",
		})
	if err != nil {
		return game.LaunchSpec{}, err
	}

	consoleParam, err := game.NewLaunchParam(
		game.LaunchParamInput{
			Type:              game.LaunchParamFlag,
			Name:              "Console",
			Description:       "Whether to open the console",
			Argument:          "-console",
			Enabled:           true,
			ArgumentSeparator: " ",
		})
	if err != nil {
		return game.LaunchSpec{}, err
	}

	netgraphParam, err := game.NewLaunchParam(
		game.LaunchParamInput{
			Type:              game.LaunchParamEnum,
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
		return game.LaunchSpec{}, err
	}

	spec, err := game.NewLaunchSpec("left4dead.exe", []game.LaunchParam{
		*connectParam,
		*nointroParam,
		*consoleParam,
		*netgraphParam,
	})
	if err != nil {
		return game.LaunchSpec{}, err
	}

	return *spec, nil
}

func (client *gameMockAPIClient) getHostLaunchSpec(ctx context.Context) (game.LaunchSpec, error) {
	minint := int64(2)
	maxint := int64(64)

	minfloat := float64(0.1)
	maxfloat := float64(10)

	floatprecision := int64(10)

	serverbaseParam, err := game.NewLaunchParam(
		game.LaunchParamInput{
			Type:              game.LaunchParamFlag,
			Name:              "Server Base",
			Description:       "The server base to use",
			Argument:          "-serverbase",
			Enabled:           true,
			Required:          true,
			ArgumentSeparator: " ",
		})
	if err != nil {
		return game.LaunchSpec{}, err
	}

	nointroParam, err := game.NewLaunchParam(
		game.LaunchParamInput{
			Type:              game.LaunchParamFlag,
			Name:              "No Intro",
			Description:       "Whether to skip the intro",
			Argument:          "-novid",
			Enabled:           true,
			ArgumentSeparator: " ",
		})
	if err != nil {
		return game.LaunchSpec{}, err
	}

	mapParam, err := game.NewLaunchParam(
		game.LaunchParamInput{
			Type:              game.LaunchParamEnum,
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
		return game.LaunchSpec{}, err
	}

	maxPlayersParam, err := game.NewLaunchParam(
		game.LaunchParamInput{
			Type:              game.LaunchParamInt,
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
		return game.LaunchSpec{}, err
	}

	gravityParam, err := game.NewLaunchParam(
		game.LaunchParamInput{
			Type:              game.LaunchParamFloat,
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
		return game.LaunchSpec{}, err
	}

	friendlyFireParam, err := game.NewLaunchParam(
		game.LaunchParamInput{
			Type:        game.LaunchParamEnum,
			Name:        "Friendly Fire",
			Description: "Whether friendly fire is enabled",
			Argument:    "+friendlyfire",
			Value:       "false",
			EnumValues:  []game.EnumValueOption{{Value: "false", Label: "Disable"}, {Value: "true", Label: "Enable"}},
		})
	if err != nil {
		return game.LaunchSpec{}, err
	}

	spec, err := game.NewLaunchSpec("unrealtournament.exe", []game.LaunchParam{
		*serverbaseParam,
		*nointroParam,
		*mapParam,
		*maxPlayersParam,
		*gravityParam,
		*friendlyFireParam,
	})
	if err != nil {
		return game.LaunchSpec{}, err
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
	c.catalog[slug] = game.CatalogItem{
		Slug:              slug,
		Name:              name,
		IconHash:          "mock-icon-" + slug,
		InstalledFileSize: blobSize,
		DetectionHints:    game.InstallationDetectionHints{},
		Capabilities: game.Capabilities{
			JoiningMultiplayer: supportsJoiningMultiplayer,
			HostingServer:      supportsHostingServer,
		},
	}
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

// gameExtensionMock implements extension service interfaces without conflicting with gameMockAPIClient.Download.
type gameExtensionMock struct {
	gameExtensions map[string][]game.Extension
	mu             sync.RWMutex
}

func NewGameExtensionMock() *gameExtensionMock {
	mock := &gameExtensionMock{
		gameExtensions: make(map[string][]game.Extension),
	}
	mock.gameExtensions[extensionMockGameSlug] = []game.Extension{
		game.NewExtension(extensionMockGameSlug, "custom_maps_v2.vpk", "Custom Maps Pack", 15728640, true),
		game.NewExtension(extensionMockGameSlug, "sound_overhaul.vpk", "Sound Overhaul", 8388608, false),
		game.NewExtension(extensionMockGameSlug, "ui_skins.vpk", "UI Skins", 2097152, true),
	}
	return mock
}

func (m *gameExtensionMock) ListForGame(ctx context.Context, slug string) ([]game.Extension, error) {
	if slug != extensionMockGameSlug {
		return nil, nil
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	extensions, ok := m.gameExtensions[slug]
	if !ok {
		return nil, nil
	}
	return slices.Clone(extensions), nil
}

func (m *gameExtensionMock) Upload(ctx context.Context, slug string, localFilePath string) error {
	return nil
}

func (m *gameExtensionMock) Download(ctx context.Context, slug string, filename string, destDir string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	extensions, ok := m.gameExtensions[slug]
	if !ok {
		return errors.New("game not found")
	}
	for _, ext := range extensions {
		if ext.Filename == filename {
			ext.MarkDownloaded()
			return nil
		}
	}
	return errors.New("extension not found")
}
