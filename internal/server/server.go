package server

import (
	"github.com/Nexadis/Storage/internal/config"
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
	return nil
}
