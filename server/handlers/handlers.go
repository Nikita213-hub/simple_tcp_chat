package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Nikita213-hub/simple_tcp_chat/server/cmd/state"

	message "github.com/Nikita213-hub/simple_tcp_chat/pb"
	"github.com/Nikita213-hub/simple_tcp_chat/server/room"
	"github.com/Nikita213-hub/simple_tcp_chat/server/user"
	"github.com/Nikita213-hub/simple_tcp_chat/server/util"
	"google.golang.org/protobuf/proto"
)

// TODO: add state with USERS, ROOMS, and mutexes
// TODO: put simple_tcp_chat to Nikita213-hub folder
// TODO: think about leave room, create room and some els funcs where to place them
// TODO: create mb entity fold for User and Room entities
// TODO: using channels add functionality for broadcasting some warnings for users
// TODO: complete othre todos

func UserMessageHandler(user *user.User, state *state.GlobalState) {
	defer user.Conn.Close()
	defer room.LeaveRoom(user, &state.ROOMS, state.RoomsMx, state.UsersMx)
	defer delete(state.USERS, user.Id)
	for {
		msgBuffer := make([]byte, 1024)
		user.Conn.Write([]byte(">> "))
		n, err := user.Conn.Read(msgBuffer)
		if err != nil {
			fmt.Println("AAAAAAAAAA")
			return
			// idk what to do in that case
		}
		msgp := message.ChatMessage{}
		proto.Unmarshal(msgBuffer[:n], &msgp)
		fmt.Println(msgp)
		message := strings.TrimSuffix(string(msgBuffer[:n]), "\n")
		switch message {
		case "/new_room":
			roomPassword, err := util.GetPassword(user.Conn)
			if err != nil {
				fmt.Println(err)
				return
			}
			newRoom, err := room.CreateRoom(user, roomPassword, &state.ROOMS, state.RoomsMx)
			if err != nil {
				user.Conn.Write([]byte(err.Error()))
			}
			state.UsersMx.Lock()
			user.Current_room = newRoom.Id
			state.UsersMx.Unlock()
			user.Conn.Write([]byte("Room (id: " + strconv.FormatInt(int64(newRoom.Id), 10) + ") was successfully created\n"))
			fmt.Println(user)
		case "/exit":
			user.Conn.Write([]byte("Exit from app...\n"))
			return
		case "/connect":
			//TODO: move that shit in a function
			user.Conn.Write([]byte("Insert chat id:\n"))
			msgBuffer := make([]byte, 16)
			n, err := user.Conn.Read(msgBuffer)
			if err != nil {
				fmt.Println(err)
				return
			}
			roomId := strings.TrimSuffix(string(msgBuffer[:n]), "\n")
			roomIdInt, err := strconv.Atoi(roomId)
			if err != nil {
				fmt.Println(err)
				return
			}
			user.Conn.Write([]byte("Insert chat password:\n"))
			pswdBuffer := make([]byte, 16)
			n, err = user.Conn.Read(pswdBuffer)
			if err != nil {
				fmt.Println(err)
				return
			}
			roomPswd := strings.TrimSuffix(string(pswdBuffer[:n]), "\n")
			err = room.ConnectToRoom(user, roomPswd, roomIdInt, &state.ROOMS, state.RoomsMx)
			if err != nil {
				user.Conn.Write([]byte(err.Error()))
			} else {
				state.UsersMx.Lock()
				user.Current_room = roomIdInt
				state.UsersMx.Unlock()
			}
		case "/leave_room":
			err = room.LeaveRoom(user, &state.ROOMS, state.RoomsMx, state.UsersMx)
			if err != nil {
				fmt.Println(err)
			}
			user.Conn.Write([]byte("You have leaved room"))
		default:
			if user.Current_room == -1 {
				user.Conn.Write([]byte("Incorrect command\n"))
			} else {
				r, ok := state.ROOMS[user.Current_room]
				if !ok {
					fmt.Println(err)
				} else {
					r.SendMessage(user, message)
				}
			}
		}
	}
}
