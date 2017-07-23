package tftputils

import (
	"net"

	"github.com/sirupsen/logrus"
)

// UDPUtils abstracts the ability to read and write from the UDP connection.
type UDPUtils struct {
	addr       *net.UDPAddr
	connection *net.UDPConn
	data       []byte
	remoteAddr *net.UDPAddr
}

func NewUDPUtils(initAddr string, remoteAddr string) (*UDPUtils, error) {
	if initAddr == "" {
		initAddr = ":0"
	}

	localAddr, err := net.ResolveUDPAddr("udp", initAddr)
	if err != nil {
		logrus.Errorf("Cannot resolve UDP address: %v", err)
		return nil, err
	}

	var connection *net.UDPConn
	var remoteUDPAddr *net.UDPAddr

	if remoteAddr != "" {
		remoteUDPAddr, err := net.ResolveUDPAddr("udp", remoteAddr)
		connection, err = net.DialUDP("udp", localAddr, remoteUDPAddr)
		if err != nil {
			logrus.Errorf("Cannot dial to UDP: %v", err)
			return nil, err
		}
		logrus.Infof("Dialing: %v", remoteAddr)
	} else {
		connection, err = net.ListenUDP("udp", localAddr)
		if err != nil {
			logrus.Errorf("Cannot listen to UDP: %v", err)
			return nil, err
		}
		logrus.Infof("Listening UDP at %v", connection.LocalAddr())
	}

	return &UDPUtils{
		addr:       localAddr,
		remoteAddr: remoteUDPAddr,
		connection: connection,
		data:       make([]byte, 1024),
	}, nil
}

func (udp *UDPUtils) CloseConnection() {
	udp.connection.Close()
}

func (udp *UDPUtils) WriteToConn(data []byte) error {
	_, err := udp.connection.Write(data)
	if err != nil {
		logrus.Errorf("Error writing to udp: %v", err)
		return err
	}
	return nil
}

func (udp *UDPUtils) ReadFromConn() ([]byte, *net.UDPAddr, error) {
	length, addr, err := udp.connection.ReadFromUDP(udp.data)

	if err != nil {
		logrus.Errorf("Cannot read from UDP: %v", err)
		return []byte{}, nil, err
	}

	newData := make([]byte, length)
	copy(newData, udp.data[0:length])
	return newData, addr, nil
}
