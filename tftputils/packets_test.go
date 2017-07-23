package tftputils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAckPackets(t *testing.T) {
	packet := createAckPacket(2)
	assert.Equal(t, packet, []byte{0x00, 0x04, 0x00, 0x02})
}

func TestErrorPackage(t *testing.T) {
	message := "hi"
	packet := createErrorPacket(FileExistsErr, message)
	assert.Equal(t, packet, []byte{0x00, 0x05, 0x00, 0x06, 0x68, 0x69, 0x00})
}

func TestDataPackage(t *testing.T) {
	message := "hi"
	packet := createDataPacket(1, []byte(message))
	assert.Equal(t, packet, []byte{0x00, 0x03, 0x00, 0x01, 0x68, 0x69})
}
