package connection

import (
	"errors"
	"fmt"
	"sync"
)

type ConnectionTable struct {
	Name        string
	Connections map[string]Connection
	mutex       sync.Mutex
}

func NewConnectionTable(name string) *ConnectionTable {
	return &ConnectionTable{
		Name:        name,
		Connections: make(map[string]Connection),
	}
}

func (table *ConnectionTable) IsConnectionExistByID(id string) (ret bool) {
	table.mutex.Lock()
	defer func() { table.mutex.Unlock() }()
	_, ret = table.Connections[id]
	return
}

func (table *ConnectionTable) Add(conn Connection) (err error) {
	table.mutex.Lock()
	defer table.mutex.Unlock()
	err = nil
	if conn.ID == "" {
		return errors.New("ConnectionTable::Add failed to add connection with empty ID")
	}
	c, present := table.Connections[conn.ID]
	if present { //if exists same ID connection
		if c.GetStatus() == ST_CONNECTED { //if it is connected, refuse to add
			return errors.New(fmt.Sprintf("ConnectionTable::Add, failed to add connection with ID[%s], a connected connection with this ID is already existing", conn.ID))
		} else { //if it is not connected, add
			c.Close()
			table.Connections[conn.ID] = conn
		}
	} else { //if not existing, directly add
		table.Connections[conn.ID] = conn
	}
	return
}

func (table *ConnectionTable) Remove(id string) *Connection {
	table.mutex.Lock()
	defer table.mutex.Unlock()
	conn, present := table.Connections[id]
	if !present {
		return nil
	} else {
		delete(table.Connections, id)
		return &conn
	}
}

func (table *ConnectionTable) GetConnByID(id string) *Connection {
	table.mutex.Lock()
	defer table.mutex.Unlock()
	conn, present := table.Connections[id]
	if present {
		return &conn
	} else {
		return nil
	}
}

func (table *ConnectionTable) String() (str string) {
	table.mutex.Lock()
	defer table.mutex.Unlock()
	str = fmt.Sprintf("ConnectionTable[%s]\n", table.Name)
	count := 0
	for k, v := range table.Connections {
		count++
		str += v.String()
	}
	str += fmt.Sprintf("ConnectionTable Size[%d]\n", count)
	return str
}
