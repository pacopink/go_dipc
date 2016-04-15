package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

func TestMap() {
	mm := make(map[string]string)
	mm["Jan"] = "January"
	mm["Feb"] = "Febrary"
	mm["Mar"] = "March"

	for k, v := range mm {
		fmt.Println(k, v)
	}
	var v string
	var present bool
	v, present = mm["Jan"]
	if present {
		fmt.Println(v)
	} else {
		fmt.Println("Jan is not exist")
	}
	delete(mm, "Jan")
	v, present = mm["Jan"]
	if present {
		fmt.Println(v)
	} else {
		fmt.Println("Jan is not exist")
	}
	for k, v := range mm {
		fmt.Println(k, v)
	}

}

func ParseFullName(fn string, local_hostname string) (string, string, error) {
	l := strings.Split(fn, ".")
	var hostname, ln = "", ""
	if len(l) > 1 {
		ln = l[len(l)-1]
		hostname = strings.Join(l[0:len(l)-1], ".")
		return hostname, ln, nil
	} else {
		ln = fn
		hostname = local_hostname
		return hostname, ln, nil
	}
}

func Byte2Str(bs []byte) string {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.LittleEndian, bs)
	if err == nil {
		return buff.String()
	} else {
		return ""
	}
}

func Str2Byte(str string) []byte {
	return []byte(str)
}

func main() {
	var x [8]byte
	fmt.Println(x)
	x[0] = 0x33
	x[1] = 0x31
	x[2] = 0x34
	x[3] = 0x30
	x[4] = 0x31
	s := Byte2Str(x[:])
	fmt.Println(len(s))
	fmt.Printf("[%s]", s)
	fmt.Printf("[%v]", Str2Byte(s))

	h, l, _ := ParseFullName("www.ericsson.com.some_application", "localhost")
	fmt.Printf("[%s] [%s]\n", h, l)
	h, l, _ = ParseFullName("some_application", "localhost")
	fmt.Printf("[%s] [%s]\n", h, l)

	TestMap()
}
