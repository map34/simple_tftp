package tftputils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUDPReadWrite(t *testing.T) {
	message := "hello"

	portChan := make(chan string)
	go func() {
		conn, err := NewUDPUtils("", <-portChan)
		if err != nil {
			t.Fatal(err)
		}
		defer conn.CloseConnection()

		if err := conn.WriteToConn([]byte(message)); err != nil {
			t.Fatal(err)
		}
	}()

	readUDPUtils, err := NewUDPUtils("", "")
	if err != nil {
		t.Fatal(err)
	}
	portChan <- readUDPUtils.LocalAddress()
	defer readUDPUtils.CloseConnection()
	for {
		readMessage, addr, err := readUDPUtils.ReadFromConn()
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, addr)
		assert.Equal(t, message, string(readMessage))
		return
	}
}
