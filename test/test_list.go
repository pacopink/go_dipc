package main

import (
	"container/list"
	"fmt"
	//"net"
	//"bytes"
	"encoding/binary"
	//"strings"
	"sync"
)

const MAX_MSG_LEN = 10240
const MAX_RAW_MSG_LEN = MAX_MSG_LEN + 1024
const MAX_HOST_LEN = 32
const MAX_PEER_LEN = 32
const MAX_FULL_LEN = MAX_HOST_LEN + MAX_PEER_LEN + 1

type MyErr struct {
	errmsg string
}

func (err MyErr) Error() string {
	return err.errmsg
}

type UpStreamMsg struct {
	orig_peer [MAX_FULL_LEN]byte
	dest_peer [MAX_FULL_LEN]byte
	flag      uint8
	data_len  uint16
	data      [MAX_MSG_LEN]byte
}

func (usm *UpStreamMsg) Dump() {
	fmt.Printf("orig[%s] dest[%s] flag[%d] len[%d] data[%v]\n",
		string(usm.orig_peer[:]), string(usm.dest_peer[:]), usm.flag,
		usm.data_len, usm.data[0:usm.data_len])
}

func (usm *UpStreamMsg) DecodeFromBytes(buff []byte) error {
	buff_len := len(buff)
	offset := 0

	if buff_len < len(usm.orig_peer)+len(usm.dest_peer)+1+2 {
		return &MyErr{"bytes buffer len too short"}
	}
	copy(usm.orig_peer[0:len(usm.orig_peer)], buff[offset:offset+len(usm.orig_peer)])
	offset += len(usm.orig_peer)
	copy(usm.dest_peer[0:len(usm.dest_peer)], buff[offset:offset+len(usm.dest_peer)])
	offset += len(usm.dest_peer)
	usm.flag = (uint8)(buff[offset])
	offset += 1
	usm.data_len = binary.LittleEndian.Uint16(buff[offset : offset+2])
	offset += 2
	if buff_len != offset+int(usm.data_len) {
		fmt.Printf("Unpack [%d] [%d] [%d] Not OK\n", offset, usm.data_len, buff_len)
		return &MyErr{"data len not match"}
	}
	copy(usm.data[0:len(usm.data)], buff[offset:offset+int(usm.data_len)])
	fmt.Printf("Unpack [%d] [%d] [%d] OK\n", offset, usm.data_len, buff_len)
	return nil
}

func (usm *UpStreamMsg) Serialize() []byte {
	var b [MAX_MSG_LEN]byte
	l, e := usm.SerializeToBytes(b[0:], len(b))
	if e == nil {
		return b[0:l]
	} else {
		return nil
	}
}

func (usm *UpStreamMsg) SerializeToBytes(buff []byte, max_len int) (int, error) {
	offset := 0

	if max_len < len(usm.orig_peer)+len(usm.dest_peer)+2+2+int(usm.data_len) {
		return 0, &MyErr{"insufficient buffer length"}
	}

	copy(buff[offset:offset+len(usm.orig_peer)], usm.orig_peer[0:len(usm.orig_peer)])
	offset += len(usm.orig_peer)
	copy(buff[offset:offset+len(usm.dest_peer)], usm.dest_peer[0:len(usm.dest_peer)])
	offset += len(usm.dest_peer)
	buff[offset] = byte(usm.flag)
	offset += 1
	binary.LittleEndian.PutUint16(buff[offset:offset+2], usm.data_len)
	offset += 2
	copy(buff[offset:], usm.data[0:usm.data_len])
	return offset + int(usm.data_len), nil
}

type UpStreamPeer struct {
	id      string
	channel chan []byte
}

type DownStreamPeer struct {
}

var l = new(list.List)
var lock = new(sync.Mutex)

func InitList() {
	for i := 0; i < 100; i++ {
		l.PushBack(i)
	}
}

func IterList(wc *sync.WaitGroup) {
	defer func() {
		wc.Done()
	}()

	for {
		lock.Lock()
		e := l.Front()
		if e != nil {
			fmt.Printf("element [%d]\n", e.Value)
			l.Remove(e)
			lock.Unlock()
		} else {
			lock.Unlock()
			break
		}
	}
}

func main() {
	InitList()
	//fmt.Printf("list len[%d]\n", l.Len())
	wg := new(sync.WaitGroup)
	wg.Add(4)
	go IterList(wg)
	go IterList(wg)
	go IterList(wg)
	go IterList(wg)
	wg.Wait()

	us := &UpStreamMsg{
		flag:     1 + 2,
		data_len: 0,
	}
	copy(us.orig_peer[0:], "AAA")
	copy(us.dest_peer[0:], "BBB")
	s := "Hello World"
	copy(us.data[0:], s)
	us.data_len = uint16(len(s))
	us.Dump()
	var b [MAX_MSG_LEN]byte
	l, e := us.SerializeToBytes(b[0:], len(b))
	if e == nil {
		var us2 UpStreamMsg
		us2.DecodeFromBytes(b[0:l])
		us2.Dump()
		fmt.Println(us2.Serialize())
	} else {
		fmt.Println("Failed to SerializeToBytes:", e)
	}
}
