package parcs

import (
	b "bytes"
	"fmt"
	"net"
	"time"
)

type Side int

const (
	SleepDuration      = 500 * time.Millisecond
	Port               = 4444
	Client        Side = iota
	Server
)

var (
	SYN = []byte("SYN")
	ACK = []byte("ACK")
)

func listen() (net.Listener, error) {
	return net.Listen("tcp", fmt.Sprintf(":%d", Port))
}

func connect(serviceName string) net.Conn {
	for {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serviceName, Port))
		if err == nil {
			return conn
		}
		time.Sleep(SleepDuration)
	}
}

func send(conn net.Conn, v interface{}) error {
	return sendBytes(conn, marshal(v))
}

func sendAll(conn net.Conn, vs ...interface{}) error {
	for _, v := range vs {
		if err := send(conn, v); err != nil {
			return err
		}
	}
	return nil
}

func sendBytes(conn net.Conn, bytes []byte) error {
	err := sendAllBytes(conn, encodeUint64(uint64(len(bytes))))
	if err != nil {
		return err
	}
	return sendAllBytes(conn, bytes)
}

func sendAllBytes(conn net.Conn, bytes []byte) error {
	sent := 0
	for sent < len(bytes) {
		n, err := conn.Write(bytes[sent:])
		if err != nil {
			return err
		}
		sent += n
	}
	return nil
}

func recvBytes(conn net.Conn) ([]byte, error) {
	bytes, err := recvAllBytes(conn, BytesInInt)
	if err != nil {
		return nil, err
	}
	l := int(decodeUint64(bytes))
	return recvAllBytes(conn, l)
}

func recv(conn net.Conn, v interface{}) error {
	bytes, err := recvBytes(conn)
	if err != nil {
		return err
	}
	return unmarshal(bytes, v)
}

func recvAllBytes(conn net.Conn, n int) ([]byte, error) {
	bytes := make([]byte, n)
	received := 0
	for received < n {
		m, err := conn.Read(bytes[received:])
		if err != nil {
			return nil, err
		}
		received += m
	}
	return bytes, nil
}

func handshake(conn net.Conn, s Side) error {
	switch s {
	case Client:
		if err := sendAllBytes(conn, SYN); err != nil {
			return err
		}
		bytes, err := recvAllBytes(conn, 3)
		if err != nil {
			return err
		}
		if !b.Equal(bytes, ACK) {
			return fmt.Errorf("Expecting ACK got %v", bytes)
		}
	case Server:

		bytes, err := recvAllBytes(conn, 3)
		if err != nil {
			return err
		}
		if !b.Equal(bytes, SYN) {
			return fmt.Errorf("Expecting SYN got %v", bytes)
		}
		if err := sendAllBytes(conn, ACK); err != nil {
			return err
		}
	}
	return nil
}
