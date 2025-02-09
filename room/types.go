package room

import (
	"github.com/Nikita213-hub/simple_tcp_chat/user"
)

type Room struct {
	Id       int
	Password string
	Users    []*user.User
}

func (r *Room) SendMessage(sender *user.User, message string) {
	for _, usr := range r.Users {
		if usr.Id == sender.Id {
			continue
		}
		go func(usr *user.User) {
			_, err := usr.Conn.Write([]byte(sender.Nickname + ": " + message + "\n"))
			if err != nil {
				sender.Conn.Write([]byte("Error occured when sending message"))
			}
		}(usr)
	}
}
