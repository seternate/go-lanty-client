package game

import (
	"context"
)

type FileExplorerOpener interface {
	OpenFileExplorer(path string) error
}
type InstallationDirectoryOpener struct {
	installationRepo InstallationRepository
	explorerOpener   FileExplorerOpener
}

func NewInstallationDirectoryOpener(installationRepo InstallationRepository, explorerOpener FileExplorerOpener) *InstallationDirectoryOpener {
	return &InstallationDirectoryOpener{
		installationRepo: installationRepo,
		explorerOpener:   explorerOpener,
	}
}

func (opener *InstallationDirectoryOpener) OpenDirectory(ctx context.Context, slug string) error {
	installation, err := opener.installationRepo.GetBySlug(ctx, slug)
	if err != nil {
		return err
	}

	err = opener.explorerOpener.OpenFileExplorer(installation.Directory())
	if err != nil {
		return err
	}

	return nil
}
