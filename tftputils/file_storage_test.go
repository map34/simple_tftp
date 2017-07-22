package tftputils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func writeAFile() (*FileObject, *FileStore, error) {
	fileStorage := NewFileStore()
	file := NewFileObject("hello.txt", []byte{0x12, 0x33, 0x33})
	err := fileStorage.Put(file)
	return file, fileStorage, err
}

func TestPut(t *testing.T) {
	_, _, err := writeAFile()
	assert.Nil(t, err)
}

func TestPutExists(t *testing.T) {
	file, fileStorage, _ := writeAFile()
	err := fileStorage.Put(file)
	if assert.NotNil(t, err) {
		assert.Equal(t, errors.New("File already exists"), err)
	}
}

func TestGet(t *testing.T) {
	file, fileStorage, _ := writeAFile()
	actualFile, _ := fileStorage.Get(file.filename)
	assert.Equal(t, file, actualFile)
}

func TestGetNonExistant(t *testing.T) {
	fileStorage := NewFileStore()
	actualFile, err := fileStorage.Get("some_file")
	assert.Nil(t, actualFile)
	if assert.NotNil(t, err) {
		assert.Equal(t, errors.New("File does not exist"), err)
	}
}
