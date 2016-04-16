package msg

import (
	//"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

/* All upper level msg shall implement this interface */
type MsgInterface interface {
	Pack(b []byte) (int, error)
	Unpack(b []byte) (int, error)
	String() string
}

/* A common msg, contains a upper level msg */
type CommMsg struct {
	Flag      uint8  // bit 0: Is MNG, bit 1: Is Compressed, bit 2: Json, bit 3: Send backup, reserve for others
	Type      uint16 //
	TimeStamp uint64 //
	Src       string //PeerName@HostName
	Dst       string //PeerName@HostName
	Back      string //PeerName@HostName or ''
	Msg       MsgInterface
}

func (m *CommMsg) IsFlagOn(flag uint8) bool {
	return (m.Flag & flag) != 0
}

func (m *CommMsg) Pack(br []byte) (int, error) {
	var err error
	//evaluate the length of the Pack result
	l_eval := 1 + 2 + 8 + /*src peer*/ 2 + len(m.Src) + /* dest peer*/ 2 + len(m.Dst) + /*backup peer*/ 2 + len(m.Back) + 2 /*inner msg len*/
	if len(br) < l_eval {
		return 0, errors.New(fmt.Sprintf("Insufficient buffer for pack actual len [%d], need [%d]", len(br), l_eval))
	}
	//pack inner msg
	li := 0
	if m.Msg != nil {
		li, err = m.Msg.Pack(br[l_eval:])
		if err != nil {
			return 0, err
		}
	}
	binary.LittleEndian.PutUint16(br[l_eval-2:l_eval], uint16(li))
	//pack common part
	offset := 0
	br[offset] = m.Flag
	offset += 1
	binary.LittleEndian.PutUint16(br[offset:offset+2], m.Type)
	offset += 2
	binary.LittleEndian.PutUint64(br[offset:offset+8], m.TimeStamp)
	offset += 8
	l, _ := PackString(br[offset:], m.Src)
	offset += l
	l, _ = PackString(br[offset:], m.Dst)
	offset += l
	l, _ = PackString(br[offset:], m.Back)
	offset += l
	offset += li
	return l_eval + li, nil
}

func (m *CommMsg) Unpack(b []byte) (int, error) {
	l := len(b)
	eval_l := 1 + 2 + 8 + 2 + 2 + 2 + 2 //at this moment the len at least
	if l < eval_l {
		return 0, errors.New("Too short bytes for unpack")
	}
	offset := 0
	m.Flag = b[offset]
	offset += 1
	m.Type = binary.LittleEndian.Uint16(b[offset : offset+2])
	offset += 2
	m.TimeStamp = binary.LittleEndian.Uint64(b[offset : offset+8])
	offset += 8
	//unpack src
	l_str := 0
	l_str = int(binary.LittleEndian.Uint16(b[offset : offset+2]))
	offset += 2
	eval_l += l_str
	if l < eval_l {
		return offset, errors.New("Too short bytes for unpack")
	}
	m.Src = string(b[offset : offset+l_str])
	offset += l_str
	//unpack dst
	l_str = int(binary.LittleEndian.Uint16(b[offset : offset+2]))
	offset += 2
	eval_l += l_str
	if l < eval_l {
		return offset, errors.New("Too short bytes for unpack")
	}
	m.Dst = string(b[offset : offset+l_str])
	offset += l_str
	//unpack back
	l_str = int(binary.LittleEndian.Uint16(b[offset : offset+2]))
	offset += 2
	eval_l += l_str
	if l < eval_l {
		return offset, errors.New("Too short bytes for unpack")
	}
	m.Back = string(b[offset : offset+l_str])
	offset += l_str
	//unpack inner msg len
	l_str = int(binary.LittleEndian.Uint16(b[offset : offset+2]))
	offset += 2
	eval_l += l_str
	if l < eval_l {
		return offset, errors.New("Too short bytes for unpack")
	}
	fmt.Printf("offset[%d] l_str[%d] [%v]\n", offset, l_str, b[offset:offset+l_str])

	m.Msg = MakeMsgByFlagAndMsgType(m.Flag, m.Type)
	if l_str > 0 {
		if m.Msg != nil {
			l_unpack, err := m.Msg.Unpack(b[offset : offset+l_str])
			if err != nil {
				return offset, err
			}
			offset += l_unpack
		} else {
			return offset, errors.New(fmt.Sprintf("Failed to MakeMsg struct from flag[%d][%s] msg_type[%d][%s]", m.Flag, Flag2Str(m.Flag), m.Type, MsgType2Str(m.Type)))
		}
	}
	return offset, nil
}

func (m *CommMsg) String() string {
	flag := Flag2Str(m.Flag)

	inner := ""
	if m.Msg != nil {
		inner = m.Msg.String()
	}

	if m.IsFlagOn(MNG_FLAG) {
		msg_type := MsgType2Str(m.Type)
		return fmt.Sprintf("flag[%s] type[%s] src[%s] dst[%s] back[%s] payload[%s]",
			flag, msg_type, m.Src, m.Dst, m.Back, inner)
	} else {
		return fmt.Sprintf("flag[%s] type[%d] src[%s] dst[%s] back[%s] payload[%s]",
			flag, m.Type, m.Src, m.Dst, m.Back, inner)
	}
}
