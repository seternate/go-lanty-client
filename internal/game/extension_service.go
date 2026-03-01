package game

import (
	"context"

	eventbus "github.com/asaskevich/EventBus"
)

type ExtensionListProvider interface {
	ListForGame(ctx context.Context, slug string) ([]Extension, error)
}

type ExtensionUploader interface {
	Upload(ctx context.Context, slug string, localFilePath string) error
}

type ExtensionDownloader interface {
	Download(ctx context.Context, slug string, filename string, destDir string) error
}

type ExtensionService struct {
	bus              eventbus.Bus
	extensionRepo    ExtensionRepository
	installationRepo InstallationRepository
	listProvider     ExtensionListProvider
	uploader         ExtensionUploader
	downloader       ExtensionDownloader
}

func NewExtensionService(
	bus eventbus.Bus,
	extensionRepo ExtensionRepository,
	installationRepo InstallationRepository,
	listProvider ExtensionListProvider,
	uploader ExtensionUploader,
	downloader ExtensionDownloader,
) *ExtensionService {
	return &ExtensionService{
		bus:              bus,
		extensionRepo:    extensionRepo,
		installationRepo: installationRepo,
		listProvider:     listProvider,
		uploader:         uploader,
		downloader:       downloader,
	}
}

func (s *ExtensionService) ListExtensions(ctx context.Context, slug string) ([]Extension, error) {
	fromServer, err := s.listProvider.ListForGame(ctx, slug)
	if err != nil {
		return nil, err
	}
	stored, err := s.extensionRepo.ListBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	downloadedByFilename := make(map[string]bool)
	for _, ext := range stored {
		downloadedByFilename[ext.Filename] = ext.Downloaded
	}
	result := make([]Extension, 0, len(fromServer))
	for _, ext := range fromServer {
		downloaded := downloadedByFilename[ext.Filename]
		result = append(result, Extension{
			Slug:       slug,
			Filename:   ext.Filename,
			Name:       ext.Name,
			Size:       ext.Size,
			Downloaded: downloaded,
		})
	}
	return result, nil
}

func (s *ExtensionService) UploadExtension(ctx context.Context, slug string, localFilePath string) error {
	if err := s.uploader.Upload(ctx, slug, localFilePath); err != nil {
		return err
	}
	fromServer, err := s.listProvider.ListForGame(ctx, slug)
	if err != nil {
		return err
	}
	stored, _ := s.extensionRepo.ListBySlug(ctx, slug)
	downloadedByFilename := make(map[string]bool)
	for _, ext := range stored {
		downloadedByFilename[ext.Filename] = ext.Downloaded
	}
	for _, ext := range fromServer {
		downloaded := downloadedByFilename[ext.Filename]
		toStore := Extension{
			Slug:       slug,
			Filename:   ext.Filename,
			Name:       ext.Name,
			Size:       ext.Size,
			Downloaded: downloaded,
		}
		if err := s.extensionRepo.Store(ctx, toStore); err != nil {
			return err
		}
	}
	s.bus.Publish(ExtensionListRefreshedEvent, ExtensionEvent{Slug: slug})
	return nil
}

func (s *ExtensionService) DownloadExtension(ctx context.Context, slug string, filename string) error {
	installation, err := s.installationRepo.GetBySlug(ctx, slug)
	if err != nil {
		return err
	}
	if !installation.IsInstalled() {
		return nil
	}
	if err := s.downloader.Download(ctx, slug, filename, installation.Directory()); err != nil {
		return err
	}
	stored, _ := s.extensionRepo.ListBySlug(ctx, slug)
	var toStore Extension
	found := false
	for _, ext := range stored {
		if ext.Filename == filename {
			toStore = ext
			toStore.MarkDownloaded()
			found = true
			break
		}
	}
	if !found {
		toStore = NewExtension(slug, filename, "", 0, true)
	}
	if err := s.extensionRepo.Store(ctx, toStore); err != nil {
		return err
	}
	s.bus.Publish(ExtensionDownloadedEvent, ExtensionEvent{Slug: slug, Filename: filename})
	return nil
}

func (s *ExtensionService) HasNewExtensions(ctx context.Context, slug string) (bool, error) {
	list, err := s.ListExtensions(ctx, slug)
	if err != nil {
		return false, err
	}
	for _, ext := range list {
		if !ext.Downloaded {
			return true, nil
		}
	}
	return false, nil
}

func (s *ExtensionService) ExtensionsVisible(slug string) (bool, error) {
	return s.HasNewExtensions(context.Background(), slug)
}
