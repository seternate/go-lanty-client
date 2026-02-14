package service

import (
	eventbus "github.com/asaskevich/EventBus"
	"github.com/seternate/go-lanty-client/internal/application/user/event"
	"github.com/seternate/go-lanty-client/internal/domain/user"
)

type UserCatalogRefreshService interface {
	Refresh() error
}

type UserCatalogSource interface {
	FetchCatalog() ([]user.UserCatalogItem, error)
}

var _ UserCatalogRefreshService = (*CatalogRefresher)(nil)

type CatalogRefresher struct {
	eventbus          eventbus.Bus
	userCatalogSource UserCatalogSource
	userCatalogRepo   user.UserCatalogRepository
}

func NewUserCatalogRefreshService(eventbus eventbus.Bus, userCatalogSource UserCatalogSource, userCatalogRepo user.UserCatalogRepository) *CatalogRefresher {
	return &CatalogRefresher{
		eventbus:          eventbus,
		userCatalogSource: userCatalogSource,
		userCatalogRepo:   userCatalogRepo,
	}
}

func (service *CatalogRefresher) Refresh() error {
	oldCatalog, err := service.userCatalogRepo.GetAll()
	if err != nil {
		return err
	}
	newCatalog, err := service.userCatalogSource.FetchCatalog()
	if err != nil {
		return err
	}

	changes := diff(oldCatalog, newCatalog)

	if len(changes.Added) == 0 && len(changes.Updated) == 0 && len(changes.Removed) == 0 {
		return nil
	}

	err = service.userCatalogRepo.ReplaceAll(newCatalog...)
	if err != nil {
		return err
	}

	for _, item := range changes.Added {
		service.eventbus.Publish(event.CatalogAddedEvent, event.CatalogEvent{IP: item.IP})
	}

	for _, item := range changes.Updated {
		service.eventbus.Publish(event.CatalogUpdatedEvent, event.CatalogEvent{IP: item.IP})
	}

	for _, item := range changes.Removed {
		service.eventbus.Publish(event.CatalogRemovedEvent, event.CatalogEvent{IP: item.IP})
	}

	return nil
}

type diffResult struct {
	Added   []user.UserCatalogItem
	Updated []user.UserCatalogItem
	Removed []user.UserCatalogItem
}

func diff(old []user.UserCatalogItem, new []user.UserCatalogItem) diffResult {
	oldByIP := make(map[string]user.UserCatalogItem, len(old))
	for _, item := range old {
		oldByIP[item.IP] = item
	}

	newByIP := make(map[string]user.UserCatalogItem, len(new))
	for _, item := range new {
		newByIP[item.IP] = item
	}

	added := make([]user.UserCatalogItem, 0, len(new))
	updated := make([]user.UserCatalogItem, 0, len(new))
	for _, newItem := range new {
		oldItem, found := oldByIP[newItem.IP]
		if !found {
			added = append(added, newItem)
		} else if found && !newItem.Equal(oldItem) {
			updated = append(updated, newItem)
		}
	}

	removed := make([]user.UserCatalogItem, 0, len(old))
	for _, oldItem := range oldByIP {
		if _, ok := newByIP[oldItem.IP]; !ok {
			removed = append(removed, oldItem)
		}
	}

	return diffResult{
		Added:   added,
		Updated: updated,
		Removed: removed,
	}
}
