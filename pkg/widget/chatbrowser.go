package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty-client/pkg/controller"
	"github.com/seternate/go-lanty-client/pkg/theme"
)

type ChatBrowser struct {
	widget.BaseWidget
	controller   *controller.Controller
	window       fyne.Window
	messageboard *MessageBoard
	scroll       *ScrollWithState
	messageentry *widget.Entry
	sendbutton   *widget.Button
	filebutton   *widget.Button
	newmessage   chan struct{}
}

func NewChatBrowser(controller *controller.Controller, window fyne.Window) (chatbrowser *ChatBrowser) {
	chatbrowser = &ChatBrowser{
		controller:   controller,
		window:       window,
		messageboard: NewMessageBoard(controller),
		messageentry: widget.NewMultiLineEntry(),
		newmessage:   make(chan struct{}, 50),
		sendbutton:   widget.NewButtonWithIcon("", fynetheme.MailSendIcon(), nil),
		filebutton:   widget.NewButtonWithIcon("", fynetheme.MailAttachmentIcon(), nil),
	}
	chatbrowser.ExtendBaseWidget(chatbrowser)

	chatbrowser.scroll = NewVScrollWithState(chatbrowser.messageboard)

	chatbrowser.messageentry.SetMinRowsVisible(2)
	chatbrowser.messageentry.Wrapping = fyne.TextWrapWord
	chatbrowser.messageentry.SetPlaceHolder("Type your message here!")
	chatbrowser.messageentry.OnSubmitted = func(s string) {
		controller.Chat.SendTextMessage(s)
		chatbrowser.messageentry.SetText("")
	}

	chatbrowser.sendbutton.OnTapped = func() {
		controller.Chat.SendTextMessage(chatbrowser.messageentry.Text)
		chatbrowser.messageentry.SetText("")
	}

	chatbrowser.filebutton.OnTapped = chatbrowser.fileUploadCallback

	chatbrowser.messageboard.Subscribe(chatbrowser.newmessage)
	chatbrowser.run()

	return chatbrowser
}

func (widget *ChatBrowser) fileUploadCallback() {
	folderdialog := dialog.NewFileOpen(func(uri fyne.URIReadCloser, err error) {
		if uri == nil || err != nil {
			return
		}
		widget.controller.Chat.SendFileMessage(uri.URI().Path())
	}, widget.window)

	dialogStartURI, err := storage.ListerForURI(storage.NewFileURI(widget.controller.Settings.Settings().GameDirectory))
	if err == nil {
		folderdialog.SetLocation(dialogStartURI)
	}

	//This will make the folderopen dialog to be "fullscreen" inside the app
	folderdialog.Resize(fyne.NewSize(10000, 10000))
	folderdialog.Show()
}

func (widget *ChatBrowser) run() {
	widget.controller.WaitGroup().Add(1)
	go widget.messageUpdater()
}

func (widget *ChatBrowser) messageUpdater() {
	defer widget.controller.WaitGroup().Done()
	for {
		select {
		case <-widget.controller.Context().Done():
			log.Trace().Msg("exiting chatbrowser messageUpdater()")
			return
		case <-widget.newmessage:
			widget.Refresh()
			if widget.scroll.IsNearBottom(0.80) {
				widget.scroll.ScrollToBottom()
			}
		}
	}
}

func (browser *ChatBrowser) CreateRenderer() fyne.WidgetRenderer {
	return newChatBrowserRenderer(browser)
}

type chatBrowserRenderer struct {
	widget *ChatBrowser
}

func newChatBrowserRenderer(widget *ChatBrowser) fyne.WidgetRenderer {
	renderer := &chatBrowserRenderer{
		widget: widget,
	}
	return renderer
}

func (renderer *chatBrowserRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{
		renderer.widget.scroll,
		renderer.widget.sendbutton,
		renderer.widget.filebutton,
		renderer.widget.messageentry,
	}
	return objects
}

func (renderer *chatBrowserRenderer) Layout(size fyne.Size) {
	bottomHeight := fyne.Max(renderer.widget.messageentry.MinSize().Height, renderer.widget.sendbutton.MinSize().Height)
	renderer.widget.messageboard.Resize(fyne.NewSize(size.Width-2*theme.InnerPadding(), renderer.widget.messageboard.MinSize().Height))
	renderer.widget.scroll.Move(fyne.NewPos(theme.InnerPadding(), theme.InnerPadding()))
	renderer.widget.scroll.Resize(fyne.NewSize(size.Width-2*theme.InnerPadding(), size.Height-3*theme.InnerPadding()-bottomHeight))
	renderer.widget.sendbutton.Resize(fyne.NewSize(bottomHeight, bottomHeight))
	renderer.widget.sendbutton.Move(fyne.NewPos(size.Width-theme.InnerPadding()-renderer.widget.sendbutton.Size().Width, renderer.widget.scroll.Position().Y+renderer.widget.scroll.Size().Height+theme.InnerPadding()))
	renderer.widget.filebutton.Resize(fyne.NewSize(bottomHeight, bottomHeight))
	renderer.widget.filebutton.Move(fyne.NewPos(renderer.widget.sendbutton.Position().X-theme.InnerPadding()-renderer.widget.filebutton.Size().Width, renderer.widget.scroll.Position().Y+renderer.widget.scroll.Size().Height+theme.InnerPadding()))
	renderer.widget.messageentry.Move(fyne.NewPos(renderer.widget.scroll.Position().X, renderer.widget.scroll.Position().Y+renderer.widget.scroll.Size().Height+theme.InnerPadding()))
	renderer.widget.messageentry.Resize(fyne.NewSize(renderer.widget.filebutton.Position().X-2*theme.InnerPadding(), bottomHeight))
}

func (renderer *chatBrowserRenderer) MinSize() fyne.Size {
	bottomHeight := fyne.Max(renderer.widget.messageentry.MinSize().Height, renderer.widget.sendbutton.MinSize().Height)
	minWidth := fyne.Max(4*theme.InnerPadding()+renderer.widget.messageentry.MinSize().Width+renderer.widget.sendbutton.MinSize().Width+renderer.widget.filebutton.MinSize().Width, 2*theme.InnerPadding()+renderer.widget.scroll.MinSize().Width)
	minHeight := 3*theme.InnerPadding() + bottomHeight
	return fyne.NewSize(minWidth, minHeight)
}

func (renderer *chatBrowserRenderer) Refresh() {
	renderer.widget.messageboard.Refresh()
	renderer.widget.scroll.Refresh()
	renderer.widget.messageentry.Refresh()
	renderer.widget.sendbutton.Refresh()
	renderer.widget.filebutton.Refresh()
	//Without first messagetiles will be never be shown until "resize" event happens
	//then all messagetile are shown correctly
	canvas.Refresh(renderer.widget)
}

func (renderer *chatBrowserRenderer) Destroy() {}
