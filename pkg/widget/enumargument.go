package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/seternate/go-lanty/pkg/game/argument"
)

type EnumArgument struct {
	widget.BaseWidget
	Argument *argument.Enum
	selectw  *widget.Select
}

func NewEnumArgument(arg *argument.Enum) *EnumArgument {
	item := &EnumArgument{
		Argument: arg,
		selectw:  widget.NewSelect([]string{""}, nil),
	}
	item.ExtendBaseWidget(item)

	options := []string{}
	option := ""
	for _, i := range item.Argument.Items {
		options = append(options, i.Name)
		if i.Value == item.Argument.Value {
			option = i.Name
		}
	}
	item.selectw.SetOptions(options)
	item.selectw.OnChanged = func(s string) {
		for _, i := range item.Argument.Items {
			if i.Name == s {
				item.Argument.Value = i.Value
			}
		}
	}
	item.selectw.SetSelected(option)

	return item
}

func (item *EnumArgument) Refresh() {
	option := ""
	for _, i := range item.Argument.Items {
		if i.Value == item.Argument.Value {
			option = i.Name
		}
	}
	item.selectw.SetSelected(option)
	item.BaseWidget.Refresh()
}

func (item *EnumArgument) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(item.selectw)
}
