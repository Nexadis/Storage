package client

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Nexadis/Storage/internal/client/config"
	pb "github.com/Nexadis/Storage/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	c    *config.Config
	conn grpc.ClientConnInterface
}

func New(c *config.Config) *Client {
	return &Client{
		c: c,
	}
}

func (c *Client) Open() error {
	creds := insecure.NewCredentials()
	conn, err := grpc.Dial(c.c.GRPC, grpc.WithBlock(), grpc.WithTransportCredentials(creds))
	if err != nil {
		return err
	}
	c.conn = conn
	log.Printf("Connect to %s", c.c.GRPC)
	return nil
}

func (c *Client) Close() {
}

func (c *Client) DoCmd(args []string) error {
	client := pb.NewKeyValueClient(c.conn)
	if len(args) < 1 {
		return errors.New("no command")
	}
	cmd := args[0]
	switch cmd {
	case "get":
		if err := checkLen(args, 2); err != nil {
			return err
		}
		key := args[1]
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := client.Get(ctx, &pb.GetRequest{Key: key})
		if err != nil {
			return err
		}
		log.Printf("Get %s=%s", key, r.Value)

	case "put":
		if err := checkLen(args, 3); err != nil {
			return err
		}
		key := args[1]
		value := args[2]
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_, err := client.Put(ctx, &pb.PutRequest{Key: key, Value: value})
		if err != nil {
			return err
		}
		log.Printf("Put %s=%s", key, value)

	case "delete":
		if err := checkLen(args, 2); err != nil {
			return err
		}
		key := args[1]
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_, err := client.Delete(ctx, &pb.DeleteRequest{Key: key})
		if err != nil {
			return err
		}
		log.Printf("Delete %s", key)
	}

	return nil
}

func checkLen(args []string, size int) error {
	if len(args) < size {
		return fmt.Errorf("invalid command, need %d args", size)
	}
	return nil
}
