package tftputils

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
)

type RequestInfo struct {
	filename string
	mode     string
}

func sendAckPacket(blockLoc uint16, udpUtils *UDPUtils) error {
	packet := createAckPacket(blockLoc)
	return udpUtils.WriteToConn(packet)
}

func sendErrorPacket(errCode uint8, errMessage string, udpUtils *UDPUtils) error {
	packet := createErrorPacket(errCode, errMessage)
	return udpUtils.WriteToConn(packet)
}

func sendDataPacket(block uint16, data []byte, udpUtils *UDPUtils) error {
	packet := createDataPacket(block, data)
	return udpUtils.WriteToConn(packet)
}

// We should be able to tolerate an error coming from client
func handleError(packet []byte) error {
	msg, err := getErrorMessage(packet)
	if err != nil {
		return err
	}
	logrus.Error(msg)
	return nil
}

//  2 bytes     string    1 byte     string   1 byte
// ------------------------------------------------
// | Opcode |  Filename  |   0  |    Mode    |   0  |
// ------------------------------------------------
func getOpCode(input []byte) (uint16, error) {
	if len(input) < 2 {
		return UNKNOWNOP, errors.New("Not enough bytes to get the opcode")
	}
	opInt, err := bytesToUint64(input[:2])
	if err != nil {
		return UNKNOWNOP, err
	}
	return uint16(opInt), nil
}

func getRequestInfo(input []byte) (string, string, error) {
	if len(input) < 2 {
		return "", "", errors.New("Not enough bytes to get filename and mode")
	}

	zeroByteTh := 0
	zeroByteIndex := -1

	reqInput := input[2:]
	filename := ""
	mode := ""
	for i, byteVal := range reqInput {
		if byteVal == 0 {
			isFileNameByte := zeroByteTh == 0
			isModeByte := zeroByteTh == 1 && zeroByteIndex > -1

			if isFileNameByte {
				filename = string(reqInput[:i])

				// Increment the proposed locations of the next 0 byte
				zeroByteTh++
				zeroByteIndex = i + 1
				if len(reqInput) <= zeroByteIndex {
					return filename, "", errors.New("Not enough bytes to get \"mode\"")
				}

				// Go straight to the next byte
				continue
			}

			if isModeByte {
				mode = string(reqInput[zeroByteIndex:i])
				return filename, mode, nil
			}
		}
	}
	return "", "", errors.New("Cannot parse input")
}

func createRequestInfo(packetBytes []byte) (*RequestInfo, error) {
	filename, mode, err := getRequestInfo(packetBytes)
	return &RequestInfo{
		filename: filename,
		mode:     mode,
	}, err
}

// 2 bytes     2 bytes
// ---------------------
// | Opcode |   Block #  |
// ---------------------
func getAck(input []byte) (uint16, error) {
	if len(input) < 4 {
		return 0, errors.New("Cannot get block. Data lengths mismatch")
	}
	block, err := bytesToUint64(input[2:4])
	if err != nil {
		return 0, err
	}
	return uint16(block), nil
}

// 2 bytes     2 bytes      n bytes
// ----------------------------------
// | Opcode |   Block #  |   Data     |
// ----------------------------------
func getData(input []byte) ([]byte, error) {
	if len(input) < 4 {
		return []byte{}, errors.New("Not enough bytes to get file data")
	}
	return input[4:], nil
}

//  2 bytes     2 bytes      string    1 byte
// -----------------------------------------
// | Opcode |  ErrorCode |   ErrMsg   |   0  |
// -----------------------------------------
func getErrorMessage(input []byte) (string, error) {
	if len(input) < 5 {
		return "", errors.New("Not enough bytes to get error message")
	}
	errCode, err := bytesToUint64(input[2:4])
	if err != nil {
		return "", err
	}
	errMessage := string(input[4 : len(input)-1])
	return fmt.Sprintf("Error from client. Code: %v, Message: %v", errCode, string(errMessage)), nil
}
