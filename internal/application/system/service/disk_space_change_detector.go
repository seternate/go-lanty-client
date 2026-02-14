package service

import (
	"math"

	eventbus "github.com/asaskevich/EventBus"
	"github.com/seternate/go-lanty-client/internal/application/system/event"
)

type DiskSpaceProvider interface {
	GetFreeDiskSpace() (uint64, error)
}

type DiskSpaceChangeDetector interface {
	DetectChange() (changed bool, freeDiskSpace uint64, err error)
}

var _ DiskSpaceChangeDetector = (*DiskChangeDetector)(nil)

type DiskChangeDetector struct {
	eventbus       eventbus.Bus
	provider       DiskSpaceProvider
	lastFreeBytes  uint64
	thresholdBytes uint64
}

func NewDiskSpaceChangeDetector(bus eventbus.Bus, thresholdBytes uint64, provider DiskSpaceProvider) *DiskChangeDetector {
	return &DiskChangeDetector{
		thresholdBytes: thresholdBytes,
		provider:       provider,
		eventbus:       bus,
	}
}

func (detector *DiskChangeDetector) DetectChange() (bool, uint64, error) {
	currentFreeBytes, err := detector.provider.GetFreeDiskSpace()
	if err != nil {
		return false, 0, err
	}

	lastFreeBytes := detector.lastFreeBytes

	delta := math.Abs(float64(currentFreeBytes - lastFreeBytes))
	if delta <= float64(detector.thresholdBytes) {
		return false, lastFreeBytes, nil
	}

	detector.lastFreeBytes = currentFreeBytes

	detector.eventbus.Publish(event.DiskSpaceChangedEvent, event.DiskSpaceEvent{FreeDiskSpace: currentFreeBytes})

	return true, currentFreeBytes, nil
}
