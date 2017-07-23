package tftputils

import (
	"fmt"
	"sync"
)

// FileObject holds filename and data of a file
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

// FileStore holds a dictionary of fileObjects and
// a mutex to protect from concurrent access
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

	if fs.DoesFileExist(file.filename) {
		return fmt.Errorf("%v exists", file.filename)
	}
	fs.fileMap[file.filename] = file
	return nil
}

func (fs *FileStore) Get(filename string) (*FileObject, error) {
	if !fs.DoesFileExist(filename) {
		return nil, fmt.Errorf("%v does not exist", filename)
	}
	file := fs.fileMap[filename]
	return file, nil
}

func (fs *FileStore) DoesFileExist(filename string) bool {
	_, ok := fs.fileMap[filename]
	return ok
}
