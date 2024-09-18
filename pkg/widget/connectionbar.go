package widget

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty-client/pkg/controller"
	"github.com/seternate/go-lanty-client/pkg/setting"
	"github.com/seternate/go-lanty-client/pkg/theme"
)

type Connectionbar struct {
	widget.BaseWidget

	controller    *controller.Controller
	statusupdated chan struct{}
	statustext    string
}

func NewConnectionbar(controller *controller.Controller) *Connectionbar {
	connectionbar := &Connectionbar{
		controller:    controller,
		statusupdated: make(chan struct{}, 50),
		statustext:    "UNKNOWN",
	}
	connectionbar.ExtendBaseWidget(connectionbar)

	controller.Connection.Subscribe(connectionbar.statusupdated)
	controller.WaitGroup().Add(1)
	go connectionbar.statusUpdater()

	return connectionbar
}

func (widget *Connectionbar) statusUpdater() {
	defer widget.controller.WaitGroup().Done()
	widget.updateStatus()
	for {
		select {
		case <-widget.controller.Context().Done():
			log.Trace().Msg("exiting connectionbar statusUpdater()")
			return
		case <-widget.statusupdated:
			widget.updateStatus()
			widget.Refresh()
		}
	}
}

func (widget *Connectionbar) updateStatus() {
	status := widget.controller.Connection.Status
	if status == controller.Connected {
		widget.statustext = fmt.Sprintf("Connected to server: %s", widget.controller.Settings.Settings().ServerURL)
	} else if status == controller.Disconnected {
		widget.statustext = fmt.Sprintf("Error connecting to server: %s", widget.controller.Settings.Settings().ServerURL)
	}
}

func (widget *Connectionbar) CreateRenderer() fyne.WidgetRenderer {
	return newConnectionbarRenderer(widget)
}

type connectionbarRenderer struct {
	widget  *Connectionbar
	line    *canvas.Line
	status  *canvas.Text
	version *canvas.Text
}

func newConnectionbarRenderer(widget *Connectionbar) *connectionbarRenderer {
	renderer := &connectionbarRenderer{
		widget:  widget,
		line:    canvas.NewLine(fynetheme.InputBorderColor()),
		status:  canvas.NewText(widget.statustext, theme.ForegroundColor()),
		version: canvas.NewText(setting.VERSION, theme.ForegroundColor()),
	}
	renderer.line.StrokeWidth = 2.0
	return renderer
}

func (renderer *connectionbarRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{
		renderer.line,
		renderer.status,
		renderer.version,
	}
	return objects
}

func (renderer *connectionbarRenderer) Layout(size fyne.Size) {
	renderer.line.Resize(fyne.NewSize(size.Width, 0))

	versiontextsize := fyne.MeasureText(renderer.version.Text, renderer.version.TextSize, renderer.version.TextStyle)
	statustextsize := fyne.MeasureText(renderer.status.Text, renderer.status.TextSize, renderer.status.TextStyle)

	renderer.version.Move(fyne.NewPos(theme.InnerPadding(), (size.Height-versiontextsize.Height)/2))
	renderer.status.Move(fyne.NewPos(size.Width-statustextsize.Width-theme.InnerPadding(), (size.Height-statustextsize.Height)/2))
}

func (renderer *connectionbarRenderer) MinSize() fyne.Size {
	versiontextsize := fyne.MeasureText(renderer.version.Text, renderer.version.TextSize, renderer.version.TextStyle)
	statustextsize := fyne.MeasureText(renderer.status.Text, renderer.status.TextSize, renderer.status.TextStyle)

	minHeight := fyne.Max(versiontextsize.Height, statustextsize.Height) + theme.InnerPadding()
	minWidth := versiontextsize.Width + statustextsize.Width + 3*theme.InnerPadding()

	return fyne.NewSize(minWidth, minHeight)
}

func (renderer *connectionbarRenderer) Refresh() {
	renderer.line.Refresh()
	renderer.status.Text = renderer.widget.statustext
	renderer.status.Refresh()
	renderer.version.Refresh()
}

func (renderer *connectionbarRenderer) Destroy() {}
