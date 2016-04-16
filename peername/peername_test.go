package peername

import (
	//"fmt"
	"testing"
)

func TestParseFullPeerName(t *testing.T) {
	var ep, eh, fn, p, h string
	var err error
	ep = "peername.xxxx.yyyyy"
	eh = "ericsson.com"
	fn = ep + "@" + eh
	p, h, err = ParseFullPeerName(fn)
	if err != nil {
		t.Fatalf("ParseFullPeerName[%s] get error: %v\n", fn, err)
	} else if p != ep || h != eh {
		t.Fatalf("ParseFullPeerName[%s] get [%s][%s] expecting [%s][%s]\n", p, h, ep, eh)
	}
	fn = "kdjfkdfjdkfjdkfjdkfj"
	p, h, err = ParseFullPeerName(fn)
	if err == nil {
		t.Fatalf("ParseFullPeerName[%s] should get err but not\n", fn)
	}
	fn = "aaa@bbb@ccc"
	p, h, err = ParseFullPeerName(fn)
	if err == nil {
		t.Fatalf("ParseFullPeerName[%s] should get err but not\n", fn)
	}
}

func TestPeerPrefix(t *testing.T) {
	var ep, eh, fn, p string
	//var err error
	ep = "peername.xxxx.yyyyy"
	eh = "ericsson.com"
	fn = ep + "@" + eh
	p = "peername.xx"
	if !PeerPrefix(fn, p) {
		t.Fatalf("PeerPrefix [%s] [%s] should match but not matched\n", fn, p)
	}
	p = "xpeername."
	if PeerPrefix(fn, p) {
		t.Fatalf("PeerPrefix [%s] [%s] should not match but not matched\n", fn, p)
	}
}

func TestPeerPattern(t *testing.T) {
	var ep, eh, fn, p string
	//var err error
	ep = "peername.xxxx.yyyyy"
	eh = "ericsson.com"
	fn = ep + "@" + eh
	p = "^peer.*\\.xxxx\\.*"
	if !PeerPattern(fn, p) {
		t.Fatalf("PeerPattern [%s] [%s] should match but not matched\n", fn, p)
	}
	p = "\\.xxxx\\."
	if !PeerPattern(fn, p) {
		t.Fatalf("PeerPattern [%s] [%s] should match but not matched\n", fn, p)
	}
	p = "peername.yyy.*"
	if PeerPattern(fn, p) {
		t.Fatalf("PeerPattern [%s] [%s] should not match but not matched\n", fn, p)
	}
}
