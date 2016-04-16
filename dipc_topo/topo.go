package dipc_topo

import (
	"fmt"
	"go_dipc/msg"
	"sort"
	"sync"
)

/* Application level peer status */
const (
	PST_UNKNOWN   = 0
	PST_CONCERNED = 1
	PST_BROKEN    = 2
)

type Peer struct {
	PeerName string
	Status   int
}

type Topo struct {
	LocalName string
	Peers     map[string]Peer
	Backup    map[string]string
	mutex     sync.Mutex
}

func (topo *Topo) GetBackup(id string) string {
	topo.mutex.Lock()
	defer func() { topo.mutex.Unlock() }()

	back, present := topo.Backup[id]
	if present {
		return back
	} else {
		return ""
	}
}

func (topo *Topo) GetMyBackup() string {
	topo.mutex.Lock()
	defer func() { topo.mutex.Unlock() }()

	back, present := topo.Backup[topo.LocalName]
	if present {
		return back
	} else {
		return ""
	}
}

func (topo *Topo) UpdatePeers(peers []string) (added []string, removed []string) {
	topo.mutex.Lock()
	defer func() { topo.mutex.Unlock() }()

	for _, v := range peers {
		fmt.Println(v)
		_, present := topo.Peers[v]
		if present {
			continue
		}
		if added == nil {
			added = make([]string, 1)
		}
		added = append(added, v)
		topo.Peers[v] = Peer{
			PeerName: v,
			Status:   PST_UNKNOWN,
		}
	}
	for k, _ := range topo.Peers {
		found := false
		for _, v := range peers {
			if k == v {
				found = true
				break
			}
		}
		if !found {
			if removed == nil {
				removed = make([]string, 1)
			}
			removed = append(removed, k)
			delete(topo.Peers, k)
		}
	}
	return
}

func (topo *Topo) UpdateBackup(backups []string) bool {
	topo.mutex.Lock()
	defer func() { topo.mutex.Unlock() }()

	l := len(backups)
	if l%2 != 0 {
		l = l - 1
	}
	if l < 2 {
		return false
	}
	topo.Backup = make(map[string]string)
	for i, j := 0, 1; j < l; i, j = i+2, j+2 {
		topo.Backup[backups[i]] = backups[j]
		topo.Backup[backups[j]] = backups[i]
	}
	return true
}

func (topo *Topo) GenTopoMsg() *msg.TopoMsg {
	topo.mutex.Lock()
	defer func() { topo.mutex.Unlock() }()

	var copy_peer []string
	for k, _ := range topo.Peers {
		if copy_peer == nil {
			copy_peer = make([]string, 1)
		}
		copy_peer = append(copy_peer, k)
	}

	var copy_backup []string
	for k, v := range topo.Backup {
		if copy_backup == nil {
			copy_backup = make([]string, 1)
		}
		//a peername shall appear twice, don't add duplicate
		for _, j := range copy_backup {
			if k == j {
				continue
			}
		}
		copy_backup = append(copy_backup, k)
		copy_backup = append(copy_backup, v)
	}
	if copy_peer != nil {
		sort.Strings(copy_peer)
	}
	return &msg.TopoMsg{
		Peers:       copy_peer,
		BackupPeers: copy_backup,
	}
}

var GlobalTopo = &Topo{
	Peers:  make(map[string]Peer),
	Backup: make(map[string]string),
}
