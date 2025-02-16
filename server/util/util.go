package util

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

// func GetNickname(conn *net.TCPConn) (string, error) {
// 	readBuffer := make([]byte, 30)

// 	n, err := conn.Read(readBuffer)
// 	if err != nil {
// 		return "", err
// 	}
// 	nickname := strings.TrimSuffix(string(readBuffer[:n]), "\n")

// 	if len(nickname) < 1 {
// 		return "", errors.New("insert smth lol")
// 	}
// 	return nickname, nil
// }

// func GetPassword(conn *net.TCPConn) (string, error) {
// 	if len(password) < 1 {
// 		return "", errors.New("insert smth lol\n")
// 	}
// 	return password, nil
// }

func GetPort() (string, error) {
	ioStream := os.Stdin
	reader := bufio.NewReader(ioStream)
	fmt.Println("insert port for serving")
	var port string
	port, err := reader.ReadString('\n')
	port = strings.TrimSuffix(port, "\n")
	if err != nil {
		return "", err
	}
	if len(port) < 1 {
		return "", errors.New("insert smth lol")
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
