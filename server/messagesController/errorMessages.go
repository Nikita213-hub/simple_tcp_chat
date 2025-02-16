package messagesController

import (
	"errors"
	"net"

	communicationEntities "github.com/Nikita213-hub/chat_proto"
	proto "google.golang.org/protobuf/proto"
)

func SendErrorMessage(conn *net.TCPConn, content string) error {
	msg := &communicationEntities.ErrorMessage{
		Content: content,
	}
	msgWrapper := communicationEntities.WrapperMessage{
		Msg: &communicationEntities.WrapperMessage_Em{
			Em: msg,
		},
	}
	messageb, err := proto.Marshal(&msgWrapper)
	if err != nil {
		return errors.New("something went wrong while processing message")
	}
	conn.Write(messageb)
	return nil
}
