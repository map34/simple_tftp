package tftputils

import "errors"

func sendAckPacket(blockLoc uint16, udpUtils *UDPUtils) error {
	packet := NewAckPacket(blockLoc)
	return udpUtils.WriteToConn(packet.packetBytes)
}

func sendErrorPacket(errCode uint8, errMessage string, udpUtils *UDPUtils) error {
	errorPacket := NewErrorPacket(errCode, errMessage)
	return udpUtils.WriteToConn(errorPacket.packetBytes)
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
	zeroByteTh := 0
	zeroByteIndex := -1

	reqInput := input[2:]
	filename := ""
	mode := ""
	for i, e := range reqInput {
		if e == 0 {
			isFileNameByte := zeroByteTh == 0
			isModeByte := zeroByteTh == 1 && zeroByteIndex > -1

			if isFileNameByte {
				filename = string(reqInput[:i])

				// Increment the proposed locations of the next 0 byte
				zeroByteTh++
				zeroByteIndex = i + 1
				if len(reqInput) < zeroByteIndex {
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

// 2 bytes     2 bytes
// ---------------------
// | Opcode |   Block #  |
// ---------------------
func getAck(input []byte) (uint16, error) {
	if len(input) < 4 {
		return 0, errors.New("Not enough bytes to get block")
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
