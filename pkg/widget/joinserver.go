package widget

import (
	"fyne.io/fyne/v2/widget"
	"github.com/seternate/go-lanty/pkg/game"
)

type JoinServer struct {
	widget.Form
	game game.Game
}

func NewJoinServer() *JoinServer {
	serverstart := &JoinServer{}
	serverstart.ExtendBaseWidget(serverstart)

	serverstart.Form.Append("TODO", widget.NewLabel("TODO"))

	return serverstart
}

func (widget *JoinServer) Update(game game.Game) {
	widget.game = game
}
