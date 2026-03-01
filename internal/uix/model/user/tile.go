package model

import (
	"fyne.io/fyne/v2/data/binding"
	viewmodel "github.com/seternate/go-lanty-client/internal/ui/model"
)

var _ viewmodel.ChangeNotifier = (*UserTileModel)(nil)

type UserTileModel struct {
	IPAddress string
	Name      binding.String
}

func NewUserTileModel(ipAddress string) *UserTileModel {
	tile := &UserTileModel{
		IPAddress: ipAddress,
		Name:      binding.NewString(),
	}

	tile.Name.Set("N/A")

	return tile
}

func (tile *UserTileModel) SetName(name string) {
	tile.Name.Set(name)
}

func (tile *UserTileModel) AddChangeListener(fn func()) {
	l := binding.NewDataListener(fn)
	tile.Name.AddListener(l)
}
