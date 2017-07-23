package tftputils

import (
	"errors"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
)

// WriteSession holds necessary info about how to receive file data
// from the client
type WriteSession struct {
	udpUtils    *UDPUtils
	fileStorage *FileStore
	tempBuf     []byte
	reqInfo     *RequestInfo
	blockLoc    uint16
}

func NewWriteSession(fileS *FileStore, reqInfo *RequestInfo, remoteAddr *net.UDPAddr) (*WriteSession, error) {
	udpUtils, err := NewUDPUtils("", remoteAddr.String())
	if err != nil {
		return nil, err
	}

	ok, err := validateWriteRequest(fileS, reqInfo, udpUtils)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errors.New("Cannot continue protocol")
	}
	return &WriteSession{
		udpUtils:    udpUtils,
		fileStorage: fileS,
		tempBuf:     []byte{},
		reqInfo:     reqInfo,
		blockLoc:    0,
	}, nil
}

// SpawnWriteSession dials up to the address provided and
// starts sending ack packets when file data is received block by block.
func SpawnWriteSession(fileS *FileStore, reqInfo *RequestInfo, remoteAddr *net.UDPAddr) error {
	writer, err := NewWriteSession(fileS, reqInfo, remoteAddr)

	if err != nil {
		return err
	}
	defer writer.udpUtils.CloseConnection()

	logrus.Infof("W: Starting a writing session for %v in %v mode",
		reqInfo.filename, reqInfo.mode)

	sendAckPacket(writer.blockLoc, writer.udpUtils)
	for {
		data, _, err := writer.udpUtils.ReadFromConn()
		if err != nil {
			return err
		}
		done, err := writer.ResolvePacket(data)
		if err != nil {
			return err
		}
		if done {
			// Once all data are received, store it on the file stoarge.
			err := writer.storeFile()
			logrus.Infof("W: Done transferring data from %v to server", reqInfo.filename)
			return err
		}
	}
}

// validateReadRequest validates if we can store file in the server (file hasn't existed) in the server
// and the mode is supported, otherwise send an error message to the client.
func validateWriteRequest(fileS *FileStore, reqInfo *RequestInfo, udpUtils *UDPUtils) (bool, error) {
	if fileS.DoesFileExist(reqInfo.filename) {
		msg := fmt.Sprintf("W: File %v exists in the server", reqInfo.filename)
		logrus.Error(msg)
		return false, sendErrorPacket(FileExistsErr, msg, udpUtils)
	}
	return validateModeAndNotify(reqInfo.mode, udpUtils)
}

// ResolvePacket determines from initial request info
// what to do (send ack packet when file data is received, or handles error)
func (ws *WriteSession) ResolvePacket(packet []byte) (bool, error) {
	opCode, err := getOpCode(packet)
	if err != nil {
		return false, err
	}

	switch opCode {
	case DATA:
		return ws.handleData(packet)
	case ERROR:
		err := handleError(packet)
		if err != nil {
			return false, err
		}
		return true, nil
	default:
		return false, fmt.Errorf("W: Opcode unknown or currently unsupported: %v", opCode)
	}
}

// handleData parses file data and store it on the temporary buffer
// block by block, once a data block is less than 512, we know that
// that is the last block of the file, returns a done flag and an error.
func (ws *WriteSession) handleData(packet []byte) (bool, error) {
	blockFromClient, err := getAck(packet)
	if err != nil {
		return false, err
	}

	if blockFromClient == ws.blockLoc+1 {
		ws.blockLoc++
	} else {
		return false,
			fmt.Errorf("W: Error reading the next block: %v", blockFromClient)
	}

	data, err := getData(packet)
	if err != nil {
		return false, err
	}

	ws.storeData(data)

	err = sendAckPacket(ws.blockLoc, ws.udpUtils)
	if err != nil {
		return false, err
	}

	if len(packet) < SmallestBlockSize {
		return true, nil
	}
	return false, nil
}

// storeData stores block data from client
// to the temporary buffer, append sto the previously
// received blocks in the temp buffer
func (ws *WriteSession) storeData(data []byte) {
	newData := append(ws.tempBuf, data...)
	ws.tempBuf = newData
}

// storeFile stores temporary buffer data
// to the file storage,
func (ws *WriteSession) storeFile() error {
	newFile := NewFileObject(ws.reqInfo.filename, ws.tempBuf)
	return ws.fileStorage.Put(newFile)
}
