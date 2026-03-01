package diskspace

type DiskSpaceEvent struct {
	FreeDiskSpace uint64
}
type DiskSpaceEventHandler func(e DiskSpaceEvent)

const DiskSpaceChangedEvent string = "disk_space:changed"
