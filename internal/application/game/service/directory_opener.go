package service

import (
	"path/filepath"
	"sync"

	"github.com/seternate/go-lanty-client/internal/application/setting"
	settingsappsrv "github.com/seternate/go-lanty-client/internal/application/setting/service"
	"github.com/seternate/go-lanty-client/internal/domain/game"
)

type GameExplorerOpener interface {
	OpenFileExplorer(path string) error
}

type GameDirectoryOpener interface {
	OpenDirectory(slug string) error
}

var _ GameDirectoryOpener = (*DirectoryOpener)(nil)
var _ settingsappsrv.SettingsListener = (*DirectoryOpener)(nil)

type DirectoryOpener struct {
	installationDirectory string
	installationRepo      game.GameInstallationRepository
	explorerOpener        GameExplorerOpener

	mu sync.RWMutex
}

func NewGameDirectoryOpener(installationDirectory string, installationRepo game.GameInstallationRepository, explorerOpener GameExplorerOpener) *DirectoryOpener {
	return &DirectoryOpener{
		installationDirectory: installationDirectory,
		installationRepo:      installationRepo,
		explorerOpener:        explorerOpener,
	}
}

func (service *DirectoryOpener) OpenDirectory(slug string) error {
	installation, err := service.installationRepo.GetBySlug(slug)
	if err != nil {
		return err
	}

	service.mu.RLock()
	installDir := service.installationDirectory
	service.mu.RUnlock()

	installationPathRelative := filepath.Join(installDir, installation.InstallationDirectoryRelative)
	installationPathAbsolute, err := filepath.Abs(installationPathRelative)
	if err != nil {
		return err
	}

	err = service.explorerOpener.OpenFileExplorer(installationPathAbsolute)
	if err != nil {
		return err
	}

	return nil
}

func (service *DirectoryOpener) ApplySettings(settings setting.Settings) {
	service.mu.Lock()
	defer service.mu.Unlock()
	service.installationDirectory = settings.GameDirectory
}
