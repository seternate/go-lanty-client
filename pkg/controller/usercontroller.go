package controller

import (
	"slices"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty/pkg/user"
	"github.com/seternate/go-lanty/pkg/util"
)

type UserController struct {
	parent          *Controller
	users           []user.User
	subscriber      []chan struct{}
	refreshinterval time.Duration
	err             error
	mutex           sync.RWMutex
}

func NewUserController(parent *Controller, refreshinteval time.Duration) (controller *UserController) {
	controller = &UserController{
		parent:          parent,
		users:           make([]user.User, 0),
		subscriber:      make([]chan struct{}, 0),
		refreshinterval: refreshinteval,
	}
	go controller.run()
	return
}

func (controller *UserController) GetUsers() []user.User {
	defer controller.mutex.RUnlock()
	controller.mutex.RLock()
	return controller.users
}

func (controller *UserController) Err() error {
	defer controller.mutex.RUnlock()
	controller.mutex.RLock()
	return controller.err
}

func (controller *UserController) run() {
	for {
		log.Trace().Msg("usercontroller update loop")
		controller.update()
		time.Sleep(controller.refreshinterval)
	}
}

func (controller *UserController) update() {
	names, err := controller.parent.client.User.GetUsers()
	if err != nil {
		controller.mutex.Lock()
		controller.err = err
		controller.mutex.Unlock()
		log.Error().Err(err).Msg("error retrieving userlist from server")
		return
	}

	users := make([]user.User, 0)
	for _, name := range names {
		user, err := controller.parent.client.User.GetUser(name)
		if err != nil {
			controller.mutex.Lock()
			controller.err = err
			controller.mutex.Unlock()
			log.Error().Err(err).Str("user", name).Msg("error retrieving user from server")
			return
		}
		users = append(users, user)
	}
	log.Debug().Interface("users", users).Msg("updated users in usercontroller")
	controller.mutex.Lock()
	controller.users = users
	controller.mutex.Unlock()
	controller.notifySubcriber()
}

func (controller *UserController) Subscribe(subscriber chan struct{}) {
	defer controller.mutex.Unlock()
	controller.mutex.Lock()
	log.Trace().Msg("new subscriber to usercontroller")
	controller.subscriber = append(controller.subscriber, subscriber)
}

func (controller *UserController) Unsubscribe(subscriber chan struct{}) {
	defer controller.mutex.Unlock()
	controller.mutex.Lock()
	index := slices.Index(controller.subscriber, subscriber)
	slices.Delete(controller.subscriber, index, index+1)
}

func (controller *UserController) notifySubcriber() {
	defer controller.mutex.RUnlock()
	controller.mutex.RLock()
	log.Trace().Msg("notify subscriber of usercontroller")
	for _, subscriber := range controller.subscriber {
		util.ChannelWriteNonBlocking(subscriber, struct{}{})
	}
}
