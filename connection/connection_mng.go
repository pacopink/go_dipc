package connection

import (
	"errors"
	"fmt"
	//"go_dipc/dipc_topo"
	"go_dipc/msg"
	//"strings"
	"time"
)

func (conn *Connection) MngRoutine() {
	defer conn.Close()
	err := HandShake(conn)
	if err != nil {
		return
	}
	//after handshake well, add to proper map change status to CONNECTED
	switch conn.ConnType {
	case CONN_TYPE_A2A_CONNECT:
		fallthrough
	case CONN_TYPE_P2P_CONNECT:
		fallthrough
	case CONN_TYPE_A2A_ACCEPT:
		fallthrough
	case CONN_TYPE_P2P_ACCEPT:
		conn.mgr.SiblingConnMap.Add(conn)
	case CONN_TYPE_P2A_CONNECT:
		fallthrough
	case CONN_TYPE_P2A_ACCEPT:
		conn.mgr.PeerConnMap.Add(conn)
	}
	conn.SetStatus(ST_CONNECTED)

	serve_forever := true
	for serve_forever {
		if conn.GetStatus() != ST_CONNECTED {
			serve_forever = false
		}
		for i := 0; i < PROCESS_MSG_PER_CYCLE; i++ {
			select {
			case m := <-conn.chanIn:
				if m.IsFlagOn(msg.MNG_FLAG) {
					/* TODO if mng msg, process it*/
				} else {
					conn.chanApp <- m
				}
			case <-time.After(time.Millisecond):
				continue
			}
		}
	}
}

func HandShake(conn *Connection) (err error) {
	switch conn.ConnType {
	case CONN_TYPE_A2A_CONNECT:
		fallthrough
	case CONN_TYPE_P2A_CONNECT:
		fallthrough
	case CONN_TYPE_P2P_CONNECT:
		err = ValidateSrcDst(conn.LocalID, conn.ID, conn.ConnType) //for connect type, can be validate here
		if err != nil {
			return err
		}
		err = connecterHandShake(conn)
	case CONN_TYPE_A2A_ACCEPT:
		fallthrough
	case CONN_TYPE_P2P_ACCEPT:
		fallthrough
	case CONN_TYPE_P2A_ACCEPT:
		err = acceptorHandShake(conn)
	default:
		err = errors.New(fmt.Sprintf("HandShake conn_type [%s] unknow", conn.ConnType))
	}
	return
}

func connecterHandShake(conn *Connection) (err error) {
	/*send a register first*/
	err = SendRegisterOrAck(conn, false, 0)
	if err != nil {
		return
	}
	/*wait for register_ack and check the result*/
	m, err := WaitForMsg(conn, msg.REGISTER_ACK, 5*time.Second)
	if err != nil {
		return
	}
	if m == nil {
		err = errors.New("ConnecterHandShake get nil msg")
		return
	}
	if m.Dst != conn.LocalID {
		err = errors.New(fmt.Sprintf("ConnecterHandShake recv reg_ack dst [%s] is not connection local id [%s]", m.Dst, conn.LocalID))
		return
	}
	if m.Src != conn.ID {
		err = errors.New(fmt.Sprintf("ConnecterHandShake recv reg_ack src [%s] is not connection id [%s]", m.Src, conn.ID))
		return
	}
	r := m.GetRegisterMsg()
	if r == nil {
		err = errors.New("ConnecterHandShake recv reg_ack but failed to decode")
		return
	}
	if r.Action != 0 {
		err = errors.New("ConnecterHandShake recv reg_ack ")
		return
	}
	/* Handshaked well*/
	return nil
}

func acceptorHandShake(conn *Connection) (err error) {
	/*wait for register_ack and check the result*/
	m, err := WaitForMsg(conn, msg.REGISTER, 5*time.Second)
	if err != nil {
		return
	}
	if m == nil {
		err = errors.New("AcceptorHandShake get nil msg")
		return
	}
	if m.Src == "" {
		err = errors.New("AcceptorHandShake recv reg src empty")
		return
	}
	/* now we begin to know who is the connector so save it*/
	conn.ID = m.Src
	/*and we can response with a register_ack after the rest check, don't return directly*/
	err = ValidateSrcDst(conn.LocalID, conn.ID, conn.ConnType)
	if err != nil {
		return
	}
	if conn.mgr.IsConnectionExistByID(m.Src) {
		err = errors.New(fmt.Sprintf("AcceptorHandShake recv reg src [%s] is invalid or already exist connection", m.Src))
	}
	r := m.GetRegisterMsg()
	if r == nil {
		err = errors.New("AcceptorHandShake recv reg but failed to decode")
	}
	if r.Action != 0 {
		err = errors.New("AcceptorHandShake recv reg but action not register")
	}
	/* Send Response */
	action := uint8(0) //by default response register_ack with successful
	if err != nil {
		action = uint8(1) //if any error, response with unsuccess
	}
	SendRegisterOrAck(conn, true, action)
	return
}

func SendRegisterOrAck(conn *Connection, is_ack bool, action uint8) (err error) {
	if conn.LocalID == conn.ID || conn.LocalID == "" || conn.ID == "" {
		errors.New("SendRegisterOrAck failed, invalid ID")
	} else {
		msg_type := msg.REGISTER
		if is_ack {
			msg_type = msg.REGISTER_ACK
		}
		m := msg.GenerateMngMsg(conn.LocalID, conn.ID, "", msg_type, &msg.RegisterMsg{Action: action})
		if !conn.SendMsg(m) {
			err = errors.New(fmt.Sprintf("SendRegister failed to send reg [%s]", m.String()))
		}
	}
	return
}

func WaitForMsg(conn *Connection, msg_type uint16, t time.Duration) (m *msg.CommMsg, err error) {
	m = nil
	err = nil
	select {
	case m = <-conn.chanIn:
		if m.IsFlagOn(msg.MNG_FLAG) && m.Type == msg_type {
			return
		} else {
			err = errors.New(fmt.Sprintf("WaitForMsg: expecting msg type [%s]  but receive [%s][%d]", msg.MsgType2Str(msg_type), msg.MsgType2Str(m.Type), m.Type))
			return
		}
	case <-time.After(t):
		err = errors.New(fmt.Sprintf("WaitForMsg: Timeout to wait for recv [%s] failed", msg.MsgType2Str(msg_type)))
		return
	}
	return
}
