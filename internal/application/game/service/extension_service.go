package service

import (
	"context"
	"path/filepath"

	eventbus "github.com/asaskevich/EventBus"
	"github.com/seternate/go-lanty-client/internal/application/game/event"
	"github.com/seternate/go-lanty-client/internal/domain/game"
)

// ExtensionListProvider returns extensions for a game from the server.
// Downloaded is overwritten by the service using local repository state.
type ExtensionListProvider interface {
	ListForGame(ctx context.Context, slug string) ([]game.Extension, error)
}

// ExtensionUploader uploads a file from the installation folder to the server.
type ExtensionUploader interface {
	Upload(ctx context.Context, slug string, localFilePath string) error
}

// ExtensionDownloader downloads one extension file into destDir (game installation directory).
type ExtensionDownloader interface {
	Download(ctx context.Context, slug string, filename string, destDir string) error
}

type ExtensionService interface {
	ListExtensions(ctx context.Context, slug string) ([]game.Extension, error)
	UploadExtension(ctx context.Context, slug string, localFilePath string) error
	DownloadExtension(ctx context.Context, slug string, filename string) error
	HasNewExtensions(ctx context.Context, slug string) (bool, error)
	ExtensionsVisible(slug string) (bool, error)
}

var _ ExtensionService = (*ExtensionServiceImpl)(nil)

type ExtensionServiceImpl struct {
	bus              eventbus.Bus
	extensionRepo    game.GameExtensionRepository
	installationRepo game.GameInstallationRepository
	installationRoot string
	listProvider     ExtensionListProvider
	uploader         ExtensionUploader
	downloader       ExtensionDownloader
}

func NewExtensionService(
	bus eventbus.Bus,
	extensionRepo game.GameExtensionRepository,
	installationRepo game.GameInstallationRepository,
	installationRoot string,
	listProvider ExtensionListProvider,
	uploader ExtensionUploader,
	downloader ExtensionDownloader,
) *ExtensionServiceImpl {
	return &ExtensionServiceImpl{
		bus:              bus,
		extensionRepo:    extensionRepo,
		installationRepo: installationRepo,
		installationRoot: installationRoot,
		listProvider:     listProvider,
		uploader:         uploader,
		downloader:       downloader,
	}
}

func (s *ExtensionServiceImpl) ListExtensions(ctx context.Context, slug string) ([]game.Extension, error) {
	fromServer, err := s.listProvider.ListForGame(ctx, slug)
	if err != nil {
		return nil, err
	}
	stored, err := s.extensionRepo.ListByGameSlug(slug)
	if err != nil {
		return nil, err
	}
	downloadedByFilename := make(map[string]bool)
	for _, ext := range stored {
		downloadedByFilename[ext.Filename] = ext.Downloaded
	}
	result := make([]game.Extension, 0, len(fromServer))
	for _, ext := range fromServer {
		downloaded := downloadedByFilename[ext.Filename]
		result = append(result, game.Extension{
			GameSlug:   slug,
			Filename:   ext.Filename,
			Name:       ext.Name,
			Size:       ext.Size,
			Downloaded: downloaded,
		})
	}
	return result, nil
}

func (s *ExtensionServiceImpl) UploadExtension(ctx context.Context, slug string, localFilePath string) error {
	if err := s.uploader.Upload(ctx, slug, localFilePath); err != nil {
		return err
	}
	fromServer, err := s.listProvider.ListForGame(ctx, slug)
	if err != nil {
		return err
	}
	stored, _ := s.extensionRepo.ListByGameSlug(slug)
	downloadedByFilename := make(map[string]bool)
	for _, ext := range stored {
		downloadedByFilename[ext.Filename] = ext.Downloaded
	}
	for _, ext := range fromServer {
		downloaded := downloadedByFilename[ext.Filename]
		toStore := game.Extension{
			GameSlug:   slug,
			Filename:   ext.Filename,
			Name:       ext.Name,
			Size:       ext.Size,
			Downloaded: downloaded,
		}
		if err := s.extensionRepo.Store(toStore); err != nil {
			return err
		}
	}
	s.bus.Publish(event.ExtensionListRefreshedEvent, event.ExtensionEvent{Slug: slug})
	return nil
}

func (s *ExtensionServiceImpl) DownloadExtension(ctx context.Context, slug string, filename string) error {
	installation, err := s.installationRepo.GetBySlug(slug)
	if err != nil {
		return err
	}
	if !installation.IsInstalled() {
		return nil
	}
	destDir, err := filepath.Abs(filepath.Join(s.installationRoot, installation.InstallationDirectoryRelative))
	if err != nil {
		return err
	}
	if err := s.downloader.Download(ctx, slug, filename, destDir); err != nil {
		return err
	}
	stored, _ := s.extensionRepo.ListByGameSlug(slug)
	var toStore game.Extension
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
		toStore = game.NewExtension(slug, filename, "", 0, true)
	}
	if err := s.extensionRepo.Store(toStore); err != nil {
		return err
	}
	s.bus.Publish(event.ExtensionDownloadedEvent, event.ExtensionEvent{Slug: slug, Filename: filename})
	return nil
}

func (s *ExtensionServiceImpl) HasNewExtensions(ctx context.Context, slug string) (bool, error) {
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

func (s *ExtensionServiceImpl) ExtensionsVisible(slug string) (bool, error) {
	return s.HasNewExtensions(context.Background(), slug)
}
