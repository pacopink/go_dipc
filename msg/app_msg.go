package msg

import (
	"errors"
	"fmt"
)

type DataMsg struct {
	DataBuffer []byte
}

func (m *DataMsg) Pack(b []byte) (int, error) {
	if m.DataBuffer == nil {
		return 0, nil
	}
	if len(b) < len(m.DataBuffer) {
		return 0, errors.New(fmt.Sprintf("Insufficient buffer for pack actual len [%d], need [%d]", len(b), len(m.DataBuffer)))
	}
	copy(b, m.DataBuffer)
	return len(m.DataBuffer), nil
}

func (m *DataMsg) Unpack(b []byte) (int, error) {
	m.DataBuffer = b
	return len(b), nil
}

func (m *DataMsg) String() string {
	return fmt.Sprintf("[%v]", m.DataBuffer)
}
