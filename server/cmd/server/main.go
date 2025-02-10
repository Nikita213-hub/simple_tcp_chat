package main

import (
	"fmt"
	"net"
	"sync"

	"github.com/Nikita213-hub/simple_tcp_chat/server/cmd/state"
	"github.com/Nikita213-hub/simple_tcp_chat/server/handlers"
	"github.com/Nikita213-hub/simple_tcp_chat/server/room"
	"github.com/Nikita213-hub/simple_tcp_chat/server/user"
	"github.com/Nikita213-hub/simple_tcp_chat/server/util"
)

func main() {

	usersMx := &sync.RWMutex{}
	roomsMx := &sync.RWMutex{}
	s := state.GlobalState{
		USERS:   map[int]*user.User{},
		ROOMS:   map[int]*room.Room{},
		UsersMx: usersMx,
		RoomsMx: roomsMx,
	}

	fmt.Println("It's a simple_tcp_chat server")
	port, err := util.GetPort()
	if err != nil {
		fmt.Println(err)
		return
	}
	strAddr := "localhost:" + port
	fmt.Println(strAddr)
	addr, err := net.ResolveTCPAddr("tcp", strAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	server, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Server is listening " + strAddr + " ...")
	for {
		conn, err := server.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			return
		}
		conn.Write([]byte("Insert your nickname\n"))
		newUser, err := user.CreateUser(conn, &s.USERS, s.UsersMx)
		if err != nil {
			fmt.Println(err)
		}
		conn.Write([]byte("Hello, " + newUser.Nickname + "\n"))
		go handlers.UserMessageHandler(newUser, &s)
	}
}
