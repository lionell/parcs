package parcs

import (
	b "bytes"
	"fmt"
	"log"
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
		log.Printf("Trying to connect to %s: %v", serviceName, err)
		if err == nil {
			log.Printf("Successfully connected to %s", serviceName)
			return conn
		}
		time.Sleep(SleepDuration)
	}
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
	log.Printf("Called with %v %v", conn, n)
	bytes := make([]byte, n)
	received := 0
	for received < n {
		log.Printf("received=%v", received)
		log.Printf("bytes[received:]=%v", bytes[received:])
		m, err := conn.Read(bytes[received:])
		log.Printf("m=%v", m)
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
		bytes, err := recvAllBytes(conn, 3)
		if err != nil {
			return err
		}
		if !b.Equal(bytes, SYN) {
			return fmt.Errorf("Expecting SYN got %v", bytes)
		}
		log.Printf("Client received SYN")
		if err := sendAllBytes(conn, ACK); err != nil {
			return err
		}
		log.Printf("Client sent ACK")
	case Server:
		if err := sendAllBytes(conn, SYN); err != nil {
			return err
		}
		log.Printf("Server sent SYN")
		bytes, err := recvAllBytes(conn, 3)
		if err != nil {
			return err
		}
		if !b.Equal(bytes, ACK) {
			return fmt.Errorf("Expecting ACK got %v", bytes)
		}
		log.Printf("Server received ACK")
	}
	return nil
}
