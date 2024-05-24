package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/seternate/go-lanty-client/pkg/theme"
)

type Form struct {
	widget.BaseWidget

	OnSubmit func()
	OnCancel func()

	items  []*FormItem
	submit *widget.Button
	cancel *widget.Button
}

func NewForm() (form *Form) {
	form = &Form{
		items:  make([]*FormItem, 0),
		submit: widget.NewButton("Submit", nil),
		cancel: widget.NewButton("Cancel", nil),
	}
	form.ExtendBaseWidget(form)
	form.submit.OnTapped = func() {
		if form.OnSubmit != nil {
			form.OnSubmit()
		}
	}
	form.cancel.OnTapped = func() {
		if form.OnCancel != nil {
			form.OnCancel()
		}
	}
	form.submit.Importance = widget.HighImportance
	return
}

func (widget *Form) AppendItem(item *FormItem) {
	widget.items = append(widget.items, item)
	if w, ok := item.Widget.(fyne.Validatable); ok {
		w.SetOnValidationChanged(func(err error) {
			if err != nil {
				widget.submit.Disable()
				return
			}
			widget.submit.Enable()
		})
	}
	widget.Refresh()
}

func (widget *Form) SetSubmitText(text string) {
	widget.submit.SetText(text)
	widget.Refresh()
}

func (widget *Form) HideSubmit() {
	widget.submit.Hide()
	widget.Refresh()
}

func (widget *Form) SetCancelText(text string) {
	widget.cancel.SetText(text)
	widget.Refresh()
}

func (widget *Form) CreateRenderer() fyne.WidgetRenderer {
	return newFormRenderer(widget)
}

type formRenderer struct {
	widget *Form
}

func newFormRenderer(widget *Form) *formRenderer {
	return &formRenderer{
		widget: widget,
	}
}

func (renderer *formRenderer) Objects() (objects []fyne.CanvasObject) {
	for _, item := range renderer.widget.items {
		objects = append(objects, item)
	}
	objects = append(objects, renderer.widget.submit, renderer.widget.cancel)
	return
}

func (renderer *formRenderer) Layout(size fyne.Size) {
	lastposition := float32(0)
	for _, item := range renderer.widget.items {
		item.Move(fyne.NewPos(theme.InnerPadding(), lastposition+theme.InnerPadding()))
		item.Resize(fyne.NewSize(size.Width-2*theme.InnerPadding(), item.MinSize().Height))
		lastposition = item.Position().Y + item.Size().Height
	}
	if renderer.widget.submit.Hidden {
		buttonSize := fyne.NewSize(renderer.widget.cancel.MinSize().Width, renderer.widget.cancel.MinSize().Height)
		renderer.widget.cancel.Resize(buttonSize)
		renderer.widget.cancel.Move(fyne.NewPos(size.Width-theme.InnerPadding()-buttonSize.Width, size.Height-theme.InnerPadding()-buttonSize.Height))
	} else {
		buttonSize := fyne.NewSize(fyne.Max(renderer.widget.cancel.MinSize().Width, renderer.widget.submit.MinSize().Width), renderer.widget.submit.MinSize().Height)
		renderer.widget.submit.Resize(buttonSize)
		renderer.widget.submit.Move(fyne.NewPos(size.Width-theme.InnerPadding()-buttonSize.Width, size.Height-theme.InnerPadding()-buttonSize.Height))
		renderer.widget.cancel.Resize(buttonSize)
		renderer.widget.cancel.Move(fyne.NewPos(renderer.widget.submit.Position().X-theme.InnerPadding()-buttonSize.Width, renderer.widget.submit.Position().Y))
	}
}

func (renderer *formRenderer) MinSize() fyne.Size {
	minwidth := float32(0)
	minheight := theme.InnerPadding()
	for _, item := range renderer.widget.items {
		minheight += item.MinSize().Height + theme.InnerPadding()
		minwidth = fyne.Max(minwidth, item.MinSize().Width+2*theme.InnerPadding())
	}
	minheight += renderer.widget.submit.MinSize().Height + theme.InnerPadding()
	minwidth = fyne.Max(minwidth, renderer.widget.submit.MinSize().Width+renderer.widget.cancel.MinSize().Width+3*theme.InnerPadding())
	return fyne.NewSize(minwidth, minheight)
}

func (renderer *formRenderer) Refresh() {
	for _, item := range renderer.widget.items {
		item.Refresh()
	}
	renderer.widget.submit.Refresh()
	renderer.widget.cancel.Refresh()
}

func (renderer *formRenderer) Destroy() {}
