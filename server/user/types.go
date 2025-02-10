package user

import "net"

type User struct {
	Id           int
	Nickname     string
	Current_room int
	Conn         *net.TCPConn
}
