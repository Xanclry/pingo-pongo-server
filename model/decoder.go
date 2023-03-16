package model

import (
	"encoding/binary"
	"fmt"
)

func DecodeMessage(binaryData []byte) RawMessage {
	lengthBinary := binaryData[0:2]
	length := binary.LittleEndian.Uint16(lengthBinary)
	fmt.Printf("decoder length: %v\n", length)

	seqNum := binary.LittleEndian.Uint32(binaryData[2:6])

	payload := string(binaryData[6:length])

	return RawMessage{
		ParsedMessage: ParsedMessage{
			Seq_num: seqNum,
			Payload: payload,
		},
		Length: length,
	}
}

func DecodeResponse(binaryData []byte) Response {
	lengthBinary := binaryData[0:2]
	length := binary.LittleEndian.Uint16(lengthBinary)

	seqNum := binary.LittleEndian.Uint32(binaryData[2:6])

	return Response{
		Seq_num: seqNum,
		Length:  length,
	}
}
