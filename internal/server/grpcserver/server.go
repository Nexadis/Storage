package grpcserver

import (
	"context"
	"log"
	"net"

	"github.com/Nexadis/Storage/internal/server/config"
	"github.com/Nexadis/Storage/internal/server/httpserver"
	"github.com/Nexadis/Storage/internal/storage"
	pb "github.com/Nexadis/Storage/protos"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedKeyValueServer
	c  *config.Config
	db storage.Storage
	l  storage.TransactionLogger
}

func (s *Server) Get(ctx context.Context, r *pb.GetRequest) (*pb.GetResponse, error) {
	log.Printf("Received GET key=%v", r.Key)

	value, err := s.db.Get(storage.DefaultUser, r.Key)

	return &pb.GetResponse{Value: value}, err
}

func New(c *config.Config) *Server {
	l := httpserver.NewLogger(c)
	db := httpserver.NewStorage(c, l)

	return &Server{
		c:  c,
		db: db,
		l:  l,
	}
}

func NewWithDB(c *config.Config, db storage.Storage, l storage.TransactionLogger) *Server {
	return &Server{
		c:  c,
		db: db,
		l:  l,
	}
}

func (s *Server) Run() error {
	gs := grpc.NewServer()

	pb.RegisterKeyValueServer(gs, s)
	l, err := net.Listen("tcp", s.c.GRPC)
	if err != nil {
		return err
	}
	log.Printf("Run grpc server started on %s", s.c.GRPC)

	return gs.Serve(l)
}
