package nets

import (
	"fmt"
)

type InitialFunc func(conn *Conn) interface{}

type SimpleConnStorage struct {
	ConnMap     map[*Conn]interface{}
	InitialFunc InitialFunc
}

func NewSimpleConnStorage(initialFunc InitialFunc) *SimpleConnStorage {
	return &SimpleConnStorage{
		ConnMap:     make(map[*Conn]interface{}),
		InitialFunc: initialFunc,
	}
}

func (t *SimpleConnStorage) AddConnData(conn *Conn, data interface{}) {
	if _, ok := t.ConnMap[conn]; ok {
		return
	}
	t.ConnMap[conn] = data
	fmt.Println("AddConn", len(t.ConnMap), data)
}

func (t *SimpleConnStorage) AddConn(conn *Conn) {
	var data interface{}
	if t.InitialFunc != nil {
		data = t.InitialFunc(conn)
	}
	t.AddConnData(conn, data)
}

func (t *SimpleConnStorage) RemoveConn(conn *Conn) interface{} {
	// remove from ConnMap
	data, ok := t.ConnMap[conn]
	if !ok {
		return nil
	}
	delete(t.ConnMap, conn)
	fmt.Println("RemoveConn", len(t.ConnMap), data)
	return data
}

func (t *SimpleConnStorage) DefaultHandshakeDoneHandler() HandshakeDoneHandler {
	return func(ctx *Context) {
		t.AddConn(ctx.Conn)
	}
}

func (t *SimpleConnStorage) DefaultClosedHandler() ClosedHandler {
	return func(ctx *Context) {
		t.RemoveConn(ctx.Conn)
	}
}
