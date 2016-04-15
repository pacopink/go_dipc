package msg

import (
	"encoding/binary"
	"errors"
	"fmt"
)

func PackString(b []byte, str string) (int, error) {
	if len(b) < len(str)+2 {
		return 0, errors.New(fmt.Sprintf("PackString insufficient buffer [%d] need [%d]", len(b), len(str)+2))
	}
	binary.LittleEndian.PutUint16(b[0:2], uint16(len(str)))
	copy(b[2:], []byte(str))
	return 2 + len(str), nil
}

func UnpackString(b []byte) (string, int, error) {
	if len(b) < 2 {
		return "", 0, errors.New(fmt.Sprintf("UnpackString insufficient buffer [%d]", len(b)))
	}
	l := int(binary.LittleEndian.Uint16(b[0:2]))
	if len(b) < l+2 {
		fmt.Println(b)
		return "", 0, errors.New(fmt.Sprintf("UnpackString insufficient buffer [%d] expect [%d]", len(b), l+2))
	}
	return string(b[2 : 2+l]), 2 + l, nil
}

func PackBytes(b []byte, b2pack []byte) (int, error) {
	if len(b) < len(b2pack)+2 {
		return 0, errors.New(fmt.Sprintf("PackBytes insufficient buffer [%d] need [%d]", len(b), len(b2pack)+2))
	}
	binary.LittleEndian.PutUint16(b[0:2], uint16(len(b2pack)))
	copy(b[2:], b2pack)
	return 2 + len(b2pack), nil
}

func UnpackBytes(b []byte) ([]byte, int, error) {
	if len(b) < 2 {
		return nil, 0, errors.New(fmt.Sprintf("UnpackBytes insufficient buffer [%d]", len(b)))
	}
	l := int(binary.LittleEndian.Uint16(b[0:2]))
	if len(b) < l+2 {
		return nil, 0, errors.New(fmt.Sprintf("UnpackBytes insufficient buffer [%d] expect [%d]", len(b), l+2))
	}
	return b[2 : 2+l], 2 + l, nil
}
