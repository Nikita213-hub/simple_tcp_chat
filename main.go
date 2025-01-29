package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

var USERS []User

var ROOMS map[int]Room

type Room struct {
	id       int
	password string
	users    []int
}

type User struct {
	id           int
	nickname     string
	current_room int
}

func createRoom(roomInitiator int) (Room, error) {
	users := make([]int, 5)
	users = append(users, roomInitiator)
	password, err := getPassword()
	if err != nil {
		fmt.Println(err)
		return Room{}, err
	}

	return Room{
		id:       int(time.Now().Unix()),
		password: password,
		users:    users,
	}, nil
}

func createUser(nickname string) User {
	return User{
		id:           int(time.Now().Unix()),
		nickname:     nickname,
		current_room: -1,
	}
}

func getPassword() (string, error) {
	ioStream := os.Stdin
	reader := bufio.NewReader(ioStream)
	fmt.Println("Insert room's password")
	password, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	if len(password) < 1 {
		return "", errors.New("Insert smth lol")
	}
	return password, nil
}

func getPort() (string, error) {
	ioStream := os.Stdin
	reader := bufio.NewReader(ioStream)
	fmt.Println("Insert port for serving")
	var port string
	port, err := reader.ReadString('\n')
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

func main() {
	fmt.Println("It's a simple_tcp_chat server")
	port, err := getPort()
	if err != nil {
		fmt.Println(err)
		return
	}
	strAddr := "localhost:" + strings.TrimSuffix(string(port), "\n")
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
		conn, err := server.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		defer conn.Close()
		conn.Write([]byte("Hello niga"))
	}
}
