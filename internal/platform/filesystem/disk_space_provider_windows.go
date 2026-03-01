//go:build windows
// +build windows

package filesystem

import (
	"context"

	"golang.org/x/sys/windows"
)

type WindowsDiskSpaceProvider struct{}

func NewDiskSpaceProvider() *WindowsDiskSpaceProvider {
	return &WindowsDiskSpaceProvider{}
}

func (provider *WindowsDiskSpaceProvider) GetFreeDiskSpace(ctx context.Context, path string) (uint64, error) {
	var freeBytesAvailable uint64
	var totalBytes uint64
	var totalFreeBytes uint64

	pathPtr, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return 0, err
	}

	err = windows.GetDiskFreeSpaceEx(
		pathPtr,
		&freeBytesAvailable,
		&totalBytes,
		&totalFreeBytes,
	)
	if err != nil {
		return 0, err
	}

	return freeBytesAvailable, nil
}
