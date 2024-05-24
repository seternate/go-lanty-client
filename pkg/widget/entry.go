package widget

import "fyne.io/fyne/v2/widget"

type Entry struct {
	widget.Entry

	OnFocusChanged func(bool)
}

func NewEntry() *Entry {
	entry := &Entry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (widget *Entry) FocusGained() {
	widget.Entry.FocusGained()
	if widget.OnFocusChanged != nil {
		widget.OnFocusChanged(true)
	}
}

func (widget *Entry) FocusLost() {
	widget.Entry.FocusLost()
	if widget.OnFocusChanged != nil {
		widget.OnFocusChanged(false)
	}
}