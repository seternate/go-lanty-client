package user

type UserCatalogItem struct {
	IP   string
	Name string
}

func (i UserCatalogItem) Equal(other UserCatalogItem) bool {
	return i.IP == other.IP && i.Name == other.Name
}
