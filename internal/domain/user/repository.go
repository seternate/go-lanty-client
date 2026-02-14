package user

type UserCatalogRepository interface {
	GetAll() ([]UserCatalogItem, error)
	GetByIP(ip string) (UserCatalogItem, error)
	ReplaceAll(metadata ...UserCatalogItem) error
}
