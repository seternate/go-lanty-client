package api

import (
	"sync"

	userappsrv "github.com/seternate/go-lanty-client/internal/application/user/service"
	"github.com/seternate/go-lanty-client/internal/domain/user"
	"github.com/seternate/go-lanty/pkg/api"
)

var _ userappsrv.UserCatalogSource = (*UserAPIClient)(nil)

type UserAPIClient struct {
	apiclient *api.Client

	mu sync.RWMutex
}

func NewUserAPIClient(apiclient *api.Client) *UserAPIClient {
	return &UserAPIClient{
		apiclient: apiclient,
	}
}

func (client *UserAPIClient) FetchCatalog() ([]user.UserCatalogItem, error) {
	panic("not implemented")
}
