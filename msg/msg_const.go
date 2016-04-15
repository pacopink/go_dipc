package msg

//uint8
const (
	MNG_FLAG      = uint8(1)
	COMPRESS_FLAG = uint8(1 << 1)
	JSON_FLAG     = uint8(1 << 2)
	BACKUP_FLAG   = uint8(1 << 3)
)

//uint16 message type definiation
const (
	/*MNG*/
	REGISTER       = uint16(1)
	REGISTER_ACK   = uint16(2)
	HEART_BEAT     = uint16(3)
	HEART_BEAT_ACK = uint16(4)
	TOPO_REQ       = uint16(5)
	TOPO           = uint16(6)
)

func Flag2Str(flag uint8) string {
	str := ""
	if (flag & MNG_FLAG) != 0 {
		str += "mng-"
	}
	if (flag & COMPRESS_FLAG) != 0 {
		str += "comp-"
	}
	if (flag & JSON_FLAG) != 0 {
		str += "json-"
	}
	if (flag & BACKUP_FLAG) != 0 {
		str += "back-"
	}
	return str
}

func MsgType2Str(msg_type uint16) string {
	str := "UNKNOWN"
	switch {
	case msg_type == REGISTER:
		str = "REGISTER"
	case msg_type == REGISTER_ACK:
		str = "REGISTER_ACK"
	case msg_type == HEART_BEAT:
		str = "HEART_BEAT"
	case msg_type == HEART_BEAT_ACK:
		str = "HEART_BEAT_ACK"
	case msg_type == TOPO:
		str = "TOPO"
	case msg_type == TOPO_REQ:
		str = "TOPO_REQ"
	default:
		str = "UNKNOWN"
	}
	return str
}
