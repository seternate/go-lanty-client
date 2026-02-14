//go:build windows
// +build windows

package detector

import (
	"sync"

	"github.com/seternate/go-lanty-client/internal/application/setting"
	settingappsrv "github.com/seternate/go-lanty-client/internal/application/setting/service"
	"github.com/seternate/go-lanty-client/internal/application/system/service"
	"golang.org/x/sys/windows"
)

var _ service.DiskSpaceProvider = (*WindowsDiskSpaceProvider)(nil)
var _ settingappsrv.SettingsListener = (*WindowsDiskSpaceProvider)(nil)

type WindowsDiskSpaceProvider struct {
	path string
	mu   sync.RWMutex
}

func NewDiskSpaceProvider(path string) *WindowsDiskSpaceProvider {
	return &WindowsDiskSpaceProvider{
		path: path,
	}
}

func (provider *WindowsDiskSpaceProvider) GetFreeDiskSpace() (uint64, error) {
	provider.mu.RLock()
	path := provider.path
	provider.mu.RUnlock()

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

func (provider *WindowsDiskSpaceProvider) ApplySettings(settings setting.Settings) {
	provider.mu.Lock()
	provider.path = settings.GameDirectory
	provider.mu.Unlock()
}
