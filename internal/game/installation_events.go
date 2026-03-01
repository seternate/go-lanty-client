package game

type InstallationEvent struct {
	Slug string
}

type InstallationEventHandler func(event InstallationEvent)

const InstallationStartedEvent = "game_installation:started"
const InstallationProgressedEvent = "game_installation:progressed"
const InstallationFinishedEvent = "game_installation:finished"
const InstallationCancelledEvent = "game_installation:cancelled"
const InstallationFailedEvent = "game_installation:failed"
const InstallationDetectedEvent = "game_installation:detected"
const InstallationRemovedEvent = "game_installation:removed"
