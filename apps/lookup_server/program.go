package main

import (
	"os"
	"flag"
	"runtime"
	"path/filepath"
	"log"
	"github.com/judwhite/go-svc/svc"
	"core/lg"
	"lookupd"
)

var logger *lg.Logger

type Program struct {
	lookupd *lookupd.Lookupd
}

func (this *Program) Init(env svc.Environment) error {
	if env.IsWindowsService() {
		dir := filepath.Dir(os.Args[0])
		return os.Chdir(dir)
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	
	logger = &lg.Logger{}
	logger.SetFlags(log.Ldate | log.Lshortfile)
	logger.SetOutput(os.Stdout)

	var err error
	port := "16688"
	handler := MsgHandler{
		logger: logger,
	}
	this.lookupd, err = lookupd.NewLookupd(port, logger, handler)
	if err != nil {
		logger.Error("new lookupd failed! %s", err.Error())
		return err
	}

	this.lookupd.Init()
	return nil
}

func (this *Program) Start() error {
	var err error

	go func() {
		err = this.lookupd.Accept()
		if err != nil {
			logger.Error("run net accept failed!", err.Error())
			os.Exit(1)
		}
	}()

	logger.Debug("Program start...")
	return nil
}

func (this *Program) Stop() error {
	this.lookupd.Exit()
	return nil
}