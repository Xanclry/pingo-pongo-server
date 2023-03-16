package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"pingo-pongo-server/controller"
	"sync"
)

const (
	LOGGING_PATH = "./transaction-log/"
)

type Logger struct {
	NewMessageChannel       chan (controller.RawMessageFromClient)
	ClientDisconnectChannel chan controller.ClientId
	clientToChannelMap      *sync.Map
	clientsWaitGroup        *sync.WaitGroup
}

func NewLogger(
	newMessageChannel chan (controller.RawMessageFromClient),
	clientDisconnectChannel chan controller.ClientId,
	mainWaitGroup *sync.WaitGroup,
) Logger {
	return Logger{
		NewMessageChannel:       newMessageChannel,
		ClientDisconnectChannel: clientDisconnectChannel,
		clientToChannelMap:      &sync.Map{},
		clientsWaitGroup:        &sync.WaitGroup{},
	}
}

func (l *Logger) Start() {
	go l.listenToDisconnects()
	for newMessage := range l.NewMessageChannel {
		l.handleNewMessage(newMessage)
	}
}

func (l *Logger) handleNewMessage(
	newMessage controller.RawMessageFromClient,
) {
	channelForThisClient, loaded := l.clientToChannelMap.LoadOrStore(newMessage.ClientId, make(chan controller.RawMessageFromClient, 1000))
	castedChannel := channelForThisClient.(chan controller.RawMessageFromClient)

	if !loaded {
		// new client
		l.clientsWaitGroup.Add(1)
		go l.startLoggerForClient(castedChannel, newMessage.ClientId)
	}
	castedChannel <- newMessage
}

func (l *Logger) startLoggerForClient(channel chan controller.RawMessageFromClient, clientId controller.ClientId) {
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
	l.clientsWaitGroup.Done()
}

func logMessageToFile(message controller.RawMessageFromClient, file *os.File) {
	fmt.Printf("Received \"%v\"\n", message.Data)
	file.WriteString(string(message.Data) + "\n")
}

func (l *Logger) listenToDisconnects() {
	for disconnect := range l.ClientDisconnectChannel {
		channel, loaded := l.clientToChannelMap.LoadAndDelete(disconnect)
		if loaded {
			close(channel.(chan controller.RawMessageFromClient))
		}
	}
}

func (l *Logger) Stop(waitGroup *sync.WaitGroup) {
	close(l.NewMessageChannel)
	l.clientToChannelMap.Range(func(key, value any) bool {
		castedChannel := value.(chan controller.RawMessageFromClient)
		close(castedChannel)
		return true
	})
	l.clientsWaitGroup.Wait()
	log.Println("Logger finished")
	waitGroup.Done()
}
