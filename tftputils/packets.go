package tftputils

func createAckPacket(block uint16) []byte {
	blockBytes := uint64ToBytes(uint64(block))
	return []byte{0x0, 0x4,
		blockBytes[len(blockBytes)-2], blockBytes[len(blockBytes)-1]}
}

func createDataPacket(block uint16, data []byte) []byte {
	blockBytes := uint64ToBytes(uint64(block))
	headerPacket := []byte{0x0, 0x3,
		blockBytes[len(blockBytes)-2], blockBytes[len(blockBytes)-1]}
	return append(headerPacket, data...)
}

func createErrorPacket(errCode uint8, errMessage string) []byte {
	errCodeBytes := uint64ToBytes(uint64(errCode))
	messageBytes := []byte(errMessage)

	packet := []byte{0x0, 0x5,
		errCodeBytes[len(errCodeBytes)-2], errCodeBytes[len(errCodeBytes)-1]}
	packet = append(packet, messageBytes...)
	return append(packet, []byte{0}...)
}
