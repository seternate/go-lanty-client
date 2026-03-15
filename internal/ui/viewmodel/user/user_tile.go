package userviewmodel

import (
	"fmt"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"github.com/seternate/go-lanty-client/internal/user"
)

type UserTile struct {
	ipAddress string
	name      binding.String

	mu sync.RWMutex
}

func NewUserTileFromModel(model user.CatalogItem) (*UserTile, error) {
	vm := &UserTile{
		ipAddress: model.IP,
		name:      binding.NewString(),
	}

	err := vm.UpdateFromModel(model)
	if err != nil {
		return nil, fmt.Errorf("could not update user tile from model: %w", err)
	}

	return vm, nil
}

func (vm *UserTile) UpdateFromModel(update user.CatalogItem) error {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	if vm.ipAddress != update.IP {
		return fmt.Errorf("IP address mismatch: want %s - got %s", vm.ipAddress, update.IP)
	}

	vm.name.Set(update.Name)

	return nil
}

func (vm *UserTile) CopyIPAddressToClipboard() {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	fyne.CurrentApp().Clipboard().SetContent(vm.ipAddress)
}

func (vm *UserTile) AddChangeListener(fn func()) {
	l := binding.NewDataListener(fn)

	vm.name.AddListener(l)
}

func (vm *UserTile) GetIPAddress() string {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	return vm.ipAddress
}

func (vm *UserTile) GetName() string {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	name, err := vm.name.Get()
	if err != nil {
		return "N/A"
	}

	return name
}
