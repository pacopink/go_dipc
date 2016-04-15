package msg

import (
	"errors"
	"fmt"
)

type RegisterMsg struct {
	Action uint8 //0: Register, non-0:Deregister, 0: Ack OK, 1: Ack Not OK
}

func (m *RegisterMsg) Pack(b []byte) (int, error) {
	if len(b) < 1 {
		return 0, errors.New(fmt.Sprintf("Insufficient buffer for pack actual len [%d], need [%d]", len(b), 1))
	}
	b[0] = m.Action
	return 1, nil
}

func (m *RegisterMsg) Unpack(b []byte) (int, error) {
	if len(b) < 1 {
		return 0, errors.New(fmt.Sprintf("Insufficient buffer for unpack actual len [%d], need [%d]", len(b), 1))
	}
	m.Action = b[0]
	return 1, nil
}

func (m *RegisterMsg) String() string {
	return fmt.Sprintf("[%d]", m.Action)
}
