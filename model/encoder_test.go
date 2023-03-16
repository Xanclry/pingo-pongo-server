package model

import (
	"testing"
)

var encodeMessageTests = []ParsedMessage{
	{
		Seq_num: 0,
		Payload: "test with 0",
	},
	{
		Seq_num: 0xFFFFFFFF,
		Payload: "test max value",
	},
	{
		Seq_num: 523905,
		Payload: "!@#$%^&*()_-+='''\n\"\":{}[]<>,./?",
	},
	{
		Seq_num: 0,
		Payload: "",
	},
}

var encodeResponseTests = []Response{
	{
		Seq_num: 0,
		Length:  6,
	},
	{
		Seq_num: 0xFFFFFFFF,
		Length:  6,
	},
}

func TestEncodeMessage(t *testing.T) {

	for _, test := range encodeMessageTests {
		encoded := EncodeMessage(test)
		decoded := DecodeMessage(encoded)
		if test != decoded.ParsedMessage {
			t.Errorf("Decoded message %q not equal to expected %q", decoded, test)
		}
	}

}

func TestEncodeResponse(t *testing.T) {
	for _, test := range encodeResponseTests {
		encoded := EncodeResponse(test)
		decoded := DecodeResponse(encoded)
		if test != decoded {
			t.Errorf("Decoded response %q not equal to expected %q", decoded, test)
		}
	}
}
