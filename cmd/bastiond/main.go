package main

import (
	"flag"
	"log"

	"github.com/yankeguo/bastion/daemon"
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
	d := daemon.New(options.Daemon)

	// run the signalHandler
	go signalHandler(d)

	// run the daemon
	if err = d.Run(); err != nil {
		log.Println("daemon exited", err)
		os.Exit(1)
		return
	}
}

func signalHandler(d *daemon.Daemon) {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)
	s := <-shutdown
	log.Println("signal:", s)
	d.Stop()
}
