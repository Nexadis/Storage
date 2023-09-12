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
	log.Printf("Received Get %s", r.Key)

	value, err := s.db.Get(storage.DefaultUser, r.Key)

	return &pb.GetResponse{Value: value}, err
}

func (s *Server) Put(ctx context.Context, r *pb.PutRequest) (*pb.PutResponse, error) {
	log.Printf("Received Put %s=%s", r.Key, r.Value)

	err := s.db.Put(storage.DefaultUser, r.Key, r.Value)
	s.l.WritePut(storage.DefaultUser, r.Key, r.Value)

	return &pb.PutResponse{}, err
}

func (s *Server) Delete(ctx context.Context, r *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	log.Printf("Received Delete %s", r.Key)

	err := s.db.Delete(storage.DefaultUser, r.Key)
	s.l.WriteDelete(storage.DefaultUser, r.Key)

	return &pb.DeleteResponse{}, err
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
