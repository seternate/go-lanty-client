package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/seternate/go-lanty/pkg/game/argument"
)

type StringArgument struct {
	widget.BaseWidget
	Argument *argument.String
	entry    *widget.Entry
}

func NewStringArgument(arg *argument.String) *StringArgument {
	item := &StringArgument{
		Argument: arg,
		entry:    widget.NewEntry(),
	}
	item.ExtendBaseWidget(item)

	item.entry.SetText(item.Argument.Value)
	item.entry.OnChanged = func(s string) {
		item.Argument.Value = s
	}

	return item
}

func (item *StringArgument) Refresh() {
	item.entry.SetText(item.Argument.Value)
	item.BaseWidget.Refresh()
}

func (item *StringArgument) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(item.entry)
}
