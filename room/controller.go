package room

import (
	"errors"
	"sync"
	"time"

	"github.com/Nikita213-hub/simple_tcp_chat/user"
)

func CreateRoom(roomInitiator *user.User, password string, ROOMS *map[int]*Room, roomsMutex *sync.RWMutex) (*Room, error) {
	users := make([]*user.User, 0, 5)
	users = append(users, roomInitiator)
	roomId := int(time.Now().Unix())
	newRoom := &Room{
		Id:       roomId,
		Password: password,
		Users:    users,
	}
	roomsMutex.Lock()
	(*ROOMS)[roomId] = newRoom
	roomsMutex.Unlock()
	return newRoom, nil
}

func ConnectToRoom(user *user.User, password string, roomId int, ROOMS *map[int]*Room, roomsMutex *sync.RWMutex) error {
	roomsMutex.Lock()
	defer roomsMutex.Unlock()
	if user.Current_room != -1 {
		return errors.New("You are already in a chat")
	}
	room, ok := (*ROOMS)[roomId]
	if !ok {
		return errors.New("Incorrect chat id")
	}
	if room.Password != password {
		return errors.New("Incorrect password")
	}
	room.Users = append(room.Users, user)
	return nil
}

func LeaveRoom(user *user.User, ROOMS *map[int]*Room, roomsMutex, usersMutex *sync.RWMutex) error {
	if user.Current_room == -1 {
		return errors.New("You are not in any chat")
	}
	roomsMutex.Lock()
	usersMutex.Lock()
	room, _ := (*ROOMS)[user.Current_room]
	var user_ind int
	for i, usr := range room.Users {
		if usr.Id == user.Id {
			user_ind = i
		}
	}
	if len(room.Users)-1 < user_ind {
		room.Users = append(room.Users[:user_ind], room.Users[user_ind+1:]...)
	} else if len(room.Users)-1 == user_ind {
		room.Users = room.Users[:user_ind]
	}
	roomsMutex.Unlock()
	user.Current_room = -1
	usersMutex.Unlock()
	return nil
}
