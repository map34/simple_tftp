package tftputils

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
)

type SpawnerFunction func(*FileStore, *RequestInfo, *net.UDPAddr) error

// ServeSession holds the udp read/write utils
// and a file storage reference
type ServeSession struct {
	udpUtils    *UDPUtils
	fileStorage *FileStore
}

// SpawnServeSession reads from socket and resolve the initial request from client
// and spawn a server in the main goroutine
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

// ResolvePacket determines from initial request info
// what to do (spawn a write/read session or handles the error)
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

// StartSession starts a session in a new goroutine.
// It handles error by repoting to the console to avoid
// affecting other goroutines.
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
	logrus.Errorf("%v", <-chanErr)
	return nil
}
