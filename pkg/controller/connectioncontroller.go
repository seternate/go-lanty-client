package controller

import (
	"slices"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty/pkg/util"
)

type ConnectionStatus int

const (
	Connected ConnectionStatus = iota
	Disconnected
)

type ConnectionController struct {
	parent          *Controller
	subscriber      []chan struct{}
	refreshinterval time.Duration
	mutex           sync.RWMutex
	Status          ConnectionStatus
}

func NewConnectionController(parent *Controller, refreshinterval time.Duration) (controller *ConnectionController) {
	controller = &ConnectionController{
		parent:          parent,
		subscriber:      make([]chan struct{}, 0, 50),
		refreshinterval: refreshinterval,
		Status:          Disconnected,
	}
	parent.WaitGroup().Add(1)
	go controller.run()
	return
}

func (controller *ConnectionController) run() {
	defer controller.parent.WaitGroup().Done()
	ticker := time.NewTicker(controller.refreshinterval)
	controller.updateStatus()
	for {
		select {
		case <-controller.parent.Context().Done():
			log.Trace().Err(controller.parent.Context().Err()).Msg("exiting connectioncontroller run()")
			return
		case <-ticker.C:
			controller.updateStatus()
		}
	}
}

func (controller *ConnectionController) updateStatus() {
	controller.mutex.Lock()
	oldstatus := controller.Status
	if controller.parent.client.Health.Health() != nil {
		controller.Status = Disconnected
	} else {
		controller.Status = Connected
	}
	newstatus := controller.Status
	controller.mutex.Unlock()
	if newstatus != oldstatus {
		controller.notifySubcriber()
	}
}

func (controller *ConnectionController) Subscribe(subscriber chan struct{}) {
	defer controller.mutex.Unlock()
	controller.mutex.Lock()
	controller.subscriber = append(controller.subscriber, subscriber)
}

func (controller *ConnectionController) Unsubscribe(subscriber chan struct{}) {
	defer controller.mutex.Unlock()
	controller.mutex.Lock()
	index := slices.Index(controller.subscriber, subscriber)
	_ = slices.Delete(controller.subscriber, index, index+1)
}

func (controller *ConnectionController) notifySubcriber() {
	defer controller.mutex.RUnlock()
	controller.mutex.RLock()
	for _, subscriber := range controller.subscriber {
		util.ChannelWriteNonBlocking(subscriber, struct{}{})
	}
}
