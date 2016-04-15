package msg

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type TopoMsg struct {
	Peers       []string
	BackupPeers []string
}

func (m *TopoMsg) Pack(b []byte) (int, error) {
	peer_num := 0
	if m.Peers != nil {
		peer_num = len(m.Peers)
	}
	back_num := 0
	if m.BackupPeers != nil {
		back_num = len(m.BackupPeers)
	}
	l := len(b)
	l_eval := 2 + 2 + peer_num*2 + back_num*2 //the at least len we can evaluate at this moment

	if l < l_eval {
		return 0, errors.New(fmt.Sprintf("Insufficient buffer for pack actual len [%d], need [%d]", l, l_eval))
	}
	offset := 0
	//pack peers
	binary.LittleEndian.PutUint16(b[offset:offset+2], uint16(peer_num))
	offset += 2
	l_str := 0
	for i := 0; i < peer_num; i++ {
		l_str = len(m.Peers[i])
		if l < (offset + 2 + l_str) {
			return offset, errors.New(fmt.Sprintf("Insufficient buffer for pack actual len [%d], need [%d]", l, l_eval))
		}
		PackString(b[offset:offset+2+l_str], m.Peers[i])
		offset += 2 + l_str
	}
	//pack backups
	binary.LittleEndian.PutUint16(b[offset:offset+2], uint16(back_num))
	offset += 2
	for i := 0; i < back_num; i++ {
		l_str = len(m.BackupPeers[i])
		if l < (offset + 2 + l_str) {
			return offset, errors.New(fmt.Sprintf("Insufficient buffer for pack actual len [%d], need [%d]", l, l_eval))
		}
		PackString(b[offset:offset+2+l_str], m.BackupPeers[i])
		offset += 2 + l_str
	}
	return offset, nil
}

func (m *TopoMsg) Unpack(b []byte) (int, error) {
	l := len(b)
	l_eval := 2 + 2

	if l < l_eval {
		return 0, errors.New(fmt.Sprintf("Insufficient buffer for unpack actual len [%d], need [%d]", len(b), 1))
	}
	offset := 0
	l_unpack := 0
	var err error
	//unpack peers
	peer_num := int(binary.LittleEndian.Uint16(b[offset : offset+2]))
	offset += 2
	m.Peers = make([]string, peer_num)
	for i := 0; i < peer_num; i++ {
		m.Peers[i], l_unpack, err = UnpackString(b[offset:])
		if err != nil {
			return offset, err
		}
		offset += l_unpack
	}
	//unpack backups
	backup_num := int(binary.LittleEndian.Uint16(b[offset : offset+2]))
	offset += 2
	m.BackupPeers = make([]string, backup_num)
	for i := 0; i < backup_num; i++ {
		m.BackupPeers[i], l_unpack, err = UnpackString(b[offset:])
		if err != nil {
			return offset, err
		}
		offset += l_unpack
	}
	return offset, nil
}

func (m *TopoMsg) String() string {
	str := ""
	n := 0
	if m.Peers != nil {
		n = len(m.Peers)
	}
	str += fmt.Sprintf("PeerNum[%d]\n", n)
	for i := 0; i < n; i++ {
		str += fmt.Sprintf("peer[%d]:[%s]\n", i, m.Peers[i])
	}

	n = 0
	if m.BackupPeers != nil {
		n = len(m.BackupPeers)
	}
	str += fmt.Sprintf("BackupPeerNum[%d]\n", n)
	for i := 0; i < n; i++ {
		str += fmt.Sprintf("backup_peer[%d]:[%s]\n", i, m.BackupPeers[i])
	}
	return str
}
