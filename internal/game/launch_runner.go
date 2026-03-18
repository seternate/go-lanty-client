package game

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/seternate/go-lanty-client/internal/user"
)

type ProcessLauncher interface {
	Launch(ctx context.Context, executable string, args []string, workingDir string) error
}

type LaunchSpecSource interface {
	GetLaunchSpec(ctx context.Context, slug string, mode LaunchSpecMode) (LaunchSpec, error)
}

type LaunchRunner struct {
	installationRepo InstallationRepository
	userRepo         user.CatalogRepository
	processLauncher  ProcessLauncher
	specSource       LaunchSpecSource
}

func NewLaunchRunner(installationRepo InstallationRepository, userRepo user.CatalogRepository, processLauncher ProcessLauncher, specSource LaunchSpecSource) *LaunchRunner {
	return &LaunchRunner{
		installationRepo: installationRepo,
		userRepo:         userRepo,
		processLauncher:  processLauncher,
		specSource:       specSource,
	}
}

func (runner *LaunchRunner) GetLaunchSpec(ctx context.Context, slug string, mode LaunchSpecMode) (LaunchSpec, error) {
	launchSpec, err := runner.specSource.GetLaunchSpec(ctx, slug, mode)
	if err != nil {
		return LaunchSpec{}, err
	}

	if mode == LaunchSpecModeJoin {
		users, err := runner.userRepo.GetAll(ctx)
		if err != nil {
			return LaunchSpec{}, fmt.Errorf("failed to get users: %w", err)
		}

		userOptions := make([]EnumValueOption, 0, len(users))
		for _, user := range users {
			userOptions = append(userOptions, EnumValueOption{Value: user.IP, Label: user.Name})
		}

		if len(userOptions) == 0 {
			return launchSpec, nil
		}

		err = launchSpec.UpdateConnectParam(userOptions, userOptions[0].Value)
		if err != nil {
			return LaunchSpec{}, fmt.Errorf("failed to update Connect param: %w", err)
		}
	}

	return launchSpec, nil
}

func (runner *LaunchRunner) StartSingleplayer(ctx context.Context, slug string, launchArgs []LaunchArg) error {
	installation, err := runner.installationRepo.GetBySlug(ctx, slug)
	if err != nil {
		return err
	}

	if !installation.IsInstalled() {
		return fmt.Errorf("game is not installed")
	}

	launchSpec, err := runner.GetLaunchSpec(ctx, slug, LaunchSpecModePlay)
	if err != nil {
		return err
	}

	err = launchSpec.UpdateParams(launchArgs)
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

func (runner *LaunchRunner) JoinMultiplayer(ctx context.Context, slug string, launchArgs []LaunchArg) error {
	installation, err := runner.installationRepo.GetBySlug(ctx, slug)
	if err != nil {
		return err
	}

	if !installation.IsInstalled() {
		return fmt.Errorf("game is not installed")
	}

	launchSpec, err := runner.GetLaunchSpec(ctx, slug, LaunchSpecModeJoin)
	if err != nil {
		return err
	}

	err = launchSpec.UpdateParams(launchArgs)
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

	launchSpec, err := runner.GetLaunchSpec(ctx, slug, LaunchSpecModeHost)
	if err != nil {
		return err
	}

	err = launchSpec.UpdateParams(launchArgs)
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
