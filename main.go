package main

import (
	"fmt"
	"pingo-pongo-server/controller"
	"pingo-pongo-server/logger"
	"sync"
)

const (
	HOST = "localhost"
	PORT = "9003"
	TYPE = "tcp"
)

func main() {

	mainWaitGroup := sync.WaitGroup{}
	mainWaitGroup.Add(3)

	messagesChannel := make(chan controller.RawMessageFromClient, 1000)
	clientDisconnectChannel := make(chan controller.ClientId, 1000)

	logger := logger.NewLogger(messagesChannel, clientDisconnectChannel, &mainWaitGroup)
	go logger.Start()
	// go logger.StartLogger(messagesChannel, clientDisconnectChannel, &mainWaitGroup)

	listeningServer := controller.New(HOST+":"+PORT, messagesChannel, clientDisconnectChannel, &mainWaitGroup)
	go serverControl(messagesChannel, &mainWaitGroup, &listeningServer, &logger)
	listeningServer.Start()

	mainWaitGroup.Wait()

}

func serverControl(
	messagesChannel chan controller.RawMessageFromClient,
	mainWaitGroup *sync.WaitGroup,
	listeningServer *controller.ListenerServer,
	logger *logger.Logger,
) {
	var command string
	for {
		fmt.Scan(&command)
		switch command {
		case "exit":
			{
				logger.Stop(mainWaitGroup)
				// close(messagesChannel)
				(*listeningServer).Stop()
				mainWaitGroup.Done()
				break
			}
		}
	}
}
