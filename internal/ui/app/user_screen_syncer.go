package app

import (
	"context"
	"sync"

	"github.com/rs/zerolog"
	userviewmodel "github.com/seternate/go-lanty-client/internal/ui/viewmodel/user"
	"github.com/seternate/go-lanty-client/internal/user"
)

type UserCatalogReader interface {
	GetAll(ctx context.Context) ([]user.CatalogItem, error)
	GetByIP(ctx context.Context, ip string) (user.CatalogItem, error)
}

type UserScreenSyncer struct {
	logger     zerolog.Logger
	mu         sync.RWMutex
	userScreen *userviewmodel.UserScreen

	catalogReader UserCatalogReader
}

func NewUserScreenSyncer(logger zerolog.Logger, userScreen *userviewmodel.UserScreen, catalogReader UserCatalogReader) *UserScreenSyncer {
	return &UserScreenSyncer{
		logger:        logger,
		userScreen:    userScreen,
		catalogReader: catalogReader,
	}
}

func (syncer *UserScreenSyncer) OnCatalogItemAdded(e user.CatalogEvent) {
	syncer.mu.Lock()
	defer syncer.mu.Unlock()
	ctx := context.Background()
	const handler = "OnCatalogItemAdded"

	catalogItem, err := syncer.catalogReader.GetByIP(ctx, e.IP)
	if err != nil {
		syncer.logger.Warn().Err(err).Str("handler", handler).Str("ip", e.IP).Msg("catalogReader.GetByIP failed")
		return
	}

	if err := syncer.userScreen.AddUserTile(catalogItem); err != nil {
		syncer.logger.Error().Err(err).Str("handler", handler).Str("ip", e.IP).Msg("AddUserTile failed")
	}

	catalogItems, err := syncer.catalogReader.GetAll(ctx)
	if err != nil {
		syncer.logger.Warn().Err(err).Str("handler", handler).Str("ip", e.IP).Msg("catalogReader.GetAll failed")
		return
	}
	if err := syncer.userScreen.UpdateAvailableUsers(len(catalogItems)); err != nil {
		syncer.logger.Error().Err(err).Str("handler", handler).Str("ip", e.IP).Msg("UpdateAvailableUsers failed")
	}
}

func (syncer *UserScreenSyncer) OnCatalogItemUpdated(e user.CatalogEvent) {
	syncer.mu.Lock()
	defer syncer.mu.Unlock()
	ctx := context.Background()
	const handler = "OnCatalogItemUpdated"

	catalogItem, err := syncer.catalogReader.GetByIP(ctx, e.IP)
	if err != nil {
		syncer.logger.Warn().Err(err).Str("handler", handler).Str("ip", e.IP).Msg("catalogReader.GetByIP failed")
		return
	}

	if err := syncer.userScreen.UpdateUserTile(catalogItem); err != nil {
		syncer.logger.Error().Err(err).Str("handler", handler).Str("ip", e.IP).Msg("UpdateUserTile failed")
	}

	catalogItems, err := syncer.catalogReader.GetAll(ctx)
	if err != nil {
		syncer.logger.Warn().Err(err).Str("handler", handler).Str("ip", e.IP).Msg("catalogReader.GetAll failed")
		return
	}
	if err := syncer.userScreen.UpdateAvailableUsers(len(catalogItems)); err != nil {
		syncer.logger.Error().Err(err).Str("handler", handler).Str("ip", e.IP).Msg("UpdateAvailableUsers failed")
	}
}

func (syncer *UserScreenSyncer) OnCatalogItemRemoved(e user.CatalogEvent) {
	syncer.mu.Lock()
	defer syncer.mu.Unlock()
	ctx := context.Background()
	const handler = "OnCatalogItemRemoved"

	if err := syncer.userScreen.RemoveUserTile(e.IP); err != nil {
		syncer.logger.Error().Err(err).Str("handler", handler).Str("ip", e.IP).Msg("RemoveUserTile failed")
	}

	catalogItems, err := syncer.catalogReader.GetAll(ctx)
	if err != nil {
		syncer.logger.Warn().Err(err).Str("handler", handler).Str("ip", e.IP).Msg("catalogReader.GetAll failed")
		return
	}
	if err := syncer.userScreen.UpdateAvailableUsers(len(catalogItems)); err != nil {
		syncer.logger.Error().Err(err).Str("handler", handler).Str("ip", e.IP).Msg("UpdateAvailableUsers failed")
	}
}
