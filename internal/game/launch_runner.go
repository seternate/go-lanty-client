package game

import (
	"context"
	"fmt"
	"path/filepath"
)

type ProcessLauncher interface {
	Launch(ctx context.Context, executable string, args []string, workingDir string) error
}

type LaunchSpecSource interface {
	GetLaunchSpec(ctx context.Context, slug string, mode LaunchSpecMode) (LaunchSpec, error)
}

type LaunchRunner struct {
	installationRepo InstallationRepository
	processLauncher  ProcessLauncher
	specSource       LaunchSpecSource
}

func NewLaunchRunner(installationRepo InstallationRepository, processLauncher ProcessLauncher, specSource LaunchSpecSource) *LaunchRunner {
	return &LaunchRunner{
		installationRepo: installationRepo,
		processLauncher:  processLauncher,
		specSource:       specSource,
	}
}

func (runner *LaunchRunner) GetLaunchSpec(ctx context.Context, slug string, mode LaunchSpecMode) (LaunchSpec, error) {
	return runner.specSource.GetLaunchSpec(ctx, slug, mode)
}

func (runner *LaunchRunner) StartSingleplayer(ctx context.Context, slug string) error {
	installation, err := runner.installationRepo.GetBySlug(ctx, slug)
	if err != nil {
		return err
	}

	if !installation.IsInstalled() {
		return fmt.Errorf("game is not installed")
	}

	launchSpec, err := runner.specSource.GetLaunchSpec(ctx, slug, LaunchSpecModePlay)
	if err != nil {
		return err
	}

	absoluteExecutablePath, err := filepath.Abs(filepath.Join(installation.Directory(), launchSpec.ExecutablePathRelative))
	if err != nil {
		return err
	}
	workingDir := filepath.Dir(absoluteExecutablePath)

	args, err := launchSpec.Render()
	if err != nil {
		return err
	}

	return runner.processLauncher.Launch(ctx, absoluteExecutablePath, args, workingDir)
}

func (runner *LaunchRunner) JoinMultiplayer(ctx context.Context, slug string, ipAddress string) error {
	installation, err := runner.installationRepo.GetBySlug(ctx, slug)
	if err != nil {
		return err
	}

	if !installation.IsInstalled() {
		return fmt.Errorf("game is not installed")
	}

	launchSpec, err := runner.specSource.GetLaunchSpec(ctx, slug, LaunchSpecModeJoin)
	if err != nil {
		return err
	}

	err = launchSpec.UpdateParamByName("Connect", ipAddress, true)
	if err != nil {
		return err
	}

	absoluteExecutablePath, err := filepath.Abs(filepath.Join(installation.Directory(), launchSpec.ExecutablePathRelative))
	if err != nil {
		return err
	}
	workingDir := filepath.Dir(absoluteExecutablePath)

	args, err := launchSpec.Render()
	if err != nil {
		return err
	}

	return runner.processLauncher.Launch(ctx, absoluteExecutablePath, args, workingDir)
}

func (runner *LaunchRunner) HostMultiplayer(ctx context.Context, slug string, launchArgs []LaunchArg) error {
	installation, err := runner.installationRepo.GetBySlug(ctx, slug)
	if err != nil {
		return err
	}

	if !installation.IsInstalled() {
		return fmt.Errorf("game is not installed")
	}

	launchSpec, err := runner.specSource.GetLaunchSpec(ctx, slug, LaunchSpecModeHost)
	if err != nil {
		return err
	}

	for _, launchArg := range launchArgs {
		err := launchSpec.UpdateParam(launchArg.Name, launchArg.Argument, launchArg.Value, launchArg.Enabled)
		if err != nil {
			return err
		}
	}

	absoluteExecutablePath, err := filepath.Abs(filepath.Join(installation.Directory(), launchSpec.ExecutablePathRelative))
	if err != nil {
		return err
	}
	workingDir := filepath.Dir(absoluteExecutablePath)

	args, err := launchSpec.Render()
	if err != nil {
		return err
	}

	return runner.processLauncher.Launch(ctx, absoluteExecutablePath, args, workingDir)
}
