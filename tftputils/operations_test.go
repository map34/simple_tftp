package tftputils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendAck(t *testing.T) {
	port := ":3000"

	go func() {
		udpUtils, err := NewUDPUtils("", port)
		if err != nil {
			t.Fatal(err)
		}
		defer udpUtils.CloseConnection()

		if err := sendAckPacket(1, udpUtils); err != nil {
			t.Fatal(err)
		}
	}()

	readUDPUtils, err := NewUDPUtils(port, "")
	if err != nil {
		t.Fatal(err)
	}
	defer readUDPUtils.CloseConnection()
	for {
		readMessage, addr, err := readUDPUtils.ReadFromConn()
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, addr)
		assert.Equal(t, []byte{0x00, 0x04, 0x00, 0x01}, readMessage)
		return
	}
}

func TestSendError(t *testing.T) {
	port := ":3000"
	message := "hi"
	go func() {
		udpUtils, err := NewUDPUtils("", port)
		if err != nil {
			t.Fatal(err)
		}
		defer udpUtils.CloseConnection()

		if err := sendErrorPacket(FileExistsErr, message, udpUtils); err != nil {
			t.Fatal(err)
		}
	}()

	readUDPUtils, err := NewUDPUtils(port, "")
	if err != nil {
		t.Fatal(err)
	}
	defer readUDPUtils.CloseConnection()
	for {
		readMessage, addr, err := readUDPUtils.ReadFromConn()
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, addr)
		assert.Equal(t, []byte{0x00, 0x05, 0x00, 0x06, 0x68, 0x69, 0x00}, readMessage)
		return
	}
}

func TestDataPacket(t *testing.T) {
	port := ":3000"
	message := "hi"
	go func() {
		udpUtils, err := NewUDPUtils("", port)
		if err != nil {
			t.Fatal(err)
		}
		defer udpUtils.CloseConnection()

		if err := sendDataPacket(1, []byte(message), udpUtils); err != nil {
			t.Fatal(err)
		}
	}()

	readUDPUtils, err := NewUDPUtils(port, "")
	if err != nil {
		t.Fatal(err)
	}
	defer readUDPUtils.CloseConnection()
	for {
		readMessage, addr, err := readUDPUtils.ReadFromConn()
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, addr)
		assert.Equal(t, []byte{0x00, 0x03, 0x00, 0x01, 0x68, 0x69}, readMessage)
		return
	}
}

func TestGetOpCode(t *testing.T) {
	bytes := []byte{0x00, 0x01, 0x68, 0x69, 0x00, 0x01, 0x00}
	opCode, err := getOpCode(bytes)
	assert.Nil(t, err)
	assert.Equal(t, uint16(RRQ), opCode)
}

func TestBadGetOpCode(t *testing.T) {
	badBytes := []byte{0x01}
	badOpCode, err := getOpCode(badBytes)
	assert.NotNil(t, err)
	assert.Equal(t, uint16(UNKNOWNOP), badOpCode)
}

func TestCreateRequestInfo(t *testing.T) {
	bytes := []byte{0x00, 0x01, 0x68, 0x69, 0x00, 0x6f, 0x63, 0x74, 0x65, 0x74, 0x00}
	reqInfo, err := createRequestInfo(bytes)
	assert.Nil(t, err)
	assert.Equal(t, "hi", reqInfo.filename)
	assert.Equal(t, "octet", reqInfo.mode)
}

func TestCreateRequestInfoNameOnly(t *testing.T) {
	badBytes := []byte{0x00, 0x01, 0x68, 0x69, 0x00}
	reqInfo, err := createRequestInfo(badBytes)
	assert.NotNil(t, err)
	assert.Equal(t, "hi", reqInfo.filename)
	assert.Equal(t, "", reqInfo.mode)
}

func TestCreateRequestInfoBad(t *testing.T) {
	badBytes := []byte{0x00, 0x01, 0x68, 0x69}
	reqInfo, err := createRequestInfo(badBytes)
	assert.NotNil(t, err)
	assert.Equal(t, "", reqInfo.filename)
	assert.Equal(t, "", reqInfo.mode)
}

func TestGetAck(t *testing.T) {
	bytes := []byte{0x00, 0x01, 0x00, 0x02}
	block, err := getAck(bytes)
	assert.Nil(t, err)
	assert.Equal(t, uint16(2), block)
}

func TestBadGetAck(t *testing.T) {
	bytes := []byte{0x00, 0x01}
	block, err := getAck(bytes)
	assert.NotNil(t, err)
	assert.Equal(t, uint16(0), block)
}

func TestGetData(t *testing.T) {
	bytes := []byte{0x00, 0x01, 0x00, 0x02, 0x69}
	block, err := getData(bytes)
	assert.Nil(t, err)
	assert.Equal(t, []byte{0x69}, block)
}

func TestBadGetData(t *testing.T) {
	bytes := []byte{0x00, 0x01}
	block, err := getData(bytes)
	assert.NotNil(t, err)
	assert.Equal(t, []byte{}, block)
}

func TestGetErrorMessage(t *testing.T) {
	bytes := []byte{0x00, 0x01, 0x00, 0x02, 0x68, 0x69, 0x00}
	actualMsg, err := getErrorMessage(bytes)
	expectedMsg := fmt.Sprintf("Error from client. Code: 2, Message: hi")
	assert.Nil(t, err)
	assert.Equal(t, actualMsg, expectedMsg)
}
