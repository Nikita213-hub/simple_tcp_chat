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

func GetPort() (string, error) {
	ioStream := os.Stdin
	reader := bufio.NewReader(ioStream)
	fmt.Println("Insert port for connetction")
	var port string
	port, err := reader.ReadString('\n')
	port = strings.TrimSuffix(port, "\n")
	if err != nil {
		return "", err
	}
	if len(port) < 1 {
		return "", errors.New("insert smth lol")
	}
	return port, nil
}

func main() {
	fmt.Println("It's a tcp chat client")
	port, err := GetPort()
	if err != nil {
		fmt.Println(err)
	}
	addr, err := net.ResolveTCPAddr("tcp", "localhost:"+port)
	if err != nil {
		fmt.Println(err)
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Println(err)
	}
	for {
		//request
		conn.Write([]byte("OLA AMIGO"))
		//response
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println(err)
		}
		msg := buffer[:n]
		fmt.Println(string(msg))
		time.Sleep(5 * time.Second)
	}

}
