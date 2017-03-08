package main

import "fmt"

// TODO: logging in file?

const (
	DBGV = 0 // iota
	DBGI = 1
	DBGE = 2
)

type Debug int
var Dbg Debug = DBGI

func (d Debug) E(s string, a ...interface{}) {
	if d <= DBGE {
		fmt.Printf(s, a...)
	}
}

func (d Debug) I(s string, a ...interface{}) {
	if d <= DBGI {
		fmt.Printf(s, a...)
	}
}

func (d Debug) V(s string, a ...interface{}) {
	if d <= DBGV {
		fmt.Printf(s, a...)
	}
}
