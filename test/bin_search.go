package main

import (
	"fmt"
)

type Searchable interface {
	Len() int
	Less(int, int) bool
	Equal(int, interface{}) bool
}

type List []int

func (l List) Len() int {
	return len(l)
}

func (l List) Less(first int, second int) bool {
	if l[first] < l[second] {
		return true
	}

	return false
}

func (l List) Equal(index int, item interface{}) bool {
	if value, ok := item.(int); ok {
		if l[index] == value {
			return true
		}
	}

	return false
}

func main() {
	list := []int{1, 2, 3, 5, 9}

	index := binSearch(list, 3)
	fmt.Printf("The index of 3 in the list is %d\n", index)

	index = binSearch(list, 4)
	fmt.Printf("The index of 4 in the list is %d\n", index)
}

func binSearch(list List, item interface{}) int {
	startFlag := 0
	stopFlag := list.Len() - 1
	middleFlag := (startFlag + stopFlag) / 2

	for (!list.Equal(middleFlag, item)) && (startFlag < stopFlag) {
		if list.Less(startFlag, stopFlag) {
			startFlag = middleFlag + 1
		} else {
			stopFlag = middleFlag - 1
		}
		middleFlag = (startFlag + stopFlag) / 2
	}

	if list.Equal(middleFlag, item) {
		return middleFlag
	} else {
		return -1
	}
}
