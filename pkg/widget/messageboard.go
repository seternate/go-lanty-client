package widget

import (
	"container/list"
	"slices"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty-client/pkg/controller"
	"github.com/seternate/go-lanty-client/pkg/theme"
	"github.com/seternate/go-lanty/pkg/chat"
	"github.com/seternate/go-lanty/pkg/util"
)

type MessageBoard struct {
	widget.BaseWidget
	controller   *controller.Controller
	messagetiles *list.List
	newmessage   chan chat.Message
	subscriber   []chan struct{}
	mutex        sync.RWMutex
}

func NewMessageBoard(controller *controller.Controller) (messageboard *MessageBoard) {
	messageboard = &MessageBoard{
		controller:   controller,
		messagetiles: list.New(),
		newmessage:   make(chan chat.Message, 50),
		subscriber:   make([]chan struct{}, 50),
	}
	messageboard.ExtendBaseWidget(messageboard)

	controller.Chat.Subscribe(messageboard.newmessage)
	messageboard.run()

	return
}

func (controller *MessageBoard) Subscribe(subscriber chan struct{}) {
	defer controller.mutex.Unlock()
	controller.mutex.Lock()
	controller.subscriber = append(controller.subscriber, subscriber)
}

func (controller *MessageBoard) Unsubscribe(subscriber chan struct{}) {
	defer controller.mutex.Unlock()
	controller.mutex.Lock()
	index := slices.Index(controller.subscriber, subscriber)
	_ = slices.Delete(controller.subscriber, index, index+1)
}

func (widget *MessageBoard) run() {
	widget.controller.WaitGroup().Add(1)
	go widget.messageUpdater()
}

func (widget *MessageBoard) messageUpdater() {
	defer widget.controller.WaitGroup().Done()
	for {
		select {
		case <-widget.controller.Context().Done():
			log.Trace().Msg("exiting messageboard messageUpdater()")
			return
		case message := <-widget.newmessage:
			log.Trace().Interface("message", message).Msg("messageboard received new message")
			widget.addMessageTile(message)
		}
	}
}

func (widget *MessageBoard) addMessageTile(message chat.Message) {
	messagetile := NewMessageTile(message)
	if widget.messagetiles.Len() >= 50 {
		widget.messagetiles.Remove(widget.messagetiles.Back())
	}
	if widget.controller.User.GetUser() == message.GetUser() {
		messagetile.SetBackgroundColor(fynetheme.PrimaryColor())
		messagetile.HideUser()
	}
	widget.messagetiles.PushFront(messagetile)
	widget.Refresh()
	widget.notifySubcriber()
}

func (controller *MessageBoard) notifySubcriber() {
	defer controller.mutex.RUnlock()
	controller.mutex.RLock()
	for _, subscriber := range controller.subscriber {
		util.ChannelWriteNonBlocking(subscriber, struct{}{})
	}
}

func (browser *MessageBoard) CreateRenderer() fyne.WidgetRenderer {
	return newMessageBoardRenderer(browser)
}

type messageBoardRenderer struct {
	widget *MessageBoard
}

func newMessageBoardRenderer(widget *MessageBoard) fyne.WidgetRenderer {
	renderer := &messageBoardRenderer{
		widget: widget,
	}
	return renderer
}

func (renderer *messageBoardRenderer) Objects() (objects []fyne.CanvasObject) {
	for e := renderer.widget.messagetiles.Front(); e != nil; e = e.Next() {
		messageTile := e.Value.(*MessageTile)
		objects = append(objects, messageTile)
	}
	return
}

func (renderer *messageBoardRenderer) Layout(size fyne.Size) {
	for e := renderer.widget.messagetiles.Back(); e != nil; e = e.Prev() {
		messagetile := (e.Value).(*MessageTile)

		//Get longest line
		rawtext := messagetile.GetMessage().GetMessage()
		normalizednewlinetext := strings.ReplaceAll(rawtext, "\r\n", "\n")
		textlines := strings.Split(normalizednewlinetext, "\n")
		longestline := ""
		for _, textline := range textlines {
			if len(textline) > len(longestline) {
				longestline = textline
			}
		}
		//Calculate Tile Width
		maxtextwidth := fyne.MeasureText(longestline, fynetheme.TextSize(), fyne.TextStyle{}).Width + 4*fynetheme.LineSpacing()
		messagetilewidth := fyne.Max(maxtextwidth, messagetile.MinSize().Width)
		if messagetilewidth > size.Width*0.6 {
			messagetilewidth = size.Width * 0.6
		}

		messagetile.Resize(fyne.NewSize(messagetilewidth, messagetile.MinSize().Height))
		posX := theme.InnerPadding()
		posY := theme.InnerPadding()
		if renderer.widget.controller.User.GetUser() == messagetile.message.GetUser() {
			posX = size.Width - messagetile.Size().Width - theme.InnerPadding()
		}
		if e.Next() != nil {
			nextMessagetile := (e.Next().Value).(*MessageTile)
			posY = nextMessagetile.Position().Y + nextMessagetile.Size().Height + theme.InnerPadding()
		}
		messagetile.Move(fyne.NewPos(posX, posY))
	}
}

func (renderer *messageBoardRenderer) MinSize() fyne.Size {
	minWidth := float32(2 * theme.InnerPadding())
	minHeight := float32(theme.InnerPadding())
	for e := renderer.widget.messagetiles.Front(); e != nil; e = e.Next() {
		messagetile := (e.Value).(*MessageTile)
		minHeight = minHeight + messagetile.MinSize().Height + theme.InnerPadding()
		minWidth = fyne.Max(minWidth, 2*theme.InnerPadding()+messagetile.MinSize().Width)
	}
	return fyne.NewSize(minWidth, minHeight)
}

func (renderer *messageBoardRenderer) Refresh() {
	for e := renderer.widget.messagetiles.Front(); e != nil; e = e.Next() {
		messagetile := (e.Value).(*MessageTile)
		messagetile.Refresh()
	}
	//Without first messagetiles will be never be shown until "resize" event happens
	//then all messagetile are shown correctly
	canvas.Refresh(renderer.widget)
}

func (renderer *messageBoardRenderer) Destroy() {}
