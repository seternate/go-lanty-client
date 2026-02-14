package event

type CatalogEvent struct {
	Slug string
}

type CatalogEventHandler func(e CatalogEvent)

const CatalogAddedEvent string = "game_catalog:added"
const CatalogUpdatedEvent string = "game_catalog:updated"
const CatalogRemovedEvent string = "game_catalog:removed"
