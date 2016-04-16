/*
A full peername is composed of PEERNAME@HOSTNAME
A PEERNAME is composed of SERVICE_NAME.PEER_IDENTIFIER
Case sensitive
*/
package peername

import (
	"errors"
	"regexp"
	"strings"
)

func ParseFullPeerName(fullname string) (peername string, hostname string, err error) {
	v := strings.Split(fullname, "@")
	if len(v) != 2 {
		err = errors.New("ParsePeerName failed")
		return
	}
	peername = v[0]
	hostname = v[1]
	return
}

func PeerPrefix(name string, prefix string) bool {
	peer, _, _ := ParseFullPeerName(name)
	return strings.Index(peer, prefix) == 0
}

func PeerPattern(name string, pattern string) bool {
	peer, _, err := ParseFullPeerName(name)
	if err != nil {
		return false
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return re.MatchString(peer)
}

func PeerService(name string, servicename string) bool {
	peer, host, err := ParseFullPeerName(name)
	if err != nil {
		return false
	}
	v := strings.Split(peer, ".")
	if len(v) < 2 {
		return false
	}
	return v[0] == servicename
}

func IsBelong2Host(name string, host string) bool {
	_, peer_host, _ := ParseFullPeerName(name)
	return host == peer_host

}

func IsAgent(name string) bool {
	peer, host, err := ParseFullPeerName(name)
	if err != nil {
		return false
	}
	return peer == "agent"
}
