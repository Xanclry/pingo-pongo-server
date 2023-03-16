package main

import (
	"fmt"
	"log"
	"net"
	"pingo-pongo-server/model"
)

const (
	DEFAULT_ADDRESS = "localhost:9003"
)

type Client struct {
	address string
	conn    *net.TCPConn
	seq_num uint32
}

func NewClient() Client {
	return Client{
		address: DEFAULT_ADDRESS,
		seq_num: 0,
	}
}

func (c *Client) Start() {
	err := c.initConnection()
	if err != nil {
		log.Fatalf("Failed to create connection: %v", err.Error())
	}
	log.Printf("Established connection to %v", c.address)

	var inputString string
	reply := make([]byte, 6)
	for {
		fmt.Scan(&inputString)
		if inputString == "exit" {
			break
		}

		err = c.send(c.buildMessage(inputString))
		if err != nil {
			log.Panicf("Write to server failed: %v", err.Error())
		}

		_, err := c.conn.Read(reply)
		if err != nil {
			log.Panicf("Write to server failed: %v", err.Error())
		}
		decodedResponse := model.DecodeResponse(reply)

		fmt.Printf("Received \"%v\"\n", reply)
		fmt.Printf("Decoded: %v\n", decodedResponse.String())
	}
}

func (c *Client) Stop() {
	err := c.conn.Close()
	if err != nil {
		log.Panicf("Failed to close connection: %v", err.Error())
	}
	log.Print("Connection closed")
}

func (c *Client) buildMessage(message string) []byte {
	encodedMessage := model.EncodeMessage(model.ParsedMessage{
		Payload: message,
		Seq_num: c.seq_num,
	})
	c.seq_num += 1
	return encodedMessage
}

func (c *Client) initConnection() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", c.address)
	if err != nil {
		return err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *Client) send(message []byte) error {
	fmt.Printf("Sending \"%v\"...\n", message)
	_, err := c.conn.Write(message)
	if err != nil {
		return err
	}
	return nil
}
