package event

type InstallationEvent struct {
	Slug string
}

type InstallationEventHandler func(event InstallationEvent)

const InstallationStartedEvent = "game_installation:started"
const InstallationProgressEvent = "game_installation:progress"
const InstallationFinishedEvent = "game_installation:finished"
const InstallationCanceledEvent = "game_installation:canceled"
const InstallationFailedEvent = "game_installation:failed"
const InstallationDetectedEvent = "game_installation:detected"
const InstallationRemovedEvent = "game_installation:removed"
