package main

import (
	"pingo-pongo-server/controller"
	"pingo-pongo-server/logger"
	"sync"
)

type Server struct {
	Address                 string
	mainWaitGroup           *sync.WaitGroup
	messagesChannel         chan controller.RawMessageFromClient
	clientDisconnectChannel chan controller.ClientId
	logger                  *logger.Logger
	listener                *controller.ListenerServer
}

func NewServer(address string) Server {
	return Server{
		Address:                 address,
		mainWaitGroup:           &sync.WaitGroup{},
		messagesChannel:         make(chan controller.RawMessageFromClient, 1000),
		clientDisconnectChannel: make(chan controller.ClientId, 1000),
	}
}

func (s *Server) Start() {
	s.mainWaitGroup.Add(3)

	localLogger := logger.NewLogger(s.messagesChannel, s.clientDisconnectChannel, s.mainWaitGroup)
	s.logger = &localLogger
	go s.logger.Start()

	localListener := controller.New(s.Address, s.messagesChannel, s.clientDisconnectChannel, s.mainWaitGroup)
	s.listener = &localListener
	go s.listener.Start()
}

func (s *Server) Stop() {
	s.logger.Stop(s.mainWaitGroup)
	s.listener.Stop()
	s.mainWaitGroup.Done()
	s.mainWaitGroup.Wait()
}
