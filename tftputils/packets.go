package tftputils

type DataPacket struct {
	block       uint16
	data        []byte
	packetBytes []byte
}

type ErrorPacket struct {
	code        uint8
	message     string
	packetBytes []byte
}

type AckPacket struct {
	block       uint16
	packetBytes []byte
}

type RequestInfo struct {
	filename string
	mode     string
}

func NewAckPacket(block uint16) *AckPacket {
	blockBytes := uint64ToBytes(uint64(block))

	packet := []byte{0x0, 0x4,
		blockBytes[len(blockBytes)-2], blockBytes[len(blockBytes)-1]}
	return &AckPacket{
		block:       block,
		packetBytes: packet,
	}
}

func NewDataPacket(block uint16, data []byte) *DataPacket {
	blockBytes := uint64ToBytes(uint64(block))
	headerPacket := []byte{0x0, 0x3,
		blockBytes[len(blockBytes)-2], blockBytes[len(blockBytes)-1]}
	packet := append(headerPacket, data...)
	return &DataPacket{
		block:       block,
		data:        data,
		packetBytes: packet,
	}
}

func NewErrorPacket(errCode uint8, errMessage string) *ErrorPacket {
	errCodeBytes := uint64ToBytes(uint64(errCode))
	messageBytes := []byte(errMessage)

	packet := []byte{0x0, 0x5,
		errCodeBytes[len(errCodeBytes)-2], errCodeBytes[len(errCodeBytes)-1]}
	packet = append(packet, messageBytes...)
	packet = append(packet, []byte{0}...)

	return &ErrorPacket{
		code:        errCode,
		message:     errMessage,
		packetBytes: packet,
	}
}

func NewRequestInfo(packetBytes []byte) (*RequestInfo, error) {
	filename, mode, err := getRequestInfo(packetBytes)
	return &RequestInfo{
		filename: filename,
		mode:     mode,
	}, err
}
