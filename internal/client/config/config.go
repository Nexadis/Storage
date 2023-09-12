package config

import "flag"

type Config struct {
	GRPC   string
	Secure bool
}

func New() *Config {
	c := &Config{}
	flag.StringVar(&c.GRPC, "g", "127.0.0.1:50051", "Address for connect to gRPC server")
	flag.BoolVar(&c.Secure, "s", false, "Use TLS for connect")
	flag.Parse()
	return c
}
