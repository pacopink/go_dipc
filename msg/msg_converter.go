package msg

func MakeMsgByFlagAndMsgType(flag uint8, msg_type uint16) MsgInterface {
	if (flag & MNG_FLAG) != 0 {
		switch msg_type {
		case REGISTER:
			return new(RegisterMsg)
		case REGISTER_ACK:
			return new(RegisterMsg)
		case TOPO:
			return new(TopoMsg)
		default:
			return nil
		}
	} else {
		return new(DataMsg)
	}
}

func (m *CommMsg) GetRegisterMsg() *RegisterMsg {
	if m.Msg == nil {
		return nil
	}
	r, ok := m.Msg.(*RegisterMsg) //Assertion
	if ok {
		return r
	} else {
		return nil
	}
}

func (m *CommMsg) GetTopoMsg() *TopoMsg {
	if m.Msg == nil {
		return nil
	}
	r, ok := m.Msg.(*TopoMsg) //Assertion
	if ok {
		return r
	} else {
		return nil
	}
}

func (m *CommMsg) GetDataMsg() *DataMsg {
	if m.Msg == nil {
		return nil
	}
	r, ok := m.Msg.(*DataMsg) //Assertion
	if ok {
		return r
	} else {
		return nil
	}
}
