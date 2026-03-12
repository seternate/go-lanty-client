package backend

import (
	"sync"

	"github.com/seternate/go-lanty-client/internal/user"
	"github.com/seternate/go-lanty/pkg/api"
)

type UserAPIClient struct {
	apiclient *api.Client

	mu sync.RWMutex
}

func NewUserAPIClient(apiclient *api.Client) *UserAPIClient {
	return &UserAPIClient{
		apiclient: apiclient,
	}
}

func (client *UserAPIClient) FetchCatalog() ([]user.CatalogItem, error) {
	panic("not implemented")
}
