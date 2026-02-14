package game

import "time"

type InstallationProgress struct {
	Slug      string
	Completed uint64
	Total     uint64
	StartedAt time.Time
	EndedAt   time.Time
	Failed    bool
	Succeeded bool
}
