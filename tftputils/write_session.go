package tftputils

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
)

type WriteSession struct {
	udpUtils    *UDPUtils
	fileStorage *FileStore
	tempBuf     []byte
	reqInfo     *RequestInfo
	blockLoc    uint16
}

func NewWriteSession(fileS *FileStore, reqInfo *RequestInfo, remoteAddr *net.UDPAddr) (*WriteSession, error) {
	udpUtils, err := NewUDPUtils(remoteAddr)
	if err != nil {
		return nil, err
	}
	return &WriteSession{
		udpUtils:    udpUtils,
		fileStorage: fileS,
		tempBuf:     []byte{},
		reqInfo:     reqInfo,
		blockLoc:    0,
	}, nil
}

func SpawnWriteSession(fileS *FileStore, reqInfo *RequestInfo, remoteAddr *net.UDPAddr) error {
	writer, err := NewWriteSession(fileS, reqInfo, remoteAddr)
	if err != nil {
		return err
	}
	sendAckPacket(writer.blockLoc, writer.udpUtils)
	for {
		data, addr, err := writer.udpUtils.ReadFromConn()
		if err != nil {
			return err
		}
		done, err := writer.ResolvePackets(data, addr)
		if err != nil {
			return err
		}
		if done {
			logrus.Infof("W: Done getting data from file")
			return writer.storeFile()
		}
	}
}

func (ws *WriteSession) ResolvePackets(packet []byte, addr *net.UDPAddr) (bool, error) {
	opCode, err := getOpCode(packet)
	if err != nil {
		return false, err
	}

	switch opCode {
	case DATA:
		return ws.handleData(packet)
	case ERROR:
		return true, nil
	default:
		return false, fmt.Errorf("W: Opcode unknown or currently unsupported: %v", opCode)
	}
}

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

	if len(packet) < SMALLESTBLOCK {
		return true, nil
	}
	return false, nil
}

func (ws *WriteSession) storeData(data []byte) {
	newData := append(ws.tempBuf, data...)
	ws.tempBuf = newData
}

func (ws *WriteSession) storeFile() error {
	newFile := NewFileObject(ws.reqInfo.filename, ws.tempBuf)
	return ws.fileStorage.Put(newFile)
}
