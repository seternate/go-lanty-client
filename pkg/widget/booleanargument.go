package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/seternate/go-lanty/pkg/game/argument"
)

type BooleanArgument struct {
	widget.BaseWidget
	Argument   *argument.Boolean
	radiogroup *widget.RadioGroup
}

func NewBooleanArgument(arg *argument.Boolean) *BooleanArgument {
	item := &BooleanArgument{
		Argument:   arg,
		radiogroup: widget.NewRadioGroup([]string{"Enable", "Disable"}, nil),
	}
	item.ExtendBaseWidget(item)

	var selected string
	if item.Argument.Value {
		selected = "Enable"
	} else {
		selected = "Disable"
	}
	item.radiogroup.SetSelected(selected)
	item.radiogroup.Horizontal = true
	item.radiogroup.OnChanged = func(s string) {
		if s == "Enable" {
			item.Argument.Value = true
		} else {
			item.Argument.Value = false
		}
	}

	return item
}

func (item *BooleanArgument) Refresh() {
	var selected string
	if item.Argument.Value {
		selected = "Enable"
	} else {
		selected = "Disable"
	}
	item.radiogroup.SetSelected(selected)
	item.BaseWidget.Refresh()
}

func (item *BooleanArgument) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(item.radiogroup)
}
