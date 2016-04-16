package common_util

import (
	"fmt"
	"testing"
	"time"
)

func Worker(rc Switcher) {
	fmt.Println("START JOB")
	rc.WaitTrigger()
	fmt.Println("START TRIGGERING JOB")
	time.Sleep(time.Second * 10)
	fmt.Println("STOP JOB")
	rc.Complete()
}

func TestSwitcher(t *testing.T) {
	rc := New()
	go Worker(*rc)
	fmt.Println("TRIGGER")
	rc.Trigger()
	fmt.Println("WAITING")
	err := rc.Wait(5)
	if err != nil {
		fmt.Println(err)
	}
	rc.Wait(0)
	fmt.Println("COMPLETE")
}
