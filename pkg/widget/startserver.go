package widget

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/seternate/go-lanty-client/pkg/controller"
	"github.com/seternate/go-lanty-client/pkg/theme"
	"github.com/seternate/go-lanty/pkg/game"
	"github.com/seternate/go-lanty/pkg/game/argument"
)

type StartServer struct {
	widget.BaseWidget
	arguments []*ArgumentWidget
	cancel    *widget.Button
	start     *widget.Button
	reset     *widget.Button

	game       game.Game
	controller *controller.Controller

	OnSubmit func(game game.Game)
	OnCancel func()
}

func NewStartServer(controller *controller.Controller) *StartServer {
	startserver := &StartServer{
		controller: controller,
		cancel:     widget.NewButton("Cancel", nil),
		start:      widget.NewButton("Start", nil),
		reset:      widget.NewButton("Reset", nil),
	}
	startserver.ExtendBaseWidget(startserver)

	startserver.start.Importance = widget.HighImportance
	startserver.start.OnTapped = func() {
		if startserver.OnSubmit != nil {
			startserver.OnSubmit(startserver.game)
		}
	}
	startserver.reset.OnTapped = func() {
		for _, arg := range startserver.arguments {
			arg.Reset()
		}
		controller.Status.Info("Reseted values to defaults", 3*time.Second)
	}
	startserver.cancel.OnTapped = func() {
		if startserver.OnCancel != nil {
			startserver.OnCancel()
		}
	}

	return startserver
}

func (widget *StartServer) SetGame(game game.Game) {
	widget.game = game
	widget.arguments = []*ArgumentWidget{}
	if game.Server.Arguments == nil {
		return
	}
	for _, arg := range game.Server.Arguments.Arguments {
		if arg.GetType() == argument.TYPE_BASE && arg.IsMandatory() {
			continue
		}
		widget.arguments = append(widget.arguments, NewBaseArgumentWidget(arg))
	}
	widget.Refresh()
}

func (w *StartServer) CreateRenderer() fyne.WidgetRenderer {
	return newStartServerRenderer(w)
}

type startServerRenderer struct {
	widget *StartServer
	name   *canvas.Text
}

func newStartServerRenderer(widget *StartServer) *startServerRenderer {
	renderer := &startServerRenderer{
		widget: widget,
		name:   canvas.NewText(widget.game.Name+" -- Start Server", theme.ForegroundColor()),
	}
	renderer.name.TextSize = 18
	renderer.name.TextStyle.Bold = true
	return renderer
}

func (renderer *startServerRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{
		renderer.name,
		renderer.widget.start,
		renderer.widget.reset,
		renderer.widget.cancel,
	}
	for _, arg := range renderer.widget.arguments {
		objects = append(objects, arg)
	}
	return objects
}

func (renderer *startServerRenderer) Layout(size fyne.Size) {
	textsize := fyne.MeasureText(renderer.name.Text, renderer.name.TextSize, renderer.name.TextStyle)
	renderer.widget.start.Resize(fyne.NewSize(100, renderer.widget.start.MinSize().Height))
	renderer.widget.start.Move(fyne.NewPos(size.Width-theme.InnerPadding()-renderer.widget.start.Size().Width, theme.InnerPadding()))
	renderer.widget.reset.Resize(fyne.NewSize(100, renderer.widget.start.MinSize().Height))
	renderer.widget.reset.Move(fyne.NewPos(renderer.widget.start.Position().X-theme.InnerPadding()-renderer.widget.reset.Size().Width, theme.InnerPadding()))
	renderer.widget.cancel.Resize(fyne.NewSize(100, renderer.widget.start.MinSize().Height))
	renderer.widget.cancel.Move(fyne.NewPos(renderer.widget.reset.Position().X-theme.InnerPadding()-renderer.widget.cancel.Size().Width, theme.InnerPadding()))
	renderer.name.Move(fyne.NewPos(theme.InnerPadding(), renderer.widget.start.Position().Y+(renderer.widget.start.Size().Height-textsize.Height)/2))
	previousPosition := fyne.NewPos(theme.InnerPadding(), renderer.widget.start.Position().Y+renderer.widget.start.Size().Height+theme.InnerPadding())
	for _, arg := range renderer.widget.arguments {
		arg.Move(previousPosition)
		arg.Resize(fyne.NewSize(size.Width-2*theme.InnerPadding(), arg.MinSize().Height))
		previousPosition = previousPosition.AddXY(0, arg.Size().Height+theme.InnerPadding())
	}
}

func (renderer *startServerRenderer) MinSize() fyne.Size {
	textsize := fyne.MeasureText(renderer.name.Text, renderer.name.TextSize, renderer.name.TextStyle)
	minsize := fyne.NewSize(2*theme.InnerPadding(), 2*theme.InnerPadding()+textsize.Height)
	for _, arg := range renderer.widget.arguments {
		minsize.Width = fyne.Max(minsize.Width, arg.MinSize().Width)
		minsize.Height = minsize.Height + arg.MinSize().Height + theme.InnerPadding()
	}
	return minsize
}

func (renderer *startServerRenderer) Refresh() {
	renderer.name.Text = renderer.widget.game.Name + " -- Start Server"
	renderer.name.Refresh()
	for _, arg := range renderer.widget.arguments {
		arg.Refresh()
	}
	renderer.widget.start.Refresh()
	renderer.widget.reset.Refresh()
	renderer.widget.cancel.Refresh()
}

func (renderer *startServerRenderer) Destroy() {}
