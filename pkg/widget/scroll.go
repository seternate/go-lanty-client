package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ScrollWithState struct {
	widget.BaseWidget
	scroll      *container.Scroll
	lastpostion float32
}

func NewVScrollWithState(content fyne.CanvasObject) (widget *ScrollWithState) {
	widget = &ScrollWithState{
		scroll: container.NewVScroll(content),
	}
	widget.ExtendBaseWidget(widget)
	content.Refresh()
	widget.Refresh()
	return
}

func (widget *ScrollWithState) Show() {
	if widget.Hidden {
		widget.scroll.Offset.Y = widget.lastpostion
	}
	widget.BaseWidget.Show()
	widget.Refresh()
}

func (widget *ScrollWithState) Hide() {
	widget.lastpostion = widget.scroll.Offset.Y
	widget.BaseWidget.Hide()
	widget.Refresh()
}

func (widget *ScrollWithState) CreateRenderer() fyne.WidgetRenderer {
	return newScrollWithStateRenderer(widget)
}

type scrollWithStateRenderer struct {
	widget *ScrollWithState
}

func newScrollWithStateRenderer(widget *ScrollWithState) *scrollWithStateRenderer {
	renderer := &scrollWithStateRenderer{
		widget: widget,
	}
	return renderer
}

func (renderer *scrollWithStateRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		renderer.widget.scroll,
	}
}

func (renderer *scrollWithStateRenderer) Layout(size fyne.Size) {
	renderer.widget.scroll.Resize(size)
	renderer.widget.scroll.Move(fyne.NewPos(0, 0))
}

func (renderer *scrollWithStateRenderer) MinSize() fyne.Size {
	return renderer.widget.scroll.MinSize()
}

func (renderer *scrollWithStateRenderer) Refresh() {
	renderer.widget.scroll.Content.Refresh()
	renderer.widget.scroll.Refresh()
}

func (renderer *scrollWithStateRenderer) Destroy() {}
