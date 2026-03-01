package gameviewmodel

// import (
// 	"fyne.io/fyne/v2/data/binding"
// 	"github.com/dustin/go-humanize"
// )

// type GameStatsModel struct {
// 	availableGames binding.Int
// 	installedGames binding.Int
// 	freeDiskSpace  binding.String
// 	StatsText      binding.String
// }

// func NewGameStatsModel() *GameStatsModel {
// 	stats := &GameStatsModel{
// 		availableGames: binding.NewInt(),
// 		installedGames: binding.NewInt(),
// 		freeDiskSpace:  binding.NewString(),
// 	}

// 	stats.StatsText = binding.NewSprintf("%d available, %d installed | %s free disk space", stats.availableGames, stats.installedGames, stats.freeDiskSpace)

// 	return stats
// }

// func (m *GameStatsModel) SetAvailableGames(count int) {
// 	m.availableGames.Set(count)
// }

// func (m *GameStatsModel) SetInstalledGames(count int) {
// 	m.installedGames.Set(count)
// }

// func (m *GameStatsModel) SetFreeDiskSpace(space uint64) {
// 	m.freeDiskSpace.Set(humanize.SIWithDigits(float64(space), 2, "B"))
// }
