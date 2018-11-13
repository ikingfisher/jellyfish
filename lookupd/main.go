package main

import (
	"fmt"
	"syscall"
	"github.com/judwhite/go-svc/svc"
)

var logger Logger

func main() {
	prg := &Program{}
	if err := svc.Run(prg, syscall.SIGINT, syscall.SIGTERM); err != nil {
		fmt.Println("program exit! err:%s", err.Error())
	}
}