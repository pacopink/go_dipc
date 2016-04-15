package main

import (
	"flag"
	"fmt"
	"runtime"
	"time"
)

var Debug = false

func LogPrinter(clog chan string) {
	ss := ""
	for {
		ss = <-clog
		fmt.Println(ss)
	}
}

func Worker(cin chan bool, cout chan bool, clog chan string) {
	if Debug {
		clog <- "Worker Begin"
	}
	defer func() {
		cout <- true
		if Debug {
			clog <- "Worker End"
		}
	}()
	bStop := false
	for !bStop {
		select {
		case x := <-cin:
			if Debug {
				clog <- "Worker Recv"
			}
			if x {
				bStop = true
			}
		}
	}
}

func main() {
	debug_mode := flag.Bool("debug", false, "is in debug mode")
	freq_gc := flag.Bool("gc", false, "frequnt gc mode")
	flag.Parse()
	Debug = *debug_mode

	clog := make(chan string, 1000)
	if Debug {
		go LogPrinter(clog)
	}
	for {
		cin := make(chan bool, 1)
		cout := make(chan bool, 1)
		if Debug {
			clog <- "Cycle Begin"
		}
		go Worker(cin, cout, clog)
		cin <- true
		bStop := false
		for !bStop {
			select {
			case x := <-cout:
				if Debug {
					clog <- "Main Recv"
				}
				if x {
					bStop = true
				}
			}
		}
		if Debug {
			clog <- "Cycle End\n======================"
			time.Sleep(time.Second)
		}
		if *freq_gc {
			runtime.GC()
		}
	}
}
