package main

import (
	//"bytes"
	//"encoding/binary"
	"flag"
	//"flex_dipc"
	"fmt"
	"os"
	//"os/signal"
	//"runtime"
	//"strconv"
	"net"
	//"runtime/pprof"
	"time"
	//"sync"
)

const PREFIX string = "FIPC"

type e2a_register struct {
	ep_name string
}

var sLn string
var sWorkPath string
var sUnixSock string

func init_arg() {
	host := flag.String("host", "", "hostname")
	work_path := flag.String("work", "./", "work path")
	flag.Parse()
	sLn = *host
	sWorkPath = *work_path
	sUnixSock = sWorkPath + "/" + sLn + ".sock"
	fmt.Printf("ln[%s] work[%s] unixsock[%s]\n", sLn, sWorkPath, sUnixSock)
}

func EchoFunc(conn net.Conn) {
	defer func() {
		fmt.Printf("Connection end: [%s]\n", conn)
		conn.Close()
	}()

	buf := make([]byte, 1024)
	fmt.Printf("Connection comes: [%s]\n", conn)
	for {
		_, err := conn.Read(buf)
		if err != nil {
			//println("Error reading:", err.Error())
			return
		}

		//send reply
		_, err = conn.Write(buf)
		if err != nil {
			//println("Error send reply:", err.Error())
			return
		}
	}
}

func Accept(c chan net.Conn, e chan bool, l net.Listener) {
	for {
		select {
		case <-e:
			break
		default:
			con, err := l.Accept()
			if err != nil {
				fmt.Printf("failed to accept socket: %s\n", err)
			} else {
				c <- con
			}
		}
	}
}

func main() {
	init_arg()
	os.Remove(sUnixSock)
	addr := net.UnixAddr{sUnixSock, "unix"}
	l, err := net.ListenUnix("unix", &addr)
	if err != nil {
		fmt.Printf("failed to open socket: %s\n", err)
		os.Exit(1)
	}
	defer l.Close()

	con_chan := make(chan net.Conn, 10)
	end_chan := make(chan bool, 1)

	go Accept(con_chan, end_chan, l)
	var c net.Conn
	for {
		select {
		case c = <-con_chan:
			go EchoFunc(c)
		default:
			time.Sleep(50 * time.Millisecond) //sleep for retry
		}
	}
}
