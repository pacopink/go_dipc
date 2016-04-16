package msg

import (
	//"bytes"
	"fmt"
	"testing"
	"time"
)

var m CommMsg
var b = make([]byte, 4096)

func TestData(t *testing.T) {
	m.Flag = (BACKUP_FLAG)
	m.Type = 105
	m.TimeStamp = uint64(time.Now().Unix())
	m.Src = "AAA@AAA"
	m.Dst = "BBB@BBB"
	m.Back = "CCC@CCC"

	x := &DataMsg{
		DataBuffer: []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
	}
	x.DataBuffer[4] = 0
	m.Msg = x

	fmt.Println("======== TestData ==========")
	fmt.Println(m)
	l, err := m.Pack(b)
	if err != nil {
		t.Fatalf("Pack Topo Failed:", err)
	}
	fmt.Println("Pack Data:", l, "  Data:", b[0:l])

	var m2 CommMsg
	fmt.Println(b[0:l])
	l_unpack, err := m2.Unpack(b[0:l])
	if err != nil {
		t.Fatalf("Unpack Data Failed:", err)
	}
	if l_unpack != l {
		t.Fatalf("Unpack Data Failed: expect uplen [%d] but get [%d]", l, l_unpack)
	}
	fmt.Println(m2)

	mm, ok := m2.Msg.(*DataMsg) //IMPORTANT: assertion and convert type
	fmt.Println(mm, ok)
	if !ok {
		t.Fatalf("Convert failed")
	}
	fmt.Println(mm.DataBuffer)

	fmt.Println("****** GenerateAppMsg Data ******")
	y := GenerateAppMsg("AAAA@AAAA", "DDDDD@DDDDD", "CCCCCC*CCCCC", 256, []byte("00000000000000000000000000"))
	fmt.Println(y)
	l, err = y.Pack(b)
	if err != nil {
		t.Fatalf("Pack Data Failed:", err)
	}
	fmt.Println("Pack Data:", l, "  Data:", b[0:l])
}

func TestTopo(t *testing.T) {
	m.Flag = (MNG_FLAG | BACKUP_FLAG)
	m.Type = TOPO
	m.TimeStamp = uint64(time.Now().Unix())
	m.Src = "AAA@AAA"
	m.Dst = "BBB@BBB"
	m.Back = "CCC@CCC"
	m.Msg = &TopoMsg{
		Peers:       []string{"AAA@AAA", "BBB@BBB", "CCC@CCC"},
		BackupPeers: []string{"AAA@AAA", "BBB@BBB", "CCC@CCC", "DDD@DDD"},
	}

	fmt.Println("======== TestTopo ==========")
	fmt.Println(m)
	l, err := m.Pack(b)
	if err != nil {
		t.Fatalf("Pack Topo Failed:", err)
	}
	fmt.Println("Pack Topo:", l, "  Data:", b[0:l])

	var m2 CommMsg
	fmt.Println(b[0:l])
	l_unpack, err := m2.Unpack(b[0:l])
	if err != nil {
		t.Fatalf("Unpack Topo Failed:", err)
	}
	if l_unpack != l {
		t.Fatalf("Unpack Topo Failed: expect uplen [%d] but get [%d]", l, l_unpack)
	}
	fmt.Println(m2)

	mm, ok := m2.Msg.(*TopoMsg) //IMPORTANT: assertion and convert type
	fmt.Println(mm, ok)
	if !ok {
		t.Fatalf("Convert failed")
	}
	fmt.Println(mm.Peers, mm.BackupPeers)

	fmt.Println("****** GenerateMngMsg Topo ******")
	x := GenerateMngMsg("AAAA@AAAA", "AGENT@DDDDD", "", TOPO, mm)
	fmt.Println(x)
	l, err = m.Pack(b)
	if err != nil {
		t.Fatalf("Pack Topo Failed:", err)
	}
	fmt.Println("Pack Topo:", l, "  Data:", b[0:l])

}

/*
func TestHB(t *testing.T) {
	m.Flag = (MNG_FLAG | BACKUP_FLAG)
	m.Type = HEART_BEAT_ACK
	m.TimeStamp = uint64(time.Now().Unix())
	m.Src = "AAA@AAA"
	m.Dst = "BBB@BBB"
	m.Back = "CCC@CCC"

	fmt.Println("========= TestHB ===========")
	l, err := m.Pack(b)
	if err != nil {
		t.Fatalf("Pack HB Failed:", err)
	}
	fmt.Println("Pack HB:", l, "  Data:", b[0:l])

	var m2 CommMsg
	fmt.Println(b[0:l])
	l_unpack, err := m2.Unpack(b[0:l])
	if err != nil {
		t.Fatalf("Unpack HB Failed:", err)
	}
	if l_unpack != l {
		t.Fatalf("Unpack HB Failed: expect uplen [%d] but get [%d]", l, l_unpack)
	}
	fmt.Println(m2)

}

func TestReg(t *testing.T) {
	var m CommMsg
	m.Flag = (MNG_FLAG | BACKUP_FLAG)
	m.Type = REGISTER
	m.TimeStamp = uint64(time.Now().Unix())
	m.Src = "AAA@AAA"
	m.Dst = "BBB@BBB"
	m.Back = "CCC@CCC"
	m.Msg = &RegisterMsg{
		8,
	}
	//fmt.Println(m)

	fmt.Println("========= TestReg ===========")
	l, err := m.Pack(b)
	if err != nil {
		t.Fatalf("Pack1 Failed:", err)
	}
	fmt.Println("Pack1:", l, "  Data:", b[0:l])

	var m2 CommMsg
	fmt.Println(b[0:l])
	m2.Unpack(b[0:l])
	fmt.Println(m2)

	mm, ok := m2.Msg.(*RegisterMsg) //IMPORTANT: assertion and convert type
	fmt.Println(mm, ok)
	if !ok {
		t.Fatalf("Convert failed")
	}
	fmt.Println(mm.Action)
}
*/
