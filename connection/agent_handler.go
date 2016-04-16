package connection

import (
	"fmt"
	"go_dipc/dipc_topo"
	"go_dipc/msg"
	"time"
)

func A2A_Accept_MsgHandler(conn *Connection) {
	var comm msg.CommMsg
	b := make([]byte, 8096)
	handshake_faild := true
	register_recv := false
	select {
	case b := <-conn.chanOUT:
		_, err := comm.Unpack(b)
		if err != nil {
			fmt.Println("A2A_Accept_MsgHandler failed to Unpack recv data")
		} else if (comm.Flag & msg.MNG_FLAG) == 0 {
			fmt.Println("A2A_Accept_MsgHandler failed to Unpack recv data")
		} else if comm.Type != msg.REGISTER {
			fmt.Println("A2A_Accept_MsgHandler failed to Unpack recv data")
		} else if comm.Dst != dipc_topo.GlobalTopo.LocalName {
			fmt.Printf("A2A_Accept_MsgHandler recv register to [%s] but not [%s]\n", comm.Dst, dipc_topo.GlobalTopo.LocalName)
		} else {
			reg_msg := comm.GetRegisterMsg()
			if reg_msg == nil {
				fmt.Println("A2A_Accept_MsgHandler failed to convert register msg")
			} else {
				register_recv = true
				if reg_msg.Action != 0 {
					fmt.Println("A2A_Accept_MsgHandler action not register")
				} else {
					handshake_failed = false
				}
			}
		}
	case <-time.After(time.Second * 5):
		fmt.Println("Connection wait for msg timeout")
	}

	if register_recv { //shall send a response
		comm.Src, comm.Dst = comm.Dst, comm.Src //swap src-dst
		comm.Back = ""
		comm.Flag = msg.MNG_FLAG
		comm.Type = msg.REGISTER_ACK
		comm.TimeStamp = time.Now().Unix()
		comm.Msg = &msg.RegisterMsg{
			Action: 0,
		}
		if handshake_faild { //set result to 1 if failed
			comm.Msg.Action = 1
		}
		l, err := comm.Pack(b)
		if err != nil {
			fmt.Println("failed to pack register ack")
		} else {
			for !conn.SendMsg(b[0:l]) {
			}
		}
	}

	if handshake_faild {
		conn.Close()
		return
	}
}

func A2A_Connect_MsgHandler(conn *Connection) {
}

func A2B_Accept_MsgHandler(conn *Connection) {
}
