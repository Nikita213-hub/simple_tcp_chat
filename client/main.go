package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"

	communicationEntities "github.com/Nikita213-hub/chat_proto"
	tea "github.com/charmbracelet/bubbletea"
	proto "google.golang.org/protobuf/proto"
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

func processResponse(conn *net.TCPConn) (string, interface{}, error) {
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println(err)
	}

	msgWrapper := communicationEntities.WrapperMessage{}
	if err := proto.Unmarshal(buffer[:n], &msgWrapper); err != nil {
		return "", struct{}{}, err
	}
	switch {
	case msgWrapper.GetCm() != nil:
		chatMsg := msgWrapper.GetCm()
		return "CM", chatMsg, nil
	case msgWrapper.GetEm() != nil:
		errMsg := msgWrapper.GetEm()
		return "EM", errMsg, nil
	case msgWrapper.GetNm() != nil:
		notification := msgWrapper.GetNm()
		return "NM", notification, nil
	default:
		return "", struct{}{}, errors.New("unknown message type")
	}
}

func receiveMessages(p *tea.Program, conn *net.TCPConn, wg *sync.WaitGroup, stop chan<- struct{}, resume chan<- struct{}) {
	defer wg.Done()
	for {
		msgType, msg, err := processResponse(conn)
		if err != nil {
			fmt.Println(err)
		}
		switch msgType {
		case "CM":
			msgTyped := msg.(*communicationEntities.ChatMessage)
			p.Send(msgTyped.Sender.Name + ": " + msgTyped.Content)
		case "EM":
			p.Send(msg.(*communicationEntities.ErrorMessage).Content)
		case "NM":
			p.Send(msg.(*communicationEntities.Notification).Content)
		}
	}
}

type model struct {
	messages    []string
	input       string
	input_chars int
	conn        *net.TCPConn
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Send message (simulate)
			m.messages = append(m.messages, "You: "+m.input)
			msgs := communicationEntities.ChatMessage{
				Content: string(m.input),
			}
			msgb, err := proto.Marshal(&msgs)
			if err != nil {
				fmt.Println(err)
			}
			m.conn.Write(msgb)
			m.input = ""
			m.input_chars = 0
			return m, nil
		case "ctrl+c":
			return m, tea.Quit
		case "backspace":
			if m.input_chars > 0 {
				m.input = string([]byte(m.input)[:m.input_chars-1])
				m.input_chars--
			}
		default:
			m.input += msg.String()
			m.input_chars++
		}
	case string: // Simulate receiving a message
		m.messages = append(m.messages, msg)
	}
	return m, nil
}

func (m model) View() string {
	s := strings.Join(m.messages, "\n") + "\n\n"
	s += "Input: " + m.input + "_"
	return s
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
	defer conn.Close()
	var wg sync.WaitGroup
	wg.Add(2)
	stopSignal := make(chan struct{})
	resumeSignal := make(chan struct{})
	p := tea.NewProgram(model{conn: conn})
	go receiveMessages(p, conn, &wg, stopSignal, resumeSignal)
	if _, err := p.Run(); err != nil {
		panic(err)
	}
	wg.Wait()
}
