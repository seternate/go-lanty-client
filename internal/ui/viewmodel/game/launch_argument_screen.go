package gameviewmodel

import (
	"fmt"

	"github.com/seternate/go-lanty-client/internal/game"
	"github.com/seternate/go-lanty-client/internal/ui/viewmodel"
)

type LaunchArgumentScreen struct {
	Header   *viewmodel.HeaderScrollScreen
	groups   []*LaunchArgumentGroup
	onSubmit func(values []game.LaunchArg)
}

func NewLaunchArgumentScreen(title string, info string, arguments []game.LaunchParam, onSubmit func(values []game.LaunchArg), onCancel func()) (*LaunchArgumentScreen, error) {
	vm := &LaunchArgumentScreen{
		Header:   viewmodel.NewHeaderScrollScreen(title),
		groups:   make([]*LaunchArgumentGroup, 0),
		onSubmit: onSubmit,
	}

	vm.Header.SetInfo(info)
	vm.Header.SetOnSubmit(func() {
		onSubmit(vm.GetArgumentValues())
	})
	vm.Header.SetOnCancel(onCancel)

	launchParamGroups := game.GroupLaunchParamsByCategory(arguments)
	showGroupHeader := len(launchParamGroups) > 1

	for _, group := range launchParamGroups {
		groupVM, err := NewLaunchArgumentGroupFromModel(group, showGroupHeader)
		if err != nil {
			return nil, fmt.Errorf("failed to create argument group: %w", err)
		}
		vm.groups = append(vm.groups, groupVM)
	}

	return vm, nil
}

func (vm *LaunchArgumentScreen) GetArgumentGroups() []*LaunchArgumentGroup {
	return vm.groups
}

func (vm *LaunchArgumentScreen) GetArgumentValues() []game.LaunchArg {
	out := make([]game.LaunchArg, 0)

	for _, group := range vm.groups {
		out = append(out, group.GetArgumentValues()...)
	}

	return out
}
