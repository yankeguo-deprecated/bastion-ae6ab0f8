package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {
	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)
	<-shutdown
}
