package main

import (
	"flag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yankeguo/bastion/daemon"
	"github.com/yankeguo/bastion/types"
	"os"
	"os/signal"
	"syscall"
)

var (
	dev         bool
	optionsFile string
	options     types.Options
)

func main() {
	var err error

	// init logger
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, NoColor: true})

	// load options from command-line arguments
	flag.StringVar(&optionsFile, "c", "/etc/bastion/bastion.yml", "bastion config file")
	flag.BoolVar(&dev, "dev", false, "enable dev mode")
	flag.Parse()

	// load options file
	log.Info().Str("file", optionsFile).Msg("load options file")
	if options, err = types.LoadOptions(optionsFile); err != nil {
		log.Error().Err(err).Str("file", optionsFile).Msg("failed to load options file")
		os.Exit(1)
		return
	}
	// merge command-line options
	if dev {
		options.Daemon.Dev = dev
	}
	log.Info().Interface("options", options.Daemon).Msg("options file loaded")

	// adjust logger
	if options.Daemon.Dev {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	// create daemon
	d := daemon.New(options.Daemon)

	// run the signalHandler
	go signalHandler(d)

	// run the daemon
	if err = d.Run(); err != nil {
		log.Info().Err(err).Msg("exited")
		os.Exit(1)
		return
	}
}

func signalHandler(d *daemon.Daemon) {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)
	s := <-shutdown
	log.Info().Str("signal", s.String()).Msg("signal received")
	d.Stop()
}
