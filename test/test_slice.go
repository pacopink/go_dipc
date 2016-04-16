package main

import (
	"fmt"
	"sort"
)

type SomeType struct {
	K string
	V int
}

type SomeInterface interface {
	Set(k string, v int)
	Get() (k string, v int)
}

func (t *SomeType) Set(k string, v int) {
	t.K = k
	t.V = v
}
func (t *SomeType) Get() (k string, v int) {
	return t.K, t.V
}

func Modify(s SomeInterface) {
	s.Set("XXXXXXXXXXX", 1000)
}

func Read(s SomeInterface) {
	fmt.Println(s.Get())
}

func main() {
	var a, b []string
	a = make([]string, 10)
	a[0] = "a"
	a[2] = "b"
	fmt.Println(a)
	b = make([]string, 10)
	copy(b, a)
	b[3] = "c"
	b[1] = "e"
	fmt.Println(b)
	fmt.Println(a)
	sort.Strings(b)
	fmt.Println(b)

	y := SomeType{
		"YYYY",
		8,
	}
	fmt.Println(y)
	Modify(&y)
	Read(&y)
	fmt.Println(y)

	var x *SomeType
	fmt.Println(x)
	x = new(SomeType)
	x.K = "ZZZZZ"
	x.V = 20
	fmt.Println(*x)
	Modify(x)
	Read(x)
	fmt.Println(x)

}
