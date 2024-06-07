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
	user            user.User
	users           user.Users
	loggedIn        bool
	subscriber      []chan struct{}
	refreshinterval time.Duration
	err             error
	usernameupdated chan struct{}
	mutex           sync.RWMutex
}

func NewUserController(parent *Controller, refreshinteval time.Duration) (controller *UserController) {
	controller = &UserController{
		parent:          parent,
		user:            user.User{Name: parent.settings.Username},
		loggedIn:        false,
		subscriber:      make([]chan struct{}, 0, 50),
		refreshinterval: refreshinteval,
		usernameupdated: make(chan struct{}, 50),
	}
	parent.Settings.Subscribe(controller.usernameupdated)
	parent.WaitGroup().Add(1)
	go controller.run()
	return
}

func (controller *UserController) GetUser() user.User {
	defer controller.mutex.RUnlock()
	controller.mutex.RLock()
	return controller.user
}

func (controller *UserController) IsLoggedIn() bool {
	defer controller.mutex.RUnlock()
	controller.mutex.RLock()
	return controller.loggedIn
}

func (controller *UserController) GetUsers() []user.User {
	defer controller.mutex.RUnlock()
	controller.mutex.RLock()
	return controller.users.Users()
}

func (controller *UserController) Err() error {
	defer controller.mutex.RUnlock()
	controller.mutex.RLock()
	return controller.err
}

func (controller *UserController) run() {
	defer controller.parent.WaitGroup().Done()
	controller.login()
	controller.updateUsers()
	ticker := time.NewTicker(controller.refreshinterval)
	updateChannel := make(chan struct{}, 10)
	for {
		select {
		case <-controller.parent.Context().Done():
			log.Trace().Err(controller.parent.Context().Err()).Msg("exiting usercontroller run()")
			return
		case <-updateChannel:
			if controller.IsLoggedIn() {
				controller.loginKeepAlive()
			} else {
				controller.login()
			}
			controller.updateUsers()
		case <-ticker.C:
			updateChannel <- struct{}{}
		case <-controller.usernameupdated:
			name := controller.GetUser().Name
			if name != controller.parent.settings.Username {
				controller.mutex.Lock()
				controller.user.Name = controller.parent.settings.Username
				controller.mutex.Unlock()
				updateChannel <- struct{}{}
				log.Debug().Msg("updated user")
			}
		}
	}
}

func (controller *UserController) updateUsers() {
	ips, err := controller.parent.client.User.GetUsers()
	if err != nil {
		controller.mutex.Lock()
		controller.err = err
		controller.mutex.Unlock()
		log.Error().Err(err).Msg("error retrieving userlist from server")
		return
	}
	users := user.Users{}
	for _, ip := range ips {
		user, err := controller.parent.client.User.GetUser(ip)
		if err != nil {
			controller.mutex.Lock()
			controller.err = err
			controller.mutex.Unlock()
			log.Error().Err(err).Str("ip", ip).Msg("error retrieving user from server")
			return
		}
		err = users.Add(user)
		if err != nil {
			controller.mutex.Lock()
			controller.err = err
			controller.mutex.Unlock()
			log.Error().Err(err).Str("ip", ip).Msg("error adding user to temporary userlist")
			return
		}
	}
	controller.mutex.RLock()
	localUsers := controller.users
	controller.mutex.RUnlock()
	if localUsers.Equal(users) {
		return
	}
	controller.mutex.Lock()
	controller.users = users
	controller.mutex.Unlock()
	log.Debug().Interface("users", users).Msg("updated users in usercontroller")
	controller.notifySubcriber()
}

func (controller *UserController) loginKeepAlive() {
	controller.mutex.RLock()
	user := controller.user
	controller.mutex.RUnlock()
	_, err := controller.parent.client.User.UpdateUser(user)
	if err != nil {
		controller.mutex.Lock()
		controller.loggedIn = false
		controller.mutex.Unlock()
		log.Error().Err(err).Msg("error updating user at server - try to login again")
	}
}

func (controller *UserController) login() {
	user, err := controller.parent.client.User.CreateNewUser(controller.user)
	if err != nil {
		log.Error().Err(err).Msgf("could not login at server - retries exceeded")
		return
	}
	log.Debug().Interface("user", user).Msg("user logged in")
	controller.mutex.Lock()
	controller.user = user
	controller.loggedIn = true
	controller.mutex.Unlock()
}

func (controller *UserController) Subscribe(subscriber chan struct{}) {
	defer controller.mutex.Unlock()
	controller.mutex.Lock()
	controller.subscriber = append(controller.subscriber, subscriber)
}

func (controller *UserController) Unsubscribe(subscriber chan struct{}) {
	defer controller.mutex.Unlock()
	controller.mutex.Lock()
	index := slices.Index(controller.subscriber, subscriber)
	_ = slices.Delete(controller.subscriber, index, index+1)
}

func (controller *UserController) notifySubcriber() {
	defer controller.mutex.RUnlock()
	controller.mutex.RLock()
	for _, subscriber := range controller.subscriber {
		util.ChannelWriteNonBlocking(subscriber, struct{}{})
	}
}
