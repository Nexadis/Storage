package server

import (
	"log"
	"sync"

	"github.com/Nexadis/Storage/internal/config"
	"github.com/Nexadis/Storage/internal/storage"
	"github.com/Nexadis/Storage/internal/storage/mem"
	"github.com/Nexadis/Storage/internal/storage/pg"
	"github.com/labstack/echo/v4"
)

type HTTPServer struct {
	*echo.Echo
	c *config.Config
	s storage.Storage
	l storage.TransactionLogger
}

func New(c *config.Config) *HTTPServer {
	var l storage.TransactionLogger
	var err error
	if c.DBURI == "" {
		l, err = mem.NewFileTransactionLogger(c.FileSave)
		log.Print("Use in file transactions")
	} else {
		l, err = pg.NewPostgreTransactionLogger(c.DBURI)
		log.Print("Use db transactions")
	}
	if err != nil {
		log.Fatal(err)
	}
	l.Run()
	s := mem.New(100)
	err = storage.RestoreTransactions(s, l)
	if err != nil {
		log.Fatal(err)
	}
	e := echo.New()
	hs := &HTTPServer{
		e,
		c,
		s,
		l,
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
