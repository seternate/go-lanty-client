package gamebrowser

import (
	"fyne.io/fyne/v2/widget"
	"github.com/seternate/go-lanty/pkg/game"
)

func NewStartServerView(game game.Game) (container *widget.Form) {
	executable := widget.NewFormItem("executable", widget.NewLabel(game.ServerExecutable))
	argument := widget.NewFormItem("argument", widget.NewEntry())

	container = widget.NewForm(executable, argument)

	return
}
