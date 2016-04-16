package common_util

import (
	"errors"
	//"fmt"
	"time"
)

type Switcher struct {
	chan_input  chan int8
	chan_output chan int8
}

func New() *Switcher {
	return &Switcher{
		chan_input:  make(chan int8, 1),
		chan_output: make(chan int8, 1),
	}
}

/* for contoller, to trigger some job */
func (rc *Switcher) Trigger() {
	rc.chan_input <- 1 //fire a signal to trigger some task
}

/* for controller, to wait til some job finish */
func (rc *Switcher) Wait(sec_to_wait int) error {
	if sec_to_wait <= 0 {
		<-rc.chan_output
		return nil
	} else {
		select {
		case <-rc.chan_output:
			return nil
		case <-time.After(time.Second * time.Duration(sec_to_wait)):
			return errors.New("Timeout")
		}
	}
	return nil
}

/* for worker, to wait for trigger from controller */
func (rc *Switcher) WaitTrigger() <-chan int8 {
	return rc.chan_input
}

/* for worker, to notify complete of the job */
func (rc *Switcher) Complete() {
	rc.chan_output <- 1
}
