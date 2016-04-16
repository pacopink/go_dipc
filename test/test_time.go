package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()
	utc := now.UTC()
	fmt.Println("Unix:", now)
	fmt.Println("UTC:", utc)

	fmt.Println(time.Now().Unix())
	x := time.Now().UnixNano()
	for i := 0; i < 5000; i++ {
		time.Sleep(100 * time.Microsecond)
	}
	y := time.Now().UnixNano()
	fmt.Println(y - x)
	fmt.Println(float64(y-x) / 1000000000)
}
