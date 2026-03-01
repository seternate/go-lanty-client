package game

import "errors"

type InstallationActivity int

const (
	InstallationActivityIdle InstallationActivity = iota
	InstallationActivityDownloading
	InstallationActivityExtracting
)

type InstallationResult int

const (
	InstallationResultNone InstallationResult = iota
	InstallationResultCancelled
	InstallationResultFailed
)

type Installation struct {
	Slug       string
	activity   InstallationActivity
	lastResult InstallationResult
	directory  string
}

func NewInstallation(slug string) *Installation {
	return &Installation{
		Slug:       slug,
		activity:   InstallationActivityIdle,
		lastResult: InstallationResultNone,
	}
}

func (i *Installation) IsIdle() bool {
	return i.activity == InstallationActivityIdle
}

func (i *Installation) IsInstalling() bool {
	return !i.IsIdle()
}

func (i *Installation) IsDownloading() bool {
	return i.activity == InstallationActivityDownloading
}

func (i *Installation) IsExtracting() bool {
	return i.activity == InstallationActivityExtracting
}

func (i *Installation) IsInstalled() bool {
	return len(i.directory) > 0
}

func (i *Installation) WasCancelled() bool {
	return i.lastResult == InstallationResultCancelled
}

func (i *Installation) HasFailed() bool {
	return i.lastResult == InstallationResultFailed
}

func (i *Installation) MarkAsDownloading() error {
	if !i.IsIdle() {
		return errors.New("cannot mark as downloading when not idling")
	}

	i.activity = InstallationActivityDownloading
	i.lastResult = InstallationResultNone

	return nil
}

func (i *Installation) MarkAsExtracting() error {
	if !i.IsDownloading() {
		return errors.New("cannot mark as extracting when not downloading")
	}

	i.activity = InstallationActivityExtracting
	i.directory = ""

	return nil
}

func (i *Installation) MarkAsCancelled() error {
	if i.IsIdle() {
		return errors.New("cannot mark as cancelled when idling")
	}

	i.activity = InstallationActivityIdle
	i.lastResult = InstallationResultCancelled

	if i.IsExtracting() {
		i.directory = ""
	}

	return nil
}

func (i *Installation) MarkAsFailed() error {
	if i.IsIdle() {
		return errors.New("cannot mark as failed when idling")
	}

	i.activity = InstallationActivityIdle
	i.lastResult = InstallationResultFailed

	if i.IsExtracting() {
		i.directory = ""
	}

	return nil
}

func (i *Installation) MarkAsCompleted(directory string) error {
	if !i.IsExtracting() {
		return errors.New("cannot mark as completed when not extracting")
	}

	i.activity = InstallationActivityIdle
	i.lastResult = InstallationResultNone
	i.directory = directory

	return nil
}

func (i *Installation) SetDirectory(directory string) {
	i.directory = directory
}

func (i *Installation) Directory() string {
	return i.directory
}

func (i *Installation) ClearDirectory() {
	i.directory = ""
}
