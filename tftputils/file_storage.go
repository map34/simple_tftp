package tftputils

import (
	"errors"
	"sync"
)

type FileObject struct {
	filename string
	data     []byte
}

func NewFileObject(filename string, data []byte) *FileObject {
	return &FileObject{
		filename: filename,
		data:     data,
	}
}

type FileStore struct {
	fileMap map[string]*FileObject
	mutex   *sync.Mutex
}

func NewFileStore() *FileStore {
	return &FileStore{
		fileMap: make(map[string]*FileObject),
		mutex:   &sync.Mutex{},
	}
}

func (fs *FileStore) Put(file *FileObject) error {
	// Protect storage from concurrent writing
	defer fs.mutex.Unlock()
	fs.mutex.Lock()

	_, ok := fs.fileMap[file.filename]
	if ok {
		return errors.New("File already exists")
	}
	fs.fileMap[file.filename] = file
	return nil
}

func (fs *FileStore) Get(filename string) (*FileObject, error) {
	file := fs.fileMap[filename]
	if file == nil {
		return file, errors.New("File does not exist")
	}
	return file, nil
}
