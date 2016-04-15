package msg

import (
	"bytes"
	"fmt"
	"testing"
)

func TestStringCodec(t *testing.T) {
	str := "PACOLI"
	b := make([]byte, 1024)
	l, _ := PackString(b[0:], str)
	if l != 2+len(str) {
		t.Fatalf("PackString failed [%d] != [%d] as expected\n", l, 2+len(str))
	}
	fmt.Println("After Pack:", b[0:l])
	str_out, ll, _ := UnpackString(b[0:l])
	if ll != l {
		t.Fatalf("UnpackString failed [%d] != [%d] as expected\n", ll, l)
	}
	if str_out != str {
		t.Fatalf("UnpackString failed [%s] != [%s] as expected\n", str_out, str)
	}
	fmt.Println("After Unpack:", str_out)
}

func TestByteCodec(t *testing.T) {
	str := []byte("PACOLI")
	b := make([]byte, 1024)
	l, _ := PackBytes(b[0:], str)
	if l != 2+len(str) {
		t.Fatalf("PackBytes failed [%d] != [%d] as expected\n", l, 2+len(str))
	}
	fmt.Println("After Pack:", b[0:l])
	str_out, ll, _ := UnpackBytes(b[0:l])
	if ll != l {
		t.Fatalf("UnpackBytes failed [%d] != [%d] as expected\n", ll, l)
	}
	if !bytes.Equal(str_out, str) {
		t.Fatalf("UnpackBytes failed [%v] != [%v] as expected\n", str_out, str)
	}
	fmt.Println("After Unpack:", str_out)
}
