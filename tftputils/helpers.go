package tftputils

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
)

func uint64ToBytes(num uint64) []byte {
	bytes := make([]byte, UINT64BYTESNUM)
	binary.BigEndian.PutUint64(bytes, num)
	return bytes
}

func bytesToUint64(bytes []byte) (uint64, error) {
	if len(bytes) > UINT64BYTESNUM {
		return 0, errors.New("Bytes length is too long for Uint64")
	}
	pad := make([]byte, UINT64BYTESNUM-len(bytes))
	completeBytes := append(pad, bytes...)
	return binary.BigEndian.Uint64(completeBytes), nil
}

func validateMode(mode string) bool {
	switch mode {
	case OCTET:
		return true
	default:
		return false
	}
}

func validateModeAndNotify(mode string, udpUtils *UDPUtils) (bool, error) {
	if !validateMode(mode) {
		msg := fmt.Sprintf("Mode %v not supported", mode)
		logrus.Error(msg)
		return false, sendErrorPacket(UnknownTransferIDErr, msg, udpUtils)
	}
	return true, nil
}
