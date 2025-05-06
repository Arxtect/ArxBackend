package ws

import (
	"github.com/Arxtect/ArxBackend/golangp/apps/arx_center/models"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/toheart/functrace"
	"golang.org/x/net/websocket"
)

var ssList subscriptionSvc

type subscriptionSvc struct {
	Subscribers map[string][]subscription
	RoomOwners  map[string][]string
	mu          sync.RWMutex
}

type subscription struct {
	UserInfo     *models.User
	LastActivity time.Time
	Connection   *websocket.Conn
}

type Subscriber interface {
	All() map[string][]subscription
	Subscribe(roomID string, s subscription) error
	Unsubscribe(roomID string, s subscription) error
	GetSubscribers(roomID string) ([]subscription, error)
	IsOwner(roomID, userID string) bool
	AddUserAsOwner(roomID, userID string) error
	RemoveUserAsOwner(roomID, userID string) error

	DeleteRoomSubscribers(roomID string) error

	BroadcastContentUpdates(roomID string, content string)
	StartSubscriberListener()
}

func GetSubscriber() Subscriber {
	defer functrace.Trace([]interface {
	}{})()
	return &ssList
}

func (subSvc *subscriptionSvc) All() map[string][]subscription {
	defer functrace.Trace([]interface {
	}{subSvc})()
	subSvc.mu.RLock()
	defer subSvc.mu.RUnlock()
	return subSvc.Subscribers
}

func (subSvc *subscriptionSvc) Subscribe(roomID string, s subscription) error {
	defer functrace.Trace([]interface {
	}{subSvc, roomID, s})()
	subSvc.mu.Lock()
	defer subSvc.mu.Unlock()
	if subSvc.Subscribers == nil {
		subSvc.Subscribers = make(map[string][]subscription)
	}
	existingSubs, ok := subSvc.Subscribers[roomID]
	if ok {

		subSvc.Subscribers[roomID] = append(existingSubs, s)
	} else {

		subSvc.Subscribers[roomID] = []subscription{s}
	}
	return nil
}

func (subSvc *subscriptionSvc) GetSubscribers(roomID string) ([]subscription, error) {
	defer functrace.Trace([]interface {
	}{subSvc, roomID})()
	subSvc.mu.RLock()
	defer subSvc.mu.RUnlock()
	if subSvc.Subscribers == nil {
		return []subscription{}, nil
	}
	existingSubs, ok := subSvc.Subscribers[roomID]
	if ok {
		return existingSubs, nil
	}

	return []subscription{}, nil
}

func (subSvc *subscriptionSvc) Unsubscribe(roomID string, s subscription) error {
	defer functrace.Trace([]interface {
	}{subSvc, roomID, s})()
	subSvc.mu.Lock()
	defer subSvc.mu.Unlock()
	if subSvc.Subscribers == nil {
		return nil
	}
	existingSubs, ok := subSvc.Subscribers[roomID]
	if !ok {
		return errors.New("room not found")
	}
	for index, ex := range existingSubs {
		if s.Connection == ex.Connection && s.UserInfo.ID.String() == ex.UserInfo.ID.String() && s.UserInfo.Name == ex.UserInfo.Name {

			subSvc.Subscribers[roomID] = append(existingSubs[:index], existingSubs[index+1:]...)
			return nil
		}
	}
	return errors.New("could not find the subscription")
}

func (subSvc *subscriptionSvc) IsOwner(roomID, userID string) bool {
	defer functrace.Trace([]interface {
	}{subSvc, roomID, userID})()
	subSvc.mu.RLock()
	defer subSvc.mu.RUnlock()
	ownerIDs, ok := subSvc.RoomOwners[roomID]
	if !ok {
		return false
	}
	for _, id := range ownerIDs {
		if id == userID {
			return true
		}
	}
	return false
}

func (subSvc *subscriptionSvc) AddUserAsOwner(roomID, userID string) error {
	defer functrace.Trace([]interface {
	}{subSvc, roomID, userID})()
	subSvc.mu.Lock()
	defer subSvc.mu.Unlock()
	if subSvc.RoomOwners == nil {
		subSvc.RoomOwners = make(map[string][]string)
	}
	ownerIDs, ok := subSvc.RoomOwners[roomID]
	if ok {
		subSvc.RoomOwners[roomID] = append(ownerIDs, userID)
	} else {
		subSvc.RoomOwners[roomID] = []string{userID}
	}
	return nil
}

func (subSvc *subscriptionSvc) RemoveUserAsOwner(roomID, userID string) error {
	defer functrace.Trace([]interface {
	}{subSvc, roomID, userID})()
	subSvc.mu.Lock()
	defer subSvc.mu.Unlock()
	ownerIDs, ok := subSvc.RoomOwners[roomID]
	if !ok {
		return errors.New("room not found")
	}
	for i, id := range ownerIDs {
		if id == userID {
			subSvc.RoomOwners[roomID] = append(ownerIDs[:i], ownerIDs[i+1:]...)
			return nil
		}
	}
	return errors.New("could not find the user as owner")
}

func (subSvc *subscriptionSvc) DeleteRoomSubscribers(roomID string) error {
	defer functrace.Trace([]interface {
	}{subSvc, roomID})()
	subSvc.mu.Lock()
	defer subSvc.mu.Unlock()

	delete(subSvc.Subscribers, roomID)
	delete(subSvc.RoomOwners, roomID)

	return nil
}

type RoomContentUpdate struct {
	RoomID  string
	Content string
}

var broadcast = make(chan RoomContentUpdate, 100)

func (subSvc *subscriptionSvc) BroadcastContentUpdates(roomID string, content string) {
	defer functrace.Trace([]interface {
	}{subSvc, roomID, content})()
	update := RoomContentUpdate{RoomID: roomID, Content: content}
	broadcast <- update
}

func (subSvc *subscriptionSvc) StartSubscriberListener() {
	defer functrace.Trace([]interface {
	}{subSvc})()
	for {
		select {
		case update := <-broadcast:
			subSvc.mu.RLock()
			subs, ok := subSvc.Subscribers[update.RoomID]
			subSvc.mu.RUnlock()
			if ok {

				for _, sub := range subs {
					go func(s subscription) {

						err := websocket.Message.Send(s.Connection, update.Content)
						if err != nil {

							fmt.Printf("Error sending update to user %s: %v\n", s.UserInfo.ID, err)
						}
					}(sub)
				}
			}
		}
	}
}
