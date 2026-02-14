package service

import (
	eventbus "github.com/asaskevich/EventBus"
	"github.com/seternate/go-lanty-client/internal/application/game/event"
	"github.com/seternate/go-lanty-client/internal/domain/game"
)

type GameCatalogRefreshService interface {
	Refresh() error
}

type GameCatalogSource interface {
	FetchCatalog() ([]game.GameCatalogItem, error)
}

var _ GameCatalogRefreshService = (*CatalogRefreshService)(nil)

type CatalogRefreshService struct {
	eventbus      eventbus.Bus
	catalogSource GameCatalogSource
	catalogRepo   game.GameCatalogRepository
}

func NewGameCatalogRefreshService(bus eventbus.Bus, catalogSource GameCatalogSource, catalogRepo game.GameCatalogRepository) *CatalogRefreshService {
	return &CatalogRefreshService{
		eventbus:      bus,
		catalogSource: catalogSource,
		catalogRepo:   catalogRepo,
	}
}

func (service *CatalogRefreshService) Refresh() error {
	oldCatalog, err := service.catalogRepo.GetAll()
	if err != nil {
		return err
	}
	newCatalog, err := service.catalogSource.FetchCatalog()
	if err != nil {
		return err
	}

	changes := diff(oldCatalog, newCatalog)

	if len(changes.Added) == 0 && len(changes.Updated) == 0 && len(changes.Removed) == 0 {
		return nil
	}

	err = service.catalogRepo.ReplaceAll(newCatalog...)
	if err != nil {
		return err
	}

	for _, item := range changes.Added {
		service.eventbus.Publish(event.CatalogAddedEvent, event.CatalogEvent{Slug: item.Slug})
	}

	for _, item := range changes.Updated {
		service.eventbus.Publish(event.CatalogUpdatedEvent, event.CatalogEvent{Slug: item.Slug})
	}

	for _, item := range changes.Removed {
		service.eventbus.Publish(event.CatalogRemovedEvent, event.CatalogEvent{Slug: item.Slug})
	}

	return nil
}

type diffResult struct {
	Added   []game.GameCatalogItem
	Updated []game.GameCatalogItem
	Removed []game.GameCatalogItem
}

func diff(old []game.GameCatalogItem, new []game.GameCatalogItem) diffResult {
	oldBySlug := make(map[string]game.GameCatalogItem, len(old))
	for _, item := range old {
		oldBySlug[item.Slug] = item
	}

	newBySlug := make(map[string]game.GameCatalogItem, len(new))
	for _, item := range new {
		newBySlug[item.Slug] = item
	}

	added := make([]game.GameCatalogItem, 0, len(new))
	updated := make([]game.GameCatalogItem, 0, len(new))
	for _, newItem := range new {
		oldItem, found := oldBySlug[newItem.Slug]
		if !found {
			added = append(added, newItem)
		} else if found && !newItem.Equal(oldItem) {
			updated = append(updated, newItem)
		}
	}

	removed := make([]game.GameCatalogItem, 0, len(old))
	for _, oldItem := range oldBySlug {
		if _, ok := newBySlug[oldItem.Slug]; !ok {
			removed = append(removed, oldItem)
		}
	}

	return diffResult{
		Added:   added,
		Updated: updated,
		Removed: removed,
	}
}
