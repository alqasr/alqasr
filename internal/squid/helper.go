package squid

import (
	"fmt"
)

const (
	OK  = "OK\n"
	ERR = "ERR\n"
)

func SendOK() {
	fmt.Print(OK)
}

func SendERR() {
	fmt.Print(ERR)
}
