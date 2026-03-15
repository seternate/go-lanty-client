package app

import (
	"context"
	"sync"

	userviewmodel "github.com/seternate/go-lanty-client/internal/ui/viewmodel/user"
	"github.com/seternate/go-lanty-client/internal/user"
)

type UserCatalogReader interface {
	GetAll(ctx context.Context) ([]user.CatalogItem, error)
	GetByIP(ctx context.Context, ip string) (user.CatalogItem, error)
}

type UserScreenSyncer struct {
	mu         sync.RWMutex
	userScreen *userviewmodel.UserScreen

	catalogReader UserCatalogReader
}

func NewUserScreenSyncer(userScreen *userviewmodel.UserScreen, catalogReader UserCatalogReader) *UserScreenSyncer {
	return &UserScreenSyncer{
		userScreen:    userScreen,
		catalogReader: catalogReader,
	}
}

func (syncer *UserScreenSyncer) OnCatalogItemAdded(e user.CatalogEvent) {
	syncer.mu.Lock()
	defer syncer.mu.Unlock()

	catalogItem, err := syncer.catalogReader.GetByIP(context.Background(), e.IP)
	if err != nil {
		return
	}

	syncer.userScreen.AddUserTile(catalogItem)

	catalogItems, err := syncer.catalogReader.GetAll(context.Background())
	if err != nil {
		return
	}
	syncer.userScreen.UpdateAvailableUsers(len(catalogItems))
}

func (syncer *UserScreenSyncer) OnCatalogItemUpdated(e user.CatalogEvent) {
	syncer.mu.Lock()
	defer syncer.mu.Unlock()

	catalogItem, err := syncer.catalogReader.GetByIP(context.Background(), e.IP)
	if err != nil {
		return
	}

	syncer.userScreen.UpdateUserTile(catalogItem)

	catalogItems, err := syncer.catalogReader.GetAll(context.Background())
	if err != nil {
		return
	}
	syncer.userScreen.UpdateAvailableUsers(len(catalogItems))
}

func (syncer *UserScreenSyncer) OnCatalogItemRemoved(e user.CatalogEvent) {
	syncer.mu.Lock()
	defer syncer.mu.Unlock()

	syncer.userScreen.RemoveUserTile(e.IP)

	catalogItems, err := syncer.catalogReader.GetAll(context.Background())
	if err != nil {
		return
	}
	syncer.userScreen.UpdateAvailableUsers(len(catalogItems))
}
