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

type FloatArgument struct {
	widget.BaseWidget
	Argument   *argument.Float
	valueentry *widget.Entry
}

func NewFloatArgument(arg *argument.Float) *FloatArgument {
	item := &FloatArgument{
		Argument:   arg,
		valueentry: widget.NewEntry(),
	}
	item.ExtendBaseWidget(item)

	item.valueentry.SetText(strconv.FormatFloat(float64(item.Argument.Value), 'f', -1, 32))
	item.valueentry.OnChanged = func(s string) {
		if item.valueentry.Validate() == nil {
			value64, err := strconv.ParseFloat(s, 32)
			if err != nil {
				return
			}
			item.Argument.Value = float32(value64)
		}
	}
	item.valueentry.Validator = func(s string) error {
		value64, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return errors.New("wrong input")
		}
		if float32(value64) < item.Argument.MinValue {
			return errors.New("input smaller than minvalue")
		}
		if float32(value64) > item.Argument.MaxValue {
			return errors.New("input greater than maxvalue")
		}
		return nil
	}

	return item
}

func (item *FloatArgument) Refresh() {
	item.valueentry.SetText(strconv.FormatFloat(float64(item.Argument.Value), 'f', -1, 32))
	item.BaseWidget.Refresh()
}

func (item *FloatArgument) CreateRenderer() fyne.WidgetRenderer {
	return newFloatArgumentRenderer(item)
}

type floatArgumentRenderer struct {
	widget   *FloatArgument
	minvalue *canvas.Text
	maxvalue *canvas.Text
}

func newFloatArgumentRenderer(widget *FloatArgument) *floatArgumentRenderer {
	renderer := &floatArgumentRenderer{
		widget:   widget,
		minvalue: canvas.NewText("Min: "+strconv.FormatFloat(float64(widget.Argument.MinValue), 'f', -1, 32), theme.ForegroundColor()),
		maxvalue: canvas.NewText("Max: "+strconv.FormatFloat(float64(widget.Argument.MaxValue), 'f', -1, 32), theme.ForegroundColor()),
	}
	return renderer
}

func (renderer *floatArgumentRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		renderer.widget.valueentry,
		renderer.minvalue,
		renderer.maxvalue,
	}
}

func (renderer *floatArgumentRenderer) Layout(size fyne.Size) {
	textsizeminvalue := fyne.MeasureText(renderer.minvalue.Text, renderer.minvalue.TextSize, renderer.minvalue.TextStyle)
	textsizemaxvalue := fyne.MeasureText(renderer.maxvalue.Text, renderer.maxvalue.TextSize, renderer.maxvalue.TextStyle)
	renderer.widget.valueentry.Move(fyne.NewPos(0, 0))
	renderer.widget.valueentry.Resize(fyne.NewSize(size.Width-textsizeminvalue.Width-textsizemaxvalue.Width-2*theme.InnerPadding(), size.Height))
	renderer.minvalue.Move(fyne.NewPos(renderer.widget.valueentry.Size().Width+theme.InnerPadding(), (size.Height-textsizeminvalue.Height)/2))
	renderer.maxvalue.Move(fyne.NewPos(renderer.minvalue.Position().X+textsizeminvalue.Width+theme.InnerPadding(), (size.Height-textsizemaxvalue.Height)/2))
}

func (renderer *floatArgumentRenderer) MinSize() fyne.Size {
	textsizeminvalue := fyne.MeasureText(renderer.minvalue.Text, renderer.minvalue.TextSize, renderer.minvalue.TextStyle)
	textsizemaxvalue := fyne.MeasureText(renderer.maxvalue.Text, renderer.maxvalue.TextSize, renderer.maxvalue.TextStyle)
	minsize := renderer.widget.valueentry.MinSize()
	minsize = minsize.AddWidthHeight(textsizeminvalue.Width+textsizemaxvalue.Width+2*theme.InnerPadding(), 0)
	return minsize
}

func (renderer *floatArgumentRenderer) Refresh() {
	renderer.widget.valueentry.Refresh()
	renderer.minvalue.Refresh()
	renderer.maxvalue.Refresh()
}

func (renderer *floatArgumentRenderer) Destroy() {}
