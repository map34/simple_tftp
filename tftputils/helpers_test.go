package tftputils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUint64ToBytes(t *testing.T) {
	bytes := uint64ToBytes(1024)
	assert.Equal(t, bytes, []byte{0, 0, 0, 0, 0, 0, 4, 0})
}

func TestBytesToUint64(t *testing.T) {
	num, _ := bytesToUint64([]byte{0, 0, 0, 0, 0, 0, 4, 0})
	assert.Equal(t, num, uint64(1024))
}

func TestBytesToUint64Big(t *testing.T) {
	_, err := bytesToUint64([]byte{0, 0, 0, 0, 0, 0, 4, 0, 8})
	assert.NotNil(t, err)
}

func TestValidateMode(t *testing.T) {
	ok := validateMode(OCTET)
	assert.True(t, ok)

	ok = validateMode("UNKNOWN")
	assert.False(t, ok)
}
