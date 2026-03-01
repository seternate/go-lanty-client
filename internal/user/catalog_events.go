package user

type CatalogEvent struct {
	IP string
}

type CatalogEventHandler func(e CatalogEvent)

const CatalogAddedEvent string = "user_catalog:added"
const CatalogUpdatedEvent string = "user_catalog:updated"
const CatalogRemovedEvent string = "user_catalog:removed"
