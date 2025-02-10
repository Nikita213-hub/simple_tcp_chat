package state

import (
	"sync"

	"github.com/Nikita213-hub/simple_tcp_chat/server/room"
	"github.com/Nikita213-hub/simple_tcp_chat/server/user"
)

type GlobalState struct {
	USERS   map[int]*user.User
	ROOMS   map[int]*room.Room
	UsersMx *sync.RWMutex
	RoomsMx *sync.RWMutex
}
