//go:build linux || darwin
// +build linux darwin

package filesystem

import (
	"context"
	"syscall"
)

type UnixDiskSpaceProvider struct{}

func NewDiskSpaceProvider() *UnixDiskSpaceProvider {
	return &UnixDiskSpaceProvider{}
}

func (provider *UnixDiskSpaceProvider) GetFreeDiskSpace(ctx context.Context, path string) (uint64, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return 0, err
	}

	return stat.Bavail * uint64(stat.Bsize), nil
}
