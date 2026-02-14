//go:build linux || darwin
// +build linux darwin

package detector

import (
	"sync"
	"syscall"

	"github.com/seternate/go-lanty-client/internal/application/setting"
	settingappsrv "github.com/seternate/go-lanty-client/internal/application/setting/service"
	"github.com/seternate/go-lanty-client/internal/application/system/service"
)

var _ service.DiskSpaceProvider = (*UnixDiskSpaceProvider)(nil)
var _ settingappsrv.SettingsListener = (*UnixDiskSpaceProvider)(nil)

type UnixDiskSpaceProvider struct {
	path string
	mu   sync.RWMutex
}

func NewDiskSpaceProvider(path string) *UnixDiskSpaceProvider {
	return &UnixDiskSpaceProvider{
		path: path,
	}
}

func (provider *UnixDiskSpaceProvider) GetFreeDiskSpace() (uint64, error) {
	provider.mu.RLock()
	path := provider.path
	provider.mu.RUnlock()

	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return 0, err
	}

	return stat.Bavail * uint64(stat.Bsize), nil
}

func (provider *UnixDiskSpaceProvider) ApplySettings(settings setting.Settings) {
	provider.mu.Lock()
	provider.path = settings.GameDirectory
	provider.mu.Unlock()
}
