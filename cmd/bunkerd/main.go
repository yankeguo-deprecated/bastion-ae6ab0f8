package main

import (
	"flag"
	"log"

	"github.com/yankeguo/bunker/types"
	"github.com/yankeguo/bunker/utils"
	"github.com/yankeguo/bunker/db"
)

var (
	optionsFile string
	options     types.Options
)

func main() {
	var err error
	defer utils.DoExit()

	// load options from command-line arguments
	flag.StringVar(&optionsFile, "c", "/etc/bunker/bunker.yml", "bunker config file")
	flag.Parse()
	log.Println("loading", optionsFile)
	if options, err = utils.LoadOptions(optionsFile); err != nil {
		log.Println("failed to load config file", err)
		utils.WillExit(1)
		return
	} else {
		utils.PrintOptions(options)
	}

	// open database
	var d *db.DB
	if d, err = db.Open(options.Daemon.DB); err != nil {
		log.Println("failed to open db file", err)
		utils.WillExit(1)
		return
	}

	// migrate database
	if err = d.Migrate(); err != nil {
		log.Println("failed to migrate", err)
		utils.WillExit(1)
		return
	}
}
