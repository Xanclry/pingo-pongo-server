package model

import (
	"encoding/binary"
)

const (
	LENGTH_SIZE  = 2
	SEQ_NUM_SIZE = 4
)

func EncodeMessage(message ParsedMessage) []byte {
	var totalSize uint16 = GetByteLengthOfEncodedMessage(message)

	binarySlice := make([]byte, LENGTH_SIZE)
	binary.LittleEndian.PutUint16(binarySlice, uint16(totalSize))
	binarySlice = binary.LittleEndian.AppendUint32(binarySlice, message.Seq_num)
	binarySlice = append(binarySlice, []byte(message.Payload)...)

	return binarySlice

}

func GetByteLengthOfEncodedMessage(message ParsedMessage) uint16 {
	payloadSize := len([]byte(message.Payload))

	var totalSize uint16 = uint16(LENGTH_SIZE + SEQ_NUM_SIZE + payloadSize)
	return totalSize
}

func EncodeResponse(response Response) []byte {
	var totalSize uint16 = uint16(LENGTH_SIZE + SEQ_NUM_SIZE)

	binarySlice := make([]byte, LENGTH_SIZE)
	binary.LittleEndian.PutUint16(binarySlice, uint16(totalSize))
	binarySlice = binary.LittleEndian.AppendUint32(binarySlice, response.Seq_num)

	return binarySlice

}
