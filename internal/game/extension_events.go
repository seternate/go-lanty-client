package game

type ExtensionEvent struct {
	Slug     string
	Filename string
}

type ExtensionEventHandler func(event ExtensionEvent)

const ExtensionListRefreshedEvent = "game_extension:list_refreshed"
const ExtensionDownloadedEvent = "game_extension:downloaded"
