package controller

import (
	"bufio"
	"encoding/binary"
	"io"
	"log"
	"net"
	"pingo-pongo-server/model"
	"sync"
)

type ClientController struct {
	Conn              net.Conn
	MessagesChannel   chan RawMessageFromClient
	DisconnectChannel chan ClientId
	quit              chan bool
	parentWaitGroup   *sync.WaitGroup
}

func NewClientController(
	conn net.Conn,
	messageChannel chan RawMessageFromClient,
	disconnectChannel chan ClientId,
	parentWG *sync.WaitGroup,
) ClientController {
	return ClientController{
		Conn:              conn,
		MessagesChannel:   messageChannel,
		DisconnectChannel: disconnectChannel,
		quit:              make(chan bool, 1),
		parentWaitGroup:   parentWG,
	}
}

func (client *ClientController) Start() {
	clientId := ClientId{client.Conn.RemoteAddr().String()}
	reader := bufio.NewReader(client.Conn)
	for {
		messageBinary, err := client.receiveMessage(reader)

		if err != nil {
			select {
			case <-client.quit:
				log.Printf("Client controller for %v finished", clientId)
				client.parentWaitGroup.Done()
			default:
				if err.Error() == "EOF" {
					client.DisconnectChannel <- clientId
					client.parentWaitGroup.Done()
					log.Printf("Client %v disconnected", clientId)
				} else {
					log.Panicf("Error: %+v", err.Error())
				}
			}
			return
		}

		decodedMessage := model.DecodeMessage(messageBinary)

		messageFromClient := RawMessageFromClient{
			messageBinary,
			clientId,
		}

		client.MessagesChannel <- messageFromClient
		client.sendResponse(decodedMessage)
	}
}

func (client *ClientController) receiveMessage(reader *bufio.Reader) ([]byte, error) {
	lengthBuffer := make([]byte, 2)
	_, err1 := io.ReadAtLeast(reader, lengthBuffer, 2)
	if err1 != nil {
		return []byte{}, err1
	}
	messageLength := binary.LittleEndian.Uint16(lengthBuffer)

	messageBuffer := make([]byte, messageLength-2)
	_, err2 := io.ReadAtLeast(reader, messageBuffer, int(messageLength)-2)
	if err2 != nil {
		return []byte{}, err2
	}

	resultBuffer := make([]byte, 0)
	resultBuffer = append(resultBuffer, lengthBuffer...)
	resultBuffer = append(resultBuffer, messageBuffer...)

	return resultBuffer, nil
}

func (client *ClientController) sendResponse(incomingMessage model.RawMessage) {
	encodedResponse := model.EncodeResponse(buildResponse(incomingMessage))
	client.Conn.Write([]byte(encodedResponse))
}

func (client *ClientController) Stop() {
	client.quit <- true
	client.Conn.Close()
}

func buildResponse(message model.RawMessage) model.Response {
	return model.Response{Seq_num: message.ParsedMessage.Seq_num}
}
