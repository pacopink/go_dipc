package connection

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type ConnectionList struct {
	conn_list list.List
	mutex     sync.Mutex
}

func NewConnectionList() *ConnectionList {
	l := new(ConnectionList)
	go l.PeriodicallyCheck() //automatically check for remove
	return l
}

func (l *ConnectionList) PeriodicallyCheck() {
	for {
		l.Check()
		time.Sleep(time.Second * 10)
	}
}

func (l *ConnectionList) IsConnectionExistByID(id string) (ret bool) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	ret = false
	e := l.conn_list.Front()
	for e != nil {
		conn, ok := e.Value.(Connection)
		if ok {
			if id == conn.ID {
				ret = true
				break
			}
			e = e.Next()
		}
	}
	return
}

func (l *ConnectionList) Add(conn *Connection) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.conn_list.PushBack(conn)
}

func (l *ConnectionList) Remove(conn *Connection) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.Remove(conn)
}

func (l *ConnectionList) Check() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	//check any DISCONNECTED, remove it
	e := l.conn_list.Front()
	for e != nil {
		conn, ok := e.Value.(Connection)
		if !ok || conn.GetStatus() != ST_UNKNOWN { //only keep unknown connection in list
			to_remove := e
			e = e.Next()
			l.conn_list.Remove(to_remove)
		} else {
			e = e.Next()
		}
	}
}

func (l *ConnectionList) String() (str string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	str = "ConnectionList\n"
	count := 0
	e := l.conn_list.Front()
	for e != nil {
		count++
		conn, ok := e.Value.(Connection)
		if !ok {
			str += "Invalid element\n"
		} else {
			str += conn.String()
		}
		e = e.Next()
	}
	str += fmt.Sprintf("ConnectionList Size[%d]\n", count)
	return
}
