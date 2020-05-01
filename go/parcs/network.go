package parcs

import (
	"fmt"
	"net"
)

const Port = 4444

func listen() (net.Listener, error) {
	return net.Listen("tcp", fmt.Sprintf(":%d", Port))
}

func connect(serviceName string) (net.Conn, error) {
	return net.Dial("tcp", fmt.Sprintf("%s:%d", serviceName, Port))
}

func send(conn net.Conn, v interface{}) error {
	return sendBytes(conn, marshal(v))
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

func recv(conn net.Conn, v interface{}) error {
	bytes, err := recvAllBytes(conn, BytesInInt)
	if err != nil {
		return err
	}
	l := int(decodeUint64(bytes))
	bytes, err = recvAllBytes(conn, l)
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
