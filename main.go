package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	usersMutex sync.RWMutex
	roomsMutex sync.RWMutex
)

var USERS = map[int]User{}

var ROOMS = map[int]Room{}

type Room struct {
	id       int
	password string
	users    []User
}

type User struct {
	id           int
	nickname     string
	current_room int
	conn         *net.TCPConn
}

func createRoom(roomInitiator User, password string) (Room, error) {
	users := make([]User, 0, 5)
	users = append(users, roomInitiator)
	roomId := int(time.Now().Unix())
	newRoom := Room{
		id:       roomId,
		password: password,
		users:    users,
	}
	roomsMutex.Lock()
	ROOMS[roomId] = newRoom
	roomsMutex.Unlock()
	return newRoom, nil
}

func createUser(conn *net.TCPConn) (User, error) {
	nickname, err := getNickname(conn)
	if err != nil {
		return User{}, err
	}
	userId := int(time.Now().Unix())
	newUser := User{
		id:           userId,
		nickname:     nickname,
		current_room: -1,
		conn:         conn,
	}
	fmt.Println("New user created")
	usersMutex.Lock()
	USERS[userId] = newUser
	usersMutex.Unlock()
	return newUser, nil
}

func connectToRoom(user User, password string, roomId int) error {
	roomsMutex.Lock()
	defer roomsMutex.Unlock()
	if user.current_room != -1 {
		return errors.New("You are already in a chat")
	}
	room, ok := ROOMS[roomId]
	if !ok {
		return errors.New("Incorrect chat id")
	}
	if room.password != password {
		return errors.New("Incorrect password")
	}
	room.users = append(room.users, user)
	ROOMS[roomId] = room
	return nil
}

func getNickname(conn *net.TCPConn) (string, error) {
	readBuffer := make([]byte, 30)

	n, err := conn.Read(readBuffer)
	if err != nil {
		return "", err
	}
	nickname := strings.TrimSuffix(string(readBuffer[:n]), "\n")

	if len(nickname) < 1 {
		return "", errors.New("Insert smth lol")
	}
	return nickname, nil
}

func getPassword(conn *net.TCPConn) (string, error) {
	conn.Write([]byte("Insert room's password:\n"))
	pswrdBuffer := make([]byte, 15)
	n, err := conn.Read(pswrdBuffer)
	if err != nil {
		return "", err
	}
	password := strings.TrimSuffix(string(pswrdBuffer[:n]), "\n")
	if len(password) < 1 {
		return "", errors.New("Insert smth lol\n")
	}
	return password, nil
}

func getPort() (string, error) {
	ioStream := os.Stdin
	reader := bufio.NewReader(ioStream)
	fmt.Println("Insert port for serving")
	var port string
	port, err := reader.ReadString('\n')
	port = strings.TrimSuffix(port, "\n")
	if err != nil {
		return "", err
	}
	if len(port) < 1 {
		return "", errors.New("Insert smth lol")
	}
	// TODO: add implementation of asking for port again if smth went wrong(use channels)
	// for {
	// 	port, err := reader.ReadString('\n')
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		continue
	// 	}
	// 	fmt.Println(len(port))
	// 	if len(port) > 0 {
	// 		ioStream.Close()
	// 		break
	// 	}
	// 	fmt.Println("Insert smth lol")
	// }
	return port, nil
}

// TODO: add errror handling
func (r *Room) SendMessage(sender User, message string) {
	for _, usr := range r.users {
		if usr.id == sender.id {
			continue
		}
		go func(usr User) {
			_, err := usr.conn.Write([]byte(sender.nickname + ": " + message + "\n"))
			if err != nil {
				sender.conn.Write([]byte("Error occured when sending message"))
			}
		}(usr)
	}
}

// TODO: add error handling
func LeaveRoom(user User) error {
	if user.current_room == -1 {
		return errors.New("You are not in any chat")
	}
	room, _ := ROOMS[user.current_room]
	var user_ind int
	for i, usr := range room.users {
		if usr.id == user.id {
			user_ind = i
		}
	}
	if len(room.users)-1 < user_ind {
		room.users = append(room.users[:user_ind], room.users[user_ind+1:]...)
	} else if len(room.users)-1 == user_ind {
		room.users = room.users[:user_ind]
	}
	roomsMutex.Lock()
	ROOMS[user.current_room] = room
	roomsMutex.Unlock()
	fmt.Println(room)
	user.current_room = -1
	usersMutex.Lock()
	USERS[user.id] = user
	usersMutex.Unlock()
	return nil
}

func userMessageHandler(user User) {
	for {
		defer user.conn.Close()
		defer LeaveRoom(user)
		msgBuffer := make([]byte, 1024)
		user.conn.Write([]byte(">> "))
		n, err := user.conn.Read(msgBuffer)
		if err != nil {
			fmt.Println("AAAAAAAAAA")
			return
			// idk what to do in that case
		}
		message := strings.TrimSuffix(string(msgBuffer[:n]), "\n")
		fmt.Println(message)
		switch message {
		case "/new_room":
			roomPassword, err := getPassword(user.conn)
			if err != nil {
				fmt.Println(err)
				return
			}
			newRoom, err := createRoom(user, roomPassword)
			if err != nil {
				user.conn.Write([]byte(err.Error()))
			}
			user.current_room = newRoom.id
			usersMutex.Lock()
			USERS[user.id] = user
			usersMutex.Unlock()
			user.conn.Write([]byte("Room (id: " + strconv.FormatInt(int64(newRoom.id), 10) + ") was successfully created\n"))
		case "/exit":
			user.conn.Write([]byte("Exit from app...\n"))
			return
		case "/connect":
			//TODO: move that shit in a function
			user.conn.Write([]byte("Insert chat id:\n"))
			msgBuffer := make([]byte, 16)
			n, err := user.conn.Read(msgBuffer)
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
			user.conn.Write([]byte("Insert chat password:\n"))
			pswdBuffer := make([]byte, 16)
			n, err = user.conn.Read(pswdBuffer)
			if err != nil {
				fmt.Println(err)
				return
			}
			roomPswd := strings.TrimSuffix(string(pswdBuffer[:n]), "\n")
			err = connectToRoom(user, roomPswd, roomIdInt)
			if err != nil {
				user.conn.Write([]byte(err.Error()))
			} else {
				user.current_room = roomIdInt
				usersMutex.Lock()
				USERS[user.id] = user
				usersMutex.Unlock()
			}
		case "/leave_room":
			err = LeaveRoom(user)
			if err != nil {
				fmt.Println(err)
			}
			user.conn.Write([]byte("You have leaved room"))
		default:
			if user.current_room == -1 {
				user.conn.Write([]byte("Incorrect command\n"))
			} else {
				r, ok := ROOMS[user.current_room]
				if !ok {
					fmt.Println(err)
				} else {
					r.SendMessage(user, message)
				}
			}
		}
	}
}

func main() {
	fmt.Println("It's a simple_tcp_chat server")
	port, err := getPort()
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
		newUser, err := createUser(conn)
		if err != nil {
			fmt.Println(err)
		}
		conn.Write([]byte("Hello, " + newUser.nickname + "\n"))
		go userMessageHandler(newUser)
	}
}
