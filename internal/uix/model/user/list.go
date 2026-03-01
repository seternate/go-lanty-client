package model

import (
	"slices"
	"strings"

	"fyne.io/fyne/v2/data/binding"
	viewmodel "github.com/seternate/go-lanty-client/internal/ui/model"
)

var _ viewmodel.ChangeNotifier = (*UserListModel)(nil)

type UserListModel struct {
	users binding.UntypedList
}

func NewUserListModel() *UserListModel {
	return &UserListModel{
		users: binding.NewUntypedList(),
	}
}

func (m *UserListModel) GetUsers() []*UserTileModel {
	userviewmodels := make([]*UserTileModel, 0)

	users, err := m.users.Get()
	if err != nil {
		return userviewmodels
	}
	for _, user := range users {
		userviewmodels = append(userviewmodels, user.(*UserTileModel))
	}
	return userviewmodels
}

func (m *UserListModel) GetUserByIPAddress(ipAddress string) *UserTileModel {
	users, err := m.users.Get()
	if err != nil {
		return nil
	}
	for _, user := range users {
		if tile, ok := user.(*UserTileModel); ok && tile.IPAddress == ipAddress {
			return tile
		}
	}
	return nil
}

func (m *UserListModel) GetUserByName(name string) *UserTileModel {
	users, err := m.users.Get()
	if err != nil {
		return nil
	}
	for _, user := range users {
		tile, ok := user.(*UserTileModel)
		if ok {
			tilename, _ := tile.Name.Get()
			if tilename == name {
				return tile
			}
		}
	}
	return nil
}

func (m *UserListModel) AddUser(user *UserTileModel) {
	m.users.Append(user)
}

func (m *UserListModel) RemoveUser(ipAddress string) {
	users, err := m.users.Get()
	if err != nil {
		return
	}
	for _, user := range users {
		if tile, ok := user.(*UserTileModel); ok && tile.IPAddress == ipAddress {
			m.users.Remove(user)
			return
		}
	}
}

func (m *UserListModel) GetUserIPAddressesOrderedByName() []string {
	userbyip := make(map[string]string, 0)
	ipAddresses := make([]string, 0)
	for _, user := range m.GetUsers() {
		name, _ := user.Name.Get()
		userbyip[user.IPAddress] = name
		ipAddresses = append(ipAddresses, user.IPAddress)
	}
	slices.SortStableFunc(ipAddresses, func(a, b string) int {
		return strings.Compare(userbyip[a], userbyip[b])
	})
	return ipAddresses
}

func (m *UserListModel) GetUserNamesOrdered() []string {
	userNames := make([]string, 0)
	for _, user := range m.GetUsers() {
		name, _ := user.Name.Get()
		userNames = append(userNames, name)
	}
	slices.SortStableFunc(userNames, func(a, b string) int {
		return strings.Compare(a, b)
	})
	return userNames
}

func (m *UserListModel) AddChangeListener(fn func()) {
	l := binding.NewDataListener(fn)
	m.users.AddListener(l)
}
