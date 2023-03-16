package main

import (
	"fmt"
	"log"
)

func main() {
	client := NewClient()
	client.Start()
	defer client.Stop()

	var inputString string
	for {
		fmt.Scan(&inputString)
		if inputString == "exit" {
			break
		}
		err := client.Send(inputString)
		if err != nil {
			log.Panicf("Sending failed: %v", err.Error())
		}

		response, err := client.Receive()
		if err != nil {
			log.Panicf("Receiving failed: %v", err.Error())
		}
		fmt.Printf("Received: %v\n", response.String())
	}

}
