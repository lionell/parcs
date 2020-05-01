package parcs

import (
	"log"
	"net"
)

type Task struct {
	serviceID   string
	serviceName string
	conn        net.Conn
	engine      *Engine
}

func NewTask(image string, engine *Engine) (*Task, error) {
	id, err := engine.createService(image)
	log.Printf("Created a service with ID '%s'", id)
	if err != nil {
		return nil, err
	}
	name, err := engine.queryServiceName(id)
	log.Printf("Service is called '%s'", name)
	if err != nil {
		return nil, err
	}
	conn, err := connect(name)
	log.Printf("conn=%v, err=%v", conn, err)
	if err != nil {
		return nil, err
	}
	handshake(conn, Client)
	log.Printf("Connection to %s established", name)
	return &Task{
		serviceID:   id,
		serviceName: name,
		conn:        conn,
		engine:      engine,
	}, nil
}

func (t *Task) Send(v interface{}) error {
	log.Printf("Trying to send %v", v)
	return send(t.conn, v)
}

func (t *Task) Recv(v interface{}) error {
	return recv(t.conn, v)
}

func (t *Task) Shutdown() error {
	if err := t.conn.Close(); err != nil {
		return err
	}
	return t.engine.removeService(t.serviceID)
}
