package main

import (
	"fmt"
	"unsafe"
)

type slice struct {
	array unsafe.Pointer
	len   int
	cap   int
}

func main() {
	x := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	s := x
	s[5] = 100
	fmt.Println(x)
}
