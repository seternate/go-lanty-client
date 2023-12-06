package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/seternate/go-lanty-client/pkg/theme"
	"github.com/seternate/go-lanty/pkg/game/argument"
)

type ArgumentWidget struct {
	widget.BaseWidget
	Name     string
	Disabled *widget.Check
	Data     fyne.CanvasObject
	Argument argument.Argument
}

func NewBaseArgumentWidget(arg argument.Argument) *ArgumentWidget {
	item := &ArgumentWidget{
		Argument: arg,
		Name:     arg.GetName(),
	}
	item.ExtendBaseWidget(item)

	item.Disabled = widget.NewCheck("", func(b bool) {
		if b {
			item.Argument.Enable()
		} else {
			item.Argument.Disable()
		}
		item.Refresh()
	})
	item.Disabled.SetChecked(!item.Argument.IsDisabled())
	if item.Argument.IsMandatory() {
		item.Disabled.Hide()
	}

	switch item.Argument.GetType() {
	case argument.TYPE_STRING:
		item.Data = NewStringArgument(item.Argument.(*argument.String))
	case argument.TYPE_BOOLEAN:
		item.Data = NewBooleanArgument(item.Argument.(*argument.Boolean))
	case argument.TYPE_INTEGER:
		item.Data = NewIntegerArgument(item.Argument.(*argument.Integer))
	case argument.TYPE_FLOAT:
		item.Data = NewFloatArgument(item.Argument.(*argument.Float))
	case argument.TYPE_ENUM:
		item.Data = NewEnumArgument(item.Argument.(*argument.Enum))
	default:
		item.Data = nil
	}

	return item
}

func (item *ArgumentWidget) Reset() {
	item.Argument.Reset()
	item.Refresh()
}

func (item *ArgumentWidget) CreateRenderer() fyne.WidgetRenderer {
	return newargumentRenderer(item)
}

type argumentRenderer struct {
	widget     *ArgumentWidget
	background *canvas.Rectangle
	name       *canvas.Text
}

func newargumentRenderer(widget *ArgumentWidget) *argumentRenderer {
	renderer := &argumentRenderer{
		widget:     widget,
		background: canvas.NewRectangle(fynetheme.InputBorderColor()),
		name:       canvas.NewText(widget.Name, theme.ForegroundColor()),
	}
	renderer.name.TextStyle.Bold = true
	renderer.background.CornerRadius = fynetheme.SelectionRadiusSize()
	return renderer
}

func (renderer *argumentRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{
		renderer.background,
		renderer.name,
		renderer.widget.Disabled,
	}
	if renderer.widget.Data != nil {
		objects = append(objects, renderer.widget.Data)
	}
	return objects
}

func (renderer *argumentRenderer) Layout(size fyne.Size) {
	renderer.background.Resize(size)
	textsize := fyne.MeasureText(renderer.name.Text, renderer.name.TextSize, renderer.name.TextStyle)
	renderer.name.Move(fyne.NewPos(theme.InnerPadding(), theme.InnerPadding()))
	renderer.widget.Disabled.Resize(renderer.widget.Disabled.MinSize())
	renderer.widget.Disabled.Move(fyne.NewPos(textsize.Width+theme.InnerPadding(), textsize.Height/2-renderer.widget.Disabled.MinSize().Height/2+theme.InnerPadding()))
	if renderer.widget.Data != nil {
		renderer.widget.Data.Move(fyne.NewPos(3*theme.InnerPadding(), textsize.Height+2*theme.InnerPadding()))
		renderer.widget.Data.Resize(fyne.NewSize(size.Width-4*theme.InnerPadding(), renderer.widget.Data.MinSize().Height))
	}
}

func (renderer *argumentRenderer) MinSize() fyne.Size {
	textsize := fyne.MeasureText(renderer.name.Text, renderer.name.TextSize, renderer.name.TextStyle)
	minsize := fyne.NewSize(textsize.Width+renderer.widget.Disabled.MinSize().Width+3*theme.InnerPadding(), textsize.Height+2*theme.InnerPadding())
	if renderer.widget.Data != nil {
		minsize.Height = minsize.Height + renderer.widget.Data.MinSize().Height + theme.InnerPadding()
		minsize.Width = fyne.Max(minsize.Width, renderer.widget.Data.MinSize().Width+4*theme.InnerPadding())
	}
	return minsize
}

func (renderer *argumentRenderer) Refresh() {
	renderer.name.Text = renderer.widget.Name
	renderer.name.Refresh()
	renderer.background.Refresh()
	renderer.widget.Disabled.Refresh()
	if renderer.widget.Data != nil {
		renderer.widget.Data.Refresh()
	}
}

func (renderer *argumentRenderer) Destroy() {}
