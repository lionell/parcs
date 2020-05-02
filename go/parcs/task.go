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
	if err != nil {
		return nil, err
	}
	name, err := engine.queryServiceName(id)
	if err != nil {
		return nil, err
	}
	conn := connect(name)
	if err != nil {
		return nil, err
	}
	if err := handshake(conn, Client); err != nil {
		return nil, err
	}
	log.Printf("Connection to %s established", name)
	return &Task{
		serviceID:   id,
		serviceName: name,
		conn:        conn,
		engine:      engine,
	}, nil
}

func (t *Task) Send(v interface{}) error {
	return send(t.conn, v)
}

func (t *Task) SendAll(vs ...interface{}) error {
	return sendAll(t.conn, vs...)
}

func (t *Task) Recv(v interface{}) error {
	return recv(t.conn, v)
}

func (t *Task) Shutdown() error {
	if err := t.conn.Close(); err != nil {
		return err
	}
	if err := t.engine.removeService(t.serviceID); err != nil {
		return err
	}
	log.Printf("Connection to %s closed", t.serviceName)
	return nil
}
