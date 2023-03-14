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
	// stopListenningChannel := make(chan bool, 1)
	go logger.StartLogger(messagesChannel, clientDisconnectChannel, &mainWaitGroup)

	listeningServer := controller.New(HOST+":"+PORT, messagesChannel, clientDisconnectChannel, &mainWaitGroup)
	go serverControl(messagesChannel, &mainWaitGroup, &listeningServer)
	listeningServer.Start()

	// go controller.StartListening(HOST+":"+PORT, messagesChannel, clientDisconnectChannel, stopListenningChannel)
	mainWaitGroup.Wait()

	// loggingChannel <- ClientTwo()
	// time.Sleep(time.Second * 1)
	// loggingChannel <- ClientOne()
	// time.Sleep(time.Second * 1)
	// loggingChannel <- ClientTwo()
	// loggingChannel <- ClientOne()

	// time.Sleep(time.Second * 1)
	// close(loggingChannel)

	// fmt.Printf("[11][1][client one]: %v\n", model.ParseRawMessage("[11][1][client one]"))
}

func serverControl(messagesChannel chan controller.RawMessageFromClient, mainWaitGroup *sync.WaitGroup, listeningServer *controller.ListenerServer) {
	var command string
	for {
		fmt.Scan(&command)
		switch command {
		case "exit":
			{
				close(messagesChannel)
				(*listeningServer).Stop()
				mainWaitGroup.Done()
				break
			}
		}
	}
}

// func ClientOne() controller.RawMessageFromClient {
// 	return controller.RawMessageFromClient{
// 		RawMessage: model.RawMessage{
// 			ParsedMessage: model.ParsedMessage{Payload: "client one", Seq_num: 1},
// 			Length:        11,
// 		},
// 		ClientId: controller.ClientId{
// 			Address: "8081",
// 		},
// 	}
// }

// func ClientTwo() controller.RawMessageFromClient {
// 	return controller.RawMessageFromClient{
// 		RawMessage: model.RawMessage{
// 			ParsedMessage: model.ParsedMessage{Payload: "client two", Seq_num: 2},
// 			Length:        22,
// 		}, ClientId: controller.ClientId{
// 			Address: "8082",
// 		},
// 	}
// }
