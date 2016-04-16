package msg

/* RAW msg parser
DELIMITER|int16 DataLen|Data Bytes
*/

import (
	"bytes"
	"encoding/binary"
)

/* Delimiter of messages in STREAM */
const DEL = "&#%$"
const DEL_LEN = len(DEL)

const (
	ST_UNKNONW      = 0
	ST_CONNECTED    = 1
	ST_DISCONNECTED = 2
)

/* Buffer for Recv */
type MsgBuff struct {
	Len    int
	Used   int
	Buffer []byte
}

func NewMsgBuff(l int) *MsgBuff {
	if l <= 0 {
		l = 1024
	}
	m := &MsgBuff{
		Len:  l,
		Used: 0,
	}
	m.Buffer = make([]byte, l)
	return m
}

func (m *MsgBuff) GetMsg(c chan *CommMsg) int {
	msg_count := 0
	offset := 0
	for {
		begin := bytes.Index(m.Buffer[offset:m.Used], []byte(DEL)) //find delimiter
		if begin < 0 {
			if m.Used == m.Len { //buffer full but no delimiter found, discard the buffer
				m.Used = 0
			}
			break
		}
		if offset == 0 && begin > 0 {
			offset = begin
		}
		begin += DEL_LEN
		if begin+2 >= m.Used { //incomplete msg in buffer
			break
		}
		l := int(binary.LittleEndian.Uint16(m.Buffer[begin : begin+2]))
		if l+begin+2 > m.Used { //incomplete msg in buffer
			break
		}
		msg := new(CommMsg)
		_, err := msg.Unpack(m.Buffer[begin:+2 : begin+2+l])
		if err != nil {
			c <- msg //send decoded msg to channel
		}
		msg_count++
		offset += begin + 2 + l
	}

	/* move forward bytes remaining in buffer*/
	if offset > 0 {
		remain_len := m.Used - offset
		if remain_len > 0 {
			copy(m.Buffer[0:], m.Buffer[offset:m.Used])
		}
		m.Used = remain_len
	}
	return msg_count
}
