package controller

import (
	"slices"
	"sync"
	"time"

	"github.com/golang-collections/collections/queue"
	"github.com/seternate/go-lanty/pkg/util"
)

type StatusLevel int

const (
	StatusLevelInfo StatusLevel = iota
	StatusLevelWarning
	StatusLevelError
)

type Status struct {
	Level    StatusLevel
	Text     string
	Duration time.Duration
}

func NewInfo(text string, duration time.Duration) Status {
	return Status{
		Level:    StatusLevelInfo,
		Text:     text,
		Duration: duration,
	}
}

func NewWarning(text string, duration time.Duration) Status {
	return Status{
		Level:    StatusLevelWarning,
		Text:     text,
		Duration: duration,
	}
}

func NewError(text string, duration time.Duration) Status {
	return Status{
		Level:    StatusLevelError,
		Text:     text,
		Duration: duration,
	}
}

type StatusController struct {
	infostatus    queue.Queue
	warningstatus queue.Queue
	errorstatus   queue.Queue
	subscriber    []chan struct{}
	mutex         sync.Mutex
}

func NewStatusController() (controller *StatusController) {
	controller = &StatusController{}
	return
}

func (controller *StatusController) Info(text string, duration time.Duration) {
	controller.mutex.Lock()
	controller.infostatus.Enqueue(NewInfo(text, duration))
	controller.mutex.Unlock()
	controller.notifySubcriber()
}

func (controller *StatusController) Warning(text string, duration time.Duration) {
	controller.mutex.Lock()
	controller.warningstatus.Enqueue(NewWarning(text, duration))
	controller.mutex.Unlock()
	controller.notifySubcriber()
}

func (controller *StatusController) Error(text string, duration time.Duration) {
	controller.mutex.Lock()
	controller.errorstatus.Enqueue(NewError(text, duration))
	controller.mutex.Unlock()
	controller.notifySubcriber()
}

func (controller *StatusController) Next() (status Status) {
	controller.mutex.Lock()
	if controller.errorstatus.Len() > 0 {
		status = controller.errorstatus.Dequeue().(Status)
	} else if controller.warningstatus.Len() > 0 {
		status = controller.warningstatus.Dequeue().(Status)
	} else if controller.infostatus.Len() > 0 {
		status = controller.infostatus.Dequeue().(Status)
	}
	controller.mutex.Unlock()
	return
}

func (controller *StatusController) Subscribe(subscriber chan struct{}) {
	defer controller.mutex.Unlock()
	controller.mutex.Lock()
	controller.subscriber = append(controller.subscriber, subscriber)
}

func (controller *StatusController) Unsubscribe(subscriber chan struct{}) {
	defer controller.mutex.Unlock()
	controller.mutex.Lock()
	index := slices.Index(controller.subscriber, subscriber)
	slices.Delete(controller.subscriber, index, index+1)
}

func (controller *StatusController) notifySubcriber() {
	defer controller.mutex.Unlock()
	controller.mutex.Lock()
	for _, subscriber := range controller.subscriber {
		util.ChannelWriteNonBlocking(subscriber, struct{}{})
	}
}
