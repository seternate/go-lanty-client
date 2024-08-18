package controller

import (
	"fmt"
	"net/url"
	"slices"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty/pkg/chat"
	"github.com/seternate/go-lanty/pkg/util"
)

type ChatController struct {
	parent          *Controller
	subscriber      []chan chat.Message
	ticker          *time.Ticker
	mutex           sync.RWMutex
	settingschanged chan struct{}
}

func NewChatController(parent *Controller) (controller *ChatController) {
	controller = &ChatController{
		parent:          parent,
		subscriber:      make([]chan chat.Message, 50),
		ticker:          time.NewTicker(time.Second),
		settingschanged: make(chan struct{}, 50),
	}

	controller.parent.Settings.Subscribe(controller.settingschanged)
	controller.run()
	return
}

func (controller *ChatController) DownloadFile(message chat.Message) {
	if message.GetType() != chat.TYPE_FILE {
		log.Warn().Interface("message", message).Msg("wrong message type provided for file download")
		return
	}
	m := message.(*chat.FileMessage)
	u, err := url.Parse(m.Message.URL)
	if err != nil {
		log.Error().Err(err).Interface("message", message).Msg("error parsing filemessage url")
		controller.parent.Status.Error("Failed to start download", 3*time.Second)
		return
	}
	download, err := controller.parent.client.File.GetFile(controller.parent.ctx, *u, controller.parent.settings.DownloadDirectory)
	if err != nil {
		log.Error().Err(err).Interface("message", message).Msg("error starting filemessage download")
		controller.parent.Status.Error(fmt.Sprintf("Failed downloading %s", download.Filename()), 3*time.Second)
		return
	}
	go func() {
		<-download.Done
		if download.Err != nil {
			log.Error().Err(err).Str("file", download.Filename()).Msg("error downloading filemessage file")
			controller.parent.Status.Error(fmt.Sprintf("Failed downloading %s", download.Filename()), 3*time.Second)
			return
		}
		log.Debug().Str("file", download.Filename()).Msg("sucessfully downloaded filemessage file")
		controller.parent.Status.Info(fmt.Sprintf("Downloaded \"%s\" to \"%s\"", download.Filename(), controller.parent.settings.DownloadDirectory), 3*time.Second)
	}()
}

func (controller *ChatController) SendTextMessage(message string) {
	err := controller.parent.client.Chat.SendMessage(chat.NewTextMessage(controller.parent.User.GetUser(), message))
	if err != nil {
		log.Error().Err(err).Msg("error sending textmessage to server")
	}
}

func (controller *ChatController) SendFileMessage(path string) {
	controller.parent.Status.Info(fmt.Sprintf("Uploading file \"%s\" ...", path), 3*time.Second)
	fileresponse, err := controller.parent.client.File.UploadFile(path)
	if err != nil {
		controller.parent.Status.Error(fmt.Sprintf("Error uploading file \"%s\"", path), 3*time.Second)
		log.Error().Err(err).Str("file", path).Msg("error uploading file to server")
		return
	}

	message := chat.NewFileMessage(controller.parent.User.GetUser(), fileresponse)
	err = controller.parent.client.Chat.SendMessage(message)
	if err != nil {
		log.Error().Err(err).Interface("message", message).Msg("error sending filemessage to server")
	}
}

func (controller *ChatController) Subscribe(subscriber chan chat.Message) {
	defer controller.mutex.Unlock()
	controller.mutex.Lock()
	controller.subscriber = append(controller.subscriber, subscriber)
}

func (controller *ChatController) Unsubscribe(subscriber chan chat.Message) {
	defer controller.mutex.Unlock()
	controller.mutex.Lock()
	index := slices.Index(controller.subscriber, subscriber)
	_ = slices.Delete(controller.subscriber, index, index+1)
}

func (controller *ChatController) run() {
	controller.parent.WaitGroup().Add(2)
	go controller.connectionWatcher()
	go controller.messageReader()
}

func (controller *ChatController) connectionWatcher() {
	defer func() { controller.parent.WaitGroup().Done(); controller.parent.client.Chat.Disconnect() }()
	_, err := controller.parent.client.Chat.Connect()
	if err != nil {
		log.Error().Err(err).Msg("error connecting to chat")
	} else {
		log.Debug().Msg("successfully connected to chat")
	}
	for {
		select {
		case <-controller.parent.Context().Done():
			log.Trace().Err(controller.parent.Context().Err()).Msg("exiting ChatController connectionWatcher()")
			return
		case <-controller.ticker.C:
			if controller.parent.client.Chat.Error != nil {
				log.Debug().Err(controller.parent.client.Chat.Error).Msg("trying to reconnect to chat due to error in chatservice")
				err = controller.parent.client.Chat.Reconnect()
				if err != nil {
					log.Error().Err(err).Msg("error reconnecting to chat")
				} else {
					log.Debug().Msg("successfully reconnected to chat")
				}
			}
		case <-controller.settingschanged:
			log.Debug().Err(controller.parent.client.Chat.Error).Msg("trying to reconnect to chat due to a settings change")
			err = controller.parent.client.Chat.Reconnect()
			if err != nil {
				log.Error().Err(err).Msg("error reconnecting to chat")
			} else {
				log.Debug().Msg("successfully reconnected to chat")
			}
		}
	}
}

func (controller *ChatController) messageReader() {
	defer controller.parent.WaitGroup().Done()
	for {
		select {
		case <-controller.parent.Context().Done():
			log.Trace().Err(controller.parent.Context().Err()).Msg("exiting ChatController messageReader()")
			return
		case message := <-controller.parent.client.Chat.Messages:
			controller.notifySubcriber(message)
		}
	}
}

func (controller *ChatController) notifySubcriber(message chat.Message) {
	defer controller.mutex.RUnlock()
	controller.mutex.RLock()
	for _, subscriber := range controller.subscriber {
		util.ChannelWriteNonBlocking(subscriber, message)
	}
}
