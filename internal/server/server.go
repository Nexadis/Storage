package server

import (
	"log"
	"sync"

	"github.com/Nexadis/Storage/internal/config"
	"github.com/Nexadis/Storage/internal/storage"
	"github.com/labstack/echo/v4"
)

type HTTPServer struct {
	*echo.Echo
	c *config.Config
	s storage.Storage
}

func New(c *config.Config) *HTTPServer {
	s := storage.New(100)
	e := echo.New()
	hs := &HTTPServer{
		e,
		c,
		s,
	}
	return hs
}

func (hs *HTTPServer) Run() error {
	var wg sync.WaitGroup
	errs := make(chan error)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(errs)
		hs.MountHandlers()
		err := hs.Start(hs.c.Addr)
		errs <- err
	}()
	wg.Wait()
	for err := range errs {
		log.Printf("Err: %v", err)
	}
	return nil
}
