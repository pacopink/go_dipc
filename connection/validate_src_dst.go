package connection

import (
	"errors"
	"fmt"
	"go_dipc/peername"
	"strings"
)

func ValidateSrcDstWhenNewConn(local_name, remote_name string, conn_type string) (err error) {
	if !ValidateConnectionType(conn_type) {
		return errors.New(fmt.Sprintf("ValidateSrcDst %s type unknown", conn_type))
	}
	switch conn_type {
	case CONN_TYPE_A2A_CONNECT:
		fallthrough
	case CONN_TYPE_P2A_CONNECT:
		fallthrough
	case CONN_TYPE_P2P_CONNECT:
		return ValidateSrcDst(local_name, remote_name, conn_type)
	case CONN_TYPE_A2A_ACCEPT:
		local_peer, _, err := peername.ParseFullPeerName(local_name)
		if err != nil {
			return err
		}
		if local_peer != "agent" {
			return errors.New(fmt.Sprintf("ValidateSrcForAccept %s invalid peer local[%s]", conn_type, local_peer))
		}
	case CONN_TYPE_P2A_ACCEPT:
		fallthrough
	case CONN_TYPE_P2P_ACCEPT:
		local_peer, _, err := peername.ParseFullPeerName(local_name)
		if err != nil {
			return err
		}
		if local_peer != "agent" {
			return errors.New(fmt.Sprintf("ValidateSrcForAccept %s invalid peer local[%s]", conn_type, local_peer))
		}
	}
	return nil
}

func ValidateSrcDst(local_name, remote_name string, conn_type string) (err error) {
	if !ValidateConnectionType(conn_type) {
		return errors.New(fmt.Sprintf("ValidateSrcDst %s type unknown", conn_type))
	}
	local_peer, local_hostname, err := peername.ParseFullPeerName(local_name)
	if err != nil {
		return err
	}
	remote_peer, remote_hostname, err := peername.ParseFullPeerName(remote_name)
	if err != nil {
		return err
	}

	switch conn_type {
	case CONN_TYPE_A2A_CONNECT:
		fallthrough
	case CONN_TYPE_A2A_ACCEPT:
		if local_peer != "agent" || remote_peer != "agent" {
			return errors.New(fmt.Sprintf("ValidateSrcDst %s invalid peer local[%s] remote[%s]", conn_type, local_peer, remote_peer))
		}
		ret := strings.Compare(local_name, remote_name)
		if conn_type == CONN_TYPE_A2A_CONNECT && ret >= 0 {
			return errors.New(fmt.Sprintf("ValidateSrcDst %s local[%s] is no lower than remote [%s]", conn_type, local_name, remote_name))
		}
		if conn_type == CONN_TYPE_A2A_ACCEPT && ret <= 0 {
			return errors.New(fmt.Sprintf("ValidateSrcDst %s local[%s] is no higher than remote [%s]", conn_type, local_name, remote_name))
		}
	case CONN_TYPE_P2A_ACCEPT:
		if local_peer != "agent" || remote_peer == "agent" {
			return errors.New(fmt.Sprintf("ValidateSrcDst %s invalid peer local[%s] remote[%s]", conn_type, local_peer, remote_peer))
		}
		if local_hostname != remote_hostname {
			return errors.New(fmt.Sprintf("ValidateSrcDst %s not at the same host local[%s] remote[%s]", conn_type, local_name, remote_name))
		}
	case CONN_TYPE_P2A_CONNECT:
		if local_peer == "agent" || remote_peer != "agent" {
			return errors.New(fmt.Sprintf("ValidateSrcDst %s invalid peer local[%s] remote[%s]", conn_type, local_peer, remote_peer))
		}
		if local_hostname != remote_hostname {
			return errors.New(fmt.Sprintf("ValidateSrcDst %s not at the same host local[%s] remote[%s]", conn_type, local_name, remote_name))
		}
	case CONN_TYPE_P2P_CONNECT:
		fallthrough
	case CONN_TYPE_P2P_ACCEPT:
		if local_peer == "agent" || remote_peer == "agent" {
			return errors.New(fmt.Sprintf("ValidateSrcDst %s invalid peer local[%s] remote[%s]", conn_type, local_peer, remote_peer))
		}
		if local_hostname != remote_hostname {
			return errors.New(fmt.Sprintf("ValidateSrcDst %s not at the same host local[%s] remote[%s]", conn_type, local_name, remote_name))
		}
		if conn_type == CONN_TYPE_P2P_CONNECT && strings.Compare(local_name, remote_name) >= 0 {
			return errors.New(fmt.Sprintf("ValidateSrcDst %s local[%s] is no lower than remote [%s]", conn_type, local_name, remote_name))
		}
		if conn_type == CONN_TYPE_P2P_ACCEPT && strings.Compare(local_name, remote_name) <= 0 {
			return errors.New(fmt.Sprintf("ValidateSrcDst %s local[%s] is no higher than remote [%s]", conn_type, local_name, remote_name))
		}
	default:
		return errors.New(fmt.Sprintf("ValidateSrcDst %s type unknown", conn_type))
	}
	return nil
}
