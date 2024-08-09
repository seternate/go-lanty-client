package controller

import (
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

func (controller *ChatController) SendTextMessage(message string) {
	err := controller.parent.client.Chat.SendMessage(chat.NewTextMessage(controller.parent.User.GetUser(), message))
	if err != nil {
		log.Error().Err(err).Msg("error sending textmessage to server")
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
