package connection

/*
A connection that send/recv message from a connection
messages are delimitered via DEL
heartbeat mechanism implemented
*/
import (
	"errors"
	"fmt"
	"msg"
	"net"
	"sync"
	"time"
)

const BUFFER_SIZE = 10000

/* b: the input pack for process
   c: the channel receiving app level msg */
type mngFunc func(b []byte, c chan []byte) error

type Connection struct {
	ID         string //the ID of this connection
	Status     int    //status 0:Unknow, 1:Conneccted, 2:Disconnected
	Conn       *net.Conn
	mutex      sync.Mutex
	chanIN     chan []byte
	chanOUT    chan []byte
	chanAPP    chan msg.CommonMsg
	LastActSec int64
}

func New(id string, conn net.Conn, mng_handler mngFunc) (*Connection, error) {
	appConn := &Connection{
		ID:         id,
		Status:     0,
		Conn:       conn,
		chanIN:     make(chan []byte, BUFFER_SIZE),
		chanOUT:    make(chan []byte, BUFFER_SIZE),
		chanAPP:    make(chan []byte, BUFFER_SIZE),
		LastActSec: time.Now().Unix(),
	}

	go appConn.RecvRoutine()
	go appConn.SendRoutine()
	go appConn.MngRoutine(mng_handler)
	return appConn, nil
}

func (conn *Conenction) UpdateLastActSec() {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	conn.LastActSec = time.Unix()
}

func (conn *Connection) GetIdleSec() {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	return time.Unix() - conn.LastActSec
}

/* Ensure the b bytes are send completely, else wait and block */
func (conn *Connection) SendBytes(b []byte) error {
	l := len(b)
	offset := 0
	x := 0
	var err error
	for offset < l {
		x, err = conn.Conn.Write(b[offset:])
		if err {
			return err
		}
		offset += x
	}
}

func (conn *Connection) SendRoutine() {
	len_byte := make([]byte, 2)
	serve_forever := true
	for serve_forever {
		select {
		case msg <- conn.chanOUT:
			//send delimitor
			err := conn.SendBytes(msg.DEL)
			if err != nil {
				fmt.Println("SendRoutine SendDEL:", err)
				return
			}
			err = conn.SendBytes(binary.LittleEndian.PutInt16(len_byte, int16(len(msg))))
			if err != nil {
				fmt.Println("SendRoutine SendLen:", err)
				return
			}
			err = conn.SendBytes(msg)
			if err != nil {
				fmt.Println("SendRoutine SendData:", err)
				return
			}
		default:
			if conn.GetStatus() != CONNECTED {
				return
			}
			time.Sleep(100)
		}
	}
}

func (conn *Connection) RecvRoutine() {
	buf := NewMsgBuff(66560) //65kb
	serve_forever := true
	for serve_forever {
		if conn.GetStatus() != CONNECTED {
			return
		}
		conn.Conn.SetReadDeadline(time.Now().Add(t.Second * 2))
		l, err := conn.Conn.Read(buf.Buffer[buf.Used:]) //read to MsgBuffer
		if err != nil {
			switch {
			case err.Timeout():
			case err.Temporary():
				continue
			default:
				fmt.Println("RecvRoutine:", err)
				serve_forever = false
			}
		} else {
			conn.UpdateLastActSec()
			buf.Used += l
			buf.GetMsg(conn.chanIN) //decode msg from MsgBuffer and send to channel
		}
	}
}

func (conn *Connection) MngRoutine(handler mngFunc) {
	for {
		select {
		case b <- conn.chanOUT:
			err = handler(b, conn.chanAPP)
		}
	}
}

func (conn *Connection) Serve() error {
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
	if conn.chanIN != nil {
		conn.chanIN <- msg
		return true
	} else {
		return false
	}
}

func (conn *Connection) GetMsg() interface{} {
	select {
	case msg <- conn.chanOUT:
		return msg
	default:
		return nil
	}
}

func main() {
	var c Connection
	fmt.Println(c)
}
