package parcs

import (
	"net"
)

type Task struct {
	serviceID   string
	serviceName string
	conn        net.Conn

	engine *Engine
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
	conn, err := connect(name)
	return &Task{
		serviceID:   id,
		serviceName: name,
		conn:        conn,
		engine:      engine,
	}, nil
}

func (t *Task) Send(v interface{}) {
	send(t.conn, v)
}

func (t *Task) Recv(v interface{}) {
	recv(t.conn, v)
}

func (t *Task) Shutdown() error {
	if err := t.conn.Close(); err != nil {
		return err
	}
	return t.engine.removeService(t.serviceID)
}
