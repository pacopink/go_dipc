package sock

import (
	"errors"
	"net"
	"os"
)

type UnixConnector struct {
	Path         string
	chan_stop    chan bool //input channel to signal stop
	chan_stopped chan bool //output channel to signal stopped
	ChanConn     chan net.Conn
	listener     *UnixListener
	IsListening  bool
}

func Accept(l net.Listener, c chan net.Conn, s chan bool, sc chan bool) {
	defer func() { sc <- true }() //signal a stop complete event
	serve_forever := true
	for serve_forever {
		select {
		case <-s:
			serve_forever = false
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

func (acpt *UnixConnector) Start() error {
	if actp.IsListening {
		return errors.New("UnixConnector is already listening, cannot Start twice")
	}
	os.Remove(acpt.path)
	addr := net.UnixAddr{acpt.path, "unix"}
	acpt.listener, err = net.ListenUnix("unix", &addr)
	if err != nil {
		return err
	}
	go Accept(acpt.listener, acpt.ChanConn, acpt.chan_stop, actp.chan_stopped)
	actp.IsListening = true
	return nil
}

func (acpt *UnixConnector) Stop() error {
	defer func() {
		if actp.listener != nil {
			actp.IsListening = false
			actp.listener.Close()
		}
	}()
	acpt.chan_stop <- true //signal the goroutine to stop
	_ <- actp.chan_stopped //wait for the stop complete signale
}

func NewUnixConnector(path string) (*UnixConnector, error) {
	if len(path) < 5 || path[len(path)-5:len(path)] != ".sock" {
		return nil, errors.New(fmt.Sprintf("NewUnixAccpetor: invalid path [%s], should end with \".sock\"", path))
	}
	acpt := &UnixConnector{
		Path:         path,
		chan_stop:    make(chan bool),
		chan_stopped: make(chan bool),
		ChanConn:     make(chan net.Conn, 10),
		IsListening:  false,
	}
	err = actp.Start()
	if err != nil {
		return nil, err
	} else {
		return acpt, nil
	}
}
