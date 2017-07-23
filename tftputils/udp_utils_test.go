package tftputils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUDPReadWrite(t *testing.T) {
	message := "hello"
	port := ":3000"

	go func() {
		conn, err := NewUDPUtils("", port)
		if err != nil {
			t.Fatal(err)
		}
		defer conn.CloseConnection()

		if err := conn.WriteToConn([]byte(message)); err != nil {
			t.Fatal(err)
		}
	}()

	readUdpUtils, err := NewUDPUtils(port, "")
	if err != nil {
		t.Fatal(err)
	}
	defer readUdpUtils.CloseConnection()
	for {
		readMessage, addr, err := readUdpUtils.ReadFromConn()
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, addr)
		assert.Equal(t, message, string(readMessage))
		return
	}
}
