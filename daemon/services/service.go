package services

import "github.com/yankeguo/bunker/daemon"

// Service service struct
type Service struct {
	// D the daemon instance
	// D contains all resources needed by a Service
	D *daemon.Daemon
}
