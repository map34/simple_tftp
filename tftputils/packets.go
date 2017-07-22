package tftputils

type DataPacket struct {
	block       uint16
	data        []byte
	packetBytes []byte
}

type ErrorPacket struct {
	code        uint16
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

func NewAckPacket(block [2]byte) (*AckPacket, error) {
	blockInt, err := bytesToUint64(block[:])
	if err != nil {
		return nil, err
	}
	packet := []byte{0x0, 0x4, block[0], block[1]}
	return &AckPacket{
		block:       uint16(blockInt),
		packetBytes: packet,
	}, nil
}

func NewDataPacket(block [2]byte, data []byte) (*DataPacket, error) {
	blockInt, err := bytesToUint64(block[:])
	if err != nil {
		return nil, err
	}
	headerPacket := []byte{0x0, 0x3, block[0], block[1]}
	packet := append(headerPacket, data...)
	return &DataPacket{
		block:       uint16(blockInt),
		data:        data,
		packetBytes: packet,
	}, nil
}

func NewRequestInfo(packetBytes []byte) (*RequestInfo, error) {
	filename, mode, err := getRequestInfo(packetBytes)
	return &RequestInfo{
		filename: filename,
		mode:     mode,
	}, err
}
