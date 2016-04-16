package connection

import (
	"go_dipc/msg"
)

type ConnectionMgr struct {
	MyName         string
	SiblingConnMap *ConnectionTable //for agent: connections with other agents via TCPSock, for Peers, connections with other peers via UnixSock
	PeerConnMap    *ConnectionTable //for agent: connections from Peers via UnixSock, for Peers, connections with Agent via UnixSock
	TmpList        *ConnectionList
	ChanSibling    chan msg.CommMsg
	ChanAgent      chan msg.CommMsg
}

func NewConnectionMgr() *ConnectionMgr {
	return &ConnectionMgr{
		SiblingConnMap: NewConnectionTable("Sibling"),
		PeerConnMap:    NewConnectionTable("Peer"),
		TmpList:        NewConnectionList(),
		ChanSibling:    make(chan msg.CommMsg, BUFFER_SIZE),
		ChanAgent:      make(chan msg.CommMsg, BUFFER_SIZE),
	}
}

func (mgr *ConnectionMgr) IsConnectionExistByID(id string) (ret bool) {
	ret = mgr.SiblingConnMap.IsConnectionExistByID(id)
	if ret {
		return
	}
	ret = mgr.PeerConnMap.IsConnectionExistByID(id)
	if ret {
		return
	}
	return mgr.TmpList.IsConnectionExistByID(id)
}
