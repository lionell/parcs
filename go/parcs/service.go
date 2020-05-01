package parcs

import (
	"log"
	"net"
)

type Service struct {
	listener net.Listener
	conn     net.Conn

	*Engine
}

func DefaultService() *Service {
	return NewService(NewEnvEngine())
}

func NewService(engine *Engine) *Service {
	l, err := listen()
	if err != nil {
		log.Fatalf("Error while listening for connections: %v", err)
	}
	return &Service{
		Engine:   engine,
		listener: l,
	}
}

func (s *Service) Init() {
	conn, err := s.listener.Accept()
	s.conn = conn
	log.Printf("s.conn=%v, err=%v", s.conn, err)
	if err != nil {
		log.Fatalf("Error while accepting a connection: %v", err)
	}
	handshake(s.conn, Server)
	log.Printf("Handshake successfull")
}

func (s *Service) Shutdown() {
	s.conn.Close()
	s.listener.Close()
}

func (s *Service) Send(v interface{}) error {
	log.Printf("Service.Send(%v)", v)
	return send(s.conn, v)
}

func (s *Service) Recv(v interface{}) error {
	return recv(s.conn, v)
}
