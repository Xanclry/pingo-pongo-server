package controller

import (
	"bufio"
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
		message := ""
		for i := 0; i < 3; i++ {
			partOfMessage, err := reader.ReadString(']')
			if err != nil {
				select {
				case <-client.quit:
					log.Printf("Client controller for %v finished", clientId)
					client.parentWaitGroup.Done()
				default:
					if err.Error() == "EOF" {
						client.DisconnectChannel <- clientId
						log.Printf("Client %v disconnected", clientId)
					} else {
						log.Panicf("Error: %+v", err.Error())
					}
				}

				return
			}
			message = message + partOfMessage
		}

		messageFromClient := RawMessageFromClient{
			model.ParseRawMessage(message),
			clientId,
		}
		client.MessagesChannel <- messageFromClient
		client.Conn.Write([]byte(buildResponse(messageFromClient)))
	}

}

func (client *ClientController) Stop() {
	client.quit <- true
	client.Conn.Close()
}

func buildResponse(message RawMessageFromClient) string {
	response := model.Response{Seq_num: message.ParsedMessage.Seq_num}
	return response.String()
}
