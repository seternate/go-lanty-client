package widget

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/seternate/go-lanty-client/pkg/controller"
	"github.com/seternate/go-lanty-client/pkg/theme"
)

func StatusIcon(level controller.StatusLevel) fyne.Resource {
	switch level {
	case controller.StatusLevelInfo:
		return fynetheme.InfoIcon()
	case controller.StatusLevelWarning:
		return fynetheme.WarningIcon()
	case controller.StatusLevelError:
		return fynetheme.CancelIcon()
	}
	return fynetheme.BrokenImageIcon()
}

type StatusBar struct {
	widget.BaseWidget

	controller *controller.Controller
	icon       fyne.Resource
	text       string

	timer *time.Timer
}

func NewStatusBar(controller *controller.Controller) (statusbar *StatusBar) {
	statusbar = &StatusBar{
		controller: controller,
		icon:       fynetheme.InfoIcon(),
		timer:      time.NewTimer(0),
		text:       "TEST",
	}
	statusbar.ExtendBaseWidget(statusbar)
	//In order to Hide() the widget at init a call of Refresh() of its parent is needed that it is redrawn to the canvas
	//(see https://github.com/fyne-io/fyne/issues/4494)
	statusbar.Hide()
	//drain timer for first time
	<-statusbar.timer.C

	return
}

func (widget *StatusBar) ShowStatus(status controller.Status) {
	widget.text = status.Text
	widget.icon = StatusIcon(status.Level)
	widget.Refresh()
	widget.Show()
	if !widget.timer.Reset(status.Duration) {
		go func() {
			<-widget.timer.C
			widget.Hide()
		}()
	}
}

func (widget *StatusBar) CreateRenderer() fyne.WidgetRenderer {
	return newStatusBarRenderer(widget)
}

type statusBarRenderer struct {
	widget     *StatusBar
	background *canvas.Rectangle
	icon       *canvas.Image
	text       *canvas.Text
}

func newStatusBarRenderer(widget *StatusBar) *statusBarRenderer {
	renderer := &statusBarRenderer{
		widget:     widget,
		background: canvas.NewRectangle(fynetheme.InputBackgroundColor()),
		icon:       canvas.NewImageFromResource(widget.icon),
		text:       canvas.NewText(widget.text, theme.ForegroundColor()),
	}
	renderer.background.CornerRadius = fynetheme.SelectionRadiusSize()
	return renderer
}

func (renderer *statusBarRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		renderer.background,
		renderer.icon,
		renderer.text,
	}
}

func (renderer *statusBarRenderer) Layout(size fyne.Size) {
	renderer.background.Move(fyne.NewPos(0, 0))
	renderer.background.Resize(size)
	renderer.icon.Resize(fyne.NewSize(size.Height-2*theme.InnerPadding(), size.Height-2*theme.InnerPadding()))
	renderer.icon.Move(fyne.NewPos(theme.InnerPadding(), (size.Height-renderer.icon.Size().Height)/2))
	textsize := fyne.MeasureText(renderer.text.Text, renderer.text.TextSize, renderer.text.TextStyle)
	renderer.text.Move(fyne.NewPos(renderer.icon.Position().X+renderer.icon.Size().Width+theme.InnerPadding(), (size.Height-textsize.Height)/2))
}

func (renderer *statusBarRenderer) MinSize() fyne.Size {
	minsize := fyne.MeasureText(renderer.text.Text, renderer.text.TextSize, renderer.text.TextStyle)
	if minsize.Height == 0 {
		minsize.Add(fyne.NewSize(0, 12))
	}
	minsize.Add(fyne.NewSize(3*theme.InnerPadding()+minsize.Height, 2*theme.InnerPadding()))
	return minsize
}

func (renderer *statusBarRenderer) Refresh() {
	renderer.text.Text = renderer.widget.text
	renderer.icon = canvas.NewImageFromResource(renderer.widget.icon)
	renderer.background.Refresh()
	renderer.icon.Refresh()
	renderer.text.Refresh()
}

func (renderer *statusBarRenderer) Destroy() {}
