package tftputils

import (
	"fmt"
	"net"
)

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
		err := s.startWriteSession(packet, addr)
		if err != nil {
			return false, err
		}
		return true, nil
	case RRQ:
		err := s.startReadSession(packet, addr)
		if err != nil {
			return false, err
		}
		return true, nil
	default:
		return false, fmt.Errorf("S: Opcode unknown or currently unsupported: %v", opCode)
	}
}

func (s *ServeSession) startWriteSession(packet []byte, addr *net.UDPAddr) error {
	reqInfo, err := NewRequestInfo(packet)
	if err != nil {
		return err
	}
	go SpawnWriteSession(s.fileStorage, reqInfo, addr)
	return nil
}

func (s *ServeSession) startReadSession(packet []byte, addr *net.UDPAddr) error {
	reqInfo, err := NewRequestInfo(packet)
	if err != nil {
		return err
	}
	go SpawnReadSession(s.fileStorage, reqInfo, addr)
	return nil
}
