package service

import (
	"fmt"
	"path/filepath"
	"sync"

	gameapp "github.com/seternate/go-lanty-client/internal/application/game"
	"github.com/seternate/go-lanty-client/internal/application/setting"
	settingsappsrv "github.com/seternate/go-lanty-client/internal/application/setting/service"
	"github.com/seternate/go-lanty-client/internal/domain/game"
)

type CommandLauncher interface {
	Run(command Command) error
}

type Command struct {
	Executable string
	Args       []string
	WorkingDir string
}

type LaunchSpecMode string

const (
	LaunchSpecModePlay LaunchSpecMode = "play"
	LaunchSpecModeJoin LaunchSpecMode = "join"
	LaunchSpecModeHost LaunchSpecMode = "host"
)

type LaunchSpecSource interface {
	GetLaunchSpec(slug string, mode LaunchSpecMode) (game.GameLaunchSpec, error)
}

type GameLauncherService interface {
	StartSingleplayer(slug string) error
	JoinMultiplayer(slug string, ipAddress string) error
	HostMultiplayer(slug string, launchArgs []gameapp.ArgInput) error
}

var _ GameLauncherService = (*LauncherService)(nil)
var _ settingsappsrv.SettingsListener = (*LauncherService)(nil)

type LauncherService struct {
	installationDirectory string
	installationRepo      game.GameInstallationRepository
	commandLauncher       CommandLauncher
	launchSpecSource      LaunchSpecSource
	mu                    sync.RWMutex
}

func NewLauncherService(installationDirectory string, installationRepo game.GameInstallationRepository, commandLauncher CommandLauncher, launchSpecSource LaunchSpecSource) *LauncherService {
	return &LauncherService{
		installationDirectory: installationDirectory,
		installationRepo:      installationRepo,
		commandLauncher:       commandLauncher,
		launchSpecSource:      launchSpecSource,
	}
}

func (service *LauncherService) StartSingleplayer(slug string) error {
	installation, err := service.installationRepo.GetBySlug(slug)
	if err != nil {
		return err
	}

	if installation.IsNotInstalled() {
		return fmt.Errorf("game is not installed")
	}

	launchSpec, err := service.launchSpecSource.GetLaunchSpec(slug, LaunchSpecModePlay)
	if err != nil {
		return err
	}

	service.mu.RLock()
	installationDirectory := service.installationDirectory
	service.mu.RUnlock()

	absoluteExecutablePath, err := filepath.Abs(filepath.Join(installationDirectory, installation.InstallationDirectoryRelative, launchSpec.ExecutablePathRelative))
	if err != nil {
		return err
	}
	workingDir := filepath.Dir(absoluteExecutablePath)

	args, err := launchSpec.Render()
	if err != nil {
		return err
	}

	command := Command{
		Executable: absoluteExecutablePath,
		Args:       args,
		WorkingDir: workingDir,
	}

	return service.commandLauncher.Run(command)
}

func (service *LauncherService) JoinMultiplayer(slug string, ipAddress string) error {
	installation, err := service.installationRepo.GetBySlug(slug)
	if err != nil {
		return err
	}

	if installation.IsNotInstalled() {
		return fmt.Errorf("game is not installed")
	}

	launchSpec, err := service.launchSpecSource.GetLaunchSpec(slug, LaunchSpecModeJoin)
	if err != nil {
		return err
	}

	err = launchSpec.UpdateParamByName("Connect", ipAddress, true)
	if err != nil {
		return err
	}

	service.mu.RLock()
	installationDirectory := service.installationDirectory
	service.mu.RUnlock()

	absoluteExecutablePath, err := filepath.Abs(filepath.Join(installationDirectory, installation.InstallationDirectoryRelative, launchSpec.ExecutablePathRelative))
	if err != nil {
		return err
	}
	workingDir := filepath.Dir(absoluteExecutablePath)

	args, err := launchSpec.Render()
	if err != nil {
		return err
	}

	command := Command{
		Executable: absoluteExecutablePath,
		Args:       args,
		WorkingDir: workingDir,
	}

	return service.commandLauncher.Run(command)
}

func (service *LauncherService) HostMultiplayer(slug string, launchArgs []gameapp.ArgInput) error {
	installation, err := service.installationRepo.GetBySlug(slug)
	if err != nil {
		return err
	}

	if installation.IsNotInstalled() {
		return fmt.Errorf("game is not installed")
	}

	launchSpec, err := service.launchSpecSource.GetLaunchSpec(slug, LaunchSpecModeHost)
	if err != nil {
		return err
	}

	for _, launchArg := range launchArgs {
		err := launchSpec.UpdateParam(launchArg.Name, launchArg.Argument, launchArg.Value, launchArg.Enabled)
		if err != nil {
			return err
		}
	}

	service.mu.RLock()
	installationDirectory := service.installationDirectory
	service.mu.RUnlock()

	absoluteExecutablePath, err := filepath.Abs(filepath.Join(installationDirectory, installation.InstallationDirectoryRelative, launchSpec.ExecutablePathRelative))
	if err != nil {
		return err
	}
	workingDir := filepath.Dir(absoluteExecutablePath)

	args, err := launchSpec.Render()
	if err != nil {
		return err
	}

	command := Command{
		Executable: absoluteExecutablePath,
		Args:       args,
		WorkingDir: workingDir,
	}

	return service.commandLauncher.Run(command)
}

func (service *LauncherService) ApplySettings(settings setting.Settings) {
	service.mu.Lock()
	defer service.mu.Unlock()
	service.installationDirectory = settings.GameDirectory
}
