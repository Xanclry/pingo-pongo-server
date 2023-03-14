package controller

import "fmt"

type ClientId struct {
	Address string
}

func (client ClientId) GenerateFilename() string {
	return fmt.Sprintf("%v.txt", client.Address)
}
