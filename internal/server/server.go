package server

import (
	"log"
	"sync"

	"github.com/Nexadis/Storage/internal/server/config"
	"github.com/Nexadis/Storage/internal/server/grpcserver"
	"github.com/Nexadis/Storage/internal/server/httpserver"
)

type Server struct {
	HTTPServer *httpserver.HTTPServer
	GRPCServer *grpcserver.Server
}

func New(c *config.Config) *Server {
	l := httpserver.NewLogger(c)
	db := httpserver.NewStorage(c, l)
	hs := httpserver.NewWithDB(c, db, l)
	gs := grpcserver.NewWithDB(c, db, l)

	return &Server{
		hs,
		gs,
	}
}

func (s *Server) Run() error {
	var wg sync.WaitGroup
	errs := make(chan error)
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := s.HTTPServer.Run()
		errs <- err
	}()
	go func() {
		defer wg.Done()
		err := s.GRPCServer.Run()
		errs <- err
	}()
	wg.Wait()
	close(errs)
	for err := range errs {
		log.Printf("Err: %v", err)
	}
	return nil
}
