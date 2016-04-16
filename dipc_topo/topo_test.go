package dipc_topo

import (
	"fmt"
	"go_dipc/msg"
	"testing"
)

func TestTopo(t *testing.T) {
	m := &msg.TopoMsg{
		Peers:       []string{"AAA@AAA", "BBB@BBB", "CCC@CCC"},
		BackupPeers: []string{"AAA@AAA", "BBB@BBB", "CCC@CCC", "DDD@DDD"},
	}
	fmt.Println(m)

	fmt.Println(GlobalTopo.UpdatePeers(m.Peers))
	fmt.Println(GlobalTopo.UpdateBackup(m.BackupPeers))
	fmt.Println(GlobalTopo)
}
