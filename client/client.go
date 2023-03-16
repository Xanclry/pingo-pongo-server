package main

import (
	"bufio"
	"fmt"
	"io"
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
	reader  *bufio.Reader
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
	c.reader = bufio.NewReader(c.conn)
	log.Printf("Established connection to %v", c.address)
}

func (c *Client) Send(message string) error {
	binaryMessage := c.buildMessage(message)
	fmt.Printf("Sending \"%v\"...\n", message)
	_, err := c.conn.Write(binaryMessage)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Receive() (model.Response, error) {
	replyBinary := make([]byte, 6)
	_, err := io.ReadAtLeast(c.reader, replyBinary, 6)
	if err != nil {
		return model.Response{}, err
	}
	decodedResponse := model.DecodeResponse(replyBinary)
	return decodedResponse, nil

}

func (c *Client) Stop() {
	err := c.conn.Close()
	if err != nil {
		log.Panicf("Failed to close connection: %v", err.Error())
	} else {
		log.Printf("Connection to %v closed", c.address)
	}
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
