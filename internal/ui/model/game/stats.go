package model

import (
	"fyne.io/fyne/v2/data/binding"
	"github.com/dustin/go-humanize"
	viewmodel "github.com/seternate/go-lanty-client/internal/ui/model"
)

var _ viewmodel.ChangeNotifier = (*GameStatsModel)(nil)

type GameStatsModel struct {
	availableGames binding.Int
	installedGames binding.Int
	freeDiskSpace  binding.String
	StatsText      binding.String
}

func NewGameStatsModel() *GameStatsModel {
	stats := &GameStatsModel{
		availableGames: binding.NewInt(),
		installedGames: binding.NewInt(),
		freeDiskSpace:  binding.NewString(),
	}

	stats.StatsText = binding.NewSprintf("%d available, %d installed | %s free disk space", stats.availableGames, stats.installedGames, stats.freeDiskSpace)

	return stats
}

func (m *GameStatsModel) SetAvailableGames(count int) {
	m.availableGames.Set(count)
}

func (m *GameStatsModel) SetInstalledGames(count int) {
	m.installedGames.Set(count)
}

func (m *GameStatsModel) SetFreeDiskSpace(space uint64) {
	m.freeDiskSpace.Set(humanize.SIWithDigits(float64(space), 2, "B"))
}

func (m *GameStatsModel) AddChangeListener(fn func()) {
	l := binding.NewDataListener(fn)
	m.StatsText.AddListener(l)
}
