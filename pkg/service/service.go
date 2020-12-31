package service

import (
	"errors"
	"log"

	"go.zoe.im/x"
	"go.zoe.im/x/version"

	"go.zoe.im/eacdn/pkg/config"
	"go.zoe.im/eacdn/pkg/provider"
)

// Service ...
type Service struct {
	Config *config.Config

	// store all supported provider
	providers map[string]provider.Object
}

// Run ...
func (s *Service) Run() error {
	log.Printf("I Welcome to have EaCDN (%s)!\n", version.GitVersion)

	err := s.Init()
	if err != nil {
		log.Println("E Init service error:", err)
		return err
	}

	// start the service
	return x.GraceStart(func(ch x.GraceSignalChan) error {

		// TODO: start other goroutine

		// serve the http and blocking at here
		return s.startHTTP()
	})
}

// Init ...
func (s *Service) Init() error {

	// loads all provider
	for _, o := range s.Config.Providers {
		if _, ok := s.providers[o.Name]; ok {
			log.Println("W Provider with name", o.Name, "exits.")
			continue
		}

		err := o.Init()
		if err != nil {
			log.Println("E Provider", o.Name, "init error:", err)
			continue
		}

		log.Println("I Init provider [", o.Name, "] success")
		s.providers[o.Name] = o
	}

	// TODO: init a internal provider

	// check if we have at least one provider
	if len(s.providers) == 0 {
		return errors.New("Must offer at least 1 provider")
	}

	// start all sites
	for _, st := range s.Config.Sites {
		// choose from providers or use provider.Manager
		p, ok := s.providers[st.Provider]
		if !ok {
			log.Println("E Can not found provider", st.Provider)
			continue
		}

		// TODO: create or update
		err := p.Create(st)
		if err != nil {
			log.Println("E Create site error", err)
			continue
		}
	}

	return nil
}

// New ...
func New() *Service {
	return &Service{
		Config: config.New(),

		providers: make(map[string]provider.Object),
	}
}
