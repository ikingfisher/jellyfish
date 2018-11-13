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

type Program struct {
	logger *lg.Logger
	lookupd *lookupd.Lookupd
}

func (this *Program) Init(env svc.Environment) error {
	if env.IsWindowsService() {
		dir := filepath.Dir(os.Args[0])
		return os.Chdir(dir)
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	
	this.logger = &lg.Logger{}
	this.logger.SetFlags(log.Ldate | log.Lshortfile)
	this.logger.SetOutput(os.Stdout)

	var err error
	port := "16688"
	this.lookupd, err = lookupd.NewLookupd(port, this.logger)
	if err != nil {
		this.logger.Error("new lookupd failed! %s", err.Error())
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
			this.logger.Error("run net accept failed!", err.Error())
			os.Exit(1)
		}
	}()

	go this.lookupd.HeartBeat()

	this.logger.Debug("Program start...")
	return nil
}

func (this *Program) Stop() error {
	this.lookupd.Exit()
	return nil
}