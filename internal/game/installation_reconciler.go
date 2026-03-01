package game

import (
	"context"
	"fmt"

	eventbus "github.com/asaskevich/EventBus"
)

type InstallationPathFinder interface {
	FindPath(ctx context.Context, baseDirectory string, slug string, detectionHints InstallationDetectionHints) (installationPath string, found bool, err error)
}

type InstallationReconciler struct {
	bus                    eventbus.Bus
	installationRepo       InstallationRepository
	catalogRepo            CatalogRepository
	installationPathFinder InstallationPathFinder
}

func NewInstallationReconciler(bus eventbus.Bus, installationRepo InstallationRepository, catalogRepo CatalogRepository, installationPathFinder InstallationPathFinder) *InstallationReconciler {
	return &InstallationReconciler{
		bus:                    bus,
		installationRepo:       installationRepo,
		catalogRepo:            catalogRepo,
		installationPathFinder: installationPathFinder,
	}
}

func (reconciler *InstallationReconciler) Reconcile(ctx context.Context, baseDirectory string) error {
	catalog, err := reconciler.catalogRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to get catalog: %w", err)
	}

	installations, err := reconciler.installationRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to get installations: %w", err)
	}

	installationsBySlug := make(map[string]Installation)
	for _, installation := range installations {
		installationsBySlug[installation.Slug] = installation
	}

	catalogBySlug := make(map[string]CatalogItem)
	for _, catalogItem := range catalog {
		catalogBySlug[catalogItem.Slug] = catalogItem
	}

	for _, catalogItem := range catalog {
		if _, found := installationsBySlug[catalogItem.Slug]; !found {
			newInstallation := NewInstallation(catalogItem.Slug)
			installationsBySlug[catalogItem.Slug] = *newInstallation
			err = reconciler.installationRepo.Store(ctx, *newInstallation)
			if err != nil {
				return fmt.Errorf("failed to store new installation: %w", err)
			}
		}
	}

	for _, installation := range installations {
		if _, found := catalogBySlug[installation.Slug]; !found {
			delete(installationsBySlug, installation.Slug)
			err = reconciler.installationRepo.Remove(ctx, installation.Slug)
			if err != nil {
				return fmt.Errorf("failed to remove installation: %w", err)
			}
		}
	}

	for slug, installation := range installationsBySlug {
		installationPath, found, err := reconciler.installationPathFinder.FindPath(ctx, baseDirectory, slug, catalogBySlug[slug].DetectionHints)
		if err != nil {
			return fmt.Errorf("failed to find installation path of game: %w", err)
		}

		oldInstallation := installation

		if found {
			installation.SetDirectory(installationPath)
		} else {
			installation.ClearDirectory()
		}

		err = reconciler.installationRepo.Store(ctx, installation)
		if err != nil {
			return fmt.Errorf("failed to store installation: %w", err)
		}

		if found && !oldInstallation.IsInstalled() {
			reconciler.bus.Publish(InstallationDetectedEvent, InstallationEvent{Slug: slug})
		} else if !found && oldInstallation.IsInstalled() {
			reconciler.bus.Publish(InstallationRemovedEvent, InstallationEvent{Slug: slug})
		}
	}

	return nil
}
