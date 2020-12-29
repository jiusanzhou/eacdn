package service

import (
	"log"

	"go.zoe.im/x"
	"go.zoe.im/x/version"

	"go.zoe.im/eacdn/pkg/config"
)

// Service ...
type Service struct {
	Config *config.Config
}

// Run ...
func (s *Service) Run() error {

	// start the service
	return x.GraceStart(func(ch x.GraceSignalChan) error {
		log.Printf("I Welcome to have EaCDN (%s)!\n", version.GitVersion)

		// TODO: start other service in goroutine

		// serve the http and blocking at here
		return s.startHTTP()
	})
}

// New ...
func New() *Service {
	return &Service{
		Config: config.New(),
	}
}
