package messagesController

import (
	"net"
	"strings"

	communicationEntities "github.com/Nikita213-hub/chat_proto"
	proto "google.golang.org/protobuf/proto"
)

func ProcessUserMessage(conn *net.TCPConn) (string, error) {
	msgBuffer := make([]byte, 1024)
	n, err := conn.Read(msgBuffer)
	if err != nil {
		return "", err
	}
	unmMessage := communicationEntities.ChatMessage{}
	err = proto.Unmarshal(msgBuffer[:n], &unmMessage)
	if err != nil {
		return "", err
	}
	msg := strings.TrimSuffix(unmMessage.Content, "\n")
	return msg, nil
}
