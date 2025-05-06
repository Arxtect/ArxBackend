package ws

import (
	"github.com/Arxtect/ArxBackend/golangp/common/constants"
	"errors"

	"github.com/toheart/functrace"
)

var roomService roomSvc

type roomSvc struct {
	Rooms map[string]*room
}

type room struct {
	ID      string
	Content string
}

type EditFileRoomAccess interface {
	NewRoom() (room, error)
	Get(id string) (string, error)
	Update(id string, content string) bool
	Delete(id string) bool
}

func (rs *roomSvc) NewRoom(fileId string) (room, error) {
	defer functrace.Trace([]interface {
	}{rs, fileId})()
	if len(rs.Rooms) >= constants.MaxRooms {
		return room{}, errors.New("too many open rooms")
	}

	newRoom := room{fileId, ""}
	if rs.Rooms == nil {
		rs.Rooms = make(map[string]*room, constants.MaxRooms)
	}
	rs.Rooms[newRoom.ID] = &newRoom
	return newRoom, nil
}

func (rs *roomSvc) Update(id string, content string) bool {
	defer functrace.Trace([]interface {
	}{rs, id, content})()
	r, ok := rs.Rooms[id]
	if !ok {
		return false
	}
	r.Content = content
	return true
}

func (rs *roomSvc) Get(id string) (string, error) {
	defer functrace.Trace([]interface {
	}{rs, id})()
	if rs.Rooms == nil {
		return "", errors.New("room is empty")
	}

	r, ok := rs.Rooms[id]
	if !ok {
		return "", errors.New("room not found")
	}
	return r.Content, nil
}

func (rs *roomSvc) Delete(id string) bool {
	defer functrace.Trace([]interface {
	}{rs, id})()
	if rs.Rooms == nil {
		return false
	}

	_, ok := rs.Rooms[id]
	if !ok {
		return false
	}

	delete(rs.Rooms, id)
	return true
}
