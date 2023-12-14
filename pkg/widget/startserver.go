package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/seternate/go-lanty-client/pkg/controller"
	"github.com/seternate/go-lanty/pkg/game"
)

type StartServer struct {
	widget.BaseWidget
	form widget.Form

	controller *controller.Controller
	game       game.Game
	name       binding.String
	executable binding.String
	argument   binding.String
}

func NewStartServer(controller *controller.Controller, parent fyne.CanvasObject) *StartServer {
	startserver := &StartServer{
		controller: controller,
		name:       binding.NewString(),
		executable: binding.NewString(),
		argument:   binding.NewString(),
	}
	startserver.ExtendBaseWidget(startserver)

	startserver.form.SubmitText = "Start"
	startserver.form.OnSubmit = func() {
		startserver.Hide()
		parent.Refresh()
		startserver.game.ServerArgument, _ = startserver.argument.Get()
		controller.Game.StartServer(startserver.game)
	}
	startserver.form.OnCancel = func() {
		startserver.Hide()
		parent.Refresh()
	}

	startserver.form.Append("Name", widget.NewLabelWithData(startserver.name))
	startserver.form.Append("Executable", widget.NewLabelWithData(startserver.executable))
	startserver.form.Append("Argument", widget.NewEntryWithData(startserver.argument))

	return startserver
}

func (widget *StartServer) Update(game game.Game) {
	widget.game = game
	widget.name.Set(game.Name)
	widget.executable.Set(game.ServerExecutable)
	widget.argument.Set(game.ServerArgument)
}

func (widget *StartServer) CreateRenderer() fyne.WidgetRenderer {
	return widget.form.CreateRenderer()
}
