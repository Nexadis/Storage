package config

import "flag"

type Config struct {
	GRPC string
}

func New() *Config {
	c := &Config{}
	flag.StringVar(&c.GRPC, "g", "127.0.0.1:50051", "Address for connect to gRPC server")
	flag.Parse()
	return c
}
