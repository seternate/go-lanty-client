package game

type InstallationStatus int

const (
	InstallationStatusNotInstalled InstallationStatus = iota
	InstallationStatusDownloading
	InstallationStatusExtracting
	InstallationStatusInstalled
	InstallationStatusCanceled
	InstallationStatusError
)

type GameInstallation struct {
	Slug                          string
	Status                        InstallationStatus
	InstallationDirectoryRelative string
}

func NewGameInstallation(slug string) *GameInstallation {
	return &GameInstallation{Slug: slug, Status: InstallationStatusNotInstalled}
}

func (i *GameInstallation) IsInstalled() bool {
	return i.Status == InstallationStatusInstalled
}

func (i *GameInstallation) IsInstalling() bool {
	return i.Status == InstallationStatusDownloading || i.Status == InstallationStatusExtracting
}

func (i *GameInstallation) IsDownloading() bool {
	return i.Status == InstallationStatusDownloading
}

func (i *GameInstallation) IsExtracting() bool {
	return i.Status == InstallationStatusExtracting
}

func (i *GameInstallation) IsNotInstalled() bool {
	return i.Status == InstallationStatusNotInstalled
}

func (i *GameInstallation) IsInstallationError() bool {
	return i.Status == InstallationStatusError
}

func (i *GameInstallation) IsInstallationCanceled() bool {
	return i.Status == InstallationStatusCanceled
}

func (i *GameInstallation) MarkInstallationStarted() {
	i.Status = InstallationStatusDownloading
}

func (i *GameInstallation) MarkExtractionStarted() {
	i.Status = InstallationStatusExtracting
	i.InstallationDirectoryRelative = ""
}

func (i *GameInstallation) MarkInstallationFinished(installationDirectoryRelative string) {
	i.Status = InstallationStatusInstalled
	i.InstallationDirectoryRelative = installationDirectoryRelative
}

func (i *GameInstallation) MarkInstallationCanceled() {
	i.Status = InstallationStatusCanceled
}

func (i *GameInstallation) MarkInstallationFailed() {
	i.Status = InstallationStatusError
}

func (i *GameInstallation) MarkInstallationDetected(installationDirectoryRelative string) {
	if !i.IsInstalling() || !i.IsInstallationCanceled() {
		i.Status = InstallationStatusInstalled
	}
	i.InstallationDirectoryRelative = installationDirectoryRelative
}

func (i *GameInstallation) MarkInstallationNotFound() {
	if !i.IsInstalling() || !i.IsInstallationCanceled() {
		i.Status = InstallationStatusNotInstalled
	}
	i.InstallationDirectoryRelative = ""
}

func (i *GameInstallation) InstallationAvailable() bool {
	return len(i.InstallationDirectoryRelative) > 0
}
