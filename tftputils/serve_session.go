package tftputils

import (
	"fmt"
	"net"
)

type SpawnerFunction func(*FileStore, *RequestInfo, *net.UDPAddr) error

type ServeSession struct {
	udpUtils    *UDPUtils
	fileStorage *FileStore
}

func SpawnServeSession() error {
	server, err := NewServeSession()
	if err != nil {
		return err
	}
	defer server.udpUtils.CloseConnection()
	for {
		data, addr, err := server.udpUtils.ReadFromConn()
		if err != nil {
			return err
		}
		_, err = server.ResolvePacket(data, addr)
		if err != nil {
			return err
		}
	}
}

func NewServeSession() (*ServeSession, error) {
	udpUtils, err := NewUDPUtils("", "")
	if err != nil {
		return nil, err
	}
	return &ServeSession{
		udpUtils:    udpUtils,
		fileStorage: NewFileStore(),
	}, nil
}

func (s *ServeSession) ResolvePacket(packet []byte, addr *net.UDPAddr) (bool, error) {
	opCode, err := getOpCode(packet)
	if err != nil {
		return false, err
	}

	switch opCode {
	case WRQ:
		err := s.StartSession(packet, addr, SpawnWriteSession)
		if err != nil {
			return false, err
		}
		return true, nil
	case RRQ:
		err := s.StartSession(packet, addr, SpawnReadSession)
		if err != nil {
			return false, err
		}
		return true, nil
	case ERROR:
		err := handleError(packet)
		if err != nil {
			return false, err
		}
		return true, nil
	default:
		return false, fmt.Errorf("S: Opcode unknown or currently unsupported: %v", opCode)
	}
}

func (s *ServeSession) StartSession(
	packet []byte,
	addr *net.UDPAddr,
	funcSig SpawnerFunction) error {

	reqInfo, err := createRequestInfo(packet)
	if err != nil {
		return err
	}
	chanErr := make(chan error)
	go func() {
		chanErr <- funcSig(s.fileStorage, reqInfo, addr)
	}()
	return <-chanErr
}
