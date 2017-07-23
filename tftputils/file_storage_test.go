package tftputils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func writeAFile() (*FileObject, *FileStore, error) {
	fileStorage := NewFileStore()
	file := NewFileObject("hello.txt", []byte{0x12, 0x33, 0x33})
	err := fileStorage.Put(file)
	return file, fileStorage, err
}

func TestPutGet(t *testing.T) {
	file, fileStorage, err := writeAFile()
	assert.Nil(t, err)
	actualFile, err := fileStorage.Get(file.filename)
	assert.Nil(t, err)
	assert.Equal(t, file, actualFile)
}
func TestPutExists(t *testing.T) {
	file, fileStorage, err := writeAFile()
	assert.Nil(t, err)
	err = fileStorage.Put(file)
	if assert.NotNil(t, err) {
		assert.Equal(t, fmt.Errorf("%v exists", file.filename), err)
	}
}

func TestGetNonExistant(t *testing.T) {
	filename := "somefile"
	fileStorage := NewFileStore()
	actualFile, err := fileStorage.Get(filename)
	assert.Nil(t, actualFile)
	if assert.NotNil(t, err) {
		assert.Equal(t, fmt.Errorf("%v does not exist", filename), err)
	}
}
