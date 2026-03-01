package user

import (
	"sync"

	userevent "github.com/seternate/go-lanty-client/internal/application/user/event"
	userdomain "github.com/seternate/go-lanty-client/internal/domain/user"
	model "github.com/seternate/go-lanty-client/internal/ui/model/user"
)

type CatalogReader interface {
	GetByIP(ip string) (userdomain.UserCatalogItem, error)
}

type UserController struct {
	UserListModel *model.UserListModel
	mu            sync.RWMutex
	catalogReader CatalogReader
}

func NewUserController(catalogReader CatalogReader) *UserController {
	return &UserController{
		UserListModel: model.NewUserListModel(),
		catalogReader: catalogReader,
	}
}

func (controller *UserController) OnUserAdded(e userevent.CatalogEvent) {
	controller.mu.Lock()
	defer controller.mu.Unlock()

	catalogItem, err := controller.catalogReader.GetByIP(e.IP)
	if err != nil {
		return
	}

	tile := model.NewUserTileModel(catalogItem.IP)
	tile.SetName(catalogItem.Name)
	controller.UserListModel.AddUser(tile)
}

func (controller *UserController) OnUserUpdated(e userevent.CatalogEvent) {
	controller.mu.Lock()
	defer controller.mu.Unlock()

	catalogItem, err := controller.catalogReader.GetByIP(e.IP)
	if err != nil {
		return
	}

	tile := controller.UserListModel.GetUserByIPAddress(e.IP)
	tile.SetName(catalogItem.Name)
}
func (controller *UserController) OnUserRemoved(e userevent.CatalogEvent) {
	controller.mu.Lock()
	defer controller.mu.Unlock()

	controller.UserListModel.RemoveUser(e.IP)
}
