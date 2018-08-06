package main

import (
	"flag"
	"github.com/yankeguo/bastion/types"
	"github.com/yankeguo/bastion/utils"
	"github.com/yankeguo/bastion/web"
	"log"
	"os"
)

var (
	optionsFile string
	dev         bool
	options     types.Options
)

func main() {
	var err error

	// load options from command-line arguments
	flag.StringVar(&optionsFile, "c", "/etc/bastion/bastion.yml", "bastion config file")
	flag.BoolVar(&dev, "dev", false, "dev flag, overriding web.dev")
	flag.Parse()
	log.Println("loading", optionsFile)
	if options, err = utils.LoadOptions(optionsFile); err != nil {
		log.Println("failed to load config file", err)
		os.Exit(1)
		return
	} else {
		log.Println(options.Daemon)
	}

	if dev {
		options.Web.Dev = true
	}

	s := web.NewServer(options.Web)
	log.Fatal(s.ListenAndServe())
}
