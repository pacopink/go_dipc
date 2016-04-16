package connection

/*
A connection that send/recv message from a connection
messages are delimitered via DEL
heartbeat mechanism implemented
*/
import (
	//"errors"
	"encoding/binary"
	"errors"
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
var PROCESS_MSG_PER_CYCLE = 2000
var IDLE_REST = 10 * time.Millisecond

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

const (
	CONN_TYPE_A2A_CONNECT = "A2A_CONNECT"
	CONN_TYPE_A2A_ACCEPT  = "A2A_ACCEPT"
	CONN_TYPE_P2A_ACCEPT  = "P2A_ACCEPT"
	CONN_TYPE_P2A_CONNECT = "P2A_CONNECT"
	CONN_TYPE_P2P_CONNECT = "P2P_CONNECT"
	CONN_TYPE_P2P_ACCEPT  = "P2P_ACCEPT"
)

func ValidateConnectionType(conn_type string) bool {
	switch conn_type {
	case CONN_TYPE_A2A_CONNECT:
		return true
	case CONN_TYPE_A2A_ACCEPT:
		return true
	case CONN_TYPE_P2A_ACCEPT:
		return true
	case CONN_TYPE_P2A_CONNECT:
		return true
	case CONN_TYPE_P2P_CONNECT:
		return true
	case CONN_TYPE_P2P_ACCEPT:
		return true
	default:
		return false
	}
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
	chanOut    chan *msg.CommMsg
	chanIn     chan *msg.CommMsg
	chanApp    chan *msg.CommMsg
	LastActSec int64
	ConnType   string
	mgr        *ConnectionMgr
}

func (conn *Connection) String() (str string) {
	c_i, l_i, c_o, l_o := 0, 0, 0, 0
	if conn.chanIn != nil {
		c_i, l_i = cap(conn.chanIn), len(conn.chanIn)
	}
	if conn.chanOut != nil {
		c_o, l_o = cap(conn.chanOut), len(conn.chanOut)
	}
	str = fmt.Sprintf("Conn [%s - %s], Status[%s], IN [%d][%d], OUT[%d][%d]\n", conn.LocalID, conn.ID, State2Str(conn.GetStatus()), c_i, l_i, c_o, l_o)
	return
}

func NewConnection(l_id string, r_id string, conn_type string, conn net.Conn, app_chan chan *msg.CommMsg, mngHandler mngFunc, mgr *ConnectionMgr) (*Connection, error) {
	if !ValidateConnectionType(conn_type) {
		return nil, errors.New(fmt.Sprintf("ValidateSrcDst %s type unknown", conn_type))
	}
	err := ValidateSrcDstWhenNewConn(l_id, r_id, conn_type)
	if err != nil {
		return nil, err
	}
	appConn := &Connection{
		LocalID:    l_id,
		ID:         r_id,
		Status:     ST_UNKNOWN,
		Conn:       conn,
		chanIn:     make(chan *msg.CommMsg, BUFFER_SIZE),
		chanOut:    make(chan *msg.CommMsg, BUFFER_SIZE),
		chanApp:    app_chan,
		LastActSec: time.Now().Unix(),
		ConnType:   conn_type,
		mgr:        mgr,
	}
	go appConn.RecvRoutine()
	go appConn.SendRoutine()
	go appConn.MngRoutine()
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
func (conn *Connection) sendBytes(b []byte) error {
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
	buf := make([]byte, 66560)                //65kb
	copy(buf[0:msg.DEL_LEN], []byte(msg.DEL)) //prepare the fixed msg delimitor

	serve_forever := true
	for serve_forever {
		for i := 0; i < PROCESS_MSG_PER_CYCLE; i++ {
			select {
			case m := <-conn.chanOut:
				l, err := m.Pack(buf[msg.DEL_LEN+2:]) //pack the data
				if err != nil {
					fmt.Println("SendPackRoutine Pack failed:", err)
				} else {
					binary.LittleEndian.PutUint16(buf[msg.DEL_LEN:msg.DEL_LEN+2], uint16(l)) //pack the length of data
					err = conn.sendBytes(buf[0 : msg.DEL_LEN+2+l])                           //sendout the msg
				}
			default:
				time.Sleep(100)
			}
		}
		if conn.GetStatus() > ST_CONNECTED {
			serve_forever = false
		}
	}
}

func (conn *Connection) RecvRoutine() {
	defer conn.Close()
	buf := msg.NewMsgBuff(66560) //65kb
	serve_forever := true
	for serve_forever {
		if conn.GetStatus() > ST_CONNECTED {
			serve_forever = false
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
			buf.GetMsg(conn.chanIn) //decode msg from MsgBuffer and send to channel
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
func (conn *Connection) SendMsg(m *msg.CommMsg) bool {
	if conn.chanIn != nil {
		select {
		case conn.chanIn <- m:
			return true
		default:
			return false
		}
	} else {
		return false
	}
}

/* GetMsg in non-blocking way */
func (conn *Connection) GetMsg() *msg.CommMsg {
	select {
	case m := <-conn.chanOut:
		return m
	default:
		return nil
	}
}
