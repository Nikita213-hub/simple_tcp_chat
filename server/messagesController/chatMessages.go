package messagesController

import (
	"errors"

	communicationEntities "github.com/Nikita213-hub/chat_proto"
	"github.com/Nikita213-hub/simple_tcp_chat/server/user"
	proto "google.golang.org/protobuf/proto"
)

func SendChatMessage(sender *user.User, receiver *user.User, content string) error {
	userMsg := &communicationEntities.User{
		Id:   int32(sender.Id),
		Name: sender.Nickname,
	}
	msg := &communicationEntities.ChatMessage{
		Sender:  userMsg,
		Content: content,
	}
	msgWrapper := communicationEntities.WrapperMessage{
		Msg: &communicationEntities.WrapperMessage_Cm{
			Cm: msg,
		},
	}
	messageb, err := proto.Marshal(&msgWrapper)
	if err != nil {
		return errors.New("something went wrong while processing message")
	}
	receiver.Conn.Write(messageb)
	return nil
}
