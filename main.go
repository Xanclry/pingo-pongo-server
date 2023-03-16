package main

import (
	"fmt"
	"log"
)

const (
	HOST = "localhost"
	PORT = "9003"
)

func main() {
	server := NewServer(HOST + ":" + PORT)
	server.Start()

	var command string
	for {
		fmt.Scan(&command)
		switch command {
		case "exit":
			{
				server.Stop()
				log.Println("Server stopped")
				return
			}
		}
	}

}
