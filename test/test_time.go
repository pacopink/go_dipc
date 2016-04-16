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
}
