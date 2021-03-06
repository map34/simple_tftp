package tftputils

import (
	"errors"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
)

// ReadSession holds necessary info about how to send
// file data to client.
type ReadSession struct {
	udpUtils *UDPUtils
	file     *FileObject
	reqInfo  *RequestInfo
	blockLoc uint16
}

func NewReadSession(fileS *FileStore, reqInfo *RequestInfo, remoteAddr *net.UDPAddr) (*ReadSession, error) {
	udpUtils, err := NewUDPUtils("", remoteAddr.String())
	if err != nil {
		return nil, err
	}
	ok, err := validateReadRequest(fileS, reqInfo, udpUtils)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errors.New("Cannot continue protocol")
	}

	file, err := fileS.Get(reqInfo.filename)
	if err != nil {
		return nil, err
	}
	return &ReadSession{
		udpUtils: udpUtils,
		file:     file,
		reqInfo:  reqInfo,
		blockLoc: 1,
	}, nil
}

// SpawnReadSession dials up to the address provided and
// starts sending necessary bytes to the client to save.
func SpawnReadSession(fileS *FileStore, reqInfo *RequestInfo, remoteAddr *net.UDPAddr) error {
	reader, err := NewReadSession(fileS, reqInfo, remoteAddr)
	if err != nil {
		return err
	}

	// close connection at the end of session
	defer reader.udpUtils.CloseConnection()

	logrus.Infof("R: Starting a reading session for %v in %v mode",
		reqInfo.filename, reqInfo.mode)

	// send first set of data for the client
	// to be acked
	reader.sendData()

	for {
		data, _, err := reader.udpUtils.ReadFromConn()
		if err != nil {
			return err
		}
		done, err := reader.ResolvePacket(data)
		if err != nil {
			return err
		}
		if done {
			logrus.Infof("R: Done transferring data from server to %v", reqInfo.filename)
			return nil
		}
	}
}

// validateReadRequest validates if a file exists in the server
// and the mode is supported, otherwise send an error message to the client.
func validateReadRequest(fileS *FileStore, reqInfo *RequestInfo, udpUtils *UDPUtils) (bool, error) {
	if !fileS.DoesFileExist(reqInfo.filename) {
		msg := fmt.Sprintf("R: File %v does not exist in the server", reqInfo.filename)
		logrus.Error(msg)
		return false, sendErrorPacket(FileNotFoundErr, msg, udpUtils)
	}
	return validateModeAndNotify(reqInfo.mode, udpUtils)
}

// ResolvePacket determines from initial request info
// what to do (send file data when ack is received, or handles error)
func (rs *ReadSession) ResolvePacket(packet []byte) (bool, error) {
	opCode, err := getOpCode(packet)
	if err != nil {
		return false, err
	}

	switch opCode {
	case ACK:
		return rs.handleAck(packet)
	case ERROR:
		err := handleError(packet)
		if err != nil {
			return false, err
		}
		return true, nil
	default:
		return false, fmt.Errorf("R: Opcode unknown or currently unsupported: %v", opCode)
	}
}

// handleAck sends data to the client when an appropriate
// ack packet is received, returns a done flag and an error
func (rs *ReadSession) handleAck(packet []byte) (bool, error) {
	blockFromClient, err := getAck(packet)
	if err != nil {
		return false, err
	}
	if blockFromClient == rs.blockLoc {
		rs.blockLoc++
	} else {
		return false, fmt.Errorf("R: Wrong expected byte, actual: %v, expected: %v", blockFromClient, rs.blockLoc)
	}
	return rs.sendData()
}

// sendData sends data to the client block by block
// returns a done flag and an error object
// done flag is true when there is no more data to send
func (rs *ReadSession) sendData() (bool, error) {
	nextBlock := rs.blockLoc * SmallestBlockSize
	prevBlock := (rs.blockLoc - 1) * SmallestBlockSize
	dataLen := uint16(len(rs.file.data))

	if prevBlock > dataLen {
		return true, nil
	}

	var blockData []byte
	if dataLen < nextBlock {
		blockData = rs.file.data[prevBlock:]
	} else {
		blockData = rs.file.data[prevBlock:nextBlock]
	}
	return false, sendDataPacket(rs.blockLoc, blockData, rs.udpUtils)
}
