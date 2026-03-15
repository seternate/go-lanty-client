package userviewmodel

import (
	"fmt"
	"slices"
	"strings"

	"fyne.io/fyne/v2/data/binding"
	"github.com/seternate/go-lanty-client/internal/ui/viewmodel"
	"github.com/seternate/go-lanty-client/internal/user"
)

type UserScreen struct {
	Header         *viewmodel.HeaderScrollScreen
	users          binding.UntypedList
	availableUsers binding.Int
}

func NewUserScreen() *UserScreen {
	availableUsers := binding.NewInt()

	info := binding.NewSprintf("%d available", availableUsers)

	vm := &UserScreen{
		Header:         viewmodel.NewHeaderScrollScreenWithData("Users", info),
		users:          binding.NewUntypedList(),
		availableUsers: availableUsers,
	}

	return vm
}

func (vm *UserScreen) GetUserTiles() []*UserTile {
	userTiles := make([]*UserTile, 0)

	users, err := vm.users.Get()
	if err != nil {
		return nil
	}

	for _, user := range users {
		if tile, ok := user.(*UserTile); ok {
			userTiles = append(userTiles, tile)
		}
	}

	slices.SortStableFunc(userTiles, func(a, b *UserTile) int {
		return strings.Compare(a.GetName(), b.GetName())
	})

	return userTiles
}

func (vm *UserScreen) getUserTileByIPAddress(ipAddress string) *UserTile {
	users, err := vm.users.Get()
	if err != nil {
		return nil
	}

	for _, user := range users {
		if tile, ok := user.(*UserTile); ok && tile.GetIPAddress() == ipAddress {
			return tile
		}
	}

	return nil
}

func (vm *UserScreen) AddUserTile(model user.CatalogItem) error {
	userTile, err := NewUserTileFromModel(model)
	if err != nil {
		return fmt.Errorf("failed to create user tile: %w", err)
	}

	err = vm.users.Append(userTile)
	if err != nil {
		return fmt.Errorf("failed to add user tile: %w", err)
	}

	return nil
}

func (vm *UserScreen) UpdateUserTile(update user.CatalogItem) error {
	userTile := vm.getUserTileByIPAddress(update.IP)
	if userTile == nil {
		return fmt.Errorf("user tile not found")
	}

	err := userTile.UpdateFromModel(update)
	if err != nil {
		return fmt.Errorf("failed to update user tile: %w", err)
	}

	return nil
}

func (vm *UserScreen) RemoveUserTile(ipAddress string) error {
	userTile := vm.getUserTileByIPAddress(ipAddress)
	if userTile == nil {
		return fmt.Errorf("user tile not found")
	}

	err := vm.users.Remove(userTile)
	if err != nil {
		return fmt.Errorf("failed to remove user tile: %w", err)
	}

	return nil
}

func (vm *UserScreen) UpdateAvailableUsers(availableUsers int) error {
	return vm.availableUsers.Set(availableUsers)
}

func (vm *UserScreen) AddChangeListener(fn func()) {
	l := binding.NewDataListener(fn)

	vm.Header.AddChangeListener(fn)
	vm.users.AddListener(l)
	vm.availableUsers.AddListener(l)
}
