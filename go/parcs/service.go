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
	var err error
	s.conn, err = s.listener.Accept()
	if err != nil {
		log.Fatalf("Error while accepting a connection: %v", err)
	}
}

func (s *Service) Shutdown() {
	s.conn.Close()
	s.listener.Close()
}

func (s *Service) Send(v interface{}) error {
	return send(s.conn, v)
}

func (s *Service) Recv(v interface{}) error {
	return recv(s.conn, v)
}
