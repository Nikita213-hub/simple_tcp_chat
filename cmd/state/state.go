package state

import (
	"sync"

	"github.com/Nikita213-hub/simple_tcp_chat/room"
	"github.com/Nikita213-hub/simple_tcp_chat/user"
)

type GlobalState struct {
	USERS   map[int]*user.User
	ROOMS   map[int]*room.Room
	UsersMx *sync.RWMutex
	RoomsMx *sync.RWMutex
}
