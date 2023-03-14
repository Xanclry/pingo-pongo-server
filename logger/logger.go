package logger

import (
	"log"
	"os"
	"path/filepath"
	"pingo-pongo-server/controller"
	"pingo-pongo-server/model"
	"sync"
)

const (
	LOGGING_PATH = "./transaction-log/"
)

func StartLogger(
	newMessageChannel chan (controller.RawMessageFromClient),
	clientDisconnectChannel chan controller.ClientId,
	mainWaitGroup *sync.WaitGroup,
) {
	clientsWaitGroup := sync.WaitGroup{}
	var clientToChannelMap = sync.Map{} // ClientId -> chan model.RawMessage
	go listenToDisconnects(clientDisconnectChannel, &clientToChannelMap)

	for newMessage := range newMessageChannel {
		handleNewMessage(&clientToChannelMap, newMessage, &clientsWaitGroup)
	}

	shutdown(&clientToChannelMap, &clientsWaitGroup, mainWaitGroup)
}

func handleNewMessage(
	clientToChannelMap *sync.Map,
	newMessage controller.RawMessageFromClient,
	clientsWaitGroup *sync.WaitGroup,
) {
	channelForThisClient, loaded := clientToChannelMap.LoadOrStore(newMessage.ClientId, make(chan model.RawMessage, 1000))

	if !loaded {
		// new client
		castedChannel := channelForThisClient.(chan model.RawMessage)
		clientsWaitGroup.Add(1)
		go startLoggerForClient(castedChannel, newMessage.ClientId, clientsWaitGroup)
	}
	castedChannel := channelForThisClient.(chan model.RawMessage)
	castedChannel <- newMessage.RawMessage
}

func startLoggerForClient(channel chan model.RawMessage, clientId controller.ClientId, waitGroup *sync.WaitGroup) {
	filename, _ := filepath.Abs(LOGGING_PATH + clientId.GenerateFilename())

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer func() {
		file.Close()
		log.Printf("File %v closed", filename)
	}()

	log.Printf("new client %v - logs saved in %v\n", clientId, filename)
	for incomingMessage := range channel {
		logMessageToFile(incomingMessage, file)
	}

	log.Printf("Logging stopped for client %v\n", clientId)
	waitGroup.Done()
}

func logMessageToFile(message model.RawMessage, file *os.File) {
	file.WriteString(message.String() + "\n")
}

func listenToDisconnects(
	clientDisconnectChannel chan controller.ClientId,
	clientsMap *sync.Map,
) {
	for disconnect := range clientDisconnectChannel {
		channel, loaded := clientsMap.LoadAndDelete(disconnect)
		if loaded {
			close(channel.(chan model.RawMessage))
		}
	}
}

func shutdown(clientsMap *sync.Map, clientsWaitGroup *sync.WaitGroup, mainWaitGroup *sync.WaitGroup) {
	clientsMap.Range(func(key, value any) bool {
		castedChannel := value.(chan model.RawMessage)
		close(castedChannel)
		return true
	})
	clientsWaitGroup.Wait()
	log.Println("Logger finished")
	mainWaitGroup.Done()

}
