package httpserver

import (
	"log"

	"github.com/Nexadis/Storage/internal/server/config"
	"github.com/Nexadis/Storage/internal/storage"
	"github.com/Nexadis/Storage/internal/storage/mem"
	"github.com/Nexadis/Storage/internal/storage/pg"
	"github.com/labstack/echo/v4"
)

type HTTPServer struct {
	*echo.Echo
	Config *config.Config
	db     storage.Storage
	l      storage.TransactionLogger
}

func New(c *config.Config) *HTTPServer {
	l := NewLogger(c)
	s := NewStorage(c, l)
	return NewWithDB(c, s, l)
}

func NewWithDB(c *config.Config, db storage.Storage, l storage.TransactionLogger) *HTTPServer {
	e := echo.New()
	hs := &HTTPServer{
		e,
		c,
		db,
		l,
	}
	return hs
}

func NewLogger(c *config.Config) storage.TransactionLogger {
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
	return l
}

func NewStorage(c *config.Config, l storage.TransactionLogger) storage.Storage {
	s := mem.New(100)
	if l != nil {
		err := storage.RestoreTransactions(s, l)
		if err != nil {
			log.Fatal(err)
		}
	}

	return s
}

func (hs *HTTPServer) Run() error {
	hs.MountHandlers()
	return hs.Start(hs.Config.Addr)
}
