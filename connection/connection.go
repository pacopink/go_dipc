package connection

/*
A connection that send/recv message from a connection
messages are delimitered via DEL
heartbeat mechanism implemented
*/
import (
	//"errors"
	"encoding/binary"
	"fmt"
	"go_dipc/msg"
	"net"
	"sync"
	"time"
)

const (
	ST_UNKNOWN      = 0
	ST_CONNECTED    = 1
	ST_DISCONNECTED = 2
)

var BUFFER_SIZE = 10000

func State2Str(st int) (str string) {
	switch st {
	case ST_CONNECTED:
		str = "CONNECTED"
	case ST_DISCONNECTED:
		str = "DISCONNECTED"
	default:
		str = "UNKNOWN"
	}
	return
}

/* b: the input pack for process
   c: the channel receiving app level msg */
type mngFunc func(conn *Connection)

type Connection struct {
	LocalID    string //the local ID of this connection
	ID         string //the ID of this connection
	Status     int    //status 0:Unknow, 1:Conneccted, 2:Disconnected
	Conn       net.Conn
	mutex      sync.Mutex
	chanIN     chan []byte
	chanOUT    chan []byte
	chanAPP    chan msg.CommMsg
	LastActSec int64
	mgr        *ConnectionMgr
}

func (conn *Connection) String() (str string) {
	c_i, l_i, c_o, l_o := 0, 0, 0, 0
	if conn.chanIN != nil {
		c_i, l_i = cap(conn.chanIN), len(conn.chanIN)
	}
	if conn.chanOUT != nil {
		c_o, l_o = cap(conn.chanOUT), len(conn.chanOUT)
	}
	str = fmt.Sprintf("Conn [%s - %s], Status[%s], IN [%d][%d], OUT[%d][%d]\n", conn.LocalID, conn.ID, State2Str(conn.GetStatus()), c_i, l_i, c_o, l_o)
	return
}

func NewConnection(l_id string, r_id string, conn net.Conn, app_chan chan msg.CommMsg, mngHandler mngFunc, mgr *ConnectionMgr) (*Connection, error) {
	appConn := &Connection{
		LocalID:    l_id,
		ID:         r_id,
		Status:     ST_UNKNOWN,
		Conn:       conn,
		chanIN:     make(chan []byte, BUFFER_SIZE),
		chanOUT:    make(chan []byte, BUFFER_SIZE),
		chanAPP:    nil, //leave it nil here, mngHandler shall set it properly
		LastActSec: time.Now().Unix(),
		mgr:        mgr,
	}
	go appConn.RecvRoutine()
	go appConn.SendRoutine()
	go mngHandler(appConn)
	return appConn, nil
}

func (conn *Connection) Close() {
	if conn.Conn != nil {
		conn.Conn.Close()
	}
	conn.SetStatus(ST_DISCONNECTED)
	//remove my self from mgr
	conn.mgr.SiblingConnMap.Remove(conn.ID)
	conn.mgr.PeerConnMap.Remove(conn.ID)
}

func (conn *Connection) UpdateLastActSec() {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	conn.LastActSec = time.Now().Unix()
}

func (conn *Connection) GetIdleSec() int64 {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	return time.Now().Unix() - conn.LastActSec
}

/* Ensure the b bytes are send completely, else wait and block */
func (conn *Connection) SendBytes(b []byte) error {
	l := len(b)
	offset := 0
	x := 0
	var err error
	for offset < l {
		x, err = conn.Conn.Write(b[offset:])
		if err != nil {
			return err
		}
		offset += x
	}
	return nil
}

func (conn *Connection) SendRoutine() {
	defer conn.Close()
	len_byte := make([]byte, 2)
	serve_forever := true
	for serve_forever {
		select {
		case b := <-conn.chanOUT:
			//send delimitor
			err := conn.SendBytes([]byte(msg.DEL))
			if err != nil {
				fmt.Println("SendRoutine SendDEL:", err)
				return
			}
			binary.LittleEndian.PutUint16(len_byte, uint16(len(b)))
			err = conn.SendBytes(len_byte)
			if err != nil {
				fmt.Println("SendRoutine SendLen:", err)
				return
			}
			err = conn.SendBytes(b)
			if err != nil {
				fmt.Println("SendRoutine SendData:", err)
				return
			}
		default:
			if conn.GetStatus() != ST_CONNECTED {
				return
			}
			time.Sleep(100)
		}
	}
}

func (conn *Connection) RecvRoutine() {
	defer conn.Close()
	buf := msg.NewMsgBuff(66560) //65kb
	serve_forever := true
	for serve_forever {
		if conn.GetStatus() != ST_CONNECTED {
			return
		}
		conn.Conn.SetReadDeadline(time.Now().Add(time.Second * 1))
		l, err := conn.Conn.Read(buf.Buffer[buf.Used:]) //read to MsgBuffer
		if err != nil {
			op_err, ok := err.(*net.OpError)
			if ok {
				switch {
				case op_err.Timeout():
					continue
				case op_err.Temporary():
					continue
				default:
					fmt.Println("RecvRoutine:", op_err)
					serve_forever = false
				}
			} else {
				fmt.Println("RecvRoutine:", op_err)
				serve_forever = false
			}

		} else {
			conn.UpdateLastActSec()
			buf.Used += l
			buf.GetMsg(conn.chanIN) //decode msg from MsgBuffer and send to channel
		}
	}
}

/* GetStatus thread-safe */
func (conn *Connection) GetStatus() int {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	return conn.Status
}

/* SetStatus thread-safe */
func (conn *Connection) SetStatus(status int) {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	conn.Status = status
}

/* SendMsg in non-blocking way */
func (conn *Connection) SendMsg(b []byte) bool {
	if conn.chanIN != nil {
		select {
		case conn.chanIN <- b:
			return true
		default:
			return false
		}
	} else {
		return false
	}
}

/* GetMsg in non-blocking way */
func (conn *Connection) GetMsg() []byte {
	select {
	case b := <-conn.chanOUT:
		return b
	default:
		return nil
	}
}
