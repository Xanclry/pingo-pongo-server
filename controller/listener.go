package controller

import (
	"log"
	"net"
	"sync"
)

const (
	TYPE = "tcp"
)

type RawMessageFromClient struct {
	Data []byte
	ClientId
}

type ListenerServer struct {
	Address                 string
	LoggingChannel          chan RawMessageFromClient
	ClientDisconnectChannel chan ClientId
	handlers                *[]ClientController
	handlersWaitGroup       *sync.WaitGroup
	mainWaitGroup           *sync.WaitGroup
	listener                *net.Listener
	quit                    chan bool
}

func New(
	address string,
	loggingChannel chan RawMessageFromClient,
	clientDisconnectChannel chan ClientId,
	mainWaitGroup *sync.WaitGroup,
) ListenerServer {
	return ListenerServer{
		Address:                 address,
		LoggingChannel:          loggingChannel,
		ClientDisconnectChannel: clientDisconnectChannel,
		handlersWaitGroup:       &sync.WaitGroup{},
		quit:                    make(chan bool, 1),
		mainWaitGroup:           mainWaitGroup,
	}
}

func (listenerServer *ListenerServer) Start() {
	l, err := net.Listen(TYPE, listenerServer.Address)
	if err != nil {
		log.Fatalf("Failed to create listener: %v", err.Error())
	}
	listenerServer.listener = &l

	handlersList := []ClientController{}
	listenerServer.handlers = &handlersList

	for {

		conn, err := (*listenerServer.listener).Accept()
		if err != nil {
			select {
			case <-listenerServer.quit:
				log.Println("Connection listener finished")
			default:
				log.Fatalf("Failed to accept connection: %v", err.Error())
			}
			return

		}
		log.Printf("Accepted connection from %v", conn.RemoteAddr())
		newClientController := NewClientController(conn, listenerServer.LoggingChannel, listenerServer.ClientDisconnectChannel, listenerServer.handlersWaitGroup)
		listenerServer.handlersWaitGroup.Add(1)
		go newClientController.Start()
		handlersList = append(handlersList, newClientController)
	}

}

func (listenerServer *ListenerServer) Stop() {
	log.Println("Shutting down listener")

	for _, c := range *(listenerServer).handlers {
		c.Stop()
	}

	listenerServer.handlersWaitGroup.Wait()
	listenerServer.quit <- true
	(*listenerServer.listener).Close()
	listenerServer.mainWaitGroup.Done()
	log.Print("Listener finished")
}
