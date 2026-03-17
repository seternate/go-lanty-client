package gameviewmodel

import (
	"fmt"

	"fyne.io/fyne/v2/data/binding"
	"github.com/seternate/go-lanty-client/internal/game"
	"github.com/seternate/go-lanty-client/internal/ui/viewmodel"
)

type LaunchArgumentScreen struct {
	Header    *viewmodel.HeaderScrollScreen
	arguments binding.UntypedList
	onSubmit  func(values []game.LaunchArg)
}

func NewLaunchArgumentScreen(title string, info string, arguments []game.LaunchParam, onSubmit func(values []game.LaunchArg)) (*LaunchArgumentScreen, error) {
	vm := &LaunchArgumentScreen{
		Header:    viewmodel.NewHeaderScrollScreen(title),
		arguments: binding.NewUntypedList(),
		onSubmit:  onSubmit,
	}

	vm.Header.SetInfo(info)

	for _, argument := range arguments {
		argumentTile, err := NewLaunchArgumentTileFromModel(argument)
		if err != nil {
			return nil, fmt.Errorf("failed to create argument tile: %w", err)
		}

		err = vm.arguments.Append(argumentTile)
		if err != nil {
			return nil, fmt.Errorf("failed to add argument tile: %w", err)
		}
	}

	return vm, nil
}

func (vm *LaunchArgumentScreen) GetArgumentTiles() []*LaunchArgumentTile {
	argumentTiles := make([]*LaunchArgumentTile, 0)

	arguments, err := vm.arguments.Get()
	if err != nil {
		return nil
	}

	for _, argument := range arguments {
		if tile, ok := argument.(*LaunchArgumentTile); ok {
			argumentTiles = append(argumentTiles, tile)
		}
	}

	return argumentTiles
}
