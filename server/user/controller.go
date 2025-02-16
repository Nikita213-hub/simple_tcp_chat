package user

import (
	"fmt"
	"net"
	"sync"
	"time"
	// "github.com/Nikita213-hub/simple_tcp_chat/server/util"
)

func CreateUser(conn *net.TCPConn, USERS *map[int]*User, usersMutex *sync.RWMutex, nickname string) (*User, error) {
	userId := int(time.Now().Unix())
	// escapes
	newUser := &User{
		Id:           userId,
		Nickname:     nickname,
		Current_room: -1,
		Conn:         conn,
	}
	fmt.Println("New user created")
	usersMutex.Lock()
	(*USERS)[userId] = newUser
	usersMutex.Unlock()
	return newUser, nil
}
