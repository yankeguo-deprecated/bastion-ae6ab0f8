package main

import (
	"flag"
	"log"

	"github.com/yankeguo/bastion/sshd"
	"github.com/yankeguo/bastion/types"
	"github.com/yankeguo/bastion/utils"
	"os"
	"os/signal"
	"syscall"
)

var (
	optionsFile string
	options     types.Options
)

func main() {
	var err error

	// load options from command-line arguments
	flag.StringVar(&optionsFile, "c", "/etc/bastion/bastion.yml", "bastion config file")
	flag.Parse()
	log.Println("loading", optionsFile)
	if options, err = utils.LoadOptions(optionsFile); err != nil {
		log.Println("failed to load config file", err)
		os.Exit(1)
		return
	} else {
		log.Println(options.Daemon)
	}

	// create daemon
	d := sshd.New(options.SSHD)

	// run the signalHandler
	go signalHandler(d)

	// run the sshd
	if err = d.Run(); err != nil {
		log.Println("sshd exited", err)
		os.Exit(1)
		return
	}
}

func signalHandler(d *sshd.SSHD) {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)
	<-shutdown
	d.Shutdown()
}
