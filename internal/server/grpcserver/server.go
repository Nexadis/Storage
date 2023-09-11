package grpcserver

import (
	"context"
	"log"

	"github.com/Nexadis/Storage/internal/config"
	"github.com/Nexadis/Storage/internal/storage"
	pb "github.com/Nexadis/Storage/protos"
)

type Server struct {
	pb.UnimplementedKeyValueServer
	db storage.Storage
	l  storage.TransactionLogger
}

func (s *Server) Get(ctx context.Context, r *pb.GetRequest) (*pb.GetResponse, error) {
	log.Printf("Received GET key=%v", r.Key)

	value, err := s.db.Get(storage.DefaultUser, r.Key)

	return &pb.GetResponse{Value: value}, err
}

func New(c *config.Config) *Server {
	return nil
}

func NewWithDB(c *config.Config, db storage.Storage, l storage.TransactionLogger) *Server {
	return nil
}
