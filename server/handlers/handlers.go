package handlers

import (
	"fmt"
	"io"
	"strconv"

	"github.com/Nikita213-hub/simple_tcp_chat/server/cmd/state"
	"github.com/Nikita213-hub/simple_tcp_chat/server/room"
	"github.com/Nikita213-hub/simple_tcp_chat/server/user"

	"github.com/Nikita213-hub/simple_tcp_chat/server/messagesController"
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
		message, err := messagesController.ProcessUserMessage(user.Conn)
		if err != nil {
			if err == io.EOF {
				fmt.Println("User has left")
				return
			}
			fmt.Println(err.Error())
		}
		fmt.Println(message)
		switch message {
		case "/new_room":
			messagesController.SendNotificationMessage(user.Conn, "insert room password\n")
			roomPassword, err := messagesController.ProcessUserMessage(user.Conn)
			if err != nil {
				messagesController.SendErrorMessage(user.Conn, "error occured while creating password\n")
			}
			if len(roomPassword) < 1 {
				messagesController.SendErrorMessage(user.Conn, "insert at least 1 character\n")
			}
			newRoom, err := room.CreateRoom(user, roomPassword, &state.ROOMS, state.RoomsMx)
			if err != nil {
				user.Conn.Write([]byte(err.Error()))
			}
			state.UsersMx.Lock()
			user.Current_room = newRoom.Id
			state.UsersMx.Unlock()
			// user.Conn.Write([]byte("Room (id: " + strconv.FormatInt(int64(newRoom.Id), 10) + ") was successfully created\n"))
			messagesController.SendNotificationMessage(user.Conn, "Room (id: "+strconv.FormatInt(int64(newRoom.Id), 10)+") was successfully created\n")
			fmt.Println(user)
		case "/exit":
			messagesController.SendNotificationMessage(user.Conn, "Exit from app...\n")
			return
			// user.Conn.Write([]byte("Exit from app...\n"))
		case "/connect":
			//TODO: move that shit in a function
			// user.Conn.Write([]byte("Insert chat id:\n"))
			messagesController.SendNotificationMessage(user.Conn, "Insert chat id:\n")
			roomId, err := messagesController.ProcessUserMessage(user.Conn)
			if err != nil {
				messagesController.SendErrorMessage(user.Conn, err.Error())
			}
			roomIdInt, err := strconv.Atoi(roomId)
			if err != nil {
				messagesController.SendErrorMessage(user.Conn, err.Error())
			}
			// user.Conn.Write([]byte("Insert chat password:\n"))
			messagesController.SendNotificationMessage(user.Conn, "Insert chat password:\n")
			roomPswd, err := messagesController.ProcessUserMessage(user.Conn)
			if err != nil {
				messagesController.SendErrorMessage(user.Conn, err.Error())
			}
			err = room.ConnectToRoom(user, roomPswd, roomIdInt, &state.ROOMS, state.RoomsMx)
			if err != nil {
				messagesController.SendErrorMessage(user.Conn, err.Error())
				// user.Conn.Write([]byte(err.Error()))
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
			// user.Conn.Write([]byte("You have leaved room"))
			messagesController.SendNotificationMessage(user.Conn, "You have leaved room\n")
		default:
			if user.Current_room == -1 {
				// user.Conn.Write([]byte("Incorrect command\n"))
				messagesController.SendErrorMessage(user.Conn, "Incorrect command\n")
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
