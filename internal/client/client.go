package client

import (
	"context"
	"log"
	"time"

	"github.com/Nexadis/Storage/internal/client/config"
	pb "github.com/Nexadis/Storage/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	c *config.Config
}

func New(c *config.Config) *Client {
	return &Client{
		c: c,
	}
}

func (c *Client) Run() error {
	var creds credentials.TransportCredentials
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if !c.c.Secure {
		creds = insecure.NewCredentials()
	}
	conn, err := grpc.DialContext(ctx, c.c.GRPC, grpc.WithBlock(), grpc.WithTransportCredentials(creds))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := pb.NewKeyValueClient(conn)
	log.Printf("Connect to %s", c.c.GRPC)
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.Get(ctx, &pb.GetRequest{Key: "some"})
	if err != nil {
		return err
	}
	log.Printf("Value=%s", r.Value)
	return nil
}
