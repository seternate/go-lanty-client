package gameviewmodel

import (
	"fmt"

	"github.com/seternate/go-lanty-client/internal/game"
)

type LaunchArgumentGroup struct {
	Header     string
	Tiles      []*LaunchArgumentTile
	ShowHeader bool
}

func NewLaunchArgumentGroupFromModel(group game.LaunchParamGroup, showHeader bool) (*LaunchArgumentGroup, error) {
	vm := &LaunchArgumentGroup{
		Header:     group.HeaderText(),
		Tiles:      make([]*LaunchArgumentTile, 0, len(group.Params)),
		ShowHeader: showHeader,
	}

	for _, param := range group.Params {
		tile, err := NewLaunchArgumentTileFromModel(param)
		if err != nil {
			return nil, fmt.Errorf("failed to create argument tile: %w", err)
		}
		vm.Tiles = append(vm.Tiles, tile)
	}

	return vm, nil
}

func (vm *LaunchArgumentGroup) GetArgumentValues() []game.LaunchArg {
	out := make([]game.LaunchArg, 0, len(vm.Tiles))

	for _, tile := range vm.Tiles {
		out = append(out, tile.GetLaunchArg())
	}

	return out
}
