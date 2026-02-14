package event

type DiskSpaceEvent struct {
	FreeDiskSpace uint64
}
type DiskSpaceEventHandler func(e DiskSpaceEvent)

const DiskSpaceChangedEvent string = "system:disk_space_changed"
