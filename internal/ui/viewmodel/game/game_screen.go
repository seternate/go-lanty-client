package gameviewmodel

import (
	"fmt"
	"slices"
	"strings"

	"fyne.io/fyne/v2/data/binding"
	"github.com/dustin/go-humanize"
	"github.com/seternate/go-lanty-client/internal/game"
	"github.com/seternate/go-lanty-client/internal/ui/viewmodel"
)

type GameScreen struct {
	Header         *viewmodel.HeaderScrollScreen
	games          binding.UntypedList
	availableGames binding.Int
	installedGames binding.Int
	freeDiskSpace  binding.String

	launchRunner                *game.LaunchRunner
	installationRunner          *game.InstallationRunner
	installationDirectoryOpener *game.InstallationDirectoryOpener
}

func NewGameScreen(launchRunner *game.LaunchRunner, installationRunner *game.InstallationRunner, installationDirectoryOpener *game.InstallationDirectoryOpener) *GameScreen {
	availableGames := binding.NewInt()
	installedGames := binding.NewInt()
	freeDiskSpace := binding.NewString()

	info := binding.NewSprintf("%d available, %d installed | %s free disk space", availableGames, installedGames, freeDiskSpace)

	vm := &GameScreen{
		Header:                      viewmodel.NewHeaderScrollScreenWithData("Games", info),
		games:                       binding.NewUntypedList(),
		availableGames:              availableGames,
		installedGames:              installedGames,
		freeDiskSpace:               freeDiskSpace,
		launchRunner:                launchRunner,
		installationRunner:          installationRunner,
		installationDirectoryOpener: installationDirectoryOpener,
	}

	return vm
}

func (vm *GameScreen) GetGameTiles() []*GameTile {
	gametiles := make([]*GameTile, 0)

	games, err := vm.games.Get()
	if err != nil {
		return nil
	}

	for _, game := range games {
		if tile, ok := game.(*GameTile); ok {
			gametiles = append(gametiles, tile)
		}
	}

	slices.SortStableFunc(gametiles, func(a, b *GameTile) int {
		return strings.Compare(a.GetSlug(), b.GetSlug())
	})

	return gametiles
}

func (vm *GameScreen) getGameTileBySlug(slug string) *GameTile {
	games, err := vm.games.Get()
	if err != nil {
		return nil
	}

	for _, game := range games {
		if tile, ok := game.(*GameTile); ok && tile.GetSlug() == slug {
			return tile
		}
	}

	return nil
}

func (vm *GameScreen) AddGameTile(create GameTileCreate) error {
	gametile, err := NewGameTileFromModel(create, vm.launchRunner, vm.installationRunner, vm.installationDirectoryOpener)
	if err != nil {
		return fmt.Errorf("failed to create game tile: %w", err)
	}

	err = vm.games.Append(gametile)
	if err != nil {
		return fmt.Errorf("failed to add game tile: %w", err)
	}

	return nil
}

func (vm *GameScreen) UpdateGameTile(update GameTileUpdate) error {
	gametile := vm.getGameTileBySlug(update.CatalogItem.Slug)
	if gametile == nil {
		return fmt.Errorf("game tile not found")
	}

	err := gametile.UpdateFromModel(update)
	if err != nil {
		return fmt.Errorf("failed to update game tile: %w", err)
	}

	return nil
}

func (vm *GameScreen) RemoveGameTile(slug string) error {
	gametile := vm.getGameTileBySlug(slug)
	if gametile == nil {
		return fmt.Errorf("game tile not found")
	}

	err := vm.games.Remove(gametile)
	if err != nil {
		return fmt.Errorf("failed to remove game tile: %w", err)
	}

	return nil
}

func (vm *GameScreen) UpdateFreeDiskSpace(freeDiskSpace uint64) error {
	return vm.freeDiskSpace.Set(humanize.SIWithDigits(float64(freeDiskSpace), 2, "B"))
}

func (vm *GameScreen) UpdateAvailableGames(availableGames int) error {
	return vm.availableGames.Set(availableGames)
}

func (vm *GameScreen) UpdateInstalledGames(installedGames int) error {
	return vm.installedGames.Set(installedGames)
}

func (vm *GameScreen) AddChangeListener(fn func()) {
	l := binding.NewDataListener(fn)

	vm.Header.AddChangeListener(fn)
	vm.games.AddListener(l)
	vm.availableGames.AddListener(l)
	vm.installedGames.AddListener(l)
	vm.freeDiskSpace.AddListener(l)
}
