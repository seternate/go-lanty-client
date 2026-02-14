package game

type GameCatalogItem struct {
	Slug                       string
	Name                       string
	IconHash                   string
	BlobSize                   uint64
	ExecutableRelativePath     string
	SupportsJoiningMultiplayer bool
	SupportsHostingServer      bool
}

func NewGameCatalogItem(slug string, name string, iconHash string, blobSize uint64, executableRelativePath string, supportsJoiningMultiplayer bool, supportsHostingServer bool) GameCatalogItem {
	return GameCatalogItem{
		Slug:                       slug,
		Name:                       name,
		IconHash:                   iconHash,
		BlobSize:                   blobSize,
		ExecutableRelativePath:     executableRelativePath,
		SupportsJoiningMultiplayer: supportsJoiningMultiplayer,
		SupportsHostingServer:      supportsHostingServer,
	}
}

func (item *GameCatalogItem) Equal(other GameCatalogItem) bool {
	return item.Slug == other.Slug &&
		item.Name == other.Name &&
		item.IconHash == other.IconHash &&
		item.BlobSize == other.BlobSize &&
		item.ExecutableRelativePath == other.ExecutableRelativePath &&
		item.SupportsJoiningMultiplayer == other.SupportsJoiningMultiplayer &&
		item.SupportsHostingServer == other.SupportsHostingServer
}
