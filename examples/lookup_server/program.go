package main

import (
	"github.com/wondywang/rpclookup/core/lg"
	"github.com/wondywang/rpclookup/lookupd"
	"os"
	"flag"
	"runtime"
	"path/filepath"
	"log"
	"github.com/judwhite/go-svc/svc"
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
	logger.SetFlags(log.Ltime | log.Lshortfile)
	logger.SetOutput(os.Stdout)
	// logger.SetLevel(lg.LogLevelDebug)

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

	err = this.lookupd.Init()
	if err != nil {
		logger.Error("init lookup failed! %s", err.Error())
		return err
	}
	return nil
}

func (this *Program) Start() error {
	logger.Debug("Program start...")
	this.lookupd.Main()
	return nil
}

func (this *Program) Stop() error {
	this.lookupd.Exit()
	return nil
}