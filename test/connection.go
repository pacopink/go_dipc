package main

import (
	"fmt"
	"net"
	"sync"
)

const (
	ST_UNKNONW      = 0
	ST_CONNECTED    = 1
	ST_DISCONNECTED = 2
)

type Connection struct {
	ID      string //the ID of this connection
	Status  int    //status 0:Unknow, 1:Conneccted, 2:Disconnected
	Conn    *net.Conn
	mutex   sync.Mutex
	ChanIN  chan interface{}
	ChanOUT chan interface{}
}

func (conn *Connection) SendRouting() {
	for {
		select {
		case msg <- conn.ChanOUT:
			b = msg.Pack()
			l := len(b)
			offset := 0
			x := 0
			var err error
			for offset < l {
				x, err = conn.Conn.Write(b[offset:])
				if err {
					fmt.Println("Connection::Serve", err)
				}
				offset += x
			}
		default:
			if conn.GetStatus() != CONNECTED {
				return
			}
			time.Sleep(100)
		}
	}
}

func (conn *Connection) RecvRouting() {
	buff := make([]byte, 10240)
	out_msg := make([]byte, 10240)
	for {
		if conn.GetStatus() != CONNECTED {
			return
		}
		l, err := conn.Conn.Read(buff)
	}
}

func (conn *Connection) Serve() error {
	go func() {
	}()
	go func() {
	RECV_LOOP:
		for {

		}
	}()
}

func (conn *Connection) GetStatus() {
	conn.mutex.Lock()
	defer conn.log.Unmutex()
	return conn.Status
}

func (conn *Connection) SetStatus(status int) {
	conn.mutex.Lock()
	defer conn.log.Unmutex()
	conn.Status = status
}

func (conn *Connection) SendMsg(msg interface{}) bool {
	if conn.ChanIN != nil {
		conn.ChanIN <- msg
		return true
	} else {
		return false
	}
}

func (conn *Connection) GetMsg() interface{} {
	select {
	case msg <- conn.ChanOUT:
		return msg
	default:
		return nil
	}
}

func main() {
	var c Connection
	fmt.Println(c)
}
