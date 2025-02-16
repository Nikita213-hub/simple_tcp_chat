package room

import (
	"github.com/Nikita213-hub/simple_tcp_chat/server/messagesController"
	"github.com/Nikita213-hub/simple_tcp_chat/server/user"
)

type Room struct {
	Id       int
	Password string
	Users    []*user.User
}

func (r *Room) SendMessage(sender *user.User, messageContent string) {
	for _, usr := range r.Users {
		if usr.Id == sender.Id {
			continue
		}
		go func(usr *user.User) {
			err := messagesController.SendChatMessage(sender, usr, messageContent)
			if err != nil {
				messagesController.SendErrorMessage(sender.Conn, err.Error())
			}
		}(usr)
	}
}
