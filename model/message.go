package model

import (
	"fmt"
)

type ParsedMessage struct {
	Seq_num uint32
	Payload string
}

func (parsed ParsedMessage) String() string {
	return fmt.Sprintf("[%v][%v]", parsed.Seq_num, parsed.Payload)
}

// [msg length uint16][sequence number uint32][payload any string]
type RawMessage struct {
	ParsedMessage ParsedMessage
	Length        uint16
}

func (raw RawMessage) GetParsedMessage() ParsedMessage {
	return raw.ParsedMessage
}

func (raw RawMessage) String() string {
	return fmt.Sprintf("[%v]%v", raw.Length, raw.GetParsedMessage().String())
}

// [msg length uint16][sequence number uint32]
type Response struct {
	Length  uint16
	Seq_num uint32
}

func (response Response) String() string {
	return fmt.Sprintf("[%v][%v]", response.Length, response.Seq_num)
}
