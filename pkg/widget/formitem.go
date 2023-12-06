package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/seternate/go-lanty-client/pkg/theme"
)

type FormItem struct {
	widget.BaseWidget
	text   string
	Widget fyne.CanvasObject
}

func NewFormItem(text string, widget fyne.CanvasObject) (formitem *FormItem) {
	formitem = &FormItem{
		text:   text,
		Widget: widget,
	}
	formitem.ExtendBaseWidget(formitem)
	return
}

func (widget *FormItem) CreateRenderer() fyne.WidgetRenderer {
	widget.ExtendBaseWidget(widget)
	return newFormItemRenderer(widget)
}

type formItemRenderer struct {
	widget         *FormItem
	itemoffset     float32
	textbackground *canvas.Rectangle
	text           *canvas.Text
	itembackground *canvas.Rectangle
}

func newFormItemRenderer(widget *FormItem) (renderer *formItemRenderer) {
	renderer = &formItemRenderer{
		widget:         widget,
		itemoffset:     25,
		textbackground: canvas.NewRectangle(fynetheme.InputBorderColor()),
		text:           canvas.NewText(widget.text, theme.ForegroundColor()),
		itembackground: canvas.NewRectangle(fynetheme.InputBackgroundColor()),
	}
	renderer.text.TextSize = 18
	renderer.text.TextStyle.Bold = true
	renderer.textbackground.CornerRadius = fynetheme.SelectionRadiusSize()
	renderer.itembackground.CornerRadius = fynetheme.SelectionRadiusSize()
	return
}

func (renderer *formItemRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		renderer.textbackground,
		renderer.text,
		renderer.itembackground,
		renderer.widget.Widget,
	}
}

func (renderer *formItemRenderer) Layout(size fyne.Size) {
	textsize := fyne.MeasureText(renderer.text.Text, renderer.text.TextSize, renderer.text.TextStyle)
	renderer.textbackground.Move(fyne.NewPos(0, 0))
	renderer.textbackground.Resize(fyne.NewSize(size.Width, textsize.Height+2*theme.InnerPadding()))
	renderer.text.Move(fyne.NewPos(theme.InnerPadding(), (renderer.textbackground.Size().Height-textsize.Height)/2))

	renderer.itembackground.Move(fyne.NewPos(renderer.itemoffset, renderer.textbackground.Size().Height+theme.InnerPadding()))
	renderer.itembackground.Resize(fyne.NewSize(size.Width-renderer.itemoffset, size.Height-renderer.itembackground.Position().Y))
	renderer.widget.Widget.Move(fyne.NewPos(renderer.itembackground.Position().X+theme.InnerPadding(), renderer.itembackground.Position().Y+theme.InnerPadding()))
	renderer.widget.Widget.Resize(renderer.itembackground.Size().SubtractWidthHeight(2*theme.InnerPadding(), 2*theme.InnerPadding()))
}

func (renderer *formItemRenderer) MinSize() fyne.Size {
	minsizeheader := fyne.MeasureText(renderer.text.Text, renderer.text.TextSize, renderer.text.TextStyle)
	minsizeheader = minsizeheader.AddWidthHeight(2*theme.InnerPadding(), 2*theme.InnerPadding())

	minsizeitem := renderer.widget.Widget.MinSize()
	minsizeitem = minsizeitem.AddWidthHeight(renderer.itemoffset+2*theme.InnerPadding(), 2*theme.InnerPadding())

	return fyne.NewSize(fyne.Max(minsizeheader.Width, minsizeitem.Width), minsizeheader.Height+minsizeitem.Height+theme.InnerPadding())
}

func (renderer *formItemRenderer) Refresh() {
	renderer.textbackground.Refresh()
	renderer.text.Refresh()
	renderer.widget.Widget.Refresh()
	renderer.itembackground.Refresh()
}

func (renderer *formItemRenderer) Destroy() {}
