package api

import (
	"maps"
	"slices"
	"sync"

	userappsrv "github.com/seternate/go-lanty-client/internal/application/user/service"
	"github.com/seternate/go-lanty-client/internal/domain/user"
)

var _ userappsrv.UserCatalogSource = (*userMockAPIClient)(nil)

type userMockAPIClient struct {
	mu                 sync.RWMutex
	catalog            map[string]user.UserCatalogItem
	alternatingPresent bool
}

func NewUserMockAPIClient() *userMockAPIClient {
	client := &userMockAPIClient{
		catalog: make(map[string]user.UserCatalogItem),
	}
	client.seedExampleUsers()
	return client
}

func (client *userMockAPIClient) FetchCatalog() ([]user.UserCatalogItem, error) {
	client.mu.Lock()
	defer client.mu.Unlock()
	client.toggleAlternatingUsersLocked()

	return slices.Collect(maps.Values(client.catalog)), nil
}

func (client *userMockAPIClient) toggleAlternatingUsersLocked() {
	if client.alternatingPresent {
		delete(client.catalog, "192.168.1.100")
		delete(client.catalog, "192.168.1.101")
		client.alternatingPresent = false
		return
	}

	client.catalog["192.168.1.100"] = user.UserCatalogItem{
		IP:   "192.168.1.100",
		Name: "DDDD",
	}
	client.catalog["192.168.1.101"] = user.UserCatalogItem{
		IP:   "192.168.1.101",
		Name: "ZZZZ",
	}
	client.alternatingPresent = true
}

func (client *userMockAPIClient) seedExampleUsers() {
	exampleUsers := []user.UserCatalogItem{
		createExampleUser1(),
		createExampleUser2(),
		createExampleUser3(),
		createExampleUser4(),
		createExampleUser5(),
		createExampleUser6(),
		createExampleUser7(),
		createExampleUser8(),
		createExampleUser9(),
		createExampleUser10(),
		createExampleUser11(),
		createExampleUser12(),
		createExampleUser13(),
	}

	for _, user := range exampleUsers {
		client.catalog[user.IP] = user
	}
}

// createExampleUser1 creates the first example user
func createExampleUser1() user.UserCatalogItem {
	return user.UserCatalogItem{
		IP:   "192.168.1.10",
		Name: "Alice",
	}
}

// createExampleUser2 creates the second example user
func createExampleUser2() user.UserCatalogItem {
	return user.UserCatalogItem{
		IP:   "192.168.1.11",
		Name: "Bob",
	}
}

// createExampleUser3 creates the third example user
func createExampleUser3() user.UserCatalogItem {
	return user.UserCatalogItem{
		IP:   "192.168.1.12",
		Name: "Charlie",
	}
}

// createExampleUser4 creates the fourth example user
func createExampleUser4() user.UserCatalogItem {
	return user.UserCatalogItem{
		IP:   "192.168.1.13",
		Name: "Diana",
	}
}

// createExampleUser5 creates the fifth example user
func createExampleUser5() user.UserCatalogItem {
	return user.UserCatalogItem{
		IP:   "192.168.1.14",
		Name: "Eve",
	}
}

// createExampleUser6 creates the sixth example user
func createExampleUser6() user.UserCatalogItem {
	return user.UserCatalogItem{
		IP:   "192.168.1.15",
		Name: "Frank",
	}
}

// createExampleUser7 creates the seventh example user
func createExampleUser7() user.UserCatalogItem {
	return user.UserCatalogItem{
		IP:   "192.168.1.16",
		Name: "Grace",
	}
}

// createExampleUser8 creates the eighth example user
func createExampleUser8() user.UserCatalogItem {
	return user.UserCatalogItem{
		IP:   "192.168.1.17",
		Name: "Henry",
	}
}

// createExampleUser9 creates the ninth example user
func createExampleUser9() user.UserCatalogItem {
	return user.UserCatalogItem{
		IP:   "192.168.1.18",
		Name: "Iris",
	}
}

// createExampleUser10 creates the tenth example user
func createExampleUser10() user.UserCatalogItem {
	return user.UserCatalogItem{
		IP:   "192.168.1.19",
		Name: "Jack",
	}
}

// createExampleUser11 creates the eleventh example user
func createExampleUser11() user.UserCatalogItem {
	return user.UserCatalogItem{
		IP:   "192.168.1.20",
		Name: "Kate",
	}
}

// createExampleUser12 creates the twelfth example user
func createExampleUser12() user.UserCatalogItem {
	return user.UserCatalogItem{
		IP:   "192.168.1.21",
		Name: "Liam",
	}
}

// createExampleUser13 creates the thirteenth example user
func createExampleUser13() user.UserCatalogItem {
	return user.UserCatalogItem{
		IP:   "192.168.1.22",
		Name: "Mia",
	}
}
