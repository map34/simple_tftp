package tftputils

import (
	"encoding/binary"
	"errors"
	"fmt"
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

func validateMode(mode string) (bool, error) {
	switch mode {
	case OCTET:
		return true, nil
	default:
		return false, fmt.Errorf("mode %v is not supported", mode)
	}
}
