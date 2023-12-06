package widget

import (
	"errors"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/seternate/go-lanty-client/pkg/theme"
	"github.com/seternate/go-lanty/pkg/game/argument"
)

type IntegerArgument struct {
	widget.BaseWidget
	Argument   *argument.Integer
	valueentry *widget.Entry
}

func NewIntegerArgument(arg *argument.Integer) *IntegerArgument {
	item := &IntegerArgument{
		Argument:   arg,
		valueentry: widget.NewEntry(),
	}
	item.ExtendBaseWidget(item)

	item.valueentry.SetText(strconv.Itoa(item.Argument.Value))
	item.valueentry.OnChanged = func(s string) {
		if item.valueentry.Validate() == nil {
			value, err := strconv.Atoi(s)
			if err != nil {
				return
			}
			item.Argument.Value = value
		}
	}
	item.valueentry.Validator = func(s string) error {
		value, err := strconv.Atoi(s)
		if err != nil {
			return errors.New("wrong input")
		}
		if value < item.Argument.MinValue {
			return errors.New("input smaller than minvalue")
		}
		if value > item.Argument.MaxValue {
			return errors.New("input greater than maxvalue")
		}
		return nil
	}

	return item
}

func (item *IntegerArgument) Refresh() {
	item.valueentry.SetText(strconv.Itoa(item.Argument.Value))
	item.BaseWidget.Refresh()
}

func (item *IntegerArgument) CreateRenderer() fyne.WidgetRenderer {
	return newIntegerArgumentRenderer(item)
}

type integerArgumentRenderer struct {
	widget   *IntegerArgument
	minvalue *canvas.Text
	maxvalue *canvas.Text
}

func newIntegerArgumentRenderer(widget *IntegerArgument) *integerArgumentRenderer {
	renderer := &integerArgumentRenderer{
		widget:   widget,
		minvalue: canvas.NewText("Min: "+strconv.Itoa(widget.Argument.MinValue), theme.ForegroundColor()),
		maxvalue: canvas.NewText("Max: "+strconv.Itoa(widget.Argument.MaxValue), theme.ForegroundColor()),
	}
	return renderer
}

func (renderer *integerArgumentRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		renderer.widget.valueentry,
		renderer.minvalue,
		renderer.maxvalue,
	}
}

func (renderer *integerArgumentRenderer) Layout(size fyne.Size) {
	textsizeminvalue := fyne.MeasureText(renderer.minvalue.Text, renderer.minvalue.TextSize, renderer.minvalue.TextStyle)
	textsizemaxvalue := fyne.MeasureText(renderer.maxvalue.Text, renderer.maxvalue.TextSize, renderer.maxvalue.TextStyle)
	renderer.widget.valueentry.Move(fyne.NewPos(0, 0))
	renderer.widget.valueentry.Resize(fyne.NewSize(size.Width-textsizeminvalue.Width-textsizemaxvalue.Width-2*theme.InnerPadding(), size.Height))
	renderer.minvalue.Move(fyne.NewPos(renderer.widget.valueentry.Size().Width+theme.InnerPadding(), (size.Height-textsizeminvalue.Height)/2))
	renderer.maxvalue.Move(fyne.NewPos(renderer.minvalue.Position().X+textsizeminvalue.Width+theme.InnerPadding(), (size.Height-textsizemaxvalue.Height)/2))
}

func (renderer *integerArgumentRenderer) MinSize() fyne.Size {
	textsizeminvalue := fyne.MeasureText(renderer.minvalue.Text, renderer.minvalue.TextSize, renderer.minvalue.TextStyle)
	textsizemaxvalue := fyne.MeasureText(renderer.maxvalue.Text, renderer.maxvalue.TextSize, renderer.maxvalue.TextStyle)
	minsize := renderer.widget.valueentry.MinSize()
	minsize = minsize.AddWidthHeight(textsizeminvalue.Width+textsizemaxvalue.Width+2*theme.InnerPadding(), 0)
	return minsize
}

func (renderer *integerArgumentRenderer) Refresh() {
	renderer.widget.valueentry.Refresh()
	renderer.minvalue.Refresh()
	renderer.maxvalue.Refresh()
}

func (renderer *integerArgumentRenderer) Destroy() {}
