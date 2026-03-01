package game

import "reflect"

type CatalogItem struct {
	Slug              string
	Name              string
	IconHash          string
	InstalledFileSize uint64
	DetectionHints    InstallationDetectionHints
	Capabilities      Capabilities
}

func (c CatalogItem) Equal(other CatalogItem) bool {
	return c.Slug == other.Slug &&
		c.Name == other.Name &&
		c.IconHash == other.IconHash &&
		c.InstalledFileSize == other.InstalledFileSize &&
		c.DetectionHints.Equal(other.DetectionHints) &&
		c.Capabilities == other.Capabilities
}

type InstallationDetectionHints struct {
	FilePaths   []string
	FolderPaths []string
}

func (h InstallationDetectionHints) Equal(other InstallationDetectionHints) bool {
	return reflect.DeepEqual(h.FilePaths, other.FilePaths) &&
		reflect.DeepEqual(h.FolderPaths, other.FolderPaths)
}

type Capabilities struct {
	JoiningMultiplayer bool
	HostingServer      bool
}
