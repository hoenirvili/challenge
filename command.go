package main

import (
	"strings"
)

const (
	balance = "balance"
	pay     = "pay"
	exit    = "exit"
)

var arrCommandsString = [...]string{balance, pay, exit}

type user struct {
	balance uint64
}

func supportedCommandsComma() string {
	return strings.Join(arrCommandsString[:], ",")
}
