package model

import (
	"fmt"
	"strconv"
	"strings"
)

type ParsedMessage struct {
	Seq_num uint32
	Payload string
}

func (parsed ParsedMessage) String() string {
	return fmt.Sprintf("[%v][%v]", parsed.Seq_num, parsed.Payload)
}

func (parsed ParsedMessage) ConvertToRawMessage() RawMessage {
	length := calculateLengthOfStringWithLength(parsed.String())
	return RawMessage{ParsedMessage: parsed, Length: length}

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

func ParseRawMessage(s string) RawMessage {
	splitted := strings.Split(s, "]")
	length, _ := strconv.ParseUint(splitted[0][1:], 10, 16)
	seq_num, _ := strconv.ParseUint(splitted[1][1:], 10, 32)
	payload := splitted[2][1:]
	return RawMessage{
		ParsedMessage: ParsedMessage{
			Seq_num: uint32(seq_num),
			Payload: payload,
		},
		Length: uint16(length),
	}
}

// [msg length uint16][sequence number uint32]
type Response struct {
	Length  uint16
	Seq_num uint32
}

func (response Response) String() string {
	return BuildStringWithLength("[" + strconv.FormatUint(uint64(response.Seq_num), 10) + "]")
}

func calculateLengthOfStringWithLength(s string) uint16 {
	stringLength := len(s) + 2
	stringLengthAsString := strconv.Itoa(stringLength)
	totalLength := stringLength + len(stringLengthAsString)

	actualLength := len(fmt.Sprintf("[%v]%v", totalLength, s))
	if actualLength != totalLength {
		totalLength = actualLength
	}
	return uint16(totalLength)
}

func BuildStringWithLength(s string) string {
	stringLength := len(s) + 2
	stringLengthAsString := strconv.Itoa(stringLength)
	totalLength := stringLength + len(stringLengthAsString)

	actualLength := len(fmt.Sprintf("[%v]%v", totalLength, s))
	if actualLength != totalLength {
		totalLength = actualLength
	}
	return fmt.Sprintf("[%v]%v", totalLength, s)
}
