package widget

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/seternate/go-lanty-client/pkg/theme"
	"github.com/seternate/go-lanty/pkg/chat"
)

type MessageTile struct {
	widget.BaseWidget
	message         chat.Message
	backgroundcolor color.Color
	showuser        bool
}

func NewMessageTile(message chat.Message) (messagetile *MessageTile) {
	messagetile = &MessageTile{
		message:         message,
		backgroundcolor: fynetheme.InputBackgroundColor(),
		showuser:        true,
	}
	messagetile.ExtendBaseWidget(messagetile)

	return messagetile
}

func (widget *MessageTile) SetBackgroundColor(color color.Color) {
	widget.backgroundcolor = color
	widget.Refresh()
}

func (widget *MessageTile) ShowUser() {
	if widget.showuser {
		return
	}
	widget.showuser = true
	widget.Refresh()
}

func (widget *MessageTile) HideUser() {
	if !widget.showuser {
		return
	}
	widget.showuser = false
	widget.Refresh()
}

func (widget *MessageTile) GetMessage() chat.Message {
	return widget.message
}

func (widget *MessageTile) CreateRenderer() fyne.WidgetRenderer {
	return newMessageTileRenderer(widget)
}

type messageTileRenderer struct {
	widget     *MessageTile
	background *canvas.Rectangle
	user       *canvas.Text
	message    *widget.Label
	time       *canvas.Text
}

func newMessageTileRenderer(w *MessageTile) (renderer *messageTileRenderer) {
	renderer = &messageTileRenderer{
		widget:     w,
		background: canvas.NewRectangle(w.backgroundcolor),
		user:       canvas.NewText(w.message.GetUser().Name, theme.ForegroundColor()),
		message:    widget.NewLabel(w.message.GetMessage()),
		time:       canvas.NewText(w.message.GetTime().Format("15:04"), theme.ForegroundColor()),
	}
	renderer.background.CornerRadius = fynetheme.SelectionRadiusSize()
	renderer.user.TextSize = 12
	renderer.time.TextSize = 8
	renderer.message.Wrapping = fyne.TextWrapWord
	return
}

func (renderer *messageTileRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		renderer.background,
		renderer.user,
		renderer.message,
		renderer.time,
	}
}

func (renderer *messageTileRenderer) Layout(size fyne.Size) {
	renderer.background.Resize(size)
	timetextsize := fyne.MeasureText(renderer.time.Text, renderer.time.TextSize, renderer.time.TextStyle)
	renderer.time.Move(fyne.NewPos(size.Width-theme.InnerPadding()-timetextsize.Width, (size.Height - theme.InnerPadding()/2 - timetextsize.Height)))
	if renderer.widget.showuser {
		usertextsize := fyne.MeasureText(renderer.user.Text, renderer.user.TextSize, renderer.user.TextStyle)
		renderer.user.Move(fyne.NewPos(theme.InnerPadding(), theme.InnerPadding()/2))
		renderer.message.Resize(fyne.NewSize(size.Width, renderer.time.Position().Y-renderer.user.Position().Y-usertextsize.Height))
		renderer.message.Move(fyne.NewPos(0, renderer.user.Position().Y+usertextsize.Height))
	} else {
		renderer.message.Resize(fyne.NewSize(size.Width, renderer.time.Position().Y))
		renderer.message.Move(fyne.NewPos(0, 0))
	}
}

func (renderer *messageTileRenderer) MinSize() (minSize fyne.Size) {
	timetextsize := fyne.MeasureText(renderer.time.Text, renderer.time.TextSize, renderer.time.TextStyle)
	if renderer.widget.showuser {
		usertextsize := fyne.MeasureText(renderer.user.Text, renderer.user.TextSize, renderer.user.TextStyle)
		minSize = fyne.NewSize(
			2*theme.InnerPadding()+fyne.Max(usertextsize.Width, timetextsize.Width),
			theme.InnerPadding()+usertextsize.Height+renderer.message.MinSize().Height+timetextsize.Height,
		)
	} else {
		minSize = fyne.NewSize(
			2*theme.InnerPadding()+timetextsize.Width,
			0.5*theme.InnerPadding()+renderer.message.MinSize().Height+timetextsize.Height,
		)
	}
	return
}

func (renderer *messageTileRenderer) Refresh() {
	renderer.background.FillColor = renderer.widget.backgroundcolor
	if renderer.widget.showuser && renderer.user.Hidden {
		renderer.user.Show()
	} else if !renderer.widget.showuser && !renderer.user.Hidden {
		renderer.user.Hide()
	}
	renderer.background.Refresh()
	renderer.user.Refresh()
	renderer.message.Refresh()
	renderer.time.Refresh()
}

func (renderer *messageTileRenderer) Destroy() {}
