package msg

import (
	"time"
)

func GenerateMngMsg(src, dest, back string, msg_type uint16, msg MsgInterface) *CommMsg {
	return &CommMsg{
		Flag:      MNG_FLAG,
		Type:      msg_type,
		TimeStamp: uint64(time.Now().Unix()),
		Src:       src,
		Dst:       dest,
		Back:      back,
		Msg:       msg,
	}
}

func GenerateAppMsg(src, dest, back string, msg_type uint16, b []byte) *CommMsg {
	flag := uint8(0)
	if len(back) > 0 {
		flag = flag | BACKUP_FLAG
	}
	return &CommMsg{
		Flag:      flag,
		Type:      msg_type,
		TimeStamp: uint64(time.Now().Unix()),
		Src:       src,
		Dst:       dest,
		Back:      back,
		Msg: &DataMsg{
			DataBuffer: b,
		},
	}
}
