package diskspace

import (
	"context"
	"fmt"
	"math"

	eventbus "github.com/asaskevich/EventBus"
)

type Provider interface {
	GetFreeDiskSpace(ctx context.Context, path string) (uint64, error)
}

type Monitor struct {
	eventbus       eventbus.Bus
	provider       Provider
	lastFreeBytes  uint64
	thresholdBytes uint64
}

func NewMonitor(bus eventbus.Bus, thresholdBytes uint64, provider Provider) *Monitor {
	return &Monitor{
		thresholdBytes: thresholdBytes,
		provider:       provider,
		eventbus:       bus,
	}
}

func (monitor *Monitor) Detect(ctx context.Context, path string) error {
	currentFreeBytes, err := monitor.provider.GetFreeDiskSpace(ctx, path)
	if err != nil {
		return fmt.Errorf("failed to get free disk space: %w", err)
	}

	lastFreeBytes := monitor.lastFreeBytes

	delta := math.Abs(float64(currentFreeBytes - lastFreeBytes))
	if delta <= float64(monitor.thresholdBytes) {
		return nil
	}

	monitor.lastFreeBytes = currentFreeBytes

	monitor.eventbus.Publish(DiskSpaceChangedEvent, DiskSpaceEvent{FreeDiskSpace: currentFreeBytes})

	return nil
}
