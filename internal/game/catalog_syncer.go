package game

import (
	"context"
	"fmt"

	eventbus "github.com/asaskevich/EventBus"
)

type CatalogSource interface {
	FetchCatalog(ctx context.Context) ([]CatalogItem, error)
}

type CatalogSyncer struct {
	eventbus      eventbus.Bus
	remoteCatalog CatalogSource
	localCatalog  CatalogRepository
}

func NewCatalogSyncer(bus eventbus.Bus, remoteCatalog CatalogSource, localCatalog CatalogRepository) *CatalogSyncer {
	return &CatalogSyncer{
		eventbus:      bus,
		remoteCatalog: remoteCatalog,
		localCatalog:  localCatalog,
	}
}

func (syncer *CatalogSyncer) Sync(ctx context.Context) error {
	localCatalog, err := syncer.localCatalog.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to get local catalog: %w", err)
	}

	remoteCatalog, err := syncer.remoteCatalog.FetchCatalog(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch remote catalog: %w", err)
	}

	changes := diff(localCatalog, remoteCatalog)

	if len(changes.Added) == 0 && len(changes.Updated) == 0 && len(changes.Removed) == 0 {
		return nil
	}

	err = syncer.localCatalog.ReplaceAll(ctx, remoteCatalog...)
	if err != nil {
		return fmt.Errorf("failed to update local catalog with remote catalog: %w", err)
	}

	for _, item := range changes.Added {
		syncer.eventbus.Publish(CatalogAddedEvent, CatalogEvent{Slug: item.Slug})
	}

	for _, item := range changes.Updated {
		syncer.eventbus.Publish(CatalogUpdatedEvent, CatalogEvent{Slug: item.Slug})
	}

	for _, item := range changes.Removed {
		syncer.eventbus.Publish(CatalogRemovedEvent, CatalogEvent{Slug: item.Slug})
	}

	return nil
}

type diffResult struct {
	Added   []CatalogItem
	Updated []CatalogItem
	Removed []CatalogItem
}

func diff(old []CatalogItem, new []CatalogItem) diffResult {
	oldBySlug := make(map[string]CatalogItem, len(old))
	for _, item := range old {
		oldBySlug[item.Slug] = item
	}

	newBySlug := make(map[string]CatalogItem, len(new))
	for _, item := range new {
		newBySlug[item.Slug] = item
	}

	added := make([]CatalogItem, 0, len(new))
	updated := make([]CatalogItem, 0, len(new))
	for _, newItem := range new {
		oldItem, found := oldBySlug[newItem.Slug]
		if !found {
			added = append(added, newItem)
		} else if found && !newItem.Equal(oldItem) {
			updated = append(updated, newItem)
		}
	}

	removed := make([]CatalogItem, 0, len(old))
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
