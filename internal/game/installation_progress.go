package game

import (
	"math"
	"time"
)

type InstallationProgress struct {
	Slug      string
	Completed uint64
	Total     uint64
	StartedAt time.Time
	EndedAt   time.Time
	Failed    bool
	Succeeded bool
}

func (p *InstallationProgress) HasFinished() bool {
	return !p.EndedAt.IsZero()
}

func (p *InstallationProgress) Progress() float64 {
	if p.Total == 0 {
		return 0
	}

	if p.Completed > p.Total {
		return 1
	}

	return float64(p.Completed) / float64(p.Total)
}

func (p *InstallationProgress) Speed() uint64 {
	speed := uint64(0)

	if p.StartedAt.IsZero() {
		return speed
	}

	if p.EndedAt.IsZero() {
		speed = uint64(math.Round(float64(p.Completed) / time.Since(p.StartedAt).Seconds()))
	} else {
		speed = uint64(math.Round(float64(p.Completed) / p.EndedAt.Sub(p.StartedAt).Seconds()))
	}

	return speed
}
