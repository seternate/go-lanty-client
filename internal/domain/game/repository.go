package game

type GameCatalogRepository interface {
	GetAll() ([]GameCatalogItem, error)
	GetBySlug(slug string) (GameCatalogItem, error)
	ReplaceAll(metadata ...GameCatalogItem) error
}

type GameInstallationRepository interface {
	GetAll() ([]GameInstallation, error)
	GetBySlug(slug string) (GameInstallation, error)
	Store(installation GameInstallation) error
	Remove(slug string) error
}
